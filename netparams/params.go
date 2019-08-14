// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2016-2017 The Utopia developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package netparams

import "github.com/UtopiaCoinOrg/ucd/chaincfg"

// Params is used to group parameters for various networks such as the main
// network and test networks.
type Params struct {
	*chaincfg.Params
	JSONRPCClientPort string
	JSONRPCServerPort string
	GRPCServerPort    string
}

// MainNetParams contains parameters specific running ucwallet and
// ucd on the main network (wire.MainNet).
var MainNetParams = Params{
	Params:            chaincfg.MainNetParams(),
	JSONRPCClientPort: "10509",
	JSONRPCServerPort: "10510",
	GRPCServerPort:    "10511",
}

// TestNet3Params contains parameters specific running ucwallet and
// ucd on the test network (version 3) (wire.TestNet3).
var TestNet3Params = Params{
	Params:            chaincfg.TestNet3Params(),
	JSONRPCClientPort: "11509",
	JSONRPCServerPort: "11510",
	GRPCServerPort:    "11511",
}

// SimNetParams contains parameters specific to the simulation test network
// (wire.SimNet).
var SimNetParams = Params{
	Params:            chaincfg.SimNetParams(),
	JSONRPCClientPort: "12509",
	JSONRPCServerPort: "12510",
	GRPCServerPort:    "12511",
}
