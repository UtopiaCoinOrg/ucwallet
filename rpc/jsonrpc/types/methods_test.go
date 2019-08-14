// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2015-2019 The Utopia developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/UtopiaCoinOrg/ucd/ucjson"
	ucdtypes "github.com/UtopiaCoinOrg/ucd/rpc/jsonrpc/types"
)

// TestWalletSvrCmds tests all of the wallet server commands marshal and
// unmarshal into valid results include handling of optional fields being
// omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestWalletSvrCmds(t *testing.T) {
	t.Parallel()

	testID := int(1)
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "addmultisigaddress",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return NewAddMultisigAddressCmd(2, keys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   nil,
			},
		},
		{
			name: "addmultisigaddress optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("addmultisigaddress", 2, []string{"031234", "035678"}, "test")
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return NewAddMultisigAddressCmd(2, keys, ucjson.String("test"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"addmultisigaddress","params":[2,["031234","035678"],"test"],"id":1}`,
			unmarshalled: &AddMultisigAddressCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
				Account:   ucjson.String("test"),
			},
		},
		{
			name: "createmultisig",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("createmultisig", 2, []string{"031234", "035678"})
			},
			staticCmd: func() interface{} {
				keys := []string{"031234", "035678"}
				return NewCreateMultisigCmd(2, keys)
			},
			marshalled: `{"jsonrpc":"1.0","method":"createmultisig","params":[2,["031234","035678"]],"id":1}`,
			unmarshalled: &CreateMultisigCmd{
				NRequired: 2,
				Keys:      []string{"031234", "035678"},
			},
		},
		{
			name: "createnewaccount",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("createnewaccount", "acct")
			},
			staticCmd: func() interface{} {
				return NewCreateNewAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"createnewaccount","params":["acct"],"id":1}`,
			unmarshalled: &CreateNewAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "dumpprivkey",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("dumpprivkey", "1Address")
			},
			staticCmd: func() interface{} {
				return NewDumpPrivKeyCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpprivkey","params":["1Address"],"id":1}`,
			unmarshalled: &DumpPrivKeyCmd{
				Address: "1Address",
			},
		},
		{
			name: "estimatepriority",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("estimatepriority", 6)
			},
			staticCmd: func() interface{} {
				return NewEstimatePriorityCmd(6)
			},
			marshalled: `{"jsonrpc":"1.0","method":"estimatepriority","params":[6],"id":1}`,
			unmarshalled: &EstimatePriorityCmd{
				NumBlocks: 6,
			},
		},
		{
			name: "getaccount",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getaccount", "1Address")
			},
			staticCmd: func() interface{} {
				return NewGetAccountCmd("1Address")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccount","params":["1Address"],"id":1}`,
			unmarshalled: &GetAccountCmd{
				Address: "1Address",
			},
		},
		{
			name: "getaccountaddress",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getaccountaddress", "acct")
			},
			staticCmd: func() interface{} {
				return NewGetAccountAddressCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaccountaddress","params":["acct"],"id":1}`,
			unmarshalled: &GetAccountAddressCmd{
				Account: "acct",
			},
		},
		{
			name: "getaddressesbyaccount",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getaddressesbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return NewGetAddressesByAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"getaddressesbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &GetAddressesByAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "getbalance",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getbalance")
			},
			staticCmd: func() interface{} {
				return NewGetBalanceCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":[],"id":1}`,
			unmarshalled: &GetBalanceCmd{
				Account: nil,
				MinConf: ucjson.Int(1),
			},
		},
		{
			name: "getbalance optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getbalance", "acct")
			},
			staticCmd: func() interface{} {
				return NewGetBalanceCmd(ucjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct"],"id":1}`,
			unmarshalled: &GetBalanceCmd{
				Account: ucjson.String("acct"),
				MinConf: ucjson.Int(1),
			},
		},
		{
			name: "getbalance optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getbalance", "acct", 6)
			},
			staticCmd: func() interface{} {
				return NewGetBalanceCmd(ucjson.String("acct"), ucjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getbalance","params":["acct",6],"id":1}`,
			unmarshalled: &GetBalanceCmd{
				Account: ucjson.String("acct"),
				MinConf: ucjson.Int(6),
			},
		},
		{
			name: "getnewaddress",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getnewaddress")
			},
			staticCmd: func() interface{} {
				return NewGetNewAddressCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":[],"id":1}`,
			unmarshalled: &GetNewAddressCmd{
				Account:   nil,
				GapPolicy: nil,
			},
		},
		{
			name: "getnewaddress optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getnewaddress", "acct", "ignore")
			},
			staticCmd: func() interface{} {
				return NewGetNewAddressCmd(ucjson.String("acct"), ucjson.String("ignore"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getnewaddress","params":["acct","ignore"],"id":1}`,
			unmarshalled: &GetNewAddressCmd{
				Account:   ucjson.String("acct"),
				GapPolicy: ucjson.String("ignore"),
			},
		},
		{
			name: "getrawchangeaddress",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getrawchangeaddress")
			},
			staticCmd: func() interface{} {
				return NewGetRawChangeAddressCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":[],"id":1}`,
			unmarshalled: &GetRawChangeAddressCmd{
				Account: nil,
			},
		},
		{
			name: "getrawchangeaddress optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getrawchangeaddress", "acct")
			},
			staticCmd: func() interface{} {
				return NewGetRawChangeAddressCmd(ucjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getrawchangeaddress","params":["acct"],"id":1}`,
			unmarshalled: &GetRawChangeAddressCmd{
				Account: ucjson.String("acct"),
			},
		},
		{
			name: "getreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getreceivedbyaccount", "acct")
			},
			staticCmd: func() interface{} {
				return NewGetReceivedByAccountCmd("acct", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct"],"id":1}`,
			unmarshalled: &GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: ucjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaccount optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getreceivedbyaccount", "acct", 6)
			},
			staticCmd: func() interface{} {
				return NewGetReceivedByAccountCmd("acct", ucjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaccount","params":["acct",6],"id":1}`,
			unmarshalled: &GetReceivedByAccountCmd{
				Account: "acct",
				MinConf: ucjson.Int(6),
			},
		},
		{
			name: "getreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getreceivedbyaddress", "1Address")
			},
			staticCmd: func() interface{} {
				return NewGetReceivedByAddressCmd("1Address", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address"],"id":1}`,
			unmarshalled: &GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: ucjson.Int(1),
			},
		},
		{
			name: "getreceivedbyaddress optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("getreceivedbyaddress", "1Address", 6)
			},
			staticCmd: func() interface{} {
				return NewGetReceivedByAddressCmd("1Address", ucjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getreceivedbyaddress","params":["1Address",6],"id":1}`,
			unmarshalled: &GetReceivedByAddressCmd{
				Address: "1Address",
				MinConf: ucjson.Int(6),
			},
		},
		{
			name: "gettransaction",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("gettransaction", "123")
			},
			staticCmd: func() interface{} {
				return NewGetTransactionCmd("123", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123"],"id":1}`,
			unmarshalled: &GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "gettransaction optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("gettransaction", "123", true)
			},
			staticCmd: func() interface{} {
				return NewGetTransactionCmd("123", ucjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"gettransaction","params":["123",true],"id":1}`,
			unmarshalled: &GetTransactionCmd{
				Txid:             "123",
				IncludeWatchOnly: ucjson.Bool(true),
			},
		},
		{
			name: "importprivkey",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("importprivkey", "abc")
			},
			staticCmd: func() interface{} {
				return NewImportPrivKeyCmd("abc", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc"],"id":1}`,
			unmarshalled: &ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   nil,
				Rescan:  ucjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("importprivkey", "abc", "label")
			},
			staticCmd: func() interface{} {
				return NewImportPrivKeyCmd("abc", ucjson.String("label"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label"],"id":1}`,
			unmarshalled: &ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   ucjson.String("label"),
				Rescan:  ucjson.Bool(true),
			},
		},
		{
			name: "importprivkey optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("importprivkey", "abc", "label", false)
			},
			staticCmd: func() interface{} {
				return NewImportPrivKeyCmd("abc", ucjson.String("label"), ucjson.Bool(false), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false],"id":1}`,
			unmarshalled: &ImportPrivKeyCmd{
				PrivKey: "abc",
				Label:   ucjson.String("label"),
				Rescan:  ucjson.Bool(false),
			},
		},
		{
			name: "importprivkey optional3",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("importprivkey", "abc", "label", false, 12345)
			},
			staticCmd: func() interface{} {
				return NewImportPrivKeyCmd("abc", ucjson.String("label"), ucjson.Bool(false), ucjson.Int(12345))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importprivkey","params":["abc","label",false,12345],"id":1}`,
			unmarshalled: &ImportPrivKeyCmd{
				PrivKey:  "abc",
				Label:    ucjson.String("label"),
				Rescan:   ucjson.Bool(false),
				ScanFrom: ucjson.Int(12345),
			},
		},
		{
			name: "keypoolrefill",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("keypoolrefill")
			},
			staticCmd: func() interface{} {
				return NewKeyPoolRefillCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[],"id":1}`,
			unmarshalled: &KeyPoolRefillCmd{
				NewSize: ucjson.Uint(100),
			},
		},
		{
			name: "keypoolrefill optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("keypoolrefill", 200)
			},
			staticCmd: func() interface{} {
				return NewKeyPoolRefillCmd(ucjson.Uint(200))
			},
			marshalled: `{"jsonrpc":"1.0","method":"keypoolrefill","params":[200],"id":1}`,
			unmarshalled: &KeyPoolRefillCmd{
				NewSize: ucjson.Uint(200),
			},
		},
		{
			name: "listaccounts",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listaccounts")
			},
			staticCmd: func() interface{} {
				return NewListAccountsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[],"id":1}`,
			unmarshalled: &ListAccountsCmd{
				MinConf: ucjson.Int(1),
			},
		},
		{
			name: "listaccounts optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listaccounts", 6)
			},
			staticCmd: func() interface{} {
				return NewListAccountsCmd(ucjson.Int(6))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaccounts","params":[6],"id":1}`,
			unmarshalled: &ListAccountsCmd{
				MinConf: ucjson.Int(6),
			},
		},
		{
			name: "listlockunspent",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listlockunspent")
			},
			staticCmd: func() interface{} {
				return NewListLockUnspentCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"listlockunspent","params":[],"id":1}`,
			unmarshalled: &ListLockUnspentCmd{},
		},
		{
			name: "listreceivedbyaccount",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listreceivedbyaccount")
			},
			staticCmd: func() interface{} {
				return NewListReceivedByAccountCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[],"id":1}`,
			unmarshalled: &ListReceivedByAccountCmd{
				MinConf:          ucjson.Int(1),
				IncludeEmpty:     ucjson.Bool(false),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listreceivedbyaccount", 6)
			},
			staticCmd: func() interface{} {
				return NewListReceivedByAccountCmd(ucjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6],"id":1}`,
			unmarshalled: &ListReceivedByAccountCmd{
				MinConf:          ucjson.Int(6),
				IncludeEmpty:     ucjson.Bool(false),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listreceivedbyaccount", 6, true)
			},
			staticCmd: func() interface{} {
				return NewListReceivedByAccountCmd(ucjson.Int(6), ucjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true],"id":1}`,
			unmarshalled: &ListReceivedByAccountCmd{
				MinConf:          ucjson.Int(6),
				IncludeEmpty:     ucjson.Bool(true),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaccount optional3",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listreceivedbyaccount", 6, true, false)
			},
			staticCmd: func() interface{} {
				return NewListReceivedByAccountCmd(ucjson.Int(6), ucjson.Bool(true), ucjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaccount","params":[6,true,false],"id":1}`,
			unmarshalled: &ListReceivedByAccountCmd{
				MinConf:          ucjson.Int(6),
				IncludeEmpty:     ucjson.Bool(true),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listreceivedbyaddress")
			},
			staticCmd: func() interface{} {
				return NewListReceivedByAddressCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[],"id":1}`,
			unmarshalled: &ListReceivedByAddressCmd{
				MinConf:          ucjson.Int(1),
				IncludeEmpty:     ucjson.Bool(false),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listreceivedbyaddress", 6)
			},
			staticCmd: func() interface{} {
				return NewListReceivedByAddressCmd(ucjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6],"id":1}`,
			unmarshalled: &ListReceivedByAddressCmd{
				MinConf:          ucjson.Int(6),
				IncludeEmpty:     ucjson.Bool(false),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listreceivedbyaddress", 6, true)
			},
			staticCmd: func() interface{} {
				return NewListReceivedByAddressCmd(ucjson.Int(6), ucjson.Bool(true), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true],"id":1}`,
			unmarshalled: &ListReceivedByAddressCmd{
				MinConf:          ucjson.Int(6),
				IncludeEmpty:     ucjson.Bool(true),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listreceivedbyaddress optional3",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listreceivedbyaddress", 6, true, false)
			},
			staticCmd: func() interface{} {
				return NewListReceivedByAddressCmd(ucjson.Int(6), ucjson.Bool(true), ucjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listreceivedbyaddress","params":[6,true,false],"id":1}`,
			unmarshalled: &ListReceivedByAddressCmd{
				MinConf:          ucjson.Int(6),
				IncludeEmpty:     ucjson.Bool(true),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listsinceblock",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listsinceblock")
			},
			staticCmd: func() interface{} {
				return NewListSinceBlockCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":[],"id":1}`,
			unmarshalled: &ListSinceBlockCmd{
				BlockHash:           nil,
				TargetConfirmations: ucjson.Int(1),
				IncludeWatchOnly:    ucjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listsinceblock", "123")
			},
			staticCmd: func() interface{} {
				return NewListSinceBlockCmd(ucjson.String("123"), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123"],"id":1}`,
			unmarshalled: &ListSinceBlockCmd{
				BlockHash:           ucjson.String("123"),
				TargetConfirmations: ucjson.Int(1),
				IncludeWatchOnly:    ucjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listsinceblock", "123", 6)
			},
			staticCmd: func() interface{} {
				return NewListSinceBlockCmd(ucjson.String("123"), ucjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6],"id":1}`,
			unmarshalled: &ListSinceBlockCmd{
				BlockHash:           ucjson.String("123"),
				TargetConfirmations: ucjson.Int(6),
				IncludeWatchOnly:    ucjson.Bool(false),
			},
		},
		{
			name: "listsinceblock optional3",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listsinceblock", "123", 6, true)
			},
			staticCmd: func() interface{} {
				return NewListSinceBlockCmd(ucjson.String("123"), ucjson.Int(6), ucjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listsinceblock","params":["123",6,true],"id":1}`,
			unmarshalled: &ListSinceBlockCmd{
				BlockHash:           ucjson.String("123"),
				TargetConfirmations: ucjson.Int(6),
				IncludeWatchOnly:    ucjson.Bool(true),
			},
		},
		{
			name: "listtransactions",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listtransactions")
			},
			staticCmd: func() interface{} {
				return NewListTransactionsCmd(nil, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":[],"id":1}`,
			unmarshalled: &ListTransactionsCmd{
				Account:          nil,
				Count:            ucjson.Int(10),
				From:             ucjson.Int(0),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listtransactions", "acct")
			},
			staticCmd: func() interface{} {
				return NewListTransactionsCmd(ucjson.String("acct"), nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct"],"id":1}`,
			unmarshalled: &ListTransactionsCmd{
				Account:          ucjson.String("acct"),
				Count:            ucjson.Int(10),
				From:             ucjson.Int(0),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listtransactions", "acct", 20)
			},
			staticCmd: func() interface{} {
				return NewListTransactionsCmd(ucjson.String("acct"), ucjson.Int(20), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20],"id":1}`,
			unmarshalled: &ListTransactionsCmd{
				Account:          ucjson.String("acct"),
				Count:            ucjson.Int(20),
				From:             ucjson.Int(0),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional3",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listtransactions", "acct", 20, 1)
			},
			staticCmd: func() interface{} {
				return NewListTransactionsCmd(ucjson.String("acct"), ucjson.Int(20),
					ucjson.Int(1), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1],"id":1}`,
			unmarshalled: &ListTransactionsCmd{
				Account:          ucjson.String("acct"),
				Count:            ucjson.Int(20),
				From:             ucjson.Int(1),
				IncludeWatchOnly: ucjson.Bool(false),
			},
		},
		{
			name: "listtransactions optional4",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listtransactions", "acct", 20, 1, true)
			},
			staticCmd: func() interface{} {
				return NewListTransactionsCmd(ucjson.String("acct"), ucjson.Int(20),
					ucjson.Int(1), ucjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listtransactions","params":["acct",20,1,true],"id":1}`,
			unmarshalled: &ListTransactionsCmd{
				Account:          ucjson.String("acct"),
				Count:            ucjson.Int(20),
				From:             ucjson.Int(1),
				IncludeWatchOnly: ucjson.Bool(true),
			},
		},
		{
			name: "listunspent",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listunspent")
			},
			staticCmd: func() interface{} {
				return NewListUnspentCmd(nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[],"id":1}`,
			unmarshalled: &ListUnspentCmd{
				MinConf:   ucjson.Int(1),
				MaxConf:   ucjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listunspent", 6)
			},
			staticCmd: func() interface{} {
				return NewListUnspentCmd(ucjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6],"id":1}`,
			unmarshalled: &ListUnspentCmd{
				MinConf:   ucjson.Int(6),
				MaxConf:   ucjson.Int(9999999),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listunspent", 6, 100)
			},
			staticCmd: func() interface{} {
				return NewListUnspentCmd(ucjson.Int(6), ucjson.Int(100), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100],"id":1}`,
			unmarshalled: &ListUnspentCmd{
				MinConf:   ucjson.Int(6),
				MaxConf:   ucjson.Int(100),
				Addresses: nil,
			},
		},
		{
			name: "listunspent optional3",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("listunspent", 6, 100, []string{"1Address", "1Address2"})
			},
			staticCmd: func() interface{} {
				return NewListUnspentCmd(ucjson.Int(6), ucjson.Int(100),
					&[]string{"1Address", "1Address2"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"listunspent","params":[6,100,["1Address","1Address2"]],"id":1}`,
			unmarshalled: &ListUnspentCmd{
				MinConf:   ucjson.Int(6),
				MaxConf:   ucjson.Int(100),
				Addresses: &[]string{"1Address", "1Address2"},
			},
		},
		{
			name: "lockunspent",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("lockunspent", true, `[{"txid":"123","vout":1}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []ucdtypes.TransactionInput{
					{Txid: "123", Vout: 1},
				}
				return NewLockUnspentCmd(true, txInputs)
			},
			marshalled: `{"jsonrpc":"1.0","method":"lockunspent","params":[true,[{"txid":"123","vout":1,"tree":0}]],"id":1}`,
			unmarshalled: &LockUnspentCmd{
				Unlock: true,
				Transactions: []ucdtypes.TransactionInput{
					{Txid: "123", Vout: 1},
				},
			},
		},
		{
			name: "renameaccount",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("renameaccount", "oldacct", "newacct")
			},
			staticCmd: func() interface{} {
				return NewRenameAccountCmd("oldacct", "newacct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"renameaccount","params":["oldacct","newacct"],"id":1}`,
			unmarshalled: &RenameAccountCmd{
				OldAccount: "oldacct",
				NewAccount: "newacct",
			},
		},
		{
			name: "sendfrom",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendfrom", "from", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return NewSendFromCmd("from", "1Address", 0.5, nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5],"id":1}`,
			unmarshalled: &SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     ucjson.Int(1),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6)
			},
			staticCmd: func() interface{} {
				return NewSendFromCmd("from", "1Address", 0.5, ucjson.Int(6), nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6],"id":1}`,
			unmarshalled: &SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     ucjson.Int(6),
				Comment:     nil,
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment")
			},
			staticCmd: func() interface{} {
				return NewSendFromCmd("from", "1Address", 0.5, ucjson.Int(6),
					ucjson.String("comment"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment"],"id":1}`,
			unmarshalled: &SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     ucjson.Int(6),
				Comment:     ucjson.String("comment"),
				CommentTo:   nil,
			},
		},
		{
			name: "sendfrom optional3",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendfrom", "from", "1Address", 0.5, 6, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return NewSendFromCmd("from", "1Address", 0.5, ucjson.Int(6),
					ucjson.String("comment"), ucjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendfrom","params":["from","1Address",0.5,6,"comment","commentto"],"id":1}`,
			unmarshalled: &SendFromCmd{
				FromAccount: "from",
				ToAddress:   "1Address",
				Amount:      0.5,
				MinConf:     ucjson.Int(6),
				Comment:     ucjson.String("comment"),
				CommentTo:   ucjson.String("commentto"),
			},
		},
		{
			name: "sendmany",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendmany", "from", `{"1Address":0.5}`)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return NewSendManyCmd("from", amounts, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5}],"id":1}`,
			unmarshalled: &SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     ucjson.Int(1),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6)
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return NewSendManyCmd("from", amounts, ucjson.Int(6), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6],"id":1}`,
			unmarshalled: &SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     ucjson.Int(6),
				Comment:     nil,
			},
		},
		{
			name: "sendmany optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendmany", "from", `{"1Address":0.5}`, 6, "comment")
			},
			staticCmd: func() interface{} {
				amounts := map[string]float64{"1Address": 0.5}
				return NewSendManyCmd("from", amounts, ucjson.Int(6), ucjson.String("comment"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendmany","params":["from",{"1Address":0.5},6,"comment"],"id":1}`,
			unmarshalled: &SendManyCmd{
				FromAccount: "from",
				Amounts:     map[string]float64{"1Address": 0.5},
				MinConf:     ucjson.Int(6),
				Comment:     ucjson.String("comment"),
			},
		},
		{
			name: "sendtoaddress",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendtoaddress", "1Address", 0.5)
			},
			staticCmd: func() interface{} {
				return NewSendToAddressCmd("1Address", 0.5, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5],"id":1}`,
			unmarshalled: &SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   nil,
				CommentTo: nil,
			},
		},
		{
			name: "sendtoaddress optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sendtoaddress", "1Address", 0.5, "comment", "commentto")
			},
			staticCmd: func() interface{} {
				return NewSendToAddressCmd("1Address", 0.5, ucjson.String("comment"),
					ucjson.String("commentto"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sendtoaddress","params":["1Address",0.5,"comment","commentto"],"id":1}`,
			unmarshalled: &SendToAddressCmd{
				Address:   "1Address",
				Amount:    0.5,
				Comment:   ucjson.String("comment"),
				CommentTo: ucjson.String("commentto"),
			},
		},
		{
			name: "settxfee",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("settxfee", 0.0001)
			},
			staticCmd: func() interface{} {
				return NewSetTxFeeCmd(0.0001)
			},
			marshalled: `{"jsonrpc":"1.0","method":"settxfee","params":[0.0001],"id":1}`,
			unmarshalled: &SetTxFeeCmd{
				Amount: 0.0001,
			},
		},
		{
			name: "signmessage",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("signmessage", "1Address", "message")
			},
			staticCmd: func() interface{} {
				return NewSignMessageCmd("1Address", "message")
			},
			marshalled: `{"jsonrpc":"1.0","method":"signmessage","params":["1Address","message"],"id":1}`,
			unmarshalled: &SignMessageCmd{
				Address: "1Address",
				Message: "message",
			},
		},
		{
			name: "signrawtransaction",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("signrawtransaction", "001122")
			},
			staticCmd: func() interface{} {
				return NewSignRawTransactionCmd("001122", nil, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122"],"id":1}`,
			unmarshalled: &SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   nil,
				PrivKeys: nil,
				Flags:    ucjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional1",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("signrawtransaction", "001122", `[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]`)
			},
			staticCmd: func() interface{} {
				txInputs := []RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				}

				return NewSignRawTransactionCmd("001122", &txInputs, nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[{"txid":"123","vout":1,"tree":0,"scriptPubKey":"00","redeemScript":"01"}]],"id":1}`,
			unmarshalled: &SignRawTransactionCmd{
				RawTx: "001122",
				Inputs: &[]RawTxInput{
					{
						Txid:         "123",
						Vout:         1,
						ScriptPubKey: "00",
						RedeemScript: "01",
					},
				},
				PrivKeys: nil,
				Flags:    ucjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional2",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("signrawtransaction", "001122", `[]`, `["abc"]`)
			},
			staticCmd: func() interface{} {
				txInputs := []RawTxInput{}
				privKeys := []string{"abc"}
				return NewSignRawTransactionCmd("001122", &txInputs, &privKeys, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],["abc"]],"id":1}`,
			unmarshalled: &SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]RawTxInput{},
				PrivKeys: &[]string{"abc"},
				Flags:    ucjson.String("ALL"),
			},
		},
		{
			name: "signrawtransaction optional3",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("signrawtransaction", "001122", `[]`, `[]`, "ALL")
			},
			staticCmd: func() interface{} {
				txInputs := []RawTxInput{}
				privKeys := []string{}
				return NewSignRawTransactionCmd("001122", &txInputs, &privKeys,
					ucjson.String("ALL"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"signrawtransaction","params":["001122",[],[],"ALL"],"id":1}`,
			unmarshalled: &SignRawTransactionCmd{
				RawTx:    "001122",
				Inputs:   &[]RawTxInput{},
				PrivKeys: &[]string{},
				Flags:    ucjson.String("ALL"),
			},
		},
		{
			name: "sweepaccount - optionals provided",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sweepaccount", "default", "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu", 6, 0.05)
			},
			staticCmd: func() interface{} {
				return NewSweepAccountCmd("default", "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",
					func(i uint32) *uint32 { return &i }(6),
					func(i float64) *float64 { return &i }(0.05))
			},
			marshalled: `{"jsonrpc":"1.0","method":"sweepaccount","params":["default","DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",6,0.05],"id":1}`,
			unmarshalled: &SweepAccountCmd{
				SourceAccount:         "default",
				DestinationAddress:    "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",
				RequiredConfirmations: func(i uint32) *uint32 { return &i }(6),
				FeePerKb:              func(i float64) *float64 { return &i }(0.05),
			},
		},
		{
			name: "sweepaccount - optionals omitted",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("sweepaccount", "default", "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu")
			},
			staticCmd: func() interface{} {
				return NewSweepAccountCmd("default", "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu", nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"sweepaccount","params":["default","DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu"],"id":1}`,
			unmarshalled: &SweepAccountCmd{
				SourceAccount:      "default",
				DestinationAddress: "DsUZxxoHJSty8DCfwfartwTYbuhmVct7tJu",
			},
		},
		{
			name: "verifyseed",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("verifyseed", "abc")
			},
			staticCmd: func() interface{} {
				return NewVerifySeedCmd("abc", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifyseed","params":["abc"],"id":1}`,
			unmarshalled: &VerifySeedCmd{
				Seed:    "abc",
				Account: nil,
			},
		},
		{
			name: "verifyseed optional",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("verifyseed", "abc", 5)
			},
			staticCmd: func() interface{} {
				account := ucjson.Uint32(5)
				return NewVerifySeedCmd("abc", account)
			},
			marshalled: `{"jsonrpc":"1.0","method":"verifyseed","params":["abc",5],"id":1}`,
			unmarshalled: &VerifySeedCmd{
				Seed:    "abc",
				Account: ucjson.Uint32(5),
			},
		},
		{
			name: "walletlock",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("walletlock")
			},
			staticCmd: func() interface{} {
				return NewWalletLockCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletlock","params":[],"id":1}`,
			unmarshalled: &WalletLockCmd{},
		},
		{
			name: "walletpassphrase",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("walletpassphrase", "pass", 60)
			},
			staticCmd: func() interface{} {
				return NewWalletPassphraseCmd("pass", 60)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrase","params":["pass",60],"id":1}`,
			unmarshalled: &WalletPassphraseCmd{
				Passphrase: "pass",
				Timeout:    60,
			},
		},
		{
			name: "walletpassphrasechange",
			newCmd: func() (interface{}, error) {
				return ucjson.NewCmd("walletpassphrasechange", "old", "new")
			},
			staticCmd: func() interface{} {
				return NewWalletPassphraseChangeCmd("old", "new")
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletpassphrasechange","params":["old","new"],"id":1}`,
			unmarshalled: &WalletPassphraseChangeCmd{
				OldPassphrase: "old",
				NewPassphrase: "new",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command
		// creation function.
		marshalled, err := ucjson.MarshalCmd("1.0", testID, test.staticCmd())
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		// Ensure the command is created without error via the generic
		// new command creation function.
		cmd, err := test.newCmd()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected ucjson.NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the command as created by the generic new command
		// creation function.
		marshalled, err = ucjson.MarshalCmd("1.0", testID, cmd)
		if err != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled)
			continue
		}

		var request ucjson.Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = ucjson.ParseParams(request.Method, request.Params)
		if err != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, err)
			continue
		}

		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled))
			continue
		}
	}
}
