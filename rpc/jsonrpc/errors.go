// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2016-2018 The Utopia developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package jsonrpc

import (
	"fmt"

	"github.com/Utopia/ucd/ucjson"
	"github.com/Utopia/ucwallet/errors"
)

func convertError(err error) *ucjson.RPCError {
	if err, ok := err.(*ucjson.RPCError); ok {
		return err
	}

	code := ucjson.ErrRPCWallet
	if err, ok := err.(*errors.Error); ok {
		switch err.Kind {
		case errors.Bug:
			code = ucjson.ErrRPCInternal.Code
		case errors.Encoding:
			code = ucjson.ErrRPCInvalidParameter
		case errors.Locked:
			code = ucjson.ErrRPCWalletUnlockNeeded
		case errors.Passphrase:
			code = ucjson.ErrRPCWalletPassphraseIncorrect
		case errors.NoPeers:
			code = ucjson.ErrRPCClientNotConnected
		case errors.InsufficientBalance:
			code = ucjson.ErrRPCWalletInsufficientFunds
		}
	}
	return &ucjson.RPCError{
		Code:    code,
		Message: err.Error(),
	}
}

func rpcError(code ucjson.RPCErrorCode, err error) *ucjson.RPCError {
	return &ucjson.RPCError{
		Code:    code,
		Message: err.Error(),
	}
}

func rpcErrorf(code ucjson.RPCErrorCode, format string, args ...interface{}) *ucjson.RPCError {
	return &ucjson.RPCError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Errors variables that are defined once here to avoid duplication.
var (
	errUnloadedWallet = &ucjson.RPCError{
		Code:    ucjson.ErrRPCWallet,
		Message: "request requires a wallet but wallet has not loaded yet",
	}

	errRPCClientNotConnected = &ucjson.RPCError{
		Code:    ucjson.ErrRPCClientNotConnected,
		Message: "disconnected from consensus RPC",
	}

	errNoNetwork = &ucjson.RPCError{
		Code:    ucjson.ErrRPCClientNotConnected,
		Message: "disconnected from network",
	}

	errAccountNotFound = &ucjson.RPCError{
		Code:    ucjson.ErrRPCWalletInvalidAccountName,
		Message: "account not found",
	}

	errAddressNotInWallet = &ucjson.RPCError{
		Code:    ucjson.ErrRPCWallet,
		Message: "address not found in wallet",
	}

	errNotImportedAccount = &ucjson.RPCError{
		Code:    ucjson.ErrRPCWallet,
		Message: "imported addresses must belong to the imported account",
	}

	errNeedPositiveAmount = &ucjson.RPCError{
		Code:    ucjson.ErrRPCInvalidParameter,
		Message: "amount must be positive",
	}

	errWalletUnlockNeeded = &ucjson.RPCError{
		Code:    ucjson.ErrRPCWalletUnlockNeeded,
		Message: "enter the wallet passphrase with walletpassphrase first",
	}

	errReservedAccountName = &ucjson.RPCError{
		Code:    ucjson.ErrRPCInvalidParameter,
		Message: "account name is reserved by RPC server",
	}
)
