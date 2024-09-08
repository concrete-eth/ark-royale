// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// RowDataBoard is an auto generated low-level Go binding around an user-defined struct.
type RowDataBoard struct {
	LandObjectType uint8
	LandPlayerId   uint8
	LandObjectId   uint8
	HoverPlayerId  uint8
	HoverUnitId    uint8
	AirPlayerId    uint8
	AirUnitId      uint8
}

// RowDataBuildingPrototypes is an auto generated low-level Go binding around an user-defined struct.
type RowDataBuildingPrototypes struct {
	Width            uint8
	Height           uint8
	ResourceCost     uint16
	ResourceCapacity uint16
	ComputeCapacity  uint8
	ResourceMine     uint8
	MineTime         uint8
	MaxIntegrity     uint8
	BuildingTime     uint8
	IsArmory         bool
	IsEnvironment    bool
}

// RowDataBuildings is an auto generated low-level Go binding around an user-defined struct.
type RowDataBuildings struct {
	X            uint16
	Y            uint16
	BuildingType uint8
	State        uint8
	Integrity    uint8
	Timestamp    uint32
}

// RowDataMeta is an auto generated low-level Go binding around an user-defined struct.
type RowDataMeta struct {
	BoardWidth             uint16
	BoardHeight            uint16
	PlayerCount            uint8
	UnitPrototypeCount     uint8
	BuildingPrototypeCount uint8
	IsInitialized          bool
	HasStarted             bool
	CreationBlockNumber    uint32
}

// RowDataPlayers is an auto generated low-level Go binding around an user-defined struct.
type RowDataPlayers struct {
	SpawnAreaX                uint16
	SpawnAreaY                uint16
	SpawnAreaWidth            uint8
	SpawnAreaHeight           uint8
	WorkerPortX               uint16
	WorkerPortY               uint16
	CurResource               uint16
	MaxResource               uint16
	CurArmories               uint8
	ComputeSupply             uint8
	ComputeDemand             uint8
	UnitCount                 uint8
	BuildingCount             uint8
	BuildingPayQueuePointer   uint8
	BuildingBuildQueuePointer uint8
	UnitPayQueuePointer       uint8
	UnpurgeableUnitCount      uint8
}

// RowDataUnitPrototypes is an auto generated low-level Go binding around an user-defined struct.
type RowDataUnitPrototypes struct {
	Layer             uint8
	ResourceCost      uint16
	ComputeCost       uint8
	SpawnTime         uint8
	MaxIntegrity      uint8
	LandStrength      uint8
	HoverStrength     uint8
	AirStrength       uint8
	AttackRange       uint8
	AttackCooldown    uint8
	IsAssault         bool
	IsConfrontational bool
	IsWorker          bool
	IsPurgeable       bool
}

// RowDataUnits is an auto generated low-level Go binding around an user-defined struct.
type RowDataUnits struct {
	X            uint16
	Y            uint16
	UnitType     uint8
	State        uint8
	Load         uint8
	Integrity    uint8
	Timestamp    uint32
	Command      uint64
	CommandExtra uint64
	CommandMeta  uint8
	IsPreTicked  bool
}

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getBoardRow\",\"inputs\":[{\"name\":\"x\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"y\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRowData_Board\",\"components\":[{\"name\":\"landObjectType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"landPlayerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"landObjectId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hoverPlayerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hoverUnitId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"airPlayerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"airUnitId\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBuildingPrototypesRow\",\"inputs\":[{\"name\":\"buildingType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRowData_BuildingPrototypes\",\"components\":[{\"name\":\"width\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"height\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"resourceCost\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"resourceCapacity\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"computeCapacity\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"resourceMine\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"mineTime\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"maxIntegrity\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"buildingTime\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isArmory\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isEnvironment\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBuildingsRow\",\"inputs\":[{\"name\":\"playerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"buildingId\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRowData_Buildings\",\"components\":[{\"name\":\"x\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"y\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"buildingType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"integrity\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"timestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMetaRow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRowData_Meta\",\"components\":[{\"name\":\"boardWidth\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"boardHeight\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"playerCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"unitPrototypeCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"buildingPrototypeCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isInitialized\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"hasStarted\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"creationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPlayersRow\",\"inputs\":[{\"name\":\"playerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRowData_Players\",\"components\":[{\"name\":\"spawnAreaX\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"spawnAreaY\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"spawnAreaWidth\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"spawnAreaHeight\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"workerPortX\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"workerPortY\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"curResource\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"maxResource\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"curArmories\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"computeSupply\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"computeDemand\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"unitCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"buildingCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"buildingPayQueuePointer\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"buildingBuildQueuePointer\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"unitPayQueuePointer\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"unpurgeableUnitCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUnitPrototypesRow\",\"inputs\":[{\"name\":\"unitType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRowData_UnitPrototypes\",\"components\":[{\"name\":\"layer\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"resourceCost\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"computeCost\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"spawnTime\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"maxIntegrity\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"landStrength\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hoverStrength\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"airStrength\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"attackRange\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"attackCooldown\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isAssault\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isConfrontational\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isWorker\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isPurgeable\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUnitsRow\",\"inputs\":[{\"name\":\"playerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"unitId\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRowData_Units\",\"components\":[{\"name\":\"x\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"y\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"unitType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"load\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"integrity\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"timestamp\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"command\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commandExtra\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commandMeta\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isPreTicked\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"}]",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// Contract is an auto generated Go binding around an Ethereum contract.
type Contract struct {
	ContractCaller     // Read-only binding to the contract
	ContractTransactor // Write-only binding to the contract
	ContractFilterer   // Log filterer for contract events
}

// ContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractSession struct {
	Contract     *Contract         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractCallerSession struct {
	Contract *ContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractTransactorSession struct {
	Contract     *ContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractRaw struct {
	Contract *Contract // Generic contract binding to access the raw methods on
}

// ContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractCallerRaw struct {
	Contract *ContractCaller // Generic read-only contract binding to access the raw methods on
}

// ContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractTransactorRaw struct {
	Contract *ContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContract creates a new instance of Contract, bound to a specific deployed contract.
func NewContract(address common.Address, backend bind.ContractBackend) (*Contract, error) {
	contract, err := bindContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

// NewContractCaller creates a new read-only instance of Contract, bound to a specific deployed contract.
func NewContractCaller(address common.Address, caller bind.ContractCaller) (*ContractCaller, error) {
	contract, err := bindContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractCaller{contract: contract}, nil
}

// NewContractTransactor creates a new write-only instance of Contract, bound to a specific deployed contract.
func NewContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractTransactor, error) {
	contract, err := bindContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractTransactor{contract: contract}, nil
}

// NewContractFilterer creates a new log filterer instance of Contract, bound to a specific deployed contract.
func NewContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractFilterer, error) {
	contract, err := bindContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractFilterer{contract: contract}, nil
}

// bindContract binds a generic wrapper to an already deployed contract.
func bindContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.ContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.ContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contract *ContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contract *ContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contract *ContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contract.Contract.contract.Transact(opts, method, params...)
}

// GetBoardRow is a free data retrieval call binding the contract method 0x2f5d0157.
//
// Solidity: function getBoardRow(uint16 x, uint16 y) view returns((uint8,uint8,uint8,uint8,uint8,uint8,uint8))
func (_Contract *ContractCaller) GetBoardRow(opts *bind.CallOpts, x uint16, y uint16) (RowDataBoard, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getBoardRow", x, y)

	if err != nil {
		return *new(RowDataBoard), err
	}

	out0 := *abi.ConvertType(out[0], new(RowDataBoard)).(*RowDataBoard)

	return out0, err

}

// GetBoardRow is a free data retrieval call binding the contract method 0x2f5d0157.
//
// Solidity: function getBoardRow(uint16 x, uint16 y) view returns((uint8,uint8,uint8,uint8,uint8,uint8,uint8))
func (_Contract *ContractSession) GetBoardRow(x uint16, y uint16) (RowDataBoard, error) {
	return _Contract.Contract.GetBoardRow(&_Contract.CallOpts, x, y)
}

// GetBoardRow is a free data retrieval call binding the contract method 0x2f5d0157.
//
// Solidity: function getBoardRow(uint16 x, uint16 y) view returns((uint8,uint8,uint8,uint8,uint8,uint8,uint8))
func (_Contract *ContractCallerSession) GetBoardRow(x uint16, y uint16) (RowDataBoard, error) {
	return _Contract.Contract.GetBoardRow(&_Contract.CallOpts, x, y)
}

// GetBuildingPrototypesRow is a free data retrieval call binding the contract method 0xad986db0.
//
// Solidity: function getBuildingPrototypesRow(uint8 buildingType) view returns((uint8,uint8,uint16,uint16,uint8,uint8,uint8,uint8,uint8,bool,bool))
func (_Contract *ContractCaller) GetBuildingPrototypesRow(opts *bind.CallOpts, buildingType uint8) (RowDataBuildingPrototypes, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getBuildingPrototypesRow", buildingType)

	if err != nil {
		return *new(RowDataBuildingPrototypes), err
	}

	out0 := *abi.ConvertType(out[0], new(RowDataBuildingPrototypes)).(*RowDataBuildingPrototypes)

	return out0, err

}

// GetBuildingPrototypesRow is a free data retrieval call binding the contract method 0xad986db0.
//
// Solidity: function getBuildingPrototypesRow(uint8 buildingType) view returns((uint8,uint8,uint16,uint16,uint8,uint8,uint8,uint8,uint8,bool,bool))
func (_Contract *ContractSession) GetBuildingPrototypesRow(buildingType uint8) (RowDataBuildingPrototypes, error) {
	return _Contract.Contract.GetBuildingPrototypesRow(&_Contract.CallOpts, buildingType)
}

// GetBuildingPrototypesRow is a free data retrieval call binding the contract method 0xad986db0.
//
// Solidity: function getBuildingPrototypesRow(uint8 buildingType) view returns((uint8,uint8,uint16,uint16,uint8,uint8,uint8,uint8,uint8,bool,bool))
func (_Contract *ContractCallerSession) GetBuildingPrototypesRow(buildingType uint8) (RowDataBuildingPrototypes, error) {
	return _Contract.Contract.GetBuildingPrototypesRow(&_Contract.CallOpts, buildingType)
}

// GetBuildingsRow is a free data retrieval call binding the contract method 0xeed886d9.
//
// Solidity: function getBuildingsRow(uint8 playerId, uint8 buildingId) view returns((uint16,uint16,uint8,uint8,uint8,uint32))
func (_Contract *ContractCaller) GetBuildingsRow(opts *bind.CallOpts, playerId uint8, buildingId uint8) (RowDataBuildings, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getBuildingsRow", playerId, buildingId)

	if err != nil {
		return *new(RowDataBuildings), err
	}

	out0 := *abi.ConvertType(out[0], new(RowDataBuildings)).(*RowDataBuildings)

	return out0, err

}

// GetBuildingsRow is a free data retrieval call binding the contract method 0xeed886d9.
//
// Solidity: function getBuildingsRow(uint8 playerId, uint8 buildingId) view returns((uint16,uint16,uint8,uint8,uint8,uint32))
func (_Contract *ContractSession) GetBuildingsRow(playerId uint8, buildingId uint8) (RowDataBuildings, error) {
	return _Contract.Contract.GetBuildingsRow(&_Contract.CallOpts, playerId, buildingId)
}

// GetBuildingsRow is a free data retrieval call binding the contract method 0xeed886d9.
//
// Solidity: function getBuildingsRow(uint8 playerId, uint8 buildingId) view returns((uint16,uint16,uint8,uint8,uint8,uint32))
func (_Contract *ContractCallerSession) GetBuildingsRow(playerId uint8, buildingId uint8) (RowDataBuildings, error) {
	return _Contract.Contract.GetBuildingsRow(&_Contract.CallOpts, playerId, buildingId)
}

// GetMetaRow is a free data retrieval call binding the contract method 0x422f7e1d.
//
// Solidity: function getMetaRow() view returns((uint16,uint16,uint8,uint8,uint8,bool,bool,uint32))
func (_Contract *ContractCaller) GetMetaRow(opts *bind.CallOpts) (RowDataMeta, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getMetaRow")

	if err != nil {
		return *new(RowDataMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(RowDataMeta)).(*RowDataMeta)

	return out0, err

}

// GetMetaRow is a free data retrieval call binding the contract method 0x422f7e1d.
//
// Solidity: function getMetaRow() view returns((uint16,uint16,uint8,uint8,uint8,bool,bool,uint32))
func (_Contract *ContractSession) GetMetaRow() (RowDataMeta, error) {
	return _Contract.Contract.GetMetaRow(&_Contract.CallOpts)
}

// GetMetaRow is a free data retrieval call binding the contract method 0x422f7e1d.
//
// Solidity: function getMetaRow() view returns((uint16,uint16,uint8,uint8,uint8,bool,bool,uint32))
func (_Contract *ContractCallerSession) GetMetaRow() (RowDataMeta, error) {
	return _Contract.Contract.GetMetaRow(&_Contract.CallOpts)
}

// GetPlayersRow is a free data retrieval call binding the contract method 0x051cfce4.
//
// Solidity: function getPlayersRow(uint8 playerId) view returns((uint16,uint16,uint8,uint8,uint16,uint16,uint16,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8))
func (_Contract *ContractCaller) GetPlayersRow(opts *bind.CallOpts, playerId uint8) (RowDataPlayers, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getPlayersRow", playerId)

	if err != nil {
		return *new(RowDataPlayers), err
	}

	out0 := *abi.ConvertType(out[0], new(RowDataPlayers)).(*RowDataPlayers)

	return out0, err

}

// GetPlayersRow is a free data retrieval call binding the contract method 0x051cfce4.
//
// Solidity: function getPlayersRow(uint8 playerId) view returns((uint16,uint16,uint8,uint8,uint16,uint16,uint16,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8))
func (_Contract *ContractSession) GetPlayersRow(playerId uint8) (RowDataPlayers, error) {
	return _Contract.Contract.GetPlayersRow(&_Contract.CallOpts, playerId)
}

// GetPlayersRow is a free data retrieval call binding the contract method 0x051cfce4.
//
// Solidity: function getPlayersRow(uint8 playerId) view returns((uint16,uint16,uint8,uint8,uint16,uint16,uint16,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8))
func (_Contract *ContractCallerSession) GetPlayersRow(playerId uint8) (RowDataPlayers, error) {
	return _Contract.Contract.GetPlayersRow(&_Contract.CallOpts, playerId)
}

// GetUnitPrototypesRow is a free data retrieval call binding the contract method 0x1903dc4a.
//
// Solidity: function getUnitPrototypesRow(uint8 unitType) view returns((uint8,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,bool,bool,bool,bool))
func (_Contract *ContractCaller) GetUnitPrototypesRow(opts *bind.CallOpts, unitType uint8) (RowDataUnitPrototypes, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getUnitPrototypesRow", unitType)

	if err != nil {
		return *new(RowDataUnitPrototypes), err
	}

	out0 := *abi.ConvertType(out[0], new(RowDataUnitPrototypes)).(*RowDataUnitPrototypes)

	return out0, err

}

// GetUnitPrototypesRow is a free data retrieval call binding the contract method 0x1903dc4a.
//
// Solidity: function getUnitPrototypesRow(uint8 unitType) view returns((uint8,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,bool,bool,bool,bool))
func (_Contract *ContractSession) GetUnitPrototypesRow(unitType uint8) (RowDataUnitPrototypes, error) {
	return _Contract.Contract.GetUnitPrototypesRow(&_Contract.CallOpts, unitType)
}

// GetUnitPrototypesRow is a free data retrieval call binding the contract method 0x1903dc4a.
//
// Solidity: function getUnitPrototypesRow(uint8 unitType) view returns((uint8,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,bool,bool,bool,bool))
func (_Contract *ContractCallerSession) GetUnitPrototypesRow(unitType uint8) (RowDataUnitPrototypes, error) {
	return _Contract.Contract.GetUnitPrototypesRow(&_Contract.CallOpts, unitType)
}

// GetUnitsRow is a free data retrieval call binding the contract method 0x0077cc5a.
//
// Solidity: function getUnitsRow(uint8 playerId, uint8 unitId) view returns((uint16,uint16,uint8,uint8,uint8,uint8,uint32,uint64,uint64,uint8,bool))
func (_Contract *ContractCaller) GetUnitsRow(opts *bind.CallOpts, playerId uint8, unitId uint8) (RowDataUnits, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getUnitsRow", playerId, unitId)

	if err != nil {
		return *new(RowDataUnits), err
	}

	out0 := *abi.ConvertType(out[0], new(RowDataUnits)).(*RowDataUnits)

	return out0, err

}

// GetUnitsRow is a free data retrieval call binding the contract method 0x0077cc5a.
//
// Solidity: function getUnitsRow(uint8 playerId, uint8 unitId) view returns((uint16,uint16,uint8,uint8,uint8,uint8,uint32,uint64,uint64,uint8,bool))
func (_Contract *ContractSession) GetUnitsRow(playerId uint8, unitId uint8) (RowDataUnits, error) {
	return _Contract.Contract.GetUnitsRow(&_Contract.CallOpts, playerId, unitId)
}

// GetUnitsRow is a free data retrieval call binding the contract method 0x0077cc5a.
//
// Solidity: function getUnitsRow(uint8 playerId, uint8 unitId) view returns((uint16,uint16,uint8,uint8,uint8,uint8,uint32,uint64,uint64,uint8,bool))
func (_Contract *ContractCallerSession) GetUnitsRow(playerId uint8, unitId uint8) (RowDataUnits, error) {
	return _Contract.Contract.GetUnitsRow(&_Contract.CallOpts, playerId, unitId)
}
