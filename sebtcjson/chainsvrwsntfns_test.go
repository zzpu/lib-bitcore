// Copyright 2018 The box.la Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sebtcjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// TestChainSvrWsNtfns tests all of the chain server websocket-specific
// notifications marshal and unmarshal into valid results include handling of
// optional fields being omitted in the marshalled command, while optional
// fields with defaults have the default assigned on unmarshalled commands.
func TestChainSvrWsNtfns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		newNtfn      func() (interface{}, error)
		staticNtfn   func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "blockconnected",
			newNtfn: func() (interface{}, error) {
				return NewCmd("blockconnected", "123", 100000, 123456789)
			},
			staticNtfn: func() interface{} {
				return NewBlockConnectedNtfn("123", 100000, 123456789)
			},
			marshalled: `{"jsonrpc":"1.0","method":"blockconnected","params":["123",100000,123456789],"id":null}`,
			unmarshalled: &BlockConnectedNtfn{
				Hash:   "123",
				Height: 100000,
				Time:   123456789,
			},
		},
		{
			name: "blockdisconnected",
			newNtfn: func() (interface{}, error) {
				return NewCmd("blockdisconnected", "123", 100000, 123456789)
			},
			staticNtfn: func() interface{} {
				return NewBlockDisconnectedNtfn("123", 100000, 123456789)
			},
			marshalled: `{"jsonrpc":"1.0","method":"blockdisconnected","params":["123",100000,123456789],"id":null}`,
			unmarshalled: &BlockDisconnectedNtfn{
				Hash:   "123",
				Height: 100000,
				Time:   123456789,
			},
		},
		{
			name: "filteredblockconnected",
			newNtfn: func() (interface{}, error) {
				return NewCmd("filteredblockconnected", 100000, "header", []string{"tx0", "tx1"})
			},
			staticNtfn: func() interface{} {
				return NewFilteredBlockConnectedNtfn(100000, "header", []string{"tx0", "tx1"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"filteredblockconnected","params":[100000,"header",["tx0","tx1"]],"id":null}`,
			unmarshalled: &FilteredBlockConnectedNtfn{
				Height:        100000,
				Header:        "header",
				SubscribedTxs: []string{"tx0", "tx1"},
			},
		},
		{
			name: "filteredblockdisconnected",
			newNtfn: func() (interface{}, error) {
				return NewCmd("filteredblockdisconnected", 100000, "header")
			},
			staticNtfn: func() interface{} {
				return NewFilteredBlockDisconnectedNtfn(100000, "header")
			},
			marshalled: `{"jsonrpc":"1.0","method":"filteredblockdisconnected","params":[100000,"header"],"id":null}`,
			unmarshalled: &FilteredBlockDisconnectedNtfn{
				Height: 100000,
				Header: "header",
			},
		},
		{
			name: "recvtx",
			newNtfn: func() (interface{}, error) {
				return NewCmd("recvtx", "001122", `{"height":100000,"hash":"123","index":0,"time":12345678}`)
			},
			staticNtfn: func() interface{} {
				blockDetails := BlockDetails{
					Height: 100000,
					Hash:   "123",
					Index:  0,
					Time:   12345678,
				}
				return NewRecvTxNtfn("001122", &blockDetails)
			},
			marshalled: `{"jsonrpc":"1.0","method":"recvtx","params":["001122",{"height":100000,"hash":"123","index":0,"time":12345678}],"id":null}`,
			unmarshalled: &RecvTxNtfn{
				HexTx: "001122",
				Block: &BlockDetails{
					Height: 100000,
					Hash:   "123",
					Index:  0,
					Time:   12345678,
				},
			},
		},
		{
			name: "redeemingtx",
			newNtfn: func() (interface{}, error) {
				return NewCmd("redeemingtx", "001122", `{"height":100000,"hash":"123","index":0,"time":12345678}`)
			},
			staticNtfn: func() interface{} {
				blockDetails := BlockDetails{
					Height: 100000,
					Hash:   "123",
					Index:  0,
					Time:   12345678,
				}
				return NewRedeemingTxNtfn("001122", &blockDetails)
			},
			marshalled: `{"jsonrpc":"1.0","method":"redeemingtx","params":["001122",{"height":100000,"hash":"123","index":0,"time":12345678}],"id":null}`,
			unmarshalled: &RedeemingTxNtfn{
				HexTx: "001122",
				Block: &BlockDetails{
					Height: 100000,
					Hash:   "123",
					Index:  0,
					Time:   12345678,
				},
			},
		},
		{
			name: "rescanfinished",
			newNtfn: func() (interface{}, error) {
				return NewCmd("rescanfinished", "123", 100000, 12345678)
			},
			staticNtfn: func() interface{} {
				return NewRescanFinishedNtfn("123", 100000, 12345678)
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescanfinished","params":["123",100000,12345678],"id":null}`,
			unmarshalled: &RescanFinishedNtfn{
				Hash:   "123",
				Height: 100000,
				Time:   12345678,
			},
		},
		{
			name: "rescanprogress",
			newNtfn: func() (interface{}, error) {
				return NewCmd("rescanprogress", "123", 100000, 12345678)
			},
			staticNtfn: func() interface{} {
				return NewRescanProgressNtfn("123", 100000, 12345678)
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescanprogress","params":["123",100000,12345678],"id":null}`,
			unmarshalled: &RescanProgressNtfn{
				Hash:   "123",
				Height: 100000,
				Time:   12345678,
			},
		},
		{
			name: "txaccepted",
			newNtfn: func() (interface{}, error) {
				return NewCmd("txaccepted", "123", 1.5)
			},
			staticNtfn: func() interface{} {
				return NewTxAcceptedNtfn("123", 1.5)
			},
			marshalled: `{"jsonrpc":"1.0","method":"txaccepted","params":["123",1.5],"id":null}`,
			unmarshalled: &TxAcceptedNtfn{
				TxID:   "123",
				Amount: 1.5,
			},
		},
		{
			name: "txacceptedverbose",
			newNtfn: func() (interface{}, error) {
				return NewCmd("txacceptedverbose", `{"hex":"001122","txid":"123","version":1,"locktime":4294967295,"vin":null,"vout":null,"confirmations":0}`)
			},
			staticNtfn: func() interface{} {
				txResult := TxRawResult{
					Hex:           "001122",
					Txid:          "123",
					Version:       1,
					LockTime:      4294967295,
					Vin:           nil,
					Vout:          nil,
					Confirmations: 0,
				}
				return NewTxAcceptedVerboseNtfn(txResult)
			},
			marshalled: `{"jsonrpc":"1.0","method":"txacceptedverbose","params":[{"hex":"001122","txid":"123","version":1,"locktime":4294967295,"vin":null,"vout":null}],"id":null}`,
			unmarshalled: &TxAcceptedVerboseNtfn{
				RawTx: TxRawResult{
					Hex:           "001122",
					Txid:          "123",
					Version:       1,
					LockTime:      4294967295,
					Vin:           nil,
					Vout:          nil,
					Confirmations: 0,
				},
			},
		},
		{
			name: "relevanttxaccepted",
			newNtfn: func() (interface{}, error) {
				return NewCmd("relevanttxaccepted", "001122")
			},
			staticNtfn: func() interface{} {
				return NewRelevantTxAcceptedNtfn("001122")
			},
			marshalled: `{"jsonrpc":"1.0","method":"relevanttxaccepted","params":["001122"],"id":null}`,
			unmarshalled: &RelevantTxAcceptedNtfn{
				Transaction: "001122",
			},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the notification as created by the new static
		// creation function.  The ID is nil for notifications.
		marshalled, err := MarshalCmd(nil, test.staticNtfn())
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

		// Ensure the notification is created without error via the
		// generic new notification creation function.
		cmd, err := test.newNtfn()
		if err != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, err)
		}

		// Marshal the notification as created by the generic new
		// notification creation function.    The ID is nil for
		// notifications.
		marshalled, err = MarshalCmd(nil, cmd)
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

		var request Request
		if err := json.Unmarshal(marshalled, &request); err != nil {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, err)
			continue
		}

		cmd, err = UnmarshalCmd(&request)
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