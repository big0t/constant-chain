package metadata

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/database"
	"github.com/pkg/errors"
)

//validate by checking vout address of this tx equal to vin address of winning proposal
type RewardDCBProposalVoterMetadata struct {
	MetadataBase
}

func (rewardDCBProposalVoterMetadata *RewardDCBProposalVoterMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	// bcr.UpdateDCBFund(tx)
	return nil
}

func NewRewardDCBProposalVoterMetadata() *RewardDCBProposalVoterMetadata {
	return &RewardDCBProposalVoterMetadata{
		MetadataBase: *NewMetadataBase(RewardDCBProposalVoterMeta),
	}
}

func (rewardDCBProposalVoterMetadata *RewardDCBProposalVoterMetadata) Hash() *common.Hash {
	record := rewardDCBProposalVoterMetadata.MetadataBase.Hash().String()
	hash := common.HashH([]byte(record))
	return &hash
}

func (rewardDCBProposalVoterMetadata *RewardDCBProposalVoterMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, b byte, db database.DatabaseInterface) (bool, error) {
	return true, nil
}

func (rewardDCBProposalVoterMetadata *RewardDCBProposalVoterMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return true, true, nil
}

func (rewardDCBProposalVoterMetadata *RewardDCBProposalVoterMetadata) ValidateMetadataByItself() bool {
	return true
}

type RewardGOVProposalVoterMetadata struct {
	MetadataBase
}

func (rewardGOVProposalVoterMetadata *RewardGOVProposalVoterMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	// bcr.UpdateDCBFund(tx)
	return nil
}

func NewRewardGOVProposalVoterMetadata() *RewardGOVProposalVoterMetadata {
	return &RewardGOVProposalVoterMetadata{
		MetadataBase: *NewMetadataBase(RewardGOVProposalVoterMeta),
	}
}

func (rewardGOVProposalVoterMetadata *RewardGOVProposalVoterMetadata) Hash() *common.Hash {
	record := rewardGOVProposalVoterMetadata.MetadataBase.Hash().String()
	hash := common.HashH([]byte(record))
	return &hash
}

func (rewardGOVProposalVoterMetadata *RewardGOVProposalVoterMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, b byte, db database.DatabaseInterface) (bool, error) {
	return true, nil
}

func (rewardGOVProposalVoterMetadata *RewardGOVProposalVoterMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return true, true, nil
}

func (rewardGOVProposalVoterMetadata *RewardGOVProposalVoterMetadata) ValidateMetadataByItself() bool {
	return true
}

func (rewardGOVProposalVoterMetadata *RewardGOVProposalVoterMetadata) VerifyMinerCreatedTxBeforeGettingInBlock(
	insts [][]string,
	instUsed []int,
	shardID byte,
	tx Transaction,
	bcr BlockchainRetriever,
	accumulatedData *component.UsedInstData,
) (bool, error) {
	instIdx := -1
	var rewardGOVProposalVoterIns component.RewardProposalVoterIns
	pubkeys, amounts := tx.GetReceivers()
	if len(pubkeys) != 1 {
		return false, errors.New("One RewardGOVProposalSubmitter instruction just for one token receiver")
	}
	for i, inst := range insts {
		if instUsed[i] > 0 {
			continue
		}
		if inst[0] != strconv.Itoa(component.RewardGOVProposalSubmitterIns) {
			continue
		}
		if inst[1] != strconv.Itoa(int(shardID)) {
			continue
		}
		if inst[2] != "accepted" {
			continue
		}
		contentStr := inst[3]
		err := json.Unmarshal([]byte(contentStr), &rewardGOVProposalVoterIns)
		if err != nil {
			return false, err
		}
		if !bytes.Equal(pubkeys[0], rewardGOVProposalVoterIns.ReceiverAddress.Bytes()) {
			continue
		}
		instIdx = i
		instUsed[i]++
		break
	}
	if instIdx == -1 {
		return false, errors.Errorf("no instruction found for RewardGOVProposalSubmitter tx %s", tx.Hash().String())
	}
	if amounts[0] != rewardGOVProposalVoterIns.Amount {
		return false, errors.Errorf("Wrong reward amount. Right amount %+v, tx amount %+v\n", rewardGOVProposalVoterIns.Amount, amounts[0])
	}
	return true, nil
}

func (rewardGOVProposalVoterMetadata *RewardGOVProposalVoterMetadata) CalculateSize() uint64 {
	return calculateSize(rewardGOVProposalVoterMetadata)
}
