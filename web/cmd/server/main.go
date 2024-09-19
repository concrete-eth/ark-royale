package main

import (
	"os"

	"github.com/concrete-eth/ark-royale/web/cmd"
	"github.com/concrete-eth/ark-royale/web/server"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelDebug, true)))

	templateDir := os.Args[1]
	staticDir := os.Args[2]

	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		log.Crit("Template directory does not exist", "dir", templateDir)
	}
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Crit("Static directory does not exist", "dir", staticDir)
	}

	config := server.ServerConfig{
		BaseURL:               cmd.GetEnvWithDefault("BASE_URL", "http://localhost:3000"),
		PortStr:               cmd.GetEnvWithDefault("PORT", "3000"),
		RpcURL:                cmd.GetEnvWithDefault("WS_URL", "ws://localhost:9546"),
		GameFactoryAddressHex: cmd.GetEnvWithDefault("FACTORY_ADDRESS", ""),
		FaucetPrivateKeyHex:   cmd.GetEnvWithDefault("FAUCET_KEY", "c6fe5b29bd8a729376cb0ef97e13705723e8c2798bac39aa71071bc4089b2762"),
		TemplateDir:           templateDir,
		StaticDir:             staticDir,
	}

	wasmUrl := cmd.GetEnvWithDefault("WASM_URL", "")
	if wasmUrl == "" {
		wasmUrl = config.BaseURL + "/static/play.wasm"
	}
	config.WasmURL = wasmUrl

	if err := config.Validate(); err != nil {
		log.Crit("Invalid configuration", "err", err)
	}

	rpc, err := ethclient.Dial(config.RpcURL)
	if err != nil {
		log.Crit("Failed to connect to the Ethereum client", "err", err)
	}

	log.Debug("Connected to RPC", "url", config.RpcURL)

	s, err := server.NewServer(config, rpc)
	if err != nil {
		log.Crit("Failed to create server", "err", err)
	}

	if err := s.Run(); err != nil {
		log.Crit("Failed to run server", "err", err)
	}
}
