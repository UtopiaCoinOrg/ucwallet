// Copyright (c) 2019 The Utopia developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wallet

import (
	hdkeychain "github.com/Utopia/ucd/hdkeychain"
	"github.com/Utopia/ucd/ucutil"
)

// hd2to1 converts a hdkeychain2 extended key to the v1 API.
// An error check during string conversion is intentionally dropped for brevity.
func hd2to1(k2 *hdkeychain.ExtendedKey) *hdkeychain.ExtendedKey {
	k, _ := hdkeychain.NewKeyFromString(k2.String(), ucutil.ActiveNet)
	return k
}

// hd1to2 converts a v1 extended key to the v2 API.
// An error check during string conversion is intentionally dropped for brevity.
func hd1to2(k *hdkeychain.ExtendedKey, params hdkeychain.NetworkParams) *hdkeychain.ExtendedKey {
	k2, _ := hdkeychain.NewKeyFromString(k.String(), params)
	return k2
}
