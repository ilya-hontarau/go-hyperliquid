package sign

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/vmihailenco/msgpack/v5"
)

func SignL1Action(address string, action any, timestamp int64, isMainnet bool, keyManager KeyManager) (byte, [32]byte, [32]byte) {
	hash := buildActionHash(action, "", timestamp)
	message := buildMessage(hash.Bytes(), isMainnet)
	return SignInner(keyManager, address, message, isMainnet)
}

func GetNetSource(isMain bool) string {
	if isMain {
		return "a"
	} else {
		return "b"
	}
}

func buildMessage(hash []byte, isMain bool) apitypes.TypedDataMessage {
	source := GetNetSource(isMain)
	return apitypes.TypedDataMessage{
		"source":       source,
		"connectionId": hash,
	}
}

func buildActionHash(action any, vaultAd string, nonce int64) common.Hash {
	var (
		data []byte
	)

	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	enc.UseCompactInts(true)
	err := enc.Encode(action)
	if err != nil {
		panic(fmt.Sprintf("Failed to pack the data %s", err))
	}
	data = buf.Bytes()

	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))
	data = ArrayAppend(data, nonceBytes)

	if vaultAd == "" {
		data = ArrayAppend(data, []byte("\x00"))
	} else {
		data = ArrayAppend(data, []byte("\x01"))
		data = ArrayAppend(data, HexToBytes(vaultAd))
	}

	result := crypto.Keccak256Hash(data)
	return result
}

func HexToBytes(addr string) []byte {
	if strings.HasPrefix(addr, "0x") {
		fAddr := strings.Replace(addr, "0x", "", 1)
		b, _ := hex.DecodeString(fAddr)
		return b
	} else {
		b, _ := hex.DecodeString(addr)
		return b
	}
}

func ArrayAppend(data []byte, toAppend []byte) []byte {
	return append(data, toAppend...)
}

func SignInner(km KeyManager, address string, message apitypes.TypedDataMessage, isMainNet bool) (byte, [32]byte, [32]byte) {
	signer := NewSigner(km)
	req := SigRequest{
		PrimaryType: "Agent",
		DType: []apitypes.Type{
			{
				Name: "source",
				Type: "string",
			},
			{
				Name: "connectionId",
				Type: "bytes32",
			},
		},
		DTypeMsg:  message,
		IsMainNet: isMainNet,
	}

	v, r, s, err := signer.Sign(address, req)
	if err != nil {
		panic("Failed to sign request")
	}

	return v, r, s

}
