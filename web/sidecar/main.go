package sidecar

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/posthog/posthog-go"

	game_contract "github.com/concrete-eth/ark-rts/gogen/abigen/game"
	tables_contract "github.com/concrete-eth/ark-rts/gogen/abigen/tables"

	snapshot_types "github.com/concrete-eth/archetype/snapshot/types"

	factory_contract "github.com/concrete-eth/ark-rts/gogen/abigen/game_factory"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

type SidecarConfig struct {
	SnapshotInterval      uint64
	EvictionThreshold     uint64
	RpcURL                string
	GameFactoryAddressHex string
	PostHogApiKey         string
	PostHogEndpoint       string
	JWTSecret             [32]byte
}

func (c *SidecarConfig) Validate() error {
	if c.RpcURL == "" {
		return errors.New("RpcURL is required")
	}
	if c.GameFactoryAddressHex == "" {
		return errors.New("GameFactoryAddressHex is required")
	}
	if c.SnapshotInterval == 0 {
		return errors.New("SnapshotInterval is required")
	}
	if c.JWTSecret == [32]byte{} {
		return errors.New("JWTSecret is required")
	}
	return nil
}

type Sidecar struct {
	config SidecarConfig
	rpc    *ethclient.Client
	ph     posthog.Client
}

func NewSidecar(config SidecarConfig, rpc *ethclient.Client, ph posthog.Client) (*Sidecar, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &Sidecar{
		config: config,
		rpc:    rpc,
		ph:     ph,
	}, nil
}

func (s *Sidecar) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := s.watchGameCreated(ctx)
		if err != nil {
			log.Error("Failed to watch GameCreated event", "err", err)
			cancel()
		}
	}()

	go func() {
		defer wg.Done()
		err := s.startCleanupScheduler(ctx, 60*time.Second)
		if err != nil {
			log.Error("Failed to run periodic task", "err", err)
			cancel()
		}
	}()

	log.Info(
		"Sidecar started",
		"rpc", s.config.RpcURL,
		"factoryAddress", s.config.GameFactoryAddressHex,
		"snapshotInterval", s.config.SnapshotInterval,
		"evictionThreshold", s.config.EvictionThreshold,
	)

	wg.Wait()
	return nil
}

func (s *Sidecar) watchGameCreated(ctx context.Context) error {
	gameFactoryAddress := common.HexToAddress(s.config.GameFactoryAddressHex)
	factoryContract, err := factory_contract.NewContract(gameFactoryAddress, s.rpc)
	if err != nil {
		return fmt.Errorf("failed to connect to the factory contract: %w", err)
	}
	gameCreatedChan := make(chan *factory_contract.ContractGameCreated, 512)
	sub, err := factoryContract.WatchGameCreated(nil, gameCreatedChan)
	if err != nil {
		return fmt.Errorf("failed to watch GameCreated event: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case gameCreated := <-gameCreatedChan:
			err := s.onGameCreated(gameCreated)
			if err != nil {
				log.Error("Failed to process GameCreated event", "err", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Sidecar) onGameCreated(game *factory_contract.ContractGameCreated) error {
	log.Info("Game created", "address", game.GameAddress.Hex())

	gameContract, err := game_contract.NewContract(game.GameAddress, s.rpc)
	if err != nil {
		return fmt.Errorf("failed to connect to the game contract: %w", err)
	}
	coreAddress, err := gameContract.Proxy(nil) // TODO: rename proxy [?]
	if err != nil {
		return fmt.Errorf("failed to get core address: %w", err)
	}
	var resp snapshot_types.ScheduleResponse
	err = s.rpc.Client().Call(&resp, "ccsnap_addSchedule", snapshot_types.Schedule{
		Addresses:   []common.Address{coreAddress},
		BlockPeriod: s.config.SnapshotInterval,
		Replace:     true,
	})
	if err != nil {
		return fmt.Errorf("failed to add schedule: %w", err)
	}
	log.Info("Added schedule", "id", resp.ID)

	if s.ph != nil {
		s.ph.Enqueue(posthog.Capture{
			DistinctId: game.Origin.Hex(),
			Event:      "GameCreated",
			Properties: posthog.NewProperties().
				Set("gameAddress", game.GameAddress.Hex()).
				Set("lobbyId", game.LobbyId),
		})
	}

	return nil
}

func (s *Sidecar) startCleanupScheduler(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.cleanup()
			if err != nil {
				log.Error("Failed to run cleanup", "err", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Sidecar) cleanup() error {
	log.Info("Running cleanup")

	currentBlockNumber, err := s.rpc.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get block number: %w", err)
	}

	var schedules map[uint64]snapshot_types.Schedule
	err = s.rpc.Client().Call(&schedules, "ccsnap_getSchedules")
	if err != nil {
		return fmt.Errorf("failed to get schedules: %w", err)
	}

	for id, schedule := range schedules {
		for _, address := range schedule.Addresses {
			gameContract, err := tables_contract.NewContract(address, s.rpc)
			if err != nil {
				return fmt.Errorf("failed to connect to the game contract: %w", err)
			}

			metaRow, err := gameContract.GetMetaRow(nil)
			if err != nil {
				return fmt.Errorf("failed to get meta row: %w", err)
			}

			evict := s.config.EvictionThreshold > 0 && currentBlockNumber-uint64(metaRow.CreationBlockNumber) > s.config.EvictionThreshold
			if evict {
				err := s.rpc.Client().Call(nil, "ccsnap_deleteSchedule", id)
				if err != nil {
					return fmt.Errorf("failed to remove schedule: %w", err)
				}
			}

			log.Debug("Eviction check", "id", id, "evict", evict)
		}
	}

	return nil
}
