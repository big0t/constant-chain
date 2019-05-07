package component

import (
	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/privacy"
)

type InstructionFromBeacon interface {
	GetStringFormat() ([]string, error)
}

type AcceptDCBBoardIns struct {
	BoardPaymentAddress []privacy.PaymentAddress
	StartAmountToken    uint64
}

type AcceptGOVBoardIns struct {
	BoardPaymentAddress []privacy.PaymentAddress
	StartAmountToken    uint64
}

type ShareRewardOldBoardIns struct {
	ChairPaymentAddress privacy.PaymentAddress
	VoterPaymentAddress privacy.PaymentAddress
	BoardType           common.BoardType
	AmountOfCoin        uint64
}

type TxSendBackTokenVoteFailIns struct {
	BoardType      common.BoardType
	PaymentAddress privacy.PaymentAddress
	Amount         uint64
	PropertyID     common.Hash
}

type RewardProposalVoterIns struct {
	ReceiverAddress *privacy.PaymentAddress
	Amount          uint64
	BoardType       common.BoardType
}

type RewardProposalSubmitterIns struct {
	ReceiverAddress *privacy.PaymentAddress
	Amount          uint64
	BoardType       common.BoardType
}

type TxSendBackTokenToOldSupporterIns struct {
	BoardType      common.BoardType
	PaymentAddress privacy.PaymentAddress
	Amount         uint64
	PropertyID     common.Hash
}
