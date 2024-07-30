package main

import (
	"image"
	"os"
	"time"

	"github.com/concrete-eth/archetype/arch"
	"github.com/concrete-eth/archetype/deploy"
	"github.com/concrete-eth/archetype/kvstore"
	"github.com/concrete-eth/archetype/precompile"
	"github.com/concrete-eth/archetype/rpc"
	"github.com/concrete-eth/ark-rts/client/core"
	"github.com/concrete-eth/ark-rts/client/game"
	game_contract "github.com/concrete-eth/ark-rts/gogen/abigen/game"
	"github.com/concrete-eth/ark-rts/gogen/archmod"
	rts "github.com/concrete-eth/ark-rts/rts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	pcAddr = common.HexToAddress("0x1234")
)

func main() {
	log.SetDefault(log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelWarn, true)))

	// Create schemas from codegen
	schemas := arch.ArchSchemas{Actions: archmod.ActionSchemas, Tables: archmod.TableSchemas}

	// Create precompile
	pc := precompile.NewCorePrecompile(schemas, func() arch.Core { return &rts.Core{} })
	registry := concrete.NewRegistry()
	registry.AddPrecompile(0, pcAddr, pc)

	// Create local simulated io
	io, err := deploy.NewLocalIO(registry, schemas, func(auth *bind.TransactOpts, ethcli bind.ContractBackend) (common.Address, *types.Transaction, deploy.InitializableProxyAdmin, error) {
		return game_contract.DeployContract(auth, ethcli)
	}, pcAddr, 1*time.Second)
	if err != nil {
		panic(err)
	}
	defer io.Stop()

	io.SetTxUpdateHook(func(txUpdate *rpc.ActionTxUpdate) {
		log.Warn("Transaction "+txUpdate.Status.String(), "nonce", txUpdate.Nonce, "txHash", txUpdate.TxHash.Hex(), "err", txUpdate.Err)
	})

	// Create and start client
	kv := kvstore.NewMemoryKeyValueStore()
	hl := core.NewHeadlessClient(kv, io)
	hl.SetPlayerId(1)
	c := game.NewClient(hl, core.ClientConfig{
		ScreenSize:  image.Point{1280, 720},
		Interpolate: true,
		FixedCamera: false,
	}, true)
	w, h := c.Layout(-1, -1)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Ark RTS")
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(c); err != nil {
		panic(err)
	}
}
