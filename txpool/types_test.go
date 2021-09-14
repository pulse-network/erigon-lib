/*
   Copyright 2021 Erigon contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package txpool

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var txParseTests = []struct {
	payloadStr  string
	senderStr   string
	idHashStr   string
	signHashStr string
	nonce       uint64
}{
	// Legacy unprotected
	{payloadStr: "f86a808459682f0082520894fe3b557e8fb62b89f4916b721be55ceb828dbd73872386f26fc10000801ca0d22fc3eed9b9b9dbef9eec230aa3fb849eff60356c6b34e86155dca5c03554c7a05e3903d7375337f103cb9583d97a59dcca7472908c31614ae240c6a8311b02d6",
		senderStr: "fe3b557e8fb62b89f4916b721be55ceb828dbd73", idHashStr: "595e27a835cd79729ff1eeacec3120eeb6ed1464a04ec727aaca734ead961328",
		signHashStr: "e2b043ecdbcfed773fe7b5ffc2e23ec238081c77137134a06d71eedf9cdd81d3", nonce: 0},
	// Legacy protected (EIP-155) from calveras, with chainId 123
	{payloadStr: "f86d808459682f0082520894e80d2a018c813577f33f9e69387dc621206fb3a48856bc75e2d63100008082011aa04ae3cae463329a32573f4fbf1bd9b011f93aecf80e4185add4682a03ba4a4919a02b8f05f3f4858b0da24c93c2a65e51b2fbbecf5ffdf97c1f8cc1801f307dc107",
		idHashStr:   "f4a91979624effdb45d2ba012a7995c2652b62ebbeb08cdcab00f4923807aa8a",
		signHashStr: "ff44cf01ee9b831f09910309a689e8da83d19aa60bad325ee9154b7c25cf4de8", nonce: 0},
	{payloadStr: "f86780862d79883d2000825208945df9b87991262f6ba471f09758cde1c0fc1de734827a69801ca088ff6cf0fefd94db46111149ae4bfc179e9b94721fffd821d38d16464b3f71d0a045e0aff800961cfce805daef7016b9b675c137a6a41a548f7b60a3484c06a33a",
		senderStr: "a1e4380a3b1f749673e270229993ee55f35663b4", idHashStr: "5c504ed432cb51138bcf09aa5e8a410dd4a1e204ef84bfed1be16dfba1b22060",
		signHashStr: "19b1e28c14f33e74b96b88eba97d4a4fc8a97638d72e972310025b7e1189b049", nonce: 0},
}

// List of false-positive cases - which we didn't reject in the past
var txParseMustFailTests = []struct {
	payloadStr string
}{
	{payloadStr: "b87a02f877018302b96484b2d05e0085174876e800830186a0948437c23fdbd156c4f0e67350122fe5ec220f4c6a88019a1859dd06e93e80c001a0e7083f4aadd9bdc2988e6590e6dd2ee24d4aea58793dcf3d19755ff64744da9aa0495b617461ce01cdf86a1775ea421f45a2d82a6f0fa641e80076737d4dd37a96"},
	{payloadStr: "b86d02f86a7b80843b9aca00843b9aca0082520894e80d2a018c813577f33f9e69387dc621206fb3a48080c001a02c73a04cd144e5a84ceb6da942f83763c2682896b51f7922e2e2f9a524dd90b7a0235adda5f87a1d098e2739e40e83129ff82837c9042e6ad61d0481334dcb6f1a"},
	{payloadStr: "b903a301f9039f018218bf85105e34df0083048a949410a0847c2d170008ddca7c3a688124f49363003280b902e4c11695480000000000000000000000004b274e4a9af31c20ed4151769a88ffe63d9439960000000000000000000000008510211a852f0c5994051dd85eaef73112a82eb5000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000bad4de000000000000000000000000607816a600000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000001146aa2600000000000000000000000000000000000000000000000000000000000001bc9b000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee000000000000000000000000482579f93dc13e6b434e38b5a0447ca543d88a4600000000000000000000000000000000000000000000000000000000000000c42df546f40000000000000000000000004b274e4a9af31c20ed4151769a88ffe63d943996000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee0000000000000000000000007d93f93d41604572119e4be7757a7a4a43705f080000000000000000000000000000000000000000000000003782dace9d90000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000082b5a61569b5898ac347c82a594c86699f1981aa88ca46a6a00b8e4f27b3d17bdf3714e7c0ca6a8023b37cca556602fce7dc7daac3fcee1ab04bbb3b94c10dec301cc57266db6567aa073efaa1fa6669bdc6f0877b0aeab4e33d18cb08b8877f08931abf427f11bade042177db48ca956feb114d6f5d56d1f5889047189562ec545e1c000000000000000000000000000000000000000000000000000000000000f84ff7946856ccf24beb7ed77f1f24eee5742f7ece5557e2e1a00000000000000000000000000000000000000000000000000000000000000001d694b1dd690cc9af7bb1a906a9b5a94f94191cc553cec080a0d52f3dbcad3530e73fcef6f4a75329b569a8903bf6d8109a960901f251a37af3a00ecf570e0c0ffa6efdc6e6e49be764b6a1a77e47de7bb99e167544ffbbcd65bc"},
	{payloadStr: "b86e01f86b7b018203e882520894236ff1e97419ae93ad80cafbaa21220c5d78fb7d880de0b6b3a764000080c080a0987e3d8d0dcd86107b041e1dca2e0583118ff466ad71ad36a8465dd2a166ca2da02361c5018e63beea520321b290097cd749febc2f437c7cb41fdd085816742060"},
	{payloadStr: "b8d202f8cf7b038502540be40085174876e8008301869f94e77162b7d2ceb3625a4993bab557403a7b706f18865af3107a400080f85bf85994de0b295669a9fd93d5f28d9ec85e40f4cb697baef842a00000000000000000000000000000000000000000000000000000000000000003a0000000000000000000000000000000000000000000000000000000000000000780a0f73da48f3f5c9f324dfd28d106dcf911b53f33c92ae068cf6135352300e7291aa06ee83d0f59275d90000ac8cf912c6eb47261d244c9db19ffefc49e52869ff197"},
}

func TestParseTransactionRLPMustFail(t *testing.T) {
	ctx := NewTxParseContext()
	tx, txSender := &TxSlot{}, [20]byte{}
	for _, testCase := range txParseMustFailTests {
		payload := decodeHex(testCase.payloadStr)
		_, err := ctx.ParseTransaction(payload, 0, tx, txSender[:])
		require.Error(t, err)
	}
}
func TestParseTransactionRLP(t *testing.T) {
	ctx := NewTxParseContext()
	for i, tt := range txParseTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			require := require.New(t)
			var err error
			payload := decodeHex(tt.payloadStr)
			tx, txSender := &TxSlot{}, [20]byte{}
			parseEnd, err := ctx.ParseTransaction(payload, 0, tx, txSender[:])
			fmt.Printf("%x\n", payload)
			require.NoError(err)
			require.Equal(len(payload), parseEnd)
			if tt.signHashStr != "" {
				signHash := decodeHex(tt.signHashStr)
				if !bytes.Equal(signHash, ctx.sighash[:]) {
					t.Errorf("signHash expected %x, got %x", signHash, ctx.sighash)
				}
			}
			if tt.idHashStr != "" {
				idHash := decodeHex(tt.idHashStr)
				if !bytes.Equal(idHash, tx.idHash[:]) {
					t.Errorf("idHash expected %x, got %x", idHash, tx.idHash)
				}
			}
			if tt.senderStr != "" {
				expectSender := decodeHex(tt.senderStr)
				if !bytes.Equal(expectSender, txSender[:]) {
					t.Errorf("expectSender expected %x, got %x", expectSender, txSender)
				}
			}
			require.Equal(tt.nonce, tx.nonce)
		})
	}
}

func TestTxSlotsGrowth(t *testing.T) {
	assert := assert.New(t)
	s := &TxSlots{}
	s.Resize(11)
	assert.Equal(11, len(s.txs))
	assert.Equal(11, s.senders.Len())
	s.Resize(23)
	assert.Equal(23, len(s.txs))
	assert.Equal(23, s.senders.Len())

	s = &TxSlots{txs: make([]*TxSlot, 20), senders: make(Addresses, 20*20)}
	s.Resize(20)
	assert.Equal(20, len(s.txs))
	assert.Equal(20, s.senders.Len())
	s.Resize(23)
	assert.Equal(23, len(s.txs))
	assert.Equal(23, s.senders.Len())

	s.Resize(2)
	assert.Equal(2, len(s.txs))
	assert.Equal(2, s.senders.Len())
}
