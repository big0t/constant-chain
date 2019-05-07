package metadata

import (
	"errors"
	"fmt"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/database"
	"github.com/constant-money/constant-chain/database/lvdb"
	"github.com/constant-money/constant-chain/metadata/fromshardins"
)

func (govVoteProposalMetadata *GOVVoteProposalMetadata) GetBoardType() common.BoardType {
	return common.GOVBoard
}

type GOVVoteProposalMetadata struct {
	VoteProposalMetadata component.VoteProposalData
	MetadataBase
}

func (govVoteProposalMetadata *GOVVoteProposalMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	rightConstitutionIndex := bcr.GetConstitution(common.GOVBoard).GetConstitutionIndex() + 1
	if govVoteProposalMetadata.VoteProposalMetadata.ConstitutionIndex != rightConstitutionIndex {
		fmt.Printf("[ndh] - Wrong constitution index, right constitution index is %+v\n", rightConstitutionIndex)
		return true, false, errors.New("Wrong constitution index")
	}
	return true, true, nil
}

func (govVoteProposalMetadata *GOVVoteProposalMetadata) ValidateMetadataByItself() bool {
	//return govVoteProposalMetadata.VoteProposalMetadata.ValidateMetadataByItself()
	return true
}

func NewGOVVoteProposalMetadata(
	voteProposal component.VoteProposalData,
) *GOVVoteProposalMetadata {
	return &GOVVoteProposalMetadata{
		VoteProposalMetadata: voteProposal,
		MetadataBase:         *NewMetadataBase(GOVVoteProposalMeta),
	}
}

func (govVoteProposalMetadata *GOVVoteProposalMetadata) Hash() *common.Hash {
	record := govVoteProposalMetadata.VoteProposalMetadata.ToBytes()

	hash := common.HashH([]byte(record))
	return &hash
}

func (govVoteProposalMetadata *GOVVoteProposalMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, shardID byte, db database.DatabaseInterface) (bool, error) {
	//Validate these pubKeys are in board
	//boardType := common.GOVBoard
	//return govVoteProposalMetadata.VoteProposalMetadata.ValidateTxWithBlockChain(
	//	boardType,
	//	tx,
	//	bcr,
	//	shardID,
	//	db,
	//)
	found := false
	board := bcr.GetBoardPaymentAddress(common.GOVBoard)
	for _, payment := range board {
		if common.ByteEqual(payment.Bytes(), govVoteProposalMetadata.VoteProposalMetadata.VoterPayment.Bytes()) {
			found = true
			break
		}
	}
	if !found {
		return false, errors.New("Voter is not governor")
	}
	key := lvdb.GetKeySubmitProposal(common.GOVBoard,
		govVoteProposalMetadata.VoteProposalMetadata.ConstitutionIndex,
		govVoteProposalMetadata.VoteProposalMetadata.VoterPayment.Bytes(),
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

func (govVoteProposalMetadata *GOVVoteProposalMetadata) BuildReqActions(
	tx Transaction,
	bcr BlockchainRetriever,
	shardID byte,
) ([][]string, error) {
	voteProposal := govVoteProposalMetadata.VoteProposalMetadata
	inst := fromshardins.NewNormalVoteProposalIns(common.GOVBoard, voteProposal)

	instStr, err := inst.GetStringFormat()
	if err != nil {
		return nil, err
	}
	return [][]string{instStr}, nil
}
