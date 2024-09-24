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

// ActionDataAddBuildingPrototype is an auto generated low-level Go binding around an user-defined struct.
type ActionDataAddBuildingPrototype struct {
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

// ActionDataAddPlayer is an auto generated low-level Go binding around an user-defined struct.
type ActionDataAddPlayer struct {
	SpawnAreaX           uint16
	SpawnAreaY           uint16
	SpawnAreaWidth       uint8
	SpawnAreaHeight      uint8
	WorkerPortX          uint16
	WorkerPortY          uint16
	UnpurgeableUnitCount uint8
}

// ActionDataAddUnitPrototype is an auto generated low-level Go binding around an user-defined struct.
type ActionDataAddUnitPrototype struct {
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

// ActionDataAssignUnit is an auto generated low-level Go binding around an user-defined struct.
type ActionDataAssignUnit struct {
	PlayerId     uint8
	UnitId       uint8
	Command      uint64
	CommandExtra uint64
	CommandMeta  uint8
}

// ActionDataCreateUnit is an auto generated low-level Go binding around an user-defined struct.
type ActionDataCreateUnit struct {
	PlayerId uint8
	UnitType uint8
	X        uint16
	Y        uint16
}

// ActionDataInitialize is an auto generated low-level Go binding around an user-defined struct.
type ActionDataInitialize struct {
	Width  uint16
	Height uint16
}

// ActionDataPlaceBuilding is an auto generated low-level Go binding around an user-defined struct.
type ActionDataPlaceBuilding struct {
	PlayerId     uint8
	BuildingType uint8
	X            uint16
	Y            uint16
}

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"addBuildingPrototype\",\"inputs\":[{\"name\":\"action\",\"type\":\"tuple\",\"internalType\":\"structActionData_AddBuildingPrototype\",\"components\":[{\"name\":\"width\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"height\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"resourceCost\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"resourceCapacity\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"computeCapacity\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"resourceMine\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"mineTime\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"maxIntegrity\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"buildingTime\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isArmory\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isEnvironment\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addPlayer\",\"inputs\":[{\"name\":\"action\",\"type\":\"tuple\",\"internalType\":\"structActionData_AddPlayer\",\"components\":[{\"name\":\"spawnAreaX\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"spawnAreaY\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"spawnAreaWidth\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"spawnAreaHeight\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"workerPortX\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"workerPortY\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"unpurgeableUnitCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addUnitPrototype\",\"inputs\":[{\"name\":\"action\",\"type\":\"tuple\",\"internalType\":\"structActionData_AddUnitPrototype\",\"components\":[{\"name\":\"layer\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"resourceCost\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"computeCost\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"spawnTime\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"maxIntegrity\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"landStrength\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"hoverStrength\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"airStrength\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"attackRange\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"attackCooldown\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isAssault\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isConfrontational\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isWorker\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isPurgeable\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"assignUnit\",\"inputs\":[{\"name\":\"action\",\"type\":\"tuple\",\"internalType\":\"structActionData_AssignUnit\",\"components\":[{\"name\":\"playerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"unitId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"command\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commandExtra\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"commandMeta\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createUnit\",\"inputs\":[{\"name\":\"action\",\"type\":\"tuple\",\"internalType\":\"structActionData_CreateUnit\",\"components\":[{\"name\":\"playerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"unitType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"x\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"y\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"action\",\"type\":\"tuple\",\"internalType\":\"structActionData_Initialize\",\"components\":[{\"name\":\"width\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"height\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"placeBuilding\",\"inputs\":[{\"name\":\"action\",\"type\":\"tuple\",\"internalType\":\"structActionData_PlaceBuilding\",\"components\":[{\"name\":\"playerId\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"buildingType\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"x\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"y\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"purge\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"start\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tick\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ActionExecuted\",\"inputs\":[{\"name\":\"actionId\",\"type\":\"bytes4\",\"indexed\":false,\"internalType\":\"bytes4\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
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

// AddBuildingPrototype is a paid mutator transaction binding the contract method 0x982d778d.
//
// Solidity: function addBuildingPrototype((uint8,uint8,uint16,uint16,uint8,uint8,uint8,uint8,uint8,bool,bool) action) returns()
func (_Contract *ContractTransactor) AddBuildingPrototype(opts *bind.TransactOpts, action ActionDataAddBuildingPrototype) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addBuildingPrototype", action)
}

// AddBuildingPrototype is a paid mutator transaction binding the contract method 0x982d778d.
//
// Solidity: function addBuildingPrototype((uint8,uint8,uint16,uint16,uint8,uint8,uint8,uint8,uint8,bool,bool) action) returns()
func (_Contract *ContractSession) AddBuildingPrototype(action ActionDataAddBuildingPrototype) (*types.Transaction, error) {
	return _Contract.Contract.AddBuildingPrototype(&_Contract.TransactOpts, action)
}

// AddBuildingPrototype is a paid mutator transaction binding the contract method 0x982d778d.
//
// Solidity: function addBuildingPrototype((uint8,uint8,uint16,uint16,uint8,uint8,uint8,uint8,uint8,bool,bool) action) returns()
func (_Contract *ContractTransactorSession) AddBuildingPrototype(action ActionDataAddBuildingPrototype) (*types.Transaction, error) {
	return _Contract.Contract.AddBuildingPrototype(&_Contract.TransactOpts, action)
}

// AddPlayer is a paid mutator transaction binding the contract method 0xb44ff0c5.
//
// Solidity: function addPlayer((uint16,uint16,uint8,uint8,uint16,uint16,uint8) action) returns()
func (_Contract *ContractTransactor) AddPlayer(opts *bind.TransactOpts, action ActionDataAddPlayer) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addPlayer", action)
}

// AddPlayer is a paid mutator transaction binding the contract method 0xb44ff0c5.
//
// Solidity: function addPlayer((uint16,uint16,uint8,uint8,uint16,uint16,uint8) action) returns()
func (_Contract *ContractSession) AddPlayer(action ActionDataAddPlayer) (*types.Transaction, error) {
	return _Contract.Contract.AddPlayer(&_Contract.TransactOpts, action)
}

// AddPlayer is a paid mutator transaction binding the contract method 0xb44ff0c5.
//
// Solidity: function addPlayer((uint16,uint16,uint8,uint8,uint16,uint16,uint8) action) returns()
func (_Contract *ContractTransactorSession) AddPlayer(action ActionDataAddPlayer) (*types.Transaction, error) {
	return _Contract.Contract.AddPlayer(&_Contract.TransactOpts, action)
}

// AddUnitPrototype is a paid mutator transaction binding the contract method 0xa5e592f4.
//
// Solidity: function addUnitPrototype((uint8,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,bool,bool,bool,bool) action) returns()
func (_Contract *ContractTransactor) AddUnitPrototype(opts *bind.TransactOpts, action ActionDataAddUnitPrototype) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "addUnitPrototype", action)
}

// AddUnitPrototype is a paid mutator transaction binding the contract method 0xa5e592f4.
//
// Solidity: function addUnitPrototype((uint8,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,bool,bool,bool,bool) action) returns()
func (_Contract *ContractSession) AddUnitPrototype(action ActionDataAddUnitPrototype) (*types.Transaction, error) {
	return _Contract.Contract.AddUnitPrototype(&_Contract.TransactOpts, action)
}

// AddUnitPrototype is a paid mutator transaction binding the contract method 0xa5e592f4.
//
// Solidity: function addUnitPrototype((uint8,uint16,uint8,uint8,uint8,uint8,uint8,uint8,uint8,uint8,bool,bool,bool,bool) action) returns()
func (_Contract *ContractTransactorSession) AddUnitPrototype(action ActionDataAddUnitPrototype) (*types.Transaction, error) {
	return _Contract.Contract.AddUnitPrototype(&_Contract.TransactOpts, action)
}

// AssignUnit is a paid mutator transaction binding the contract method 0xf8613b59.
//
// Solidity: function assignUnit((uint8,uint8,uint64,uint64,uint8) action) returns()
func (_Contract *ContractTransactor) AssignUnit(opts *bind.TransactOpts, action ActionDataAssignUnit) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "assignUnit", action)
}

// AssignUnit is a paid mutator transaction binding the contract method 0xf8613b59.
//
// Solidity: function assignUnit((uint8,uint8,uint64,uint64,uint8) action) returns()
func (_Contract *ContractSession) AssignUnit(action ActionDataAssignUnit) (*types.Transaction, error) {
	return _Contract.Contract.AssignUnit(&_Contract.TransactOpts, action)
}

// AssignUnit is a paid mutator transaction binding the contract method 0xf8613b59.
//
// Solidity: function assignUnit((uint8,uint8,uint64,uint64,uint8) action) returns()
func (_Contract *ContractTransactorSession) AssignUnit(action ActionDataAssignUnit) (*types.Transaction, error) {
	return _Contract.Contract.AssignUnit(&_Contract.TransactOpts, action)
}

// CreateUnit is a paid mutator transaction binding the contract method 0x143ca15f.
//
// Solidity: function createUnit((uint8,uint8,uint16,uint16) action) returns()
func (_Contract *ContractTransactor) CreateUnit(opts *bind.TransactOpts, action ActionDataCreateUnit) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "createUnit", action)
}

// CreateUnit is a paid mutator transaction binding the contract method 0x143ca15f.
//
// Solidity: function createUnit((uint8,uint8,uint16,uint16) action) returns()
func (_Contract *ContractSession) CreateUnit(action ActionDataCreateUnit) (*types.Transaction, error) {
	return _Contract.Contract.CreateUnit(&_Contract.TransactOpts, action)
}

// CreateUnit is a paid mutator transaction binding the contract method 0x143ca15f.
//
// Solidity: function createUnit((uint8,uint8,uint16,uint16) action) returns()
func (_Contract *ContractTransactorSession) CreateUnit(action ActionDataCreateUnit) (*types.Transaction, error) {
	return _Contract.Contract.CreateUnit(&_Contract.TransactOpts, action)
}

// Initialize is a paid mutator transaction binding the contract method 0xeaba9837.
//
// Solidity: function initialize((uint16,uint16) action) returns()
func (_Contract *ContractTransactor) Initialize(opts *bind.TransactOpts, action ActionDataInitialize) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "initialize", action)
}

// Initialize is a paid mutator transaction binding the contract method 0xeaba9837.
//
// Solidity: function initialize((uint16,uint16) action) returns()
func (_Contract *ContractSession) Initialize(action ActionDataInitialize) (*types.Transaction, error) {
	return _Contract.Contract.Initialize(&_Contract.TransactOpts, action)
}

// Initialize is a paid mutator transaction binding the contract method 0xeaba9837.
//
// Solidity: function initialize((uint16,uint16) action) returns()
func (_Contract *ContractTransactorSession) Initialize(action ActionDataInitialize) (*types.Transaction, error) {
	return _Contract.Contract.Initialize(&_Contract.TransactOpts, action)
}

// PlaceBuilding is a paid mutator transaction binding the contract method 0xd74de075.
//
// Solidity: function placeBuilding((uint8,uint8,uint16,uint16) action) returns()
func (_Contract *ContractTransactor) PlaceBuilding(opts *bind.TransactOpts, action ActionDataPlaceBuilding) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "placeBuilding", action)
}

// PlaceBuilding is a paid mutator transaction binding the contract method 0xd74de075.
//
// Solidity: function placeBuilding((uint8,uint8,uint16,uint16) action) returns()
func (_Contract *ContractSession) PlaceBuilding(action ActionDataPlaceBuilding) (*types.Transaction, error) {
	return _Contract.Contract.PlaceBuilding(&_Contract.TransactOpts, action)
}

// PlaceBuilding is a paid mutator transaction binding the contract method 0xd74de075.
//
// Solidity: function placeBuilding((uint8,uint8,uint16,uint16) action) returns()
func (_Contract *ContractTransactorSession) PlaceBuilding(action ActionDataPlaceBuilding) (*types.Transaction, error) {
	return _Contract.Contract.PlaceBuilding(&_Contract.TransactOpts, action)
}

// Purge is a paid mutator transaction binding the contract method 0x70f0c351.
//
// Solidity: function purge() returns()
func (_Contract *ContractTransactor) Purge(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "purge")
}

// Purge is a paid mutator transaction binding the contract method 0x70f0c351.
//
// Solidity: function purge() returns()
func (_Contract *ContractSession) Purge() (*types.Transaction, error) {
	return _Contract.Contract.Purge(&_Contract.TransactOpts)
}

// Purge is a paid mutator transaction binding the contract method 0x70f0c351.
//
// Solidity: function purge() returns()
func (_Contract *ContractTransactorSession) Purge() (*types.Transaction, error) {
	return _Contract.Contract.Purge(&_Contract.TransactOpts)
}

// Start is a paid mutator transaction binding the contract method 0xbe9a6555.
//
// Solidity: function start() returns()
func (_Contract *ContractTransactor) Start(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "start")
}

// Start is a paid mutator transaction binding the contract method 0xbe9a6555.
//
// Solidity: function start() returns()
func (_Contract *ContractSession) Start() (*types.Transaction, error) {
	return _Contract.Contract.Start(&_Contract.TransactOpts)
}

// Start is a paid mutator transaction binding the contract method 0xbe9a6555.
//
// Solidity: function start() returns()
func (_Contract *ContractTransactorSession) Start() (*types.Transaction, error) {
	return _Contract.Contract.Start(&_Contract.TransactOpts)
}

// Tick is a paid mutator transaction binding the contract method 0x3eaf5d9f.
//
// Solidity: function tick() returns()
func (_Contract *ContractTransactor) Tick(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "tick")
}

// Tick is a paid mutator transaction binding the contract method 0x3eaf5d9f.
//
// Solidity: function tick() returns()
func (_Contract *ContractSession) Tick() (*types.Transaction, error) {
	return _Contract.Contract.Tick(&_Contract.TransactOpts)
}

// Tick is a paid mutator transaction binding the contract method 0x3eaf5d9f.
//
// Solidity: function tick() returns()
func (_Contract *ContractTransactorSession) Tick() (*types.Transaction, error) {
	return _Contract.Contract.Tick(&_Contract.TransactOpts)
}

// ContractActionExecutedIterator is returned from FilterActionExecuted and is used to iterate over the raw logs and unpacked data for ActionExecuted events raised by the Contract contract.
type ContractActionExecutedIterator struct {
	Event *ContractActionExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractActionExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractActionExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractActionExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractActionExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractActionExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractActionExecuted represents a ActionExecuted event raised by the Contract contract.
type ContractActionExecuted struct {
	ActionId [4]byte
	Data     []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterActionExecuted is a free log retrieval operation binding the contract event 0x45065f461aede1b904079823f6d858e465fa8c25fcf1654bb4a89e6dee320a1a.
//
// Solidity: event ActionExecuted(bytes4 actionId, bytes data)
func (_Contract *ContractFilterer) FilterActionExecuted(opts *bind.FilterOpts) (*ContractActionExecutedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "ActionExecuted")
	if err != nil {
		return nil, err
	}
	return &ContractActionExecutedIterator{contract: _Contract.contract, event: "ActionExecuted", logs: logs, sub: sub}, nil
}

// WatchActionExecuted is a free log subscription operation binding the contract event 0x45065f461aede1b904079823f6d858e465fa8c25fcf1654bb4a89e6dee320a1a.
//
// Solidity: event ActionExecuted(bytes4 actionId, bytes data)
func (_Contract *ContractFilterer) WatchActionExecuted(opts *bind.WatchOpts, sink chan<- *ContractActionExecuted) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "ActionExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractActionExecuted)
				if err := _Contract.contract.UnpackLog(event, "ActionExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseActionExecuted is a log parse operation binding the contract event 0x45065f461aede1b904079823f6d858e465fa8c25fcf1654bb4a89e6dee320a1a.
//
// Solidity: event ActionExecuted(bytes4 actionId, bytes data)
func (_Contract *ContractFilterer) ParseActionExecuted(log types.Log) (*ContractActionExecuted, error) {
	event := new(ContractActionExecuted)
	if err := _Contract.contract.UnpackLog(event, "ActionExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
