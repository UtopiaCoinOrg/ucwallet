// Copyright (c) 2016 The btcsuite developers
// Copyright (c) 2016 The Utopia developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package txrules

import (
	"github.com/UtopiaCoinOrg/ucd/txscript"
	"github.com/UtopiaCoinOrg/ucd/ucutil"
	"github.com/UtopiaCoinOrg/ucd/wire"
	"github.com/UtopiaCoinOrg/ucwallet/errors"
	h "github.com/UtopiaCoinOrg/ucwallet/internal/helpers"
)

// DefaultRelayFeePerKb is the default minimum relay fee policy for a mempool.
const DefaultRelayFeePerKb ucutil.Amount = 1e4

// IsDustAmount determines whether a transaction output value and script length would
// cause the output to be considered dust.  Transactions with dust outputs are
// not standard and are rejected by mempools with default policies.
func IsDustAmount(amount ucutil.Amount, scriptSize int, relayFeePerKb ucutil.Amount) bool {
	// Calculate the total (estimated) cost to the network.  This is
	// calculated using the serialize size of the output plus the serial
	// size of a transaction input which redeems it.  The output is assumed
	// to be compressed P2PKH as this is the most common script type.  Use
	// the average size of a compressed P2PKH redeem input (165) rather than
	// the largest possible (txsizes.RedeemP2PKHInputSize).
	totalSize := 8 + 2 + wire.VarIntSerializeSize(uint64(scriptSize)) +
		scriptSize + 165

	if amount > 1e12 {
		return false
	}
	// Dust is defined as an output value where the total cost to the network
	// (output size + input size) is greater than 1/3 of the relay fee.
	return int64(amount)*1000/(3*int64(totalSize)) < int64(relayFeePerKb)
}

// IsDustOutput determines whether a transaction output is considered dust.
// Transactions with dust outputs are not standard and are rejected by mempools
// with default policies.
func IsDustOutput(output *wire.TxOut, relayFeePerKb ucutil.Amount) bool {
	// Unspendable outputs which solely carry data are not checked for dust.
	if txscript.GetScriptClass(output.Version, output.PkScript) == txscript.NullDataTy {
		return false
	}

	// All other unspendable outputs are considered dust.
	if txscript.IsUnspendable(output.Value, output.PkScript) {
		return true
	}

	return IsDustAmount(ucutil.Amount(output.Value), len(output.PkScript),
		relayFeePerKb)
}

// CheckOutput performs simple consensus and policy tests on a transaction
// output.  Returns with errors.Invalid if output violates consensus rules, and
// errors.Policy if the output violates a non-consensus policy.
func CheckOutput(output *wire.TxOut, relayFeePerKb ucutil.Amount) error {
	if output.Value < 0 {
		return errors.E(errors.Invalid, "transaction output amount is negative")
	}
	if output.Value > ucutil.MaxAmount {
		return errors.E(errors.Invalid, "transaction output amount exceeds maximum value")
	}
	if IsDustOutput(output, relayFeePerKb) {
		return errors.E(errors.Policy, "transaction output is dust")
	}
	return nil
}

// FeeForSerializeSize calculates the required fee for a transaction of some
// arbitrary size given a mempool's relay fee policy.
func FeeForSerializeSize(relayFeePerKb ucutil.Amount, txSerializeSize int) ucutil.Amount {
	fee := relayFeePerKb * ucutil.Amount(txSerializeSize) / 1000

	if fee == 0 && relayFeePerKb > 0 {
		fee = relayFeePerKb
	}

	if fee < 0 || fee > ucutil.MaxAmount {
		fee = ucutil.MaxAmount
	}

	return fee
}

// PaysHighFees checks whether the signed transaction pays insanely high fees.
// Transactons are defined to have a high fee if they have pay a fee rate that
// is 1000 time higher than the default fee.
func PaysHighFees(totalInput ucutil.Amount, tx *wire.MsgTx, changeInedx int ) bool {
	fee := totalInput - h.SumOutputValues(tx.TxOut)
	if fee <= 0 {
		// Impossible to determine
		return false
	}
	if changeInedx != -1{
		isFlashTx :=false
		totalValue := int64(0)
		flashFee :=int64(0)
		for i,out:=range tx.TxOut{
			totalValue+=out.Value
			if _,has:=txscript.HaveFlashTxTag(out.PkScript);has{
				isFlashTx=true
			}
			if isFlashTx && changeInedx == i{
				totalValue -= out.Value
			}
		}
		if isFlashTx{
			flashFee = totalValue/1000.0
		}
		maxFee := FeeForSerializeSize(1000*DefaultRelayFeePerKb, tx.SerializeSize())
		return fee > maxFee + ucutil.Amount(flashFee)
	}else{
		maxFee := FeeForSerializeSize(1000*DefaultRelayFeePerKb, tx.SerializeSize())
		return fee > maxFee
	}
}
