package metadata

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/constant-money/constant-chain/blockchain/component"
	"github.com/constant-money/constant-chain/database"
	"github.com/pkg/errors"
)

type SendBackTokenVoteBoardFailMetadata struct {
	MetadataBase
}

func (SendBackTokenVoteBoardFailMetadata) ValidateTxWithBlockChain(tx Transaction, bcr BlockchainRetriever, b byte, db database.DatabaseInterface) (bool, error) {
	return true, nil
}

func (SendBackTokenVoteBoardFailMetadata) ValidateSanityData(bcr BlockchainRetriever, tx Transaction) (bool, bool, error) {
	return true, true, nil
}

func (SendBackTokenVoteBoardFailMetadata) ValidateMetadataByItself() bool {
	return true
}

func NewSendBackTokenVoteFailMetadata() *SendBackTokenVoteBoardFailMetadata {
	return &SendBackTokenVoteBoardFailMetadata{
		MetadataBase: *NewMetadataBase(SendBackTokenVoteBoardFailMeta),
	}
}

func (sendBackTokenVoteBoardFailMetadata *SendBackTokenVoteBoardFailMetadata) VerifyMinerCreatedTxBeforeGettingInBlock(insts [][]string,
	instUsed []int,
	shardID byte,
	tx Transaction,
	bcr BlockchainRetriever,
	accumulatedData *component.UsedInstData,
) (bool, error) {
	instIdx := -1
	var sendBackTokenIns component.TxSendBackTokenVoteFailIns
	pubkeys, amounts := tx.GetTokenReceivers()
	if len(pubkeys) != 1 {
		return false, errors.New("One sendbacktokenvoteboardfail just for one token receiver")
	}
	for i, inst := range insts {
		if instUsed[i] > 0 {
			continue
		}
		if inst[0] != strconv.Itoa(component.SendBackTokenVoteBoardFailIns) {
			continue
		}
		if inst[1] != strconv.Itoa(int(shardID)) {
			continue
		}
		if inst[2] != "accepted" {
			continue
		}
		contentStr := inst[3]
		err := json.Unmarshal([]byte(contentStr), &sendBackTokenIns)
		if err != nil {
			return false, err
		}
		if !bytes.Equal(pubkeys[0], sendBackTokenIns.PaymentAddress.Bytes()) {
			continue
		}
		instIdx = i
		instUsed[i]++
		break
	}
	if instIdx == -1 {
		return false, errors.Errorf("no instruction found for SendBackTokenVoteBoardFail tx %s", tx.Hash().String())
	}
	if tx.GetTokenID().Cmp(&sendBackTokenIns.PropertyID) != 0 {
		return false, errors.New("Wrong token ID")
	}
	if amounts[0] != sendBackTokenIns.Amount {
		return false, errors.Errorf("Wrong token amount. Right amount %+v, tx amount %+v\n", sendBackTokenIns.Amount, amounts[0])
	}
	return true, nil
}
