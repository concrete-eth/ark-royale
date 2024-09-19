package main

import (
	"context"
	"fmt"
	"image"
	"math/big"
	"os"
	"time"

	"github.com/concrete-eth/archetype/arch"
	"github.com/concrete-eth/archetype/kvstore"
	"github.com/concrete-eth/archetype/rpc"
	"github.com/concrete-eth/ark-royale/client/core"
	"github.com/concrete-eth/ark-royale/client/game"
	game_contract "github.com/concrete-eth/ark-royale/gogen/abigen/game"
	factory_contract "github.com/concrete-eth/ark-royale/gogen/abigen/game_factory"
	"github.com/concrete-eth/ark-royale/gogen/archmod"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	pcAddr = common.HexToAddress("0x80")
	// rpcUrl = "wss://dcp.concretelabs.dev"
	rpcUrl        = "ws://127.0.0.1:8545"
	privateKeyHex = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	blockTime     = 1 * time.Second
	// chainId       = big.NewInt(901)
	chainId = big.NewInt(1337)
)

func waitForTx(ethcli *ethclient.Client, tx *types.Transaction) {
	for {
		_, pending, err := ethcli.TransactionByHash(context.Background(), tx.Hash())
		if err != nil {
			panic(err)
		}
		if !pending {
			break
		}
		time.Sleep(1 * time.Second)
	}
	receipt, err := ethcli.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		panic(err)
	}
	if receipt.Status != 1 {
		panic("tx failed")
	}
}

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelDebug, true)))

	// Create eth client
	ethcli, err := ethclient.Dial(rpcUrl)
	if err != nil {
		panic(err)
	}

	// Load tx opts
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		panic(err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		panic(err)
	}
	// auth.Nonce = big.NewInt(0)
	auth.GasLimit = 10_000_000

	blockNum, err := ethcli.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}

	// Deploy game with factory
	implAddr, implTx, _, err := game_contract.DeployContract(auth, ethcli)
	if err != nil {
		panic(err)
	}

	fmt.Println("Game implementation deployed at", implAddr.Hex(), "with tx", implTx.Hash().Hex())

	factoryAddr, factoryTx, factoryContract, err := factory_contract.DeployContract(auth, ethcli, new(big.Int).SetInt64(25_000_000), implAddr, pcAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Factory deployed at", factoryAddr.Hex(), "with tx", factoryTx.Hash().Hex())

	fmt.Println("Waiting for game implementation transaction...")
	waitForTx(ethcli, implTx)
	fmt.Println("Waiting for factory transaction...")
	waitForTx(ethcli, factoryTx)

	gameCreatedChan := make(chan *factory_contract.ContractGameCreated, 1)
	sub, err := factoryContract.WatchGameCreated(nil, gameCreatedChan)
	if err != nil {
		panic(err)
	}

	createTx, err := factoryContract.CreateGame(auth, "0", []common.Address{auth.From, auth.From})
	if err != nil {
		panic(err)
	}

	fmt.Println("Game creation transaction sent with hash", createTx.Hash().Hex())
	fmt.Println("Waiting for game creation transaction...")
	waitForTx(ethcli, createTx)

	fmt.Println("Waiting for game created event...")

	gameCreated := <-gameCreatedChan
	sub.Unsubscribe()
	close(gameCreatedChan)

	gameAddr := gameCreated.GameAddress

	fmt.Println("Game deployed at", gameAddr.Hex(), "with tx", gameCreated.Raw.TxHash.Hex())

	// blockNum := uint64(0)
	// gameAddr := common.HexToAddress("0x856e4424f806D16E8CBC702B3c0F2ede5468eae5")

	gameContract, err := game_contract.NewContract(gameAddr, ethcli)
	if err != nil {
		panic(err)
	}

	coreAddr, err := gameContract.Proxy(nil)
	if err != nil {
		panic(err)
	}

	// Create local simulated io
	// Create schemas from codegen
	schemas := arch.ArchSchemas{Actions: archmod.ActionSchemas, Tables: archmod.TableSchemas}
	io := rpc.NewIO(ethcli, blockTime, schemas, auth, gameAddr, coreAddr, blockNum, 300*time.Millisecond)
	// io := rpc.NewIO(ethcli, blockTime, schemas, auth, gameAddr, coreAddr, 0, 0)
	defer io.Stop()

	go func() {
		for err := range io.ErrChan() {
			log.Error("IO error", "err", err)
		}
	}()

	// Create and start client
	kv := kvstore.NewMemoryKeyValueStore()
	hl := core.NewHeadlessClient(kv, io)
	hl.SetPlayerId(1)

	lastBlockNum, err := ethcli.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	hl.SyncUntil(lastBlockNum)

	c := game.NewClient(hl, core.ClientConfig{
		ScreenSize: image.Point{1280, 720},
	}, true)
	w, h := c.Layout(-1, -1)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Ark Royale")
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(c); err != nil {
		panic(err)
	}
}
