package component

import (
	"encoding/json"
	"strconv"

	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/privacy"
)

func NewTxSendBackTokenToOldSupporterIns(
	boardType common.BoardType,
	paymentAddress privacy.PaymentAddress,
	amount uint64,
	propertyID common.Hash,
) *TxSendBackTokenToOldSupporterIns {
	return &TxSendBackTokenToOldSupporterIns{
		BoardType:      boardType,
		PaymentAddress: paymentAddress,
		Amount:         amount,
		PropertyID:     propertyID,
	}
}

func NewTxSendBackTokenVoteFailIns(
	boardType common.BoardType,
	paymentAddress privacy.PaymentAddress,
	amount uint64,
	propertyID common.Hash,
) *TxSendBackTokenVoteFailIns {
	return &TxSendBackTokenVoteFailIns{
		BoardType:      boardType,
		PaymentAddress: paymentAddress,
		Amount:         amount,
		PropertyID:     propertyID,
	}
}

func NewRewardProposalSubmitterIns(
	receiverAddress *privacy.PaymentAddress,
	amount uint64,
	boardType common.BoardType,
) *RewardProposalSubmitterIns {
	return &RewardProposalSubmitterIns{
		ReceiverAddress: receiverAddress,
		Amount:          amount,
		BoardType:       boardType,
	}
}

func NewRewardProposalVoterIns(
	receiverAddress *privacy.PaymentAddress,
	amount uint64,
	boardType common.BoardType,
) *RewardProposalVoterIns {
	return &RewardProposalVoterIns{
		ReceiverAddress: receiverAddress,
		Amount:          amount,
		BoardType:       boardType,
	}
}

func NewShareRewardOldBoardMetadataIns(
	chairPaymentAddress privacy.PaymentAddress,
	voterPaymentAddress privacy.PaymentAddress,
	boardType common.BoardType,
	amountOfCoin uint64,
) *ShareRewardOldBoardIns {
	return &ShareRewardOldBoardIns{
		ChairPaymentAddress: chairPaymentAddress,
		VoterPaymentAddress: voterPaymentAddress,
		BoardType:           boardType,
		AmountOfCoin:        amountOfCoin,
	}
}

func NewAcceptBoardIns(
	boardType common.BoardType,
	boardPaymentAddress []privacy.PaymentAddress,
	startAmountToken uint64,
) InstructionFromBeacon {
	if boardType == common.DCBBoard {
		return NewAcceptDCBBoardIns(
			boardPaymentAddress,
			startAmountToken,
		)
	} else {
		return NewAcceptGOVBoardIns(
			boardPaymentAddress,
			startAmountToken,
		)
	}
}

func NewAcceptDCBBoardIns(
	boardPaymentAddress []privacy.PaymentAddress,
	startAmountToken uint64,
) *AcceptDCBBoardIns {
	return &AcceptDCBBoardIns{
		BoardPaymentAddress: boardPaymentAddress,
		StartAmountToken:    startAmountToken,
	}
}

func NewAcceptGOVBoardIns(
	boardPaymentAddress []privacy.PaymentAddress,
	startAmountToken uint64,
) *AcceptGOVBoardIns {
	return &AcceptGOVBoardIns{
		BoardPaymentAddress: boardPaymentAddress,
		StartAmountToken:    startAmountToken,
	}
}

func NewAcceptProposalIns(
	boardType common.BoardType,
	txID common.Hash,
	voters []privacy.PaymentAddress,
	shardID byte,
) *AcceptProposalIns {
	return &AcceptProposalIns{
		BoardType: boardType,
		TxID:      txID,
		Voters:    voters,
		ShardID:   shardID,
	}
}

func NewKeepOldProposalIns(
	boardType common.BoardType,
) *KeepOldProposalIns {
	return &KeepOldProposalIns{
		BoardType: boardType,
	}
}

func NewUpdateDCBConstitutionIns(
	submitProposalInfo SubmitProposalInfo,
	DCBParams DCBParams,
	voters []privacy.PaymentAddress,
) *UpdateDCBConstitutionIns {
	return &UpdateDCBConstitutionIns{
		SubmitProposalInfo: submitProposalInfo,
		DCBParams:          DCBParams,
		Voters:             voters,
	}
}

func NewUpdateGOVConstitutionIns(
	submitProposalInfo SubmitProposalInfo,
	GOVParams GOVParams,
	voters []privacy.PaymentAddress,
) *UpdateGOVConstitutionIns {
	return &UpdateGOVConstitutionIns{
		SubmitProposalInfo: submitProposalInfo,
		GOVParams:          GOVParams,
		Voters:             voters,
	}
}

func NewUpdateDCBConstitutionInsFromStr(inst []string) (*UpdateDCBConstitutionIns, error) {
	updateDCBConstitutionIns := &UpdateDCBConstitutionIns{}
	err := json.Unmarshal([]byte(inst[2]), updateDCBConstitutionIns)
	if err != nil {
		return nil, err
	}
	return updateDCBConstitutionIns, nil
}

func NewKeepOldProposalInsFromStr(inst []string) (*KeepOldProposalIns, error) {
	keepOldProposalIns := &KeepOldProposalIns{}
	err := json.Unmarshal([]byte(inst[2]), keepOldProposalIns)
	if err != nil {
		return nil, err
	}
	return keepOldProposalIns, nil
}

func NewUpdateGOVConstitutionInsFromStr(inst []string) (*UpdateGOVConstitutionIns, error) {
	updateGOVConstitutionIns := &UpdateGOVConstitutionIns{}
	err := json.Unmarshal([]byte(inst[2]), updateGOVConstitutionIns)
	if err != nil {
		return nil, err
	}
	return updateGOVConstitutionIns, nil
}

func (acceptDCBBoardIns *AcceptDCBBoardIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(acceptDCBBoardIns)
	if err != nil {
		return nil, err
	}
	return []string{
		strconv.Itoa(AcceptDCBBoardInsType),
		strconv.Itoa(-1),
		string(content),
	}, nil
}

func (acceptGOVBoardIns *AcceptGOVBoardIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(acceptGOVBoardIns)
	if err != nil {
		return nil, err
	}
	return []string{
		strconv.Itoa(AcceptGOVBoardInsType),
		strconv.Itoa(-1),
		string(content),
	}, nil
}

func (txSendBackTokenVoteFailIns *TxSendBackTokenVoteFailIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(txSendBackTokenVoteFailIns)
	if err != nil {
		return nil, err
	}
	shardID := GetShardIDFromPaymentAddressBytes(txSendBackTokenVoteFailIns.PaymentAddress)
	return []string{
		strconv.Itoa(SendBackTokenVoteBoardFailIns),
		strconv.Itoa(int(shardID)),
		string(content),
	}, nil
}

func (txSendBackTokenToOldSupporterIns *TxSendBackTokenToOldSupporterIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(txSendBackTokenToOldSupporterIns)
	if err != nil {
		return nil, err
	}
	shardID := GetShardIDFromPaymentAddressBytes(txSendBackTokenToOldSupporterIns.PaymentAddress)
	return []string{
		strconv.Itoa(SendBackTokenToOldSupporterIns),
		strconv.Itoa(int(shardID)),
		string(content),
	}, nil
}

func (shareRewardOldBoardIns *ShareRewardOldBoardIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(shareRewardOldBoardIns)
	if err != nil {
		return nil, err
	}
	shardID := GetShardIDFromPaymentAddressBytes(shareRewardOldBoardIns.VoterPaymentAddress)
	var metadataType int
	if shareRewardOldBoardIns.BoardType == common.DCBBoard {
		metadataType = ShareRewardOldDCBBoardSupportterIns
	} else {
		metadataType = ShareRewardOldGOVBoardSupportterIns
	}
	return []string{
		strconv.Itoa(metadataType),
		strconv.Itoa(int(shardID)),
		string(content),
	}, nil
}

func (rewardProposalSubmitterIns RewardProposalSubmitterIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(rewardProposalSubmitterIns)
	if err != nil {
		return nil, err
	}
	shardID := GetShardIDFromPaymentAddressBytes(*rewardProposalSubmitterIns.ReceiverAddress)
	var metadataType int
	if rewardProposalSubmitterIns.BoardType == common.DCBBoard {
		metadataType = RewardDCBProposalSubmitterIns
	} else {
		metadataType = RewardGOVProposalSubmitterIns
	}
	return []string{
		strconv.Itoa(metadataType),
		strconv.Itoa(int(shardID)),
		string(content),
	}, nil
}

func (rewardProposalVoterIns RewardProposalVoterIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(rewardProposalVoterIns)
	if err != nil {
		return nil, err
	}
	shardID := GetShardIDFromPaymentAddressBytes(*rewardProposalVoterIns.ReceiverAddress)
	var metadataType int
	if rewardProposalVoterIns.BoardType == common.DCBBoard {
		metadataType = RewardDCBProposalVoterIns
	} else {
		metadataType = RewardGOVProposalVoterIns
	}
	return []string{
		strconv.Itoa(metadataType),
		strconv.Itoa(int(shardID)),
		string(content),
	}, nil
}

func (updateDCBConstitutionIns *UpdateDCBConstitutionIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(updateDCBConstitutionIns)
	if err != nil {
		return nil, err
	}
	return []string{
		strconv.Itoa(UpdateDCBConstitutionInsType),
		strconv.Itoa(AllShards),
		string(content),
	}, nil
}

func (updateGOVConstitutionIns *UpdateGOVConstitutionIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(updateGOVConstitutionIns)
	if err != nil {
		return nil, err
	}
	return []string{
		strconv.Itoa(UpdateGOVConstitutionInsType),
		strconv.Itoa(AllShards),
		string(content),
	}, nil
}

func (acceptProposalIns AcceptProposalIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(acceptProposalIns)
	if err != nil {
		return nil, err
	}
	var t int
	if acceptProposalIns.BoardType == common.DCBBoard {
		t = AcceptDCBProposalInsType
	} else {
		t = AcceptGOVProposalInsType
	}
	return []string{
		strconv.Itoa(t),
		strconv.Itoa(int(acceptProposalIns.ShardID)),
		string(content),
	}, nil
}

func (keepOldProposalIns KeepOldProposalIns) GetStringFormat() ([]string, error) {
	content, err := json.Marshal(keepOldProposalIns)
	if err != nil {
		return nil, err
	}
	var t int
	if keepOldProposalIns.BoardType == common.DCBBoard {
		t = KeepOldDCBProposalIns
	} else {
		t = KeepOldGOVProposalIns
	}
	return []string{
		strconv.Itoa(t),
		strconv.Itoa(AllShards),
		string(content),
	}, nil
}

func GetShardIDFromPaymentAddressBytes(paymentAddress privacy.PaymentAddress) byte {
	lastByte := paymentAddress.Pk[len(paymentAddress.Pk)-1]
	return common.GetShardIDFromLastByte(lastByte)
}
