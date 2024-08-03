package main

import (
	"context"
	"os"

	"github.com/concrete-eth/ark-rts/web/cmd"
	"github.com/concrete-eth/ark-rts/web/sidecar"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/posthog/posthog-go"
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelDebug, true)))

	config := sidecar.SidecarConfig{
		SnapshotInterval:      uint64(cmd.GetIntEnvWithDefault("SNAPSHOT_INTERVAL", 60)),
		EvictionThreshold:     uint64(cmd.GetIntEnvWithDefault("EVICTION_THRESHOLD", 1800)),
		RpcURL:                cmd.GetEnvWithDefault("WS_URL", "ws://localhost:9546"),
		GameFactoryAddressHex: cmd.GetEnvWithDefault("FACTORY_ADDRESS", ""),
		PostHogApiKey:         cmd.GetEnvWithDefault("POSTHOG_API_KEY", "phc_LsLNl9nWm168DyAPmzYxLCG42C6dzWXGdSgQEmzfHx0"),
		PostHogEndpoint:       cmd.GetEnvWithDefault("POSTHOG_ENDPOINT", "https://eu.posthog.com"),
		JWTSecret:             common.HexToHash(cmd.GetEnvWithDefault("JWT_SECRET", "688f5d737bad920bdfb2fc2f488d6b6209eebda1dae949a8de91398d932c517a")),
	}
	if err := config.Validate(); err != nil {
		log.Crit("Invalid configuration", "err", err)
	}

	auth := rpc.WithHTTPAuth(node.NewJWTAuth(config.JWTSecret))
	rpcClient, err := rpc.DialOptions(context.Background(), config.RpcURL, auth)
	if err != nil {
		log.Crit("Failed to connect to RPC", "err", err)
	}

	ethClient := ethclient.NewClient(rpcClient)

	log.Debug("Connected to RPC", "url", config.RpcURL)

	var ph posthog.Client
	if config.PostHogApiKey != "" && config.PostHogEndpoint != "" {
		ph, err = posthog.NewWithConfig(config.PostHogApiKey, posthog.Config{
			Endpoint: config.PostHogEndpoint,
		})
		if err != nil {
			log.Crit("Failed to create PostHog client", "err", err)
		}
		log.Debug("Created PostHog client", "endpoint", config.PostHogEndpoint)
	}

	s, err := sidecar.NewSidecar(config, ethClient, ph)
	if err != nil {
		log.Crit("Failed to create sidecar", "err", err)
	}

	if err := s.Run(); err != nil {
		log.Crit("Failed to run sidecar", "err", err)
	}
}
