package main

import (
	"fmt"
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
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

	// Create init data
	pk, _ := crypto.HexToECDSA(deploy.LocalPrivateKeyHex)
	address := crypto.PubkeyToAddress(pk.PublicKey)
	data, err := encodeAddressArray([]common.Address{address, address})
	if err != nil {
		panic(err)
	}

	// Create local simulated io
	io, err := deploy.NewLocalIO(registry, schemas, func(auth *bind.TransactOpts, ethcli bind.ContractBackend) (addr common.Address, tx *types.Transaction, game deploy.InitializableProxyAdmin, err error) {
		auth.GasLimit = 3_500_000
		return game_contract.DeployContract(auth, ethcli)
	}, pcAddr, data, 1000*time.Millisecond)
	if err != nil {
		panic(err)
	}
	defer io.Stop()

	io.SetTxUpdateHook(func(txUpdate *rpc.ActionTxUpdate) {
		if txUpdate.Status == rpc.ActionTxStatus_Failed {
			log.Error("Failed to send transaction", "txHash", txUpdate.TxHash, "err", txUpdate.Err)
		}
	})

	// Create and start client
	kv := kvstore.NewMemoryKeyValueStore()
	hl := core.NewHeadlessClient(kv, io)
	hl.SetPlayerId(1)

	// Start game
	hl.Start()

	c := game.NewClient(hl, core.ClientConfig{
		ScreenSize: image.Point{700, 500},
	}, true)
	w, h := c.Layout(-1, -1)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle("Ark RTS")
	ebiten.SetTPS(60)
	if err := ebiten.RunGame(c); err != nil {
		panic(err)
	}
}

func encodeAddressArray(addresses []common.Address) ([]byte, error) {
	addressArrayType, err := abi.NewType("address[]", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %v", err)
	}

	arguments := abi.Arguments{{Type: addressArrayType}}

	data, err := arguments.Pack(addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to encode addresses: %v", err)
	}

	return data, nil
}
