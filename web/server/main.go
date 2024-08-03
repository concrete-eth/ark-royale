package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/google/uuid"
)

type LobbyEntry struct {
	Players      []common.Address
	LastActivity time.Time
}

func (l *LobbyEntry) GetPlayers() []common.Address {
	players := make([]common.Address, len(l.Players))
	copy(players, l.Players)
	return players
}

type LobbyRegistry struct {
	Lobbies map[string]*LobbyEntry
	lock    sync.RWMutex
}

// NewLobbyRegistry creates a new lobby registry.
func NewLobbyRegistry() *LobbyRegistry {
	return &LobbyRegistry{
		Lobbies: make(map[string]*LobbyEntry),
	}
}

// NewLobby creates a new lobby and returns its ID.
func (r *LobbyRegistry) NewLobby() string {
	r.lock.Lock()
	defer r.lock.Unlock()

	uuid := uuid.New().String()
	r.Lobbies[uuid] = &LobbyEntry{
		Players:      make([]common.Address, 0),
		LastActivity: time.Now(),
	}
	return uuid
}

// AddPlayer adds a player to the specified lobby.
// Returns an error if the lobby does not exist.
func (r *LobbyRegistry) AddPlayer(lobbyId string, player common.Address) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	lobby, ok := r.Lobbies[lobbyId]
	if !ok {
		return errors.New("lobby does not exist")
	}
	lobby.Players = append(lobby.Players, player)
	lobby.LastActivity = time.Now()
	return nil
}

// GetLobby retrieves a lobby by its ID. It returns false if the lobby does not exist.
func (r *LobbyRegistry) GetPlayerAddresses(lobbyId string) ([]common.Address, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	lobby, ok := r.Lobbies[lobbyId]
	if !ok {
		return nil, false
	}
	return lobby.GetPlayers(), true
}

func (r *LobbyRegistry) Cleanup(inactivityDuration time.Duration) {
	r.lock.Lock()
	defer r.lock.Unlock()

	now := time.Now()
	for id, lobby := range r.Lobbies {
		if now.Sub(lobby.LastActivity) > inactivityDuration {
			delete(r.Lobbies, id)
		}
	}
}

func (r *LobbyRegistry) StartCleanupScheduler(ctx context.Context, inactivityDuration time.Duration, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				r.Cleanup(inactivityDuration)
			case <-ctx.Done():
				return
			}
		}
	}()
}

type Faucet struct {
	auth   *bind.TransactOpts
	client *ethclient.Client
	lock   sync.Mutex
}

func NewFaucet(client *ethclient.Client, auth *bind.TransactOpts) *Faucet {
	return &Faucet{
		auth:   auth,
		client: client,
	}
}

func (f *Faucet) Drip(to common.Address, amount *big.Int) (common.Hash, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	nonce, err := f.client.PendingNonceAt(context.Background(), f.auth.From)
	if err != nil {
		return common.Hash{}, err
	}
	gasPrice, err := f.client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}
	tx := types.NewTransaction(nonce, to, amount, 21000, gasPrice, nil)
	signedTx, err := f.auth.Signer(f.auth.From, tx)
	if err != nil {
		return common.Hash{}, err
	}
	err = f.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, err
	}
	return signedTx.Hash(), nil
}

type ServerConfig struct {
	BaseURL               string
	PortStr               string
	RpcURL                string
	GameFactoryAddressHex string
	FaucetPrivateKeyHex   string
	TemplateDir           string
	StaticDir             string
	WasmURL               string
}

func (c *ServerConfig) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("BaseURL is required")
	}
	if c.PortStr == "" {
		return fmt.Errorf("PortStr is required")
	}
	if c.RpcURL == "" {
		return fmt.Errorf("RpcURL is required")
	}
	if c.GameFactoryAddressHex == "" {
		return fmt.Errorf("GameFactoryAddressHex is required")
	}
	if c.FaucetPrivateKeyHex == "" {
		return fmt.Errorf("FaucetPrivateKeyHex is required")
	}
	if c.TemplateDir == "" {
		return fmt.Errorf("TemplateDir is required")
	}
	if c.StaticDir == "" {
		return fmt.Errorf("StaticDir is required")
	}
	if c.WasmURL == "" {
		return fmt.Errorf("WasmUrl is required")
	}
	return nil
}

func loadTemplates(dir string) (*template.Template, error) {
	funcs := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}
	pattern := filepath.Join(dir, "*.html")
	templates, err := template.New("").Funcs(funcs).ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}
	return templates, nil
}

func loadAuth(config ServerConfig, rpc *ethclient.Client) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(config.FaucetPrivateKeyHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %w", err)
	}
	chainId, err := rpc.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}
	return auth, nil
}

type Server struct {
	config    ServerConfig
	templates *template.Template
	rpc       *ethclient.Client
	auth      *bind.TransactOpts
	faucet    *Faucet
	lobbies   *LobbyRegistry
}

func NewServer(config ServerConfig, rpc *ethclient.Client) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	templates, err := loadTemplates(config.TemplateDir)
	if err != nil {
		return nil, err
	}
	auth, err := loadAuth(config, rpc)
	if err != nil {
		return nil, err
	}
	return &Server{
		config:    config,
		templates: templates,
		rpc:       rpc,
		auth:      auth,
		faucet:    NewFaucet(rpc, auth),
		lobbies:   NewLobbyRegistry(),
	}, nil
}

func (s *Server) Run() error {
	// Start the lobby cleanup scheduler
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.lobbies.StartCleanupScheduler(ctx, 5*time.Minute, 1*time.Minute)

	// Create specific ServeMux instances
	staticMux := http.NewServeMux()
	actionsMux := http.NewServeMux()
	pagesMux := http.NewServeMux()
	dataMux := http.NewServeMux()

	// Static files
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(s.config.StaticDir)))
	staticMux.Handle("/static/", fs)

	// Actions
	actionsMux.HandleFunc("/actions/new-lobby", s.handleNewLobby)
	actionsMux.HandleFunc("/actions/join-lobby", s.handleJoinLobby)
	actionsMux.HandleFunc("/actions/request-drip", s.handleRequestDrip)

	// Pages
	pagesMux.HandleFunc("/game/", func(w http.ResponseWriter, r *http.Request) {
		s.renderTemplate(w, "game.html", nil)
	})
	pagesMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "" {
			http.NotFound(w, r)
			return
		}
		s.renderTemplate(w, "dashboard.html", nil)
	})
	pagesMux.HandleFunc("/lobby/", func(w http.ResponseWriter, r *http.Request) {
		segments := strings.Split(r.URL.Path[len("/lobby/"):], "/")
		if len(segments) != 1 {
			http.NotFound(w, r)
			return
		}
		lobbyId := segments[0]
		if _, ok := s.lobbies.GetPlayerAddresses(lobbyId); !ok {
			http.NotFound(w, r)
			return
		}
		s.renderTemplate(w, "lobby.html", nil)
	})

	// Raw data
	dataMux.HandleFunc("/data/lobby/", func(w http.ResponseWriter, r *http.Request) {
		segments := strings.Split(r.URL.Path[len("/data/lobby/"):], "/")
		if len(segments) != 1 {
			http.NotFound(w, r)
			return
		}
		lobbyId := segments[0]
		players, ok := s.lobbies.GetPlayerAddresses(lobbyId)
		if !ok {
			http.NotFound(w, r)
			return
		}
		lobbyData := struct {
			Id      string           `json:"id"`
			Players []common.Address `json:"players"`
		}{
			Id:      lobbyId,
			Players: players,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(lobbyData); err != nil {
			log.Crit("Error encoding JSON response", "err", err)
		}
	})

	// Main ServeMux to delegate to specific muxes
	mainMux := http.NewServeMux()
	mainMux.Handle("/static/", staticMux)
	mainMux.Handle("/actions/", actionsMux)
	mainMux.Handle("/data/", dataMux)
	// Since pages are potentially at the root, delegate last
	mainMux.Handle("/", pagesMux)

	log.Info(
		"Server started",
		"url", s.config.BaseURL,
		"port", s.config.PortStr,
		"rpc", s.config.RpcURL,
		"factory", s.config.GameFactoryAddressHex,
		"templates", s.config.TemplateDir,
		"static", s.config.StaticDir,
		"wasm", s.config.WasmURL,
	)

	// Start the server with the main ServeMux
	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.config.PortStr), mainMux); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func (s *Server) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	d := map[string]interface{}{
		"BaseURL":        template.JS(s.config.BaseURL),
		"WasmURL":        template.JS(s.config.WasmURL),
		"FactoryAddress": template.JS(s.config.GameFactoryAddressHex),
		"RpcURL":         template.JS(s.config.RpcURL),
		"Data":           data,
	}
	if err := s.templates.ExecuteTemplate(w, name, d); err != nil {
		log.Error("Error executing template", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type newLobbyRequestData struct {
	Players []string `json:"players"`
}

type newLobbyResponseData struct {
	LobbyId string `json:"lobbyId"`
}

func (s *Server) handleNewLobby(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData newLobbyRequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Error("Error decoding JSON request", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobbyId := s.lobbies.NewLobby()
	for _, player := range requestData.Players {
		s.lobbies.AddPlayer(lobbyId, common.HexToAddress(player))
	}

	responseData := newLobbyResponseData{
		LobbyId: lobbyId,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Crit("Error encoding JSON response", "err", err)
	}
}

type joinLobbyRequestData struct {
	LobbyId       string         `json:"lobbyId"`
	PlayerAddress common.Address `json:"playerAddress"`
}

type joinLobbyResponseData struct{}

func (s *Server) handleJoinLobby(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// If the request method is not POST, send a 405 Method Not Allowed status.
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData joinLobbyRequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Error("Error decoding JSON request", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lobbyId := requestData.LobbyId
	playerAddress := requestData.PlayerAddress
	s.lobbies.AddPlayer(lobbyId, playerAddress)

	responseData := joinLobbyResponseData{}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Crit("Error encoding JSON response", "err", err)
	}
}

type requestDripRequestData struct {
	Address string `json:"address"`
}

type requestDripResponseData struct {
	TxHash common.Hash `json:"txHash"`
	Ok     bool        `json:"ok"`
}

func (s *Server) handleRequestDrip(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData requestDripRequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Error("Error decoding JSON request", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	address := common.HexToAddress(requestData.Address)
	balance, err := s.rpc.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Error("Error getting balance", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := requestDripResponseData{}

	oneEther := new(big.Int).SetUint64(1e18)
	if balance.Cmp(oneEther) < 0 {
		txHash, err := s.faucet.Drip(address, new(big.Int).Sub(oneEther, balance))
		if err != nil {
			log.Error("Error dripping", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		responseData.Ok = true
		responseData.TxHash = txHash
	} else {
		responseData.Ok = false
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Crit("Error encoding JSON response", "err", err)
	}
}
