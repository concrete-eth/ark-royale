//go:build js
// +build js

package main

// Click outside 0,0
// Ghosts

import (
	"context"
	"fmt"
	"image"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/concrete-eth/archetype/arch"
	"github.com/concrete-eth/archetype/rpc"
	snapshot_utils "github.com/concrete-eth/archetype/snapshot/utils"
	"github.com/concrete-eth/archetype/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/lib"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"

	"github.com/concrete-eth/archetype/kvstore"
	"github.com/concrete-eth/ark-royale/client/core"
	"github.com/concrete-eth/ark-royale/client/game"
	tables_contract "github.com/concrete-eth/ark-royale/gogen/abigen/tables"
	"github.com/concrete-eth/ark-royale/gogen/archmod"
	"github.com/hajimehoshi/ebiten/v2"

	snapshot_types "github.com/concrete-eth/archetype/snapshot/types"
	game_contract "github.com/concrete-eth/ark-royale/gogen/abigen/game"
)

var (
	clientConfig = core.ClientConfig{
		ScreenSize: image.Point{900, 600},
	}
)

type URLParams struct {
	GameAddress common.Address
	WsURL       string
	Interpolate bool
	Debug       bool
	BlockTime   time.Duration
	Delay       time.Duration
}

func getURLParams() (URLParams, error) {
	window := js.Global()
	href := window.Get("location").Get("href").String()

	parsedUrl, err := url.Parse(href)
	if err != nil {
		return URLParams{}, err
	}

	path := parsedUrl.Path
	segments := strings.Split(path, "/")
	if len(segments) != 3 {
		return URLParams{}, fmt.Errorf("invalid path")
	}
	gameAddressHex := segments[2]
	if gameAddressHex == "" {
		return URLParams{}, fmt.Errorf("address parameter is required")
	}
	gameAddress := common.HexToAddress(gameAddressHex)

	queryParams := parsedUrl.Query()

	var paramValue string
	paramValue = queryParams.Get("ws")
	if paramValue == "" {
		return URLParams{}, fmt.Errorf("ws parameter is required")
	}
	wsURL := paramValue

	paramValue = queryParams.Get("interpolate")
	interpolate := strings.ToLower(paramValue) != "false"

	paramValue = queryParams.Get("debug")
	debug := strings.ToLower(paramValue) == "true"

	paramValue = queryParams.Get("blockTime")
	var blockTimeDuration time.Duration
	if paramValue == "" {
		blockTimeDuration = 1 * time.Second
	} else {
		blockTime, err := strconv.Atoi(paramValue)
		if err != nil {
			return URLParams{}, fmt.Errorf("blockTime parameter is required")
		}
		blockTimeDuration = time.Duration(blockTime) * time.Millisecond
	}

	paramValue = queryParams.Get("delay")
	var delayDuration time.Duration
	if paramValue == "" {
		delayDuration = 250 * time.Millisecond
	} else {
		delay, err := strconv.Atoi(paramValue)
		if err != nil {
			return URLParams{}, fmt.Errorf("delay parameter is required")
		}
		delayDuration = time.Duration(delay) * time.Millisecond
		if delayDuration < 0 {
			delayDuration = 0
		} else if delayDuration > blockTimeDuration {
			delayDuration = blockTimeDuration
		}
	}

	return URLParams{
		GameAddress: gameAddress,
		WsURL:       wsURL,
		Interpolate: interpolate,
		Debug:       debug,
		BlockTime:   blockTimeDuration,
		Delay:       delayDuration,
	}, nil
}

func getPrivateKey() (string, error) {
	privateKeyHex := getLocalStorage("burnerKey")
	if privateKeyHex == "" {
		return "", fmt.Errorf("private key is required")
	}
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	return privateKeyHex, nil
}

func setLocalStorage(key string, value string) {
	window := js.Global()
	window.Get("localStorage").Set(key, value)
}

func getLocalStorage(key string) string {
	window := js.Global()
	return window.Get("localStorage").Get(key).String()
}

func newGameScreenSize() image.Point {
	window := js.Global()
	gameScreenSize := image.Point{
		X: window.Get("innerWidth").Int(),
		Y: window.Get("innerHeight").Int(),
	}
	if gameScreenSize.X > 2*gameScreenSize.Y {
		// Width will be at most 2x height
		gameScreenSize.X = gameScreenSize.Y * 2
	}
	return gameScreenSize
}

func setLoadStatus(status string) {
	window := js.Global()
	if element := window.Get("document").Call("getElementById", "loader-status-main"); !element.IsNull() {
		element.Set("innerText", status)
	}
}

func hideLoadStatus() {
	window := js.Global()
	if element := window.Get("document").Call("getElementById", "loader-container-main"); !element.IsNull() {
		element.Get("parentNode").Call("removeChild", element)
	}
}

func showErrorScreen(err error) {
	body := js.Global().Get("document").Call("getElementsByTagName", "body").Index(0)
	body.Set("innerHTML", `
        <div id="error-container-main" class="error-container">
            <div id="error-status-main" class="error-status">
				<h1>Error</h1>
				<p>`+err.Error()+`</p>
			</div>
        </div>
    `)
}

func lastCriticalErrorTime() time.Time {
	lastCriticalErrorTimeStr := getLocalStorage("lastCriticalErrorTime")
	if lastCriticalErrorTimeStr == "" {
		return time.Time{}
	}
	lastCriticalErrorTime, err := time.Parse(time.RFC3339, lastCriticalErrorTimeStr)
	if err != nil {
		return time.Time{}
	}
	return lastCriticalErrorTime
}

func setLastCriticalErrorTime(t time.Time) {
	setLocalStorage("lastCriticalErrorTime", t.Format(time.RFC3339))
}

func logCrit(err error) {
	showErrorScreen(err)
	log.Error(err.Error())
	// Reload page if there was no critical error in the last 10 seconds
	reload := time.Since(lastCriticalErrorTime()) > 10*time.Second
	setLastCriticalErrorTime(time.Now())
	if reload {
		js.Global().Get("location").Call("reload")
	}
	os.Exit(0)
}

func runGameClient(clientConfig core.ClientConfig, params URLParams, privateKeyHex string) {

	// Connect to rpc
	setLoadStatus("Connecting...")
	rpcClient, err := ethclient.Dial(params.WsURL)
	if err != nil {
		logCrit(fmt.Errorf("Failed to connect to RPC: %v", err))
	}
	log.Info("Connected to RPC", "url", params.WsURL)

	// Create signer
	chainId, err := rpcClient.ChainID(context.Background())
	if err != nil {
		logCrit(fmt.Errorf("Failed to get chain ID: %v", err))
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		logCrit(fmt.Errorf("Failed to parse private key: %v", err))
	}
	opts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		logCrit(fmt.Errorf("Failed to create transactor: %v", err))
	}
	log.Info("Loaded burner wallet", "address", opts.From)
	senderAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Create instance of game and core contracts
	gameContract, err := game_contract.NewContract(params.GameAddress, rpcClient)
	if err != nil {
		logCrit(fmt.Errorf("Failed to load game contract: %v", err))
	}
	coreAddress, err := gameContract.Proxy(nil)
	if err != nil {
		logCrit(fmt.Errorf("Failed to get core address: %v", err))
	}
	tablesContract, err := tables_contract.NewContract(params.GameAddress, rpcClient)
	if err != nil {
		logCrit(fmt.Errorf("Failed to load core contract: %v", err))
	}
	metaRow, err := tablesContract.GetMetaRow(nil)
	if err != nil {
		logCrit(fmt.Errorf("Failed to get meta row: %v", err))
	}

	// Get player ID
	playerId, err := gameContract.GetPlayerId(nil, senderAddress)
	if err != nil {
		logCrit(fmt.Errorf("Failed to get player ID: %v", err))
	}
	log.Info("Fetched player ID", "playerId", playerId)

	// Connect to PostHog
	postHog := NewPostHog("https://eu.posthog.com/", "phc_LsLNl9nWm168DyAPmzYxLCG42C6dzWXGdSgQEmzfHx0", senderAddress.Hex())
	postHog.Start()
	defer postHog.Stop()

	// Check if there is a snapshot available
	var mustLoadSnapshot bool
	var mostRecentSnapshotMetadata snapshot_types.SnapshotMetadataWithStatus
	err = rpcClient.Client().Call(&mostRecentSnapshotMetadata, "arch_last", coreAddress)
	if err != nil {
		log.Error("Failed to get most recent snapshot metadata", "err", err)
	} else {
		mustLoadSnapshot = mostRecentSnapshotMetadata.Status == snapshot_types.SnapshotStatus_Done
	}

	// Instantiate the kv store
	var kv lib.KeyValueStore
	var startingBlockNumber uint64

	var loadSnapshotDuration time.Duration
	var fetchSnapshotDuration time.Duration
	var blockSyncDuration time.Duration

	if mustLoadSnapshot {
		log.Info("Syncing from snapshot", "blockNumber", mostRecentSnapshotMetadata.BlockNumber.Uint64(), "blockHash", mostRecentSnapshotMetadata.BlockHash)
		setLoadStatus("Loading snapshot...")

		startTime := time.Now()
		var mostRecentSnapshot snapshot_types.SnapshotResponse
		err = rpcClient.Client().Call(&mostRecentSnapshot, "arch_get", coreAddress, mostRecentSnapshotMetadata.BlockHash)
		if err != nil {
			logCrit(fmt.Errorf("Failed to get most recent snapshot: %v", err))
		}
		fetchSnapshotDuration = time.Since(startTime)
		snapshotSize := len(mostRecentSnapshot.Storage)
		log.Debug("Fetched snapshot", "time", fetchSnapshotDuration, "size", snapshotSize)
		postHog.Capture("SnapshotFetched", map[string]interface{}{
			"playerAddress": senderAddress.Hex(),
			"gameAddress":   params.GameAddress.Hex(),
			"coreAddress":   coreAddress.Hex(),

			"blockNumber": mostRecentSnapshotMetadata.BlockNumber.Uint64(),
			"blockHash":   mostRecentSnapshotMetadata.BlockHash.Hex(),

			"size": snapshotSize,

			"duration": fetchSnapshotDuration.Seconds(),
		})

		startTime = time.Now()
		_kv := kvstore.NewHashedMemoryKeyValueStore()
		rawBlob, err := snapshot_utils.Decompress(mostRecentSnapshot.Storage)
		if err != nil {
			logCrit(fmt.Errorf("Failed to decompress snapshot: %v", err))
		}
		storageIt := snapshot_utils.BlobToStorageIt(rawBlob)
		for storageIt.Next() {
			key := storageIt.Hash()
			enc := storageIt.Slot()
			value, err := snapshot_utils.DecodeSnapshotSlot(enc)
			if err != nil {
				logCrit(fmt.Errorf("Failed to decode snapshot slot: %v", err))
			}
			_kv.SetByKeyHash(key, value)
		}
		loadSnapshotDuration = time.Since(startTime)
		log.Debug("Loaded snapshot into KV store", "time", loadSnapshotDuration)
		postHog.Capture("SnapshotLoaded", map[string]interface{}{
			"playerAddress": senderAddress.Hex(),
			"gameAddress":   params.GameAddress.Hex(),
			"coreAddress":   coreAddress.Hex(),

			"blockNumber": mostRecentSnapshotMetadata.BlockNumber.Uint64(),
			"blockHash":   mostRecentSnapshotMetadata.BlockHash.Hex(),

			"size": snapshotSize,

			"duration": loadSnapshotDuration.Seconds(),
		})

		kv = _kv
		startingBlockNumber = mostRecentSnapshotMetadata.BlockNumber.Uint64() + 1
		log.Info("Synced from snapshot", "blockNumber", startingBlockNumber, "blockHash", mostRecentSnapshotMetadata.BlockHash)
	} else {
		kv = kvstore.NewMemoryKeyValueStore()
		startingBlockNumber = uint64(metaRow.CreationBlockNumber)
	}

	// Create chain IO
	var (
		schemas   = arch.ArchSchemas{Actions: archmod.ActionSchemas, Tables: archmod.TableSchemas}
		blockTime = params.BlockTime
	)

	io := rpc.NewIO(rpcClient, blockTime, schemas, opts, params.GameAddress, coreAddress, startingBlockNumber, params.Delay)
	defer io.Stop()

	go func() {
		for err := range io.ErrChan() {
			log.Error("IO error", "err", err)
		}
	}()

	// Determine client mode (VIEW ONLY, ACTIVE) and block to sync to
	// Set mode to view only if the game was created 1800 blocks ago

	headBlockNumber, err := rpcClient.BlockNumber(context.Background())
	if err != nil {
		logCrit(fmt.Errorf("Failed to get head block number: %v", err))
	}
	var (
		maxBlockToSyncTo = uint64(metaRow.CreationBlockNumber) + 600 // TODO: Parameterize this
		blockToSyncTo    = utils.Min(headBlockNumber, maxBlockToSyncTo)
		syncToHead       = blockToSyncTo == headBlockNumber
		// clientCanSend    = playerId != 0 && syncToHead
		clientSync = syncToHead
	)

	clientPlayerId := playerId
	if clientPlayerId == 0 {
		clientPlayerId = 1
	}

	// Create headless client
	hl := core.NewHeadlessClient(kv, io)
	hl.SetPlayerId(clientPlayerId)

	log.Info("Started headless client")

	// Sync from blocks
	setLoadStatus("Syncing...0%")

	blockSyncStartTime := time.Now()
	syncedBlockNumber := uint64(hl.Game().BlockNumber())

	for {
		// Sync until target block
		if syncedBlockNumber >= blockToSyncTo {
			break
		}
		if syncToHead {
			log.Info("Syncing to head", "blockNumber", headBlockNumber)
		} else {
			log.Info("Syncing to block", "blockNumber", blockToSyncTo)
		}

		syncTo := utils.Min(blockToSyncTo, syncedBlockNumber+512)
		hl.SyncUntil(syncTo)
		syncedBlockNumber = syncTo

		percentLoaded := 100 * float64(syncedBlockNumber-startingBlockNumber) / float64(blockToSyncTo-startingBlockNumber)
		setLoadStatus(fmt.Sprintf("Syncing...%d%%", int(percentLoaded)))

		if syncToHead {
			// Check for new head
			headBlockNumber, err = rpcClient.BlockNumber(context.Background())
			if err != nil {
				log.Crit("Failed to get head block number", "err", err)
			}
			blockToSyncTo = utils.Min(headBlockNumber, maxBlockToSyncTo)
		}
	}

	blockSyncDuration = time.Since(blockSyncStartTime)
	log.Debug("Synced to block", "blockNumber", blockToSyncTo, "duration", blockSyncDuration)
	log.Info("Synced to block", "blockNumber", blockToSyncTo)

	postHog.Capture("SyncReachedTargetBlock", map[string]interface{}{
		"playerAddress": senderAddress.Hex(),
		"gameAddress":   params.GameAddress.Hex(),
		"coreAddress":   coreAddress.Hex(),

		"headBlockNumber":    uint64(headBlockNumber),
		"initialBlockNumber": uint64(startingBlockNumber),
		"syncedBlocks":       uint64(headBlockNumber - startingBlockNumber),

		"duration": blockSyncDuration.Seconds(),
	})

	postHog.Capture("SyncCompleted", map[string]interface{}{
		"playerAddress": senderAddress.Hex(),
		"gameAddress":   params.GameAddress.Hex(),
		"coreAddress":   coreAddress.Hex(),

		"blockNumber": uint64(hl.Game().BlockNumber()),

		"fetchSnapshotTime": fetchSnapshotDuration.Seconds(),
		"loadSnapshotTime":  loadSnapshotDuration.Seconds(),
		"syncToHeadTime":    blockSyncDuration.Seconds(),
		"totalDuration":     (fetchSnapshotDuration + loadSnapshotDuration + blockSyncDuration).Seconds(),

		"loadSnapshot": mustLoadSnapshot,
	})

	io.SetTxUpdateHook(func(txUpdate *rpc.ActionTxUpdate) {
		if txUpdate.Status == rpc.ActionTxStatus_Failed {
			log.Error("Failed to send transaction", "txHash", txUpdate.TxHash, "err", txUpdate.Err)
		} else if txUpdate.Status == rpc.ActionTxStatus_Included {
			if !params.Debug {
				log.Info("Transaction included", "txHash", txUpdate.TxHash)
			}

			tx, _, err := rpcClient.TransactionByHash(context.Background(), txUpdate.TxHash)
			if err != nil {
				log.Error("Failed to retrieve transaction", "err", err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			receipt, err := bind.WaitMined(ctx, rpcClient, tx)
			cancel()
			if err != nil {
				log.Error("Failed to retrieve receipt", "err", err)
				return
			}

			var (
				gasUsed         = receipt.GasUsed
				gasLimit        = tx.Gas()
				overpaidPercent = float64(gasLimit)/float64(gasUsed)*100 - 100
			)

			log.Info("New transaction", "txHash", tx.Hash(), "status", receipt.Status, "gasUsed", gasUsed, "gasLimit", gasLimit, "overpaid", fmt.Sprintf("%.2f%%", overpaidPercent))
		}
	})

	log.Info("Starting game client")
	setLoadStatus("Starting...")

	// Create client
	cli := game.NewClient(hl, clientConfig, clientSync)

	// Start client
	hideLoadStatus()
	ebiten.SetWindowSize(clientConfig.ScreenSize.X, clientConfig.ScreenSize.Y)
	ebiten.SetWindowTitle("Game")
	if err := ebiten.RunGame(cli); err != nil && err != core.ErrQuit {
		logCrit(fmt.Errorf("Failed to run game: %v", err))
	}
}

func main() {
	// Get URL params
	params, err := getURLParams()
	if err != nil {
		logCrit(fmt.Errorf("Failed to get URL params: %v", err))
	}

	// Get private key
	privateKey, err := getPrivateKey()
	if err != nil {
		logCrit(fmt.Errorf("Failed to get burner key: %v", err))
	}

	// Set client config
	clientConfig.ScreenSize = newGameScreenSize()

	// Set log level
	minLvl := log.LvlInfo
	if params.Debug {
		minLvl = log.LvlDebug
	}
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, minLvl, true)))

	// Start
	log.Debug(
		"Starting game client",
		"gameAddress", params.GameAddress.Hex(),
		"wsURL", params.WsURL,
		"interpolate", params.Interpolate,
		"blockTime", params.BlockTime,
		"delay", params.Delay,
		"screenSize", clientConfig.ScreenSize,
	)
	runGameClient(clientConfig, params, privateKey)
}
