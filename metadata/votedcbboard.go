package metadata

import (
	"errors"

	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/database"
	"github.com/constant-money/constant-chain/metadata/fromshardins"
	"github.com/constant-money/constant-chain/privacy"
	"github.com/constant-money/constant-chain/wallet"
)

type VoteDCBBoardMetadata struct {
	VoteBoardMetadata VoteBoardMetadata

	MetadataBase
}

type GovernorInterface interface {
	GetBoardIndex() uint32
}

func NewVoteDCBBoardMetadata(candidatePaymentAddress privacy.PaymentAddress, boardIndex uint32) *VoteDCBBoardMetadata {
	return &VoteDCBBoardMetadata{
		VoteBoardMetadata: *NewVoteBoardMetadata(candidatePaymentAddress, boardIndex),
		MetadataBase:      *NewMetadataBase(VoteDCBBoardMeta),
	}
}

func NewVoteDCBBoardMetadataFromRPC(data map[string]interface{}) (Metadata, error) {
	paymentAddress := data["PaymentAddress"].(string)
	boardIndex := uint32(data["BoardIndex"].(float64))
	account, _ := wallet.Base58CheckDeserialize(paymentAddress)
	meta := NewVoteDCBBoardMetadata(account.KeySet.PaymentAddress, boardIndex)
	return meta, nil
}

func (voteDCBBoardMetadata *VoteDCBBoardMetadata) ValidateTxWithBlockChain(txr Transaction, bcr BlockchainRetriever, shardID byte, db database.DatabaseInterface) (bool, error) {
	if txr.IsPrivacy() {
		return false, errors.New("Can not vote with privacy transaction!")
	}
	voteAmount, err := txr.GetAmountOfVote(common.DCBBoard)
	if err != nil {
		return false, err
	}
	if voteAmount == 0 {
		return false, errors.New("Amount of vote must large than zero!")
	}
	return true, nil
}

func (voteDCBBoardMetadata *VoteDCBBoardMetadata) Hash() *common.Hash {
	record := string(voteDCBBoardMetadata.VoteBoardMetadata.GetBytes())
	record += voteDCBBoardMetadata.MetadataBase.Hash().String()
	hash := common.HashH([]byte(record))
	return &hash
}

func (voteDCBBoardMetadata *VoteDCBBoardMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	if voteDCBBoardMetadata.VoteBoardMetadata.BoardIndex != bcr.GetGovernor(common.DCBBoard).GetBoardIndex()+1 {
		return true, false, errors.New("Wrong board index!")
	}
	return true, true, nil
}

func (voteDCBBoardMetadata *VoteDCBBoardMetadata) ValidateMetadataByItself() bool {
	return true
}

func (voteDCBBoardMetadata *VoteDCBBoardMetadata) BuildReqActions(tx Transaction, bcr BlockchainRetriever, shardID byte) ([][]string, error) {
	voterPaymentAddress, err := tx.GetVoterPaymentAddress()
	if err != nil {
		return nil, err
	}
	amountOfVote, err := tx.GetAmountOfVote(common.DCBBoard)
	if err != nil {
		return nil, err
	}
	inst := fromshardins.NewVoteBoardIns(
		common.DCBBoard,
		voteDCBBoardMetadata.VoteBoardMetadata.CandidatePaymentAddress,
		*voterPaymentAddress,
		voteDCBBoardMetadata.VoteBoardMetadata.BoardIndex,
		amountOfVote,
	)
	instStr, err := inst.GetStringFormat()
	if err != nil {
		return nil, err
	}
	return [][]string{instStr}, nil
}
