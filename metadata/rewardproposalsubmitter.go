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
type RewardDCBProposalSubmitterMetadata struct {
	MetadataBase
}

func (rewardDCBProposalSubmitterMetadata *RewardDCBProposalSubmitterMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	// bcr.UpdateDCBFund(tx)
	return nil
}

func NewRewardDCBProposalSubmitterMetadata() *RewardDCBProposalSubmitterMetadata {
	return &RewardDCBProposalSubmitterMetadata{
		MetadataBase: *NewMetadataBase(RewardDCBProposalSubmitterMeta),
	}
}

func (rewardDCBProposalSubmitterMetadata *RewardDCBProposalSubmitterMetadata) Hash() *common.Hash {
	record := rewardDCBProposalSubmitterMetadata.MetadataBase.Hash().String()
	hash := common.HashH([]byte(record))
	return &hash
}

func (rewardDCBProposalSubmitterMetadata *RewardDCBProposalSubmitterMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, b byte, db database.DatabaseInterface) (bool, error) {
	return true, nil
}

func (rewardDCBProposalSubmitterMetadata *RewardDCBProposalSubmitterMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return true, true, nil
}

func (rewardDCBProposalSubmitterMetadata *RewardDCBProposalSubmitterMetadata) ValidateMetadataByItself() bool {
	return true
}

type RewardGOVProposalSubmitterMetadata struct {
	MetadataBase
}

func (rewardGOVProposalSubmitterMetadata *RewardGOVProposalSubmitterMetadata) ProcessWhenInsertBlockShard(tx Transaction, bcr BlockchainRetriever) error {
	// bcr.UpdateDCBFund(tx)
	return nil
}

func NewRewardGOVProposalSubmitterMetadata() *RewardGOVProposalSubmitterMetadata {
	return &RewardGOVProposalSubmitterMetadata{
		MetadataBase: *NewMetadataBase(RewardGOVProposalSubmitterMeta),
	}
}

func (rewardGOVProposalSubmitterMetadata *RewardGOVProposalSubmitterMetadata) Hash() *common.Hash {
	record := rewardGOVProposalSubmitterMetadata.MetadataBase.Hash().String()
	hash := common.HashH([]byte(record))
	return &hash
}

func (rewardGOVProposalSubmitterMetadata *RewardGOVProposalSubmitterMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, b byte, db database.DatabaseInterface) (bool, error) {
	return true, nil
}

func (rewardGOVProposalSubmitterMetadata *RewardGOVProposalSubmitterMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return true, true, nil
}

func (rewardGOVProposalSubmitterMetadata *RewardGOVProposalSubmitterMetadata) ValidateMetadataByItself() bool {
	return true
}

func (rewardGOVProposalSubmitterMetadata *RewardGOVProposalSubmitterMetadata) VerifyMinerCreatedTxBeforeGettingInBlock(
	insts [][]string,
	instUsed []int,
	shardID byte,
	tx Transaction,
	bcr BlockchainRetriever,
	accumulatedData *component.UsedInstData,
) (bool, error) {
	instIdx := -1
	var rewardProposalSubmitterIns component.RewardProposalSubmitterIns
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
		err := json.Unmarshal([]byte(contentStr), &rewardProposalSubmitterIns)
		if err != nil {
			return false, err
		}
		if !bytes.Equal(pubkeys[0], rewardProposalSubmitterIns.ReceiverAddress.Bytes()) {
			continue
		}
		instIdx = i
		instUsed[i]++
		break
	}
	if instIdx == -1 {
		return false, errors.Errorf("no instruction found for RewardGOVProposalSubmitter tx %s", tx.Hash().String())
	}
	if amounts[0] != rewardProposalSubmitterIns.Amount {
		return false, errors.Errorf("Wrong reward amount. Right amount %+v, tx amount %+v\n", rewardProposalSubmitterIns.Amount, amounts[0])
	}
	return true, nil
}

func (rewardGOVProposalSubmitterMetadata *RewardGOVProposalSubmitterMetadata) CalculateSize() uint64 {
	return calculateSize(rewardGOVProposalSubmitterMetadata)
}
