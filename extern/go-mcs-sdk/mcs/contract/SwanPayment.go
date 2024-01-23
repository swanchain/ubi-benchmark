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
)

// IPaymentMinimallockPaymentParam is an auto generated low-level Go binding around an user-defined struct.
type IPaymentMinimallockPaymentParam struct {
	Id         string
	MinPayment *big.Int
	Amount     *big.Int
	LockTime   *big.Int
	Recipient  common.Address
	Size       *big.Int
	CopyLimit  uint8
}

// SwanPaymentMetaData contains all meta data concerning the SwanPayment contract.
var SwanPaymentMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"minPayment\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockTime\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"size\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"copyLimit\",\"type\":\"uint8\"}],\"internalType\":\"structIPaymentMinimal.lockPaymentParam\",\"name\":\"param\",\"type\":\"tuple\"}],\"name\":\"lockTokenPayment\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SwanPaymentABI is the input ABI used to generate the binding from.
// Deprecated: Use SwanPaymentMetaData.ABI instead.
var SwanPaymentABI = SwanPaymentMetaData.ABI

// SwanPayment is an auto generated Go binding around an Ethereum contract.
type SwanPayment struct {
	SwanPaymentCaller     // Read-only binding to the contract
	SwanPaymentTransactor // Write-only binding to the contract
	SwanPaymentFilterer   // Log filterer for contract events
}

// SwanPaymentCaller is an auto generated read-only Go binding around an Ethereum contract.
type SwanPaymentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwanPaymentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SwanPaymentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwanPaymentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SwanPaymentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SwanPaymentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SwanPaymentSession struct {
	Contract     *SwanPayment      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SwanPaymentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SwanPaymentCallerSession struct {
	Contract *SwanPaymentCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// SwanPaymentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SwanPaymentTransactorSession struct {
	Contract     *SwanPaymentTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// SwanPaymentRaw is an auto generated low-level Go binding around an Ethereum contract.
type SwanPaymentRaw struct {
	Contract *SwanPayment // Generic contract binding to access the raw methods on
}

// SwanPaymentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SwanPaymentCallerRaw struct {
	Contract *SwanPaymentCaller // Generic read-only contract binding to access the raw methods on
}

// SwanPaymentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SwanPaymentTransactorRaw struct {
	Contract *SwanPaymentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSwanPayment creates a new instance of SwanPayment, bound to a specific deployed contract.
func NewSwanPayment(address common.Address, backend bind.ContractBackend) (*SwanPayment, error) {
	contract, err := bindSwanPayment(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SwanPayment{SwanPaymentCaller: SwanPaymentCaller{contract: contract}, SwanPaymentTransactor: SwanPaymentTransactor{contract: contract}, SwanPaymentFilterer: SwanPaymentFilterer{contract: contract}}, nil
}

// NewSwanPaymentCaller creates a new read-only instance of SwanPayment, bound to a specific deployed contract.
func NewSwanPaymentCaller(address common.Address, caller bind.ContractCaller) (*SwanPaymentCaller, error) {
	contract, err := bindSwanPayment(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SwanPaymentCaller{contract: contract}, nil
}

// NewSwanPaymentTransactor creates a new write-only instance of SwanPayment, bound to a specific deployed contract.
func NewSwanPaymentTransactor(address common.Address, transactor bind.ContractTransactor) (*SwanPaymentTransactor, error) {
	contract, err := bindSwanPayment(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SwanPaymentTransactor{contract: contract}, nil
}

// NewSwanPaymentFilterer creates a new log filterer instance of SwanPayment, bound to a specific deployed contract.
func NewSwanPaymentFilterer(address common.Address, filterer bind.ContractFilterer) (*SwanPaymentFilterer, error) {
	contract, err := bindSwanPayment(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SwanPaymentFilterer{contract: contract}, nil
}

// bindSwanPayment binds a generic wrapper to an already deployed contract.
func bindSwanPayment(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SwanPaymentABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwanPayment *SwanPaymentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwanPayment.Contract.SwanPaymentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwanPayment *SwanPaymentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwanPayment.Contract.SwanPaymentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwanPayment *SwanPaymentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwanPayment.Contract.SwanPaymentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SwanPayment *SwanPaymentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SwanPayment.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SwanPayment *SwanPaymentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SwanPayment.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SwanPayment *SwanPaymentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SwanPayment.Contract.contract.Transact(opts, method, params...)
}

// LockTokenPayment is a paid mutator transaction binding the contract method 0xf4d98717.
//
// Solidity: function lockTokenPayment((string,uint256,uint256,uint256,address,uint256,uint8) param) returns(bool)
func (_SwanPayment *SwanPaymentTransactor) LockTokenPayment(opts *bind.TransactOpts, param IPaymentMinimallockPaymentParam) (*types.Transaction, error) {
	return _SwanPayment.contract.Transact(opts, "lockTokenPayment", param)
}

// LockTokenPayment is a paid mutator transaction binding the contract method 0xf4d98717.
//
// Solidity: function lockTokenPayment((string,uint256,uint256,uint256,address,uint256,uint8) param) returns(bool)
func (_SwanPayment *SwanPaymentSession) LockTokenPayment(param IPaymentMinimallockPaymentParam) (*types.Transaction, error) {
	return _SwanPayment.Contract.LockTokenPayment(&_SwanPayment.TransactOpts, param)
}

// LockTokenPayment is a paid mutator transaction binding the contract method 0xf4d98717.
//
// Solidity: function lockTokenPayment((string,uint256,uint256,uint256,address,uint256,uint8) param) returns(bool)
func (_SwanPayment *SwanPaymentTransactorSession) LockTokenPayment(param IPaymentMinimallockPaymentParam) (*types.Transaction, error) {
	return _SwanPayment.Contract.LockTokenPayment(&_SwanPayment.TransactOpts, param)
}
