package jsonresult

import (
	"github.com/constant-money/constant-chain/privacy/zeroknowledge"
	"math/big"
)

type TransactionDetail struct {
	BlockHash   string `json:"BlockHash"`
	BlockHeight uint64 `json:"BlockHeight"`
	Index       uint64 `json:"index"`
	ShardID     byte   `json:"shardID"`
	Hash        string `json:"Hash"`
	Version     int8   `json:"Version"`
	Type        string `json:"Type"` // Transaction type
	LockTime    string `json:"LockTime"`
	Fee         uint64 `json:"Fee"` // Fee applies: always consant
	Image       string `json:"Image"`

	IsPrivacy       bool              `json:"IsPrivacy"`
	Proof           *zkp.PaymentProof `json:"Proof"`
	ProofDetail     ProofDetail       `json:"ProofDetail"`
	InputCoinPubKey string            `json:"InputCoinPubKey"`
	SigPubKey       []byte            `json:"SigPubKey,omitempty"` // 64 bytes
	Sig             []byte            `json:"Sig,omitempty"`       // 64 bytes

	Metadata               string `json:"Metadata"`
	CustomTokenData        string `json:"CustomTokenData"`
	PrivacyCustomTokenData string `json:"PrivacyCustomTokenData"`

	IsInMempool bool `json:"IsInMempool"`
	IsInBlock   bool `json:"IsInBlock"`
}

type ProofDetail struct {
	InputCoins  []*CoinDetail
	OutputCoins []*CoinDetail
}

func (proofDetail *ProofDetail) ConvertFromProof(proof *zkp.PaymentProof) {
	proofDetail.InputCoins = make([]*CoinDetail, 0)
	for _, input := range proof.InputCoins {
		in := CoinDetail{
			CoinDetails: Coin{},
		}
		if input.CoinDetails != nil {
			in.CoinDetails.Value = input.CoinDetails.Value
			in.CoinDetails.Info = input.CoinDetails.Info
			if input.CoinDetails.CoinCommitment != nil {
				in.CoinDetails.CoinCommitment = input.CoinDetails.CoinCommitment.Compress()
			}
			if input.CoinDetails.Randomness != nil {
				in.CoinDetails.Randomness = *input.CoinDetails.Randomness
			}
			if input.CoinDetails.SNDerivator != nil {
				in.CoinDetails.SNDerivator = *input.CoinDetails.SNDerivator
			}
			if input.CoinDetails.SerialNumber != nil {
				in.CoinDetails.SerialNumber = input.CoinDetails.SerialNumber.Compress()
			}
			if input.CoinDetails.PublicKey != nil {
				in.CoinDetails.PublicKey = input.CoinDetails.PublicKey.Compress()
			}
		}
		proofDetail.InputCoins = append(proofDetail.InputCoins, &in)
	}

	for _, output := range proof.OutputCoins {
		out := CoinDetail{
			CoinDetails: Coin{},
		}
		if output.CoinDetails != nil {
			out.CoinDetails.Value = output.CoinDetails.Value
			out.CoinDetails.Info = output.CoinDetails.Info
			if output.CoinDetails.CoinCommitment != nil {
				out.CoinDetails.CoinCommitment = output.CoinDetails.CoinCommitment.Compress()
			}
			if output.CoinDetails.Randomness != nil {
				out.CoinDetails.Randomness = *output.CoinDetails.Randomness
			}
			if output.CoinDetails.SNDerivator != nil {
				out.CoinDetails.SNDerivator = *output.CoinDetails.SNDerivator
			}
			if output.CoinDetails.SerialNumber != nil {
				out.CoinDetails.SerialNumber = output.CoinDetails.SerialNumber.Compress()
			}
			if output.CoinDetails.PublicKey != nil {
				out.CoinDetails.PublicKey = output.CoinDetails.PublicKey.Compress()
			}
			if output.CoinDetailsEncrypted != nil {
				out.CoinDetailsEncrypted = output.CoinDetailsEncrypted.Bytes()
			}
		}
		proofDetail.OutputCoins = append(proofDetail.OutputCoins, &out)
	}
}

type CoinDetail struct {
	CoinDetails          Coin
	CoinDetailsEncrypted []byte
}

type Coin struct {
	PublicKey      []byte
	CoinCommitment []byte
	SNDerivator    big.Int
	SerialNumber   []byte
	Randomness     big.Int
	Value          uint64
	Info           []byte
}
