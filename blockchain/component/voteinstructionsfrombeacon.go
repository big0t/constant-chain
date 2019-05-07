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

type AcceptProposalIns struct {
	BoardType common.BoardType
	TxID      common.Hash
	Voters    []privacy.PaymentAddress
	ShardID   byte
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

type KeepOldProposalIns struct {
	BoardType common.BoardType
}

type UpdateDCBConstitutionIns struct {
	SubmitProposalInfo SubmitProposalInfo
	DCBParams          DCBParams
	Voters             []privacy.PaymentAddress
}

type UpdateGOVConstitutionIns struct {
	SubmitProposalInfo SubmitProposalInfo
	GOVParams          GOVParams
	Voters             []privacy.PaymentAddress
}
