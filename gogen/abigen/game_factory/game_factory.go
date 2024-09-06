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

// ContractMetaData contains all meta data concerning the Contract contract.
var ContractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_maxGasAllocation\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_gameImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_coreImplementation\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"coreImplementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createGame\",\"inputs\":[{\"name\":\"lobbyId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_players\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"gameImplementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAddressOf\",\"inputs\":[{\"name\":\"idx\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGasAllocOf\",\"inputs\":[{\"name\":\"t\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getIndexOf\",\"inputs\":[{\"name\":\"t\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastBlock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"maxGasAllocation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nActiveTickees\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setGasAlloc\",\"inputs\":[{\"name\":\"t\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gas\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"evictionBlock\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tick\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalGasAllocation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GameCreated\",\"inputs\":[{\"name\":\"gameAddress\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"lobbyId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"origin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GasAllocSet\",\"inputs\":[{\"name\":\"tickee\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gas\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"evictionBlock\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC1167FailedCreateClone\",\"inputs\":[]}]",
	Bin: "0x60c060405234801561001057600080fd5b506040516110e53803806110e583398101604081905261002f916100a7565b600080546001600160a01b031916339081178255604051859282917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3506005556001600160a01b039182166080521660a052506100e3565b80516001600160a01b03811681146100a257600080fd5b919050565b6000806000606084860312156100bc57600080fd5b835192506100cc6020850161008b565b91506100da6040850161008b565b90509250925092565b60805160a051610fcf610116600039600081816101e201526105ac01526000818161025401526105750152610fcf6000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80639d492a2c1161008c578063b0cc1f0111610066578063b0cc1f011461022a578063dbc3352a14610233578063f2fde38b1461023c578063f6cf916e1461024f57600080fd5b80639d492a2c146101dd578063a97547d714610204578063aca113131461021757600080fd5b80633eaf5d9f116100c85780633eaf5d9f1461017657806344b920b714610180578063806b984f146101c15780638da5cb5b146101ca57600080fd5b8063017df522146100ef5780632e327fef1461013b57806332f79bd814610144575b600080fd5b6101286100fd366004610b98565b6001600160a01b0316600090815260016020526040902054600160401b90046001600160401b031690565b6040519081526020015b60405180910390f35b61012860045481565b610128610152366004610b98565b6001600160a01b03166000908152600160205260409020546001600160401b031690565b61017e610276565b005b6101a961018e366004610bba565b6000908152600260205260409020546001600160a01b031690565b6040516001600160a01b039091168152602001610132565b61012860065481565b6000546101a9906001600160a01b031681565b6101a97f000000000000000000000000000000000000000000000000000000000000000081565b61017e610212366004610bd3565b610514565b6101a9610225366004610cd2565b61056d565b61012860055481565b61012860035481565b61017e61024a366004610b98565b6106a5565b6101a97f000000000000000000000000000000000000000000000000000000000000000081565b60065443116102cc5760405162461bcd60e51b815260206004820152601f60248201527f5469636b4d61737465723a206f6e6c79206f6e63652070657220626c6f636b0060448201526064015b60405180910390fd5b4360065560035460005b8181101561039c57610307604051806040016040528060078152602001662a34b1b5b2b29d60c91b81525082610739565b620124f85a1015610348576103436040518060400160405280601081526020016f2ab73232b91033b0b99036b0b933b4b760811b81525061077e565b61039c565b6000818152600260209081526040808320546001600160a01b0316808452600190925290912054600160801b90046001600160401b031643811161039257610392826000836107c4565b50506001016102d6565b5060005b600354811015610439576103d46040518060400160405280600881526020016723b0b9b632b33a1d60c11b8152505a610739565b6103fd604051806040016040528060078152602001662a34b1b5b2b29d60c91b81525082610739565b620124f85a101561043d576104396040518060400160405280601081526020016f2ab73232b91033b0b99036b0b933b4b760811b81525061077e565b5050565b6000818152600260209081526040808320546001600160a01b03168084526001909252909120546001600160401b031661047961138882610d9f565b5a1061050a5760408051600481526024810182526020810180516001600160e01b0316633eaf5d9f60e01b17905290516000916001600160a01b0385169184916104c291610dd6565b60006040518083038160008787f1925050503d8060008114610500576040519150601f19603f3d011682016040523d82523d6000602084013e610505565b606091505b505050505b50506001016103a0565b6000546001600160a01b0316331461055d5760405162461bcd60e51b815260206004820152600c60248201526b15539055551213d49256915160a21b60448201526064016102c3565b6105688383836107c4565b505050565b6000806105997f0000000000000000000000000000000000000000000000000000000000000000610a96565b9050806001600160a01b031663d1f578947f0000000000000000000000000000000000000000000000000000000000000000856040516020016105dc9190610df2565b6040516020818303038152906040526040518363ffffffff1660e01b8152600401610608929190610e6b565b600060405180830381600087803b15801561062257600080fd5b505af1158015610636573d6000803e3d6000fd5b5050505061065f8184516207a12061064e9190610e97565b61065a61070843610d9f565b6107c4565b7f6d9f5f843298227fedb5ae27fcf3ebf729b71a00cdae9de0122e48a4aed64f17818533326040516106949493929190610eae565b60405180910390a190505b92915050565b6000546001600160a01b031633146106ee5760405162461bcd60e51b815260206004820152600c60248201526b15539055551213d49256915160a21b60448201526064016102c3565b600080546001600160a01b0319166001600160a01b0383169081178255604051909133917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a350565b610439828260405160240161074f929190610ee9565b60408051601f198184030181529190526020810180516001600160e01b0316632d839cb360e21b179052610b08565b6107c1816040516024016107929190610f0b565b60408051601f198184030181529190526020810180516001600160e01b031663104c13eb60e21b179052610b08565b50565b6107fa6040518060400160405280601281526020017129b2ba3a34b7339033b0b99030b63637b19d60711b815250848484610b11565b6001600160a01b038316600090815260016020526040902060055481546004546001600160401b0390911690610831908690610d9f565b61083b9190610f1e565b106108975760405162461bcd60e51b815260206004820152602660248201527f5469636b4d61737465723a2067617320616c6c6f636174696f6e2065786365656044820152650c8e640dac2f60d31b60648201526084016102c3565b8054600480546001600160401b03909216916000906108b7908490610f1e565b9250508190555082600460008282546108d09190610d9f565b909155505060008390036109a05780546001600160401b03166000036108f65750505050565b6003805490600061090683610f31565b909155505060035460009081526002602090815260408083205484546001600160401b03600160401b918290048116865283862080546001600160a01b0319166001600160a01b03909416938417905586549286526001909452919093208054938290049092160267ffffffffffffffff60401b1990921691909117905580546fffffffffffffffffffffffffffffffff19168155610a4c565b80546001600160401b0316600003610a14576003805490819060006109c483610f48565b9091555050815467ffffffffffffffff60401b1916600160401b6001600160401b03831602178255600090815260026020526040902080546001600160a01b0319166001600160a01b0386161790555b80546001600160401b03838116600160801b0277ffffffffffffffff0000000000000000ffffffffffffffff19909216908516171781555b60408051848152602081018490526001600160a01b038616917fde5be304e9fb13da67e61f6d156dd2aa96789f8e81a9a690e6d4e434fcb6cb35910160405180910390a250505050565b6000763d602d80600a3d3981f3363d3d373d3d3d363d730000008260601b60e81c176000526e5af43d82803e903d91602b57fd5bf38260781b17602052603760096000f090506001600160a01b038116610b03576040516330be1a3d60e21b815260040160405180910390fd5b919050565b6107c181610b60565b610b5a84848484604051602401610b2b9493929190610f61565b60408051601f198184030181529190526020810180516001600160e01b0316637c7a8d8f60e11b179052610b08565b50505050565b80516a636f6e736f6c652e6c6f67602083016000808483855afa5050505050565b80356001600160a01b0381168114610b0357600080fd5b600060208284031215610baa57600080fd5b610bb382610b81565b9392505050565b600060208284031215610bcc57600080fd5b5035919050565b600080600060608486031215610be857600080fd5b610bf184610b81565b95602085013595506040909401359392505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b0381118282101715610c4457610c44610c06565b604052919050565b600082601f830112610c5d57600080fd5b813560206001600160401b03821115610c7857610c78610c06565b8160051b610c87828201610c1c565b9283528481018201928281019087851115610ca157600080fd5b83870192505b84831015610cc757610cb883610b81565b82529183019190830190610ca7565b979650505050505050565b60008060408385031215610ce557600080fd5b82356001600160401b0380821115610cfc57600080fd5b818501915085601f830112610d1057600080fd5b8135602082821115610d2457610d24610c06565b610d36601f8301601f19168201610c1c565b8281528882848701011115610d4a57600080fd5b82828601838301376000928101820192909252909450850135915080821115610d7257600080fd5b50610d7f85828601610c4c565b9150509250929050565b634e487b7160e01b600052601160045260246000fd5b8082018082111561069f5761069f610d89565b60005b83811015610dcd578181015183820152602001610db5565b50506000910152565b60008251610de8818460208701610db2565b9190910192915050565b6020808252825182820181905260009190848201906040850190845b81811015610e335783516001600160a01b031683529284019291840191600101610e0e565b50909695505050505050565b60008151808452610e57816020860160208601610db2565b601f01601f19169290920160200192915050565b6001600160a01b0383168152604060208201819052600090610e8f90830184610e3f565b949350505050565b808202811582820484141761069f5761069f610d89565b600060018060a01b03808716835260806020840152610ed06080840187610e3f565b9481166040840152929092166060909101525092915050565b604081526000610efc6040830185610e3f565b90508260208301529392505050565b602081526000610bb36020830184610e3f565b8181038181111561069f5761069f610d89565b600081610f4057610f40610d89565b506000190190565b600060018201610f5a57610f5a610d89565b5060010190565b608081526000610f746080830187610e3f565b6001600160a01b0395909516602083015250604081019290925260609091015291905056fea2646970667358221220f09c36084d96f1a2d45fb2790c4086a0c532487672ddbbae51c271734ca5c84564736f6c63430008190033",
}

// ContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractMetaData.ABI instead.
var ContractABI = ContractMetaData.ABI

// ContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractMetaData.Bin instead.
var ContractBin = ContractMetaData.Bin

// DeployContract deploys a new Ethereum contract, binding an instance of Contract to it.
func DeployContract(auth *bind.TransactOpts, backend bind.ContractBackend, _maxGasAllocation *big.Int, _gameImplementation common.Address, _coreImplementation common.Address) (common.Address, *types.Transaction, *Contract, error) {
	parsed, err := ContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractBin), backend, _maxGasAllocation, _gameImplementation, _coreImplementation)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Contract{ContractCaller: ContractCaller{contract: contract}, ContractTransactor: ContractTransactor{contract: contract}, ContractFilterer: ContractFilterer{contract: contract}}, nil
}

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

// CoreImplementation is a free data retrieval call binding the contract method 0x9d492a2c.
//
// Solidity: function coreImplementation() view returns(address)
func (_Contract *ContractCaller) CoreImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "coreImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CoreImplementation is a free data retrieval call binding the contract method 0x9d492a2c.
//
// Solidity: function coreImplementation() view returns(address)
func (_Contract *ContractSession) CoreImplementation() (common.Address, error) {
	return _Contract.Contract.CoreImplementation(&_Contract.CallOpts)
}

// CoreImplementation is a free data retrieval call binding the contract method 0x9d492a2c.
//
// Solidity: function coreImplementation() view returns(address)
func (_Contract *ContractCallerSession) CoreImplementation() (common.Address, error) {
	return _Contract.Contract.CoreImplementation(&_Contract.CallOpts)
}

// GameImplementation is a free data retrieval call binding the contract method 0xf6cf916e.
//
// Solidity: function gameImplementation() view returns(address)
func (_Contract *ContractCaller) GameImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "gameImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GameImplementation is a free data retrieval call binding the contract method 0xf6cf916e.
//
// Solidity: function gameImplementation() view returns(address)
func (_Contract *ContractSession) GameImplementation() (common.Address, error) {
	return _Contract.Contract.GameImplementation(&_Contract.CallOpts)
}

// GameImplementation is a free data retrieval call binding the contract method 0xf6cf916e.
//
// Solidity: function gameImplementation() view returns(address)
func (_Contract *ContractCallerSession) GameImplementation() (common.Address, error) {
	return _Contract.Contract.GameImplementation(&_Contract.CallOpts)
}

// GetAddressOf is a free data retrieval call binding the contract method 0x44b920b7.
//
// Solidity: function getAddressOf(uint256 idx) view returns(address)
func (_Contract *ContractCaller) GetAddressOf(opts *bind.CallOpts, idx *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getAddressOf", idx)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAddressOf is a free data retrieval call binding the contract method 0x44b920b7.
//
// Solidity: function getAddressOf(uint256 idx) view returns(address)
func (_Contract *ContractSession) GetAddressOf(idx *big.Int) (common.Address, error) {
	return _Contract.Contract.GetAddressOf(&_Contract.CallOpts, idx)
}

// GetAddressOf is a free data retrieval call binding the contract method 0x44b920b7.
//
// Solidity: function getAddressOf(uint256 idx) view returns(address)
func (_Contract *ContractCallerSession) GetAddressOf(idx *big.Int) (common.Address, error) {
	return _Contract.Contract.GetAddressOf(&_Contract.CallOpts, idx)
}

// GetGasAllocOf is a free data retrieval call binding the contract method 0x32f79bd8.
//
// Solidity: function getGasAllocOf(address t) view returns(uint256)
func (_Contract *ContractCaller) GetGasAllocOf(opts *bind.CallOpts, t common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getGasAllocOf", t)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGasAllocOf is a free data retrieval call binding the contract method 0x32f79bd8.
//
// Solidity: function getGasAllocOf(address t) view returns(uint256)
func (_Contract *ContractSession) GetGasAllocOf(t common.Address) (*big.Int, error) {
	return _Contract.Contract.GetGasAllocOf(&_Contract.CallOpts, t)
}

// GetGasAllocOf is a free data retrieval call binding the contract method 0x32f79bd8.
//
// Solidity: function getGasAllocOf(address t) view returns(uint256)
func (_Contract *ContractCallerSession) GetGasAllocOf(t common.Address) (*big.Int, error) {
	return _Contract.Contract.GetGasAllocOf(&_Contract.CallOpts, t)
}

// GetIndexOf is a free data retrieval call binding the contract method 0x017df522.
//
// Solidity: function getIndexOf(address t) view returns(uint256)
func (_Contract *ContractCaller) GetIndexOf(opts *bind.CallOpts, t common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "getIndexOf", t)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetIndexOf is a free data retrieval call binding the contract method 0x017df522.
//
// Solidity: function getIndexOf(address t) view returns(uint256)
func (_Contract *ContractSession) GetIndexOf(t common.Address) (*big.Int, error) {
	return _Contract.Contract.GetIndexOf(&_Contract.CallOpts, t)
}

// GetIndexOf is a free data retrieval call binding the contract method 0x017df522.
//
// Solidity: function getIndexOf(address t) view returns(uint256)
func (_Contract *ContractCallerSession) GetIndexOf(t common.Address) (*big.Int, error) {
	return _Contract.Contract.GetIndexOf(&_Contract.CallOpts, t)
}

// LastBlock is a free data retrieval call binding the contract method 0x806b984f.
//
// Solidity: function lastBlock() view returns(uint256)
func (_Contract *ContractCaller) LastBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "lastBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastBlock is a free data retrieval call binding the contract method 0x806b984f.
//
// Solidity: function lastBlock() view returns(uint256)
func (_Contract *ContractSession) LastBlock() (*big.Int, error) {
	return _Contract.Contract.LastBlock(&_Contract.CallOpts)
}

// LastBlock is a free data retrieval call binding the contract method 0x806b984f.
//
// Solidity: function lastBlock() view returns(uint256)
func (_Contract *ContractCallerSession) LastBlock() (*big.Int, error) {
	return _Contract.Contract.LastBlock(&_Contract.CallOpts)
}

// MaxGasAllocation is a free data retrieval call binding the contract method 0xb0cc1f01.
//
// Solidity: function maxGasAllocation() view returns(uint256)
func (_Contract *ContractCaller) MaxGasAllocation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "maxGasAllocation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxGasAllocation is a free data retrieval call binding the contract method 0xb0cc1f01.
//
// Solidity: function maxGasAllocation() view returns(uint256)
func (_Contract *ContractSession) MaxGasAllocation() (*big.Int, error) {
	return _Contract.Contract.MaxGasAllocation(&_Contract.CallOpts)
}

// MaxGasAllocation is a free data retrieval call binding the contract method 0xb0cc1f01.
//
// Solidity: function maxGasAllocation() view returns(uint256)
func (_Contract *ContractCallerSession) MaxGasAllocation() (*big.Int, error) {
	return _Contract.Contract.MaxGasAllocation(&_Contract.CallOpts)
}

// NActiveTickees is a free data retrieval call binding the contract method 0xdbc3352a.
//
// Solidity: function nActiveTickees() view returns(uint256)
func (_Contract *ContractCaller) NActiveTickees(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "nActiveTickees")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NActiveTickees is a free data retrieval call binding the contract method 0xdbc3352a.
//
// Solidity: function nActiveTickees() view returns(uint256)
func (_Contract *ContractSession) NActiveTickees() (*big.Int, error) {
	return _Contract.Contract.NActiveTickees(&_Contract.CallOpts)
}

// NActiveTickees is a free data retrieval call binding the contract method 0xdbc3352a.
//
// Solidity: function nActiveTickees() view returns(uint256)
func (_Contract *ContractCallerSession) NActiveTickees() (*big.Int, error) {
	return _Contract.Contract.NActiveTickees(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Contract *ContractCallerSession) Owner() (common.Address, error) {
	return _Contract.Contract.Owner(&_Contract.CallOpts)
}

// TotalGasAllocation is a free data retrieval call binding the contract method 0x2e327fef.
//
// Solidity: function totalGasAllocation() view returns(uint256)
func (_Contract *ContractCaller) TotalGasAllocation(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Contract.contract.Call(opts, &out, "totalGasAllocation")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalGasAllocation is a free data retrieval call binding the contract method 0x2e327fef.
//
// Solidity: function totalGasAllocation() view returns(uint256)
func (_Contract *ContractSession) TotalGasAllocation() (*big.Int, error) {
	return _Contract.Contract.TotalGasAllocation(&_Contract.CallOpts)
}

// TotalGasAllocation is a free data retrieval call binding the contract method 0x2e327fef.
//
// Solidity: function totalGasAllocation() view returns(uint256)
func (_Contract *ContractCallerSession) TotalGasAllocation() (*big.Int, error) {
	return _Contract.Contract.TotalGasAllocation(&_Contract.CallOpts)
}

// CreateGame is a paid mutator transaction binding the contract method 0xaca11313.
//
// Solidity: function createGame(string lobbyId, address[] _players) returns(address)
func (_Contract *ContractTransactor) CreateGame(opts *bind.TransactOpts, lobbyId string, _players []common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "createGame", lobbyId, _players)
}

// CreateGame is a paid mutator transaction binding the contract method 0xaca11313.
//
// Solidity: function createGame(string lobbyId, address[] _players) returns(address)
func (_Contract *ContractSession) CreateGame(lobbyId string, _players []common.Address) (*types.Transaction, error) {
	return _Contract.Contract.CreateGame(&_Contract.TransactOpts, lobbyId, _players)
}

// CreateGame is a paid mutator transaction binding the contract method 0xaca11313.
//
// Solidity: function createGame(string lobbyId, address[] _players) returns(address)
func (_Contract *ContractTransactorSession) CreateGame(lobbyId string, _players []common.Address) (*types.Transaction, error) {
	return _Contract.Contract.CreateGame(&_Contract.TransactOpts, lobbyId, _players)
}

// SetGasAlloc is a paid mutator transaction binding the contract method 0xa97547d7.
//
// Solidity: function setGasAlloc(address t, uint256 gas, uint256 evictionBlock) returns()
func (_Contract *ContractTransactor) SetGasAlloc(opts *bind.TransactOpts, t common.Address, gas *big.Int, evictionBlock *big.Int) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "setGasAlloc", t, gas, evictionBlock)
}

// SetGasAlloc is a paid mutator transaction binding the contract method 0xa97547d7.
//
// Solidity: function setGasAlloc(address t, uint256 gas, uint256 evictionBlock) returns()
func (_Contract *ContractSession) SetGasAlloc(t common.Address, gas *big.Int, evictionBlock *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetGasAlloc(&_Contract.TransactOpts, t, gas, evictionBlock)
}

// SetGasAlloc is a paid mutator transaction binding the contract method 0xa97547d7.
//
// Solidity: function setGasAlloc(address t, uint256 gas, uint256 evictionBlock) returns()
func (_Contract *ContractTransactorSession) SetGasAlloc(t common.Address, gas *big.Int, evictionBlock *big.Int) (*types.Transaction, error) {
	return _Contract.Contract.SetGasAlloc(&_Contract.TransactOpts, t, gas, evictionBlock)
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

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Contract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Contract *ContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Contract.Contract.TransferOwnership(&_Contract.TransactOpts, newOwner)
}

// ContractGameCreatedIterator is returned from FilterGameCreated and is used to iterate over the raw logs and unpacked data for GameCreated events raised by the Contract contract.
type ContractGameCreatedIterator struct {
	Event *ContractGameCreated // Event containing the contract specifics and raw log

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
func (it *ContractGameCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGameCreated)
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
		it.Event = new(ContractGameCreated)
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
func (it *ContractGameCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGameCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGameCreated represents a GameCreated event raised by the Contract contract.
type ContractGameCreated struct {
	GameAddress common.Address
	LobbyId     string
	Sender      common.Address
	Origin      common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGameCreated is a free log retrieval operation binding the contract event 0x6d9f5f843298227fedb5ae27fcf3ebf729b71a00cdae9de0122e48a4aed64f17.
//
// Solidity: event GameCreated(address gameAddress, string lobbyId, address sender, address origin)
func (_Contract *ContractFilterer) FilterGameCreated(opts *bind.FilterOpts) (*ContractGameCreatedIterator, error) {

	logs, sub, err := _Contract.contract.FilterLogs(opts, "GameCreated")
	if err != nil {
		return nil, err
	}
	return &ContractGameCreatedIterator{contract: _Contract.contract, event: "GameCreated", logs: logs, sub: sub}, nil
}

// WatchGameCreated is a free log subscription operation binding the contract event 0x6d9f5f843298227fedb5ae27fcf3ebf729b71a00cdae9de0122e48a4aed64f17.
//
// Solidity: event GameCreated(address gameAddress, string lobbyId, address sender, address origin)
func (_Contract *ContractFilterer) WatchGameCreated(opts *bind.WatchOpts, sink chan<- *ContractGameCreated) (event.Subscription, error) {

	logs, sub, err := _Contract.contract.WatchLogs(opts, "GameCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGameCreated)
				if err := _Contract.contract.UnpackLog(event, "GameCreated", log); err != nil {
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

// ParseGameCreated is a log parse operation binding the contract event 0x6d9f5f843298227fedb5ae27fcf3ebf729b71a00cdae9de0122e48a4aed64f17.
//
// Solidity: event GameCreated(address gameAddress, string lobbyId, address sender, address origin)
func (_Contract *ContractFilterer) ParseGameCreated(log types.Log) (*ContractGameCreated, error) {
	event := new(ContractGameCreated)
	if err := _Contract.contract.UnpackLog(event, "GameCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractGasAllocSetIterator is returned from FilterGasAllocSet and is used to iterate over the raw logs and unpacked data for GasAllocSet events raised by the Contract contract.
type ContractGasAllocSetIterator struct {
	Event *ContractGasAllocSet // Event containing the contract specifics and raw log

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
func (it *ContractGasAllocSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractGasAllocSet)
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
		it.Event = new(ContractGasAllocSet)
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
func (it *ContractGasAllocSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractGasAllocSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractGasAllocSet represents a GasAllocSet event raised by the Contract contract.
type ContractGasAllocSet struct {
	Tickee        common.Address
	Gas           *big.Int
	EvictionBlock *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterGasAllocSet is a free log retrieval operation binding the contract event 0xde5be304e9fb13da67e61f6d156dd2aa96789f8e81a9a690e6d4e434fcb6cb35.
//
// Solidity: event GasAllocSet(address indexed tickee, uint256 gas, uint256 evictionBlock)
func (_Contract *ContractFilterer) FilterGasAllocSet(opts *bind.FilterOpts, tickee []common.Address) (*ContractGasAllocSetIterator, error) {

	var tickeeRule []interface{}
	for _, tickeeItem := range tickee {
		tickeeRule = append(tickeeRule, tickeeItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "GasAllocSet", tickeeRule)
	if err != nil {
		return nil, err
	}
	return &ContractGasAllocSetIterator{contract: _Contract.contract, event: "GasAllocSet", logs: logs, sub: sub}, nil
}

// WatchGasAllocSet is a free log subscription operation binding the contract event 0xde5be304e9fb13da67e61f6d156dd2aa96789f8e81a9a690e6d4e434fcb6cb35.
//
// Solidity: event GasAllocSet(address indexed tickee, uint256 gas, uint256 evictionBlock)
func (_Contract *ContractFilterer) WatchGasAllocSet(opts *bind.WatchOpts, sink chan<- *ContractGasAllocSet, tickee []common.Address) (event.Subscription, error) {

	var tickeeRule []interface{}
	for _, tickeeItem := range tickee {
		tickeeRule = append(tickeeRule, tickeeItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "GasAllocSet", tickeeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractGasAllocSet)
				if err := _Contract.contract.UnpackLog(event, "GasAllocSet", log); err != nil {
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

// ParseGasAllocSet is a log parse operation binding the contract event 0xde5be304e9fb13da67e61f6d156dd2aa96789f8e81a9a690e6d4e434fcb6cb35.
//
// Solidity: event GasAllocSet(address indexed tickee, uint256 gas, uint256 evictionBlock)
func (_Contract *ContractFilterer) ParseGasAllocSet(log types.Log) (*ContractGasAllocSet, error) {
	event := new(ContractGasAllocSet)
	if err := _Contract.contract.UnpackLog(event, "GasAllocSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Contract contract.
type ContractOwnershipTransferredIterator struct {
	Event *ContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractOwnershipTransferred)
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
		it.Event = new(ContractOwnershipTransferred)
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
func (it *ContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractOwnershipTransferred represents a OwnershipTransferred event raised by the Contract contract.
type ContractOwnershipTransferred struct {
	User     common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed user, address indexed newOwner)
func (_Contract *ContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, user []common.Address, newOwner []common.Address) (*ContractOwnershipTransferredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.FilterLogs(opts, "OwnershipTransferred", userRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ContractOwnershipTransferredIterator{contract: _Contract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed user, address indexed newOwner)
func (_Contract *ContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ContractOwnershipTransferred, user []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Contract.contract.WatchLogs(opts, "OwnershipTransferred", userRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractOwnershipTransferred)
				if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed user, address indexed newOwner)
func (_Contract *ContractFilterer) ParseOwnershipTransferred(log types.Log) (*ContractOwnershipTransferred, error) {
	event := new(ContractOwnershipTransferred)
	if err := _Contract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
