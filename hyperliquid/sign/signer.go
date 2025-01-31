package sign

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

const ChainId = 1337
const VerifyingContract = "0x0000000000000000000000000000000000000000"

type KeyManager interface {
	GetKey(address string) *ecdsa.PrivateKey
}

type Signer struct {
	manager KeyManager
	secret  string
}

func NewSigner(manager KeyManager) Signer {
	return Signer{
		manager: manager,
	}
}

type SigRequest struct {
	PrimaryType string
	DType       []apitypes.Type
	DTypeMsg    map[string]interface{}
	IsMainNet   bool
}

func (req SigRequest) GetChainId() *math.HexOrDecimal256 {
	if req.PrimaryType == "HyperliquidTransaction:Withdraw" || req.PrimaryType == "Hyperliquid:UserPoints" {
		if req.IsMainNet {
			return math.NewHexOrDecimal256(int64(42161))
		} else {
			return math.NewHexOrDecimal256(int64(421614))
		}
	} else {
		return math.NewHexOrDecimal256(int64(1337))
	}
}

func (signer Signer) Sign(address string, req SigRequest) (byte, [32]byte, [32]byte, error) {
	var (
		err error
	)

	types := GetContractTypes(req)
	domain := GetDomain(req)
	typedData := apitypes.TypedData{
		Types:       types,
		PrimaryType: req.PrimaryType,
		Domain:      domain,
		Message:     req.DTypeMsg,
	}

	key := signer.manager.GetKey(address)

	bytes, _, err := apitypes.TypedDataAndHash(typedData)

	if err != nil {
		return 0, [32]byte{}, [32]byte{}, err
	}

	sig, err := crypto.Sign(bytes, key)

	if err != nil {
		return 0, [32]byte{}, [32]byte{}, err
	}

	return SigToVRS(sig)
}
