package blockchain

import (
	"errors"
	"fmt"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/transaction"

	"github.com/constant-money/constant-chain/database"
	"github.com/constant-money/constant-chain/metadata"
	"github.com/constant-money/constant-chain/privacy"
)

var (
	mintDCBTokenParam = transaction.CustomTokenParamTx{
		PropertyID:     common.DCBTokenID.String(),
		PropertyName:   common.DCBTokenID.String(),
		PropertySymbol: common.DCBTokenID.String(),
		Amount:         0,
		TokenTxType:    transaction.CustomTokenInit,
		Receiver:       nil,
		Mintable:       true,
	}
	mintGOVTokenParam = transaction.CustomTokenParamTx{
		PropertyID:     common.GOVTokenID.String(),
		PropertyName:   common.GOVTokenID.String(),
		PropertySymbol: common.GOVTokenID.String(),
		Amount:         0,
		TokenTxType:    transaction.CustomTokenInit,
		Receiver:       nil,
		Mintable:       true,
	}
)

func (bc *BlockChain) BuildTransaction(
	metaType int,
	inst interface{},
	minerPrivateKey *privacy.PrivateKey,
	db database.DatabaseInterface,
	byteInfo byte,
) (metadata.Transaction, error) {
	fn, ok := buildTxForVotes[metaType]
	if ok {
		return fn(
			inst,
			minerPrivateKey,
			db,
			byteInfo,
		)
	} else {
		fmt.Printf("[ndh] BuildTransaction function for this metatype not found!\n")
		return nil, errors.New("BuildTransaction function for this metatype not found!")
	}
}

func buildTxSendBackTokenVoteFailIns(
	inst interface{},
	minerPrivateKey *privacy.PrivateKey,
	db database.DatabaseInterface,
	shardID byte,
) (metadata.Transaction, error) {
	sendBackTokenVoteFailIns := inst.(component.TxSendBackTokenVoteFailIns)
	//create token params
	customTokenParamTx := mintDCBTokenParam
	if sendBackTokenVoteFailIns.BoardType == common.GOVBoard {
		customTokenParamTx = mintGOVTokenParam
	}
	customTokenParamTx.Receiver = []transaction.TxTokenVout{{
		Value:          sendBackTokenVoteFailIns.Amount,
		PaymentAddress: sendBackTokenVoteFailIns.PaymentAddress,
	}}
	customTokenParamTx.Amount = sendBackTokenVoteFailIns.Amount

	//CALL DB
	//listCustomTokens, err := GetListCustomTokens(db, bcr)
	//if err != nil {
	//	return nil, err
	//}
	txCustom := &transaction.TxCustomToken{}
	err1 := txCustom.Init(
		minerPrivateKey,
		[]*privacy.PaymentInfo{},
		nil,
		0,
		&customTokenParamTx,
		//listCustomTokens,
		db,
		metadata.NewSendBackTokenVoteFailMetadata(),
		false,
		shardID,
	)
	if err1 != nil {
		return nil, err1
	}
	return txCustom, nil
}

func buildTxSendBackTokenToOldSupporterIns(
	inst interface{},
	minerPrivateKey *privacy.PrivateKey,
	db database.DatabaseInterface,
	shardID byte,
) (metadata.Transaction, error) {
	sendBackTokenToOldSupporterIns := inst.(component.TxSendBackTokenToOldSupporterIns)
	//create token params
	customTokenParamTx := mintDCBTokenParam
	if sendBackTokenToOldSupporterIns.BoardType == common.GOVBoard {
		customTokenParamTx = mintGOVTokenParam
	}
	customTokenParamTx.Receiver = []transaction.TxTokenVout{{
		Value:          sendBackTokenToOldSupporterIns.Amount,
		PaymentAddress: sendBackTokenToOldSupporterIns.PaymentAddress,
	}}
	customTokenParamTx.Amount = sendBackTokenToOldSupporterIns.Amount

	//CALL DB
	//listCustomTokens, err := GetListCustomTokens(db, bcr)
	//if err != nil {
	//	return nil, err
	//}
	txCustom := &transaction.TxCustomToken{}
	err1 := txCustom.Init(
		minerPrivateKey,
		[]*privacy.PaymentInfo{},
		nil,
		0,
		&customTokenParamTx,
		//listCustomTokens,
		db,
		metadata.NewSendBackTokenVoteFailMetadata(),
		false,
		shardID,
	)
	if err1 != nil {
		return nil, err1
	}
	return txCustom, nil
}

func buildTxShareRewardOldBoardIns(
	inst interface{},
	minerPrivateKey *privacy.PrivateKey,
	db database.DatabaseInterface,
	boardType byte,
) (metadata.Transaction, error) {
	shareRewardOldBoardIns := inst.(component.ShareRewardOldBoardIns)
	rewardShareOldBoardMeta := metadata.NewShareRewardOldBoardMetadata(
		shareRewardOldBoardIns.ChairPaymentAddress,
		shareRewardOldBoardIns.VoterPaymentAddress,
		shareRewardOldBoardIns.BoardType,
	)
	tx := transaction.Tx{}
	err := tx.InitTxSalary(
		shareRewardOldBoardIns.AmountOfCoin,
		&shareRewardOldBoardIns.VoterPaymentAddress,
		minerPrivateKey,
		db,
		rewardShareOldBoardMeta,
	)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func buildTxRewardProposalVoterIns(
	inst interface{},
	minerPrivateKey *privacy.PrivateKey,
	db database.DatabaseInterface,
	boardType byte,
) (metadata.Transaction, error) {
	var meta metadata.Metadata
	if common.BoardType(boardType) == common.DCBBoard {
		meta = metadata.NewRewardDCBProposalVoterMetadata()
	} else {
		meta = metadata.NewRewardGOVProposalVoterMetadata()
	}
	rewardProposalVoterIns := inst.(component.RewardProposalVoterIns)
	tx := transaction.Tx{}
	receiverAddress := rewardProposalVoterIns.ReceiverAddress
	amount := rewardProposalVoterIns.Amount
	err := tx.InitTxSalary(amount, receiverAddress, minerPrivateKey, db, meta)
	return &tx, err
}

func buildTxRewardProposalSubmitterIns(
	inst interface{},
	minerPrivateKey *privacy.PrivateKey,
	db database.DatabaseInterface,
	boardType byte,
) (metadata.Transaction, error) {
	var meta metadata.Metadata
	if common.BoardType(boardType) == common.DCBBoard {
		meta = metadata.NewRewardDCBProposalSubmitterMetadata()
	} else {
		meta = metadata.NewRewardGOVProposalSubmitterMetadata()
	}
	rewardProposalSubmitterIns := inst.(component.RewardProposalSubmitterIns)
	tx := transaction.Tx{}
	receiverAddress := rewardProposalSubmitterIns.ReceiverAddress
	amount := rewardProposalSubmitterIns.Amount
	err := tx.InitTxSalary(amount, receiverAddress, minerPrivateKey, db, meta)
	return &tx, err
}
