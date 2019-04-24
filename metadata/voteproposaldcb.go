package metadata

import (
	"errors"

	"github.com/constant-money/constant-chain/database/lvdb"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/database"
	"github.com/constant-money/constant-chain/metadata/fromshardins"
)

func (dcbVoteProposalMetadata *DCBVoteProposalMetadata) GetBoardType() common.BoardType {
	return common.DCBBoard
}

type DCBVoteProposalMetadata struct {
	NormalVoteProposalMetadata component.VoteProposalData
	MetadataBase
}

func (dcbVoteProposalMetadata *DCBVoteProposalMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	//return dcbVoteProposalMetadata.VoteProposalMetadata.ValidateSanityData(bcr, tx)
	return true, true, nil
}

func (dcbVoteProposalMetadata *DCBVoteProposalMetadata) ValidateMetadataByItself() bool {
	//return dcbVoteProposalMetadata.VoteProposalMetadata.ValidateMetadataByItself()
	return true
}

func NewDCBVoteProposalMetadata(
	voteProposal component.VoteProposalData,
) *DCBVoteProposalMetadata {
	return &DCBVoteProposalMetadata{
		NormalVoteProposalMetadata: voteProposal,
		MetadataBase:               *NewMetadataBase(DCBVoteProposalMeta),
	}
}

func (dcbVoteProposalMetadata *DCBVoteProposalMetadata) Hash() *common.Hash {
	record := dcbVoteProposalMetadata.NormalVoteProposalMetadata.ToBytes()

	hash := common.HashH([]byte(record))
	return &hash
}

func (dcbVoteProposalMetadata *DCBVoteProposalMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, shardID byte, db database.DatabaseInterface) (bool, error) {
	//Validate these pubKeys are in board
	//boardType := common.DCBBoard
	//return dcbVoteProposalMetadata.VoteProposalMetadata.ValidateTxWithBlockChain(
	//	boardType,
	//	tx,
	//	bcr,
	//	shardID,
	//	db,
	//)
	found := false
	board := bcr.GetBoardPaymentAddress(common.DCBBoard)
	for _, payment := range board {
		if common.ByteEqual(payment.Bytes(), dcbVoteProposalMetadata.NormalVoteProposalMetadata.VoterPayment.Bytes()) {
			found = true
			break
		}
	}
	if !found {
		return false, errors.New("Voter is not governor")
	}
	key := lvdb.GetKeySubmitProposal(common.DCBBoard,
		dcbVoteProposalMetadata.NormalVoteProposalMetadata.ConstitutionIndex,
		dcbVoteProposalMetadata.NormalVoteProposalMetadata.VoterPayment.Bytes(),
	)
	found, err := bcr.GetDatabase().HasValue(key)
	if err != nil {
		return false, err
	}
	if found {
		return false, errors.New("Just vote 1 proposal")
	}
	return true, nil
}

func (dcbVoteProposalMetadata *DCBVoteProposalMetadata) BuildReqActions(
	tx Transaction,
	bcr BlockchainRetriever,
	shardID byte,
) ([][]string, error) {
	voteProposal := dcbVoteProposalMetadata.NormalVoteProposalMetadata
	inst := fromshardins.NewNormalVoteProposalIns(common.DCBBoard, voteProposal)

	instStr, err := inst.GetStringFormat()
	if err != nil {
		return nil, err
	}
	return [][]string{instStr}, nil
}
