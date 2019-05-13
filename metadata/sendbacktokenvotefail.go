package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	fmt.Println("[ndh] - - SendBackTokenVoteBoardFail instruction type: ", component.SendBackTokenVoteBoardFailIns)
	for i, inst := range insts {
		fmt.Printf("[ndh] - - - - - instruction:%+v ", inst)
		if instUsed[i] > 0 {
			fmt.Println("is used.")
			continue
		}
		if inst[0] != strconv.Itoa(component.SendBackTokenVoteBoardFailIns) {
			fmt.Println("wrong type.")
			continue
		}
		if inst[1] != strconv.Itoa(int(shardID)) {
			fmt.Println("wrong shardID")
			continue
		}
		contentStr := inst[2]
		err := json.Unmarshal([]byte(contentStr), &sendBackTokenIns)
		if err != nil {
			return false, err
		}
		if !bytes.Equal(pubkeys[0], sendBackTokenIns.PaymentAddress.Pk) {
			fmt.Printf("tx pk: %+v, tx inst: %+v\n", pubkeys[0], sendBackTokenIns.PaymentAddress.Pk)
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

func (sendBackTokenVoteBoardFailMetadata *SendBackTokenVoteBoardFailMetadata) CheckTransactionFee(tr Transaction, minFee uint64) bool {
	// no need to have fee for this tx
	return true
}
