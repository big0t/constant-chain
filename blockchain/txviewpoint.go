package blockchain

import (
	"errors"

	"math/big"

	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/common/base58"
	"github.com/constant-money/constant-chain/database"
	"github.com/constant-money/constant-chain/privacy"
	"github.com/constant-money/constant-chain/privacy/zeroknowledge"
	"github.com/constant-money/constant-chain/transaction"
)

// TxViewPoint is used to contain data which is fetched from tx of every block
type TxViewPoint struct {
	tokenID           *common.Hash
	shardID           byte
	listSerialNumbers [][]byte // array serialNumbers
	//listSnD            []big.Int
	mapSnD         map[string][]big.Int
	mapCommitments map[string][][]byte //map[base58check.encode{pubkey}]([]([]byte-commitment))
	mapOutputCoins map[string][]privacy.OutputCoin

	// data of normal custom token
	customTokenTxs map[int32]*transaction.TxCustomToken

	// data of privacy custom token
	privacyCustomTokenViewPoint map[int32]*TxViewPoint
	privacyCustomTokenTxs       map[int32]*transaction.TxCustomTokenPrivacy
	privacyCustomTokenMetadata  *CrossShardTokenPrivacyMetaData

	//cross shard tx token
	crossTxTokenData map[int32]*transaction.TxTokenData
}

/*
ListSerialNumbers returns list nullifers which is contained in TxViewPoint
*/
// #1: joinSplitDescType is "Coin" Or "Bond" or other token
func (view *TxViewPoint) ListSerialNumbers() [][]byte {
	return view.listSerialNumbers
}

// func (view *TxViewPoint) ListSnDerivators() []big.Int {
// 	return view.listSnD
// }
func (view *TxViewPoint) MapSnDerivators() map[string][]big.Int {
	return view.mapSnD
}

func (view *TxViewPoint) ListSerialNumnbersEclipsePoint() []*privacy.EllipticPoint {
	result := []*privacy.EllipticPoint{}
	for _, commitment := range view.listSerialNumbers {
		point := &privacy.EllipticPoint{}
		point.Decompress(commitment)
		result = append(result, point)
	}
	return result
}

// fetch from desc of tx to get nullifiers and commitments
// (note: still storage full data of commitments, serialnumbers, snderivator to check double spend)
func (view *TxViewPoint) processFetchTxViewPoint(
	shardID byte,
	db database.DatabaseInterface,
	proof *zkp.PaymentProof,
	tokenID *common.Hash,
) ([][]byte, map[string][][]byte, map[string][]privacy.OutputCoin, map[string][]big.Int, error) {
	acceptedNullifiers := make([][]byte, 0)
	acceptedCommitments := make(map[string][][]byte)
	acceptedOutputcoins := make(map[string][]privacy.OutputCoin)
	acceptedSnD := make(map[string][]big.Int)
	if proof == nil {
		return acceptedNullifiers, acceptedCommitments, acceptedOutputcoins, acceptedSnD, nil
	}
	// Get data for serialnumbers
	// Process input of transaction
	// Get Serial numbers of input
	// Append into accepttedNullifiers if this serial number haven't exist yet
	for _, item := range proof.InputCoins {
		serialNum := item.CoinDetails.SerialNumber.Compress()
		ok, err := db.HasSerialNumber(tokenID, serialNum, shardID)
		if err != nil {
			return acceptedNullifiers, acceptedCommitments, acceptedOutputcoins, acceptedSnD, err
		}
		if !ok {
			acceptedNullifiers = append(acceptedNullifiers, serialNum)
		}
	}

	// Process Output Coins (just created UTXO of this transaction)
	// Proccessed variable: commitment, snd, outputcoins
	// Commitment and SND must not exist before in db
	// Outputcoins will be stored as new utxo for next transaction
	for _, item := range proof.OutputCoins {
		commitment := item.CoinDetails.CoinCommitment.Compress()
		pubkey := item.CoinDetails.PublicKey.Compress()
		pubkeyStr := base58.Base58Check{}.Encode(pubkey, common.ZeroByte)
		ok, err := db.HasCommitment(tokenID, commitment, shardID)
		if err != nil {
			return acceptedNullifiers, acceptedCommitments, acceptedOutputcoins, acceptedSnD, err
		}
		if !ok {
			pubkeyStr := base58.Base58Check{}.Encode(pubkey, common.ZeroByte)
			if acceptedCommitments[pubkeyStr] == nil {
				acceptedCommitments[pubkeyStr] = make([][]byte, 0)
			}
			// get data for commitments
			acceptedCommitments[pubkeyStr] = append(acceptedCommitments[pubkeyStr], item.CoinDetails.CoinCommitment.Compress())

			// get data for output coin
			if acceptedOutputcoins[pubkeyStr] == nil {
				acceptedOutputcoins[pubkeyStr] = make([]privacy.OutputCoin, 0)
			}
			acceptedOutputcoins[pubkeyStr] = append(acceptedOutputcoins[pubkeyStr], *item)
		}

		// get data for Snderivators
		snD := item.CoinDetails.SNDerivator
		ok, err = db.HasSNDerivator(tokenID, privacy.AddPaddingBigInt(snD, privacy.BigIntSize), shardID)
		if !ok && err == nil {
			acceptedSnD[pubkeyStr] = append(acceptedSnD[pubkeyStr], *snD)
		}
	}
	return acceptedNullifiers, acceptedCommitments, acceptedOutputcoins, acceptedSnD, nil
}

/*
fetchTxViewPointFromBlock get list serialnumber and commitments, output coins from txs in block and check if they are not in Main chain db
return a tx view point which contains list new nullifiers and new commitments from block
// (note: still storage full data of commitments, serialnumbers, snderivator to check double spend)
*/

func (view *TxViewPoint) fetchTxViewPointFromBlock(db database.DatabaseInterface, block *ShardBlock) error {
	transactions := block.Body.Transactions
	// Loop through all of the transaction descs (except for the salary tx)
	acceptedSerialNumbers := make([][]byte, 0)
	acceptedCommitments := make(map[string][][]byte)
	acceptedOutputcoins := make(map[string][]privacy.OutputCoin)
	acceptedSnD := make(map[string][]big.Int)
	constantTokenID := &common.Hash{}
	constantTokenID.SetBytes(common.ConstantID[:])
	for indexTx, tx := range transactions {
		switch tx.GetType() {
		case common.TxNormalType, common.TxRewardType, common.TxReturnStakingType:
			{
				normalTx := tx.(*transaction.Tx)
				serialNumbers, commitments, outCoins, snDs, err := view.processFetchTxViewPoint(block.Header.ShardID, db, normalTx.Proof, constantTokenID)
				if err != nil {
					return NewBlockChainError(UnExpectedError, err)
				}
				acceptedSerialNumbers = append(acceptedSerialNumbers, serialNumbers...)
				for pubkey, data := range commitments {
					if acceptedCommitments[pubkey] == nil {
						acceptedCommitments[pubkey] = make([][]byte, 0)
					}
					acceptedCommitments[pubkey] = append(acceptedCommitments[pubkey], data...)
				}
				for pubkey, data := range outCoins {
					if acceptedOutputcoins[pubkey] == nil {
						acceptedOutputcoins[pubkey] = make([]privacy.OutputCoin, 0)
					}
					acceptedOutputcoins[pubkey] = append(acceptedOutputcoins[pubkey], data...)
				}
				for pubkey, data := range snDs {
					if snDs[pubkey] == nil {
						snDs[pubkey] = make([]big.Int, 0)
					}
					snDs[pubkey] = append(snDs[pubkey], data...)
				}
				// acceptedSnD = append(acceptedSnD, snDs...)
			}
		case common.TxCustomTokenType:
			{
				tx := tx.(*transaction.TxCustomToken)
				serialNumbers, commitments, outCoins, snDs, err := view.processFetchTxViewPoint(block.Header.ShardID, db, tx.Proof, constantTokenID)
				if err != nil {
					return NewBlockChainError(UnExpectedError, err)
				}
				acceptedSerialNumbers = append(acceptedSerialNumbers, serialNumbers...)
				for pubkey, data := range commitments {
					if acceptedCommitments[pubkey] == nil {
						acceptedCommitments[pubkey] = make([][]byte, 0)
					}
					acceptedCommitments[pubkey] = append(acceptedCommitments[pubkey], data...)
				}
				for pubkey, data := range outCoins {
					if acceptedOutputcoins[pubkey] == nil {
						acceptedOutputcoins[pubkey] = make([]privacy.OutputCoin, 0)
					}
					acceptedOutputcoins[pubkey] = append(acceptedOutputcoins[pubkey], data...)
				}
				for pubkey, data := range snDs {
					if snDs[pubkey] == nil {
						snDs[pubkey] = make([]big.Int, 0)
					}
					snDs[pubkey] = append(snDs[pubkey], data...)
				}
				// acceptedSnD = append(acceptedSnD, snDs...)

				// indexTx is index of transaction in block
				view.customTokenTxs[int32(indexTx)] = tx
			}
		case common.TxCustomTokenPrivacyType:
			{
				tx := tx.(*transaction.TxCustomTokenPrivacy)
				serialNumbers, commitments, outCoins, snDs, err := view.processFetchTxViewPoint(block.Header.ShardID, db, tx.Proof, constantTokenID)
				if err != nil {
					return NewBlockChainError(UnExpectedError, err)
				}
				acceptedSerialNumbers = append(acceptedSerialNumbers, serialNumbers...)
				for pubkey, data := range commitments {
					if acceptedCommitments[pubkey] == nil {
						acceptedCommitments[pubkey] = make([][]byte, 0)
					}
					acceptedCommitments[pubkey] = append(acceptedCommitments[pubkey], data...)
				}
				for pubkey, data := range outCoins {
					if acceptedOutputcoins[pubkey] == nil {
						acceptedOutputcoins[pubkey] = make([]privacy.OutputCoin, 0)
					}
					acceptedOutputcoins[pubkey] = append(acceptedOutputcoins[pubkey], data...)
				}
				for pubkey, data := range snDs {
					if snDs[pubkey] == nil {
						snDs[pubkey] = make([]big.Int, 0)
					}
					snDs[pubkey] = append(snDs[pubkey], data...)
				}
				// acceptedSnD = append(acceptedSnD, snDs...)
				if err != nil {
					return NewBlockChainError(UnExpectedError, err)
				}

				// sub view for privacy custom token
				subView := NewTxViewPoint(block.Header.ShardID)
				subView.tokenID = &tx.TxTokenPrivacyData.PropertyID
				serialNumbersP, commitmentsP, outCoinsP, snDsP, errP := subView.processFetchTxViewPoint(subView.shardID, db, tx.TxTokenPrivacyData.TxNormal.Proof, subView.tokenID)
				if errP != nil {
					return NewBlockChainError(UnExpectedError, errP)
				}
				subView.listSerialNumbers = serialNumbersP
				for pubkey, data := range commitmentsP {
					if subView.mapCommitments[pubkey] == nil {
						subView.mapCommitments[pubkey] = make([][]byte, 0)
					}
					subView.mapCommitments[pubkey] = append(subView.mapCommitments[pubkey], data...)
				}
				for pubkey, data := range outCoinsP {
					if subView.mapOutputCoins[pubkey] == nil {
						subView.mapOutputCoins[pubkey] = make([]privacy.OutputCoin, 0)
					}
					subView.mapOutputCoins[pubkey] = append(subView.mapOutputCoins[pubkey], data...)
				}
				for pubkey, data := range snDsP {
					if subView.mapSnD[pubkey] == nil {
						subView.mapSnD[pubkey] = make([]big.Int, 0)
					}
					subView.mapSnD[pubkey] = append(subView.mapSnD[pubkey], data...)
				}
				// subView.listSnD = append(subView.listSnD, snDsP...)
				if err != nil {
					return NewBlockChainError(UnExpectedError, err)
				}

				view.privacyCustomTokenViewPoint[int32(indexTx)] = subView
				view.privacyCustomTokenTxs[int32(indexTx)] = tx
			}
		default:
			{
				return NewBlockChainError(UnExpectedError, errors.New("TxNormal type is invalid"))
			}
		}
	}

	if len(acceptedSerialNumbers) > 0 {
		view.listSerialNumbers = acceptedSerialNumbers
	}
	if len(acceptedCommitments) > 0 {
		view.mapCommitments = acceptedCommitments
	}
	if len(acceptedOutputcoins) > 0 {
		view.mapOutputCoins = acceptedOutputcoins
	}
	if len(acceptedSnD) > 0 {
		view.mapSnD = acceptedSnD
		// view.listSnD = acceptedSnD
	}
	return nil
}

/*
Create a TxNormal view point, which contains data about nullifiers and commitments
*/
func NewTxViewPoint(shardID byte) *TxViewPoint {
	result := &TxViewPoint{
		shardID:                     shardID,
		listSerialNumbers:           make([][]byte, 0),
		mapCommitments:              make(map[string][][]byte),
		mapOutputCoins:              make(map[string][]privacy.OutputCoin),
		mapSnD:                      make(map[string][]big.Int),
		customTokenTxs:              make(map[int32]*transaction.TxCustomToken),
		tokenID:                     &common.Hash{},
		privacyCustomTokenViewPoint: make(map[int32]*TxViewPoint),
		privacyCustomTokenTxs:       make(map[int32]*transaction.TxCustomTokenPrivacy),
		privacyCustomTokenMetadata:  &CrossShardTokenPrivacyMetaData{},
		crossTxTokenData:            make(map[int32]*transaction.TxTokenData),
	}
	result.tokenID.SetBytes(common.ConstantID[:])
	return result
}

/*
	fetch information from cross output coin
	- UTXO: outcoin
	- Commitment
	- snd
*/
func (view *TxViewPoint) processFetchCrossOutputViewPoint(
	shardID byte,
	db database.DatabaseInterface,
	outputCoins []privacy.OutputCoin,
	tokenID *common.Hash,
) (map[string][][]byte, map[string][]privacy.OutputCoin, map[string][]big.Int, error) {
	acceptedCommitments := make(map[string][][]byte)
	acceptedOutputcoins := make(map[string][]privacy.OutputCoin)
	acceptedSnD := make(map[string][]big.Int)
	if len(outputCoins) == 0 {
		return acceptedCommitments, acceptedOutputcoins, acceptedSnD, nil
	}

	// Process Output Coins (just created UTXO of this transaction)
	// Proccessed variable: commitment, snd, outputcoins
	// Commitment and SND must not exist before in db
	// Outputcoins will be stored as new utxo for next transaction
	for _, outputCoin := range outputCoins {
		item := &outputCoin
		commitment := item.CoinDetails.CoinCommitment.Compress()
		pubkey := item.CoinDetails.PublicKey.Compress()
		pubkeyStr := base58.Base58Check{}.Encode(pubkey, common.ZeroByte)
		ok, err := db.HasCommitment(tokenID, commitment, shardID)
		if err != nil {
			return acceptedCommitments, acceptedOutputcoins, acceptedSnD, err
		}
		if !ok {
			pubkeyStr := base58.Base58Check{}.Encode(pubkey, common.ZeroByte)
			if acceptedCommitments[pubkeyStr] == nil {
				acceptedCommitments[pubkeyStr] = make([][]byte, 0)
			}
			// get data for commitments
			acceptedCommitments[pubkeyStr] = append(acceptedCommitments[pubkeyStr], item.CoinDetails.CoinCommitment.Compress())

			// get data for output coin
			if acceptedOutputcoins[pubkeyStr] == nil {
				acceptedOutputcoins[pubkeyStr] = make([]privacy.OutputCoin, 0)
			}
			acceptedOutputcoins[pubkeyStr] = append(acceptedOutputcoins[pubkeyStr], *item)
		}

		// get data for Snderivators
		snD := item.CoinDetails.SNDerivator
		ok, err = db.HasSNDerivator(tokenID, privacy.AddPaddingBigInt(snD, privacy.BigIntSize), shardID)
		if !ok && err == nil {
			acceptedSnD[pubkeyStr] = append(acceptedSnD[pubkeyStr], *snD)
		}
	}
	return acceptedCommitments, acceptedOutputcoins, acceptedSnD, nil
}

func (view *TxViewPoint) fetchCrossTransactionViewPointFromBlock(db database.DatabaseInterface, block *ShardBlock) error {
	allShardCrossTransactions := block.Body.CrossTransactions
	// Loop through all of the transaction descs (except for the salary tx)
	acceptedOutputcoins := make(map[string][]privacy.OutputCoin)
	acceptedCommitments := make(map[string][][]byte)
	acceptedSnD := make(map[string][]big.Int)
	constantTokenID := &common.Hash{}
	constantTokenID.SetBytes(common.ConstantID[:])
	//@NOTICE: this function just work for Normal Transaction
	for _, crossTransactions := range allShardCrossTransactions {
		for _, crossTransaction := range crossTransactions {
			commitments, outCoins, snDs, err := view.processFetchCrossOutputViewPoint(block.Header.ShardID, db, crossTransaction.OutputCoin, constantTokenID)
			if err != nil {
				return NewBlockChainError(UnExpectedError, err)
			}
			for pubkey, data := range commitments {
				if acceptedCommitments[pubkey] == nil {
					acceptedCommitments[pubkey] = make([][]byte, 0)
				}
				acceptedCommitments[pubkey] = append(acceptedCommitments[pubkey], data...)
			}
			for pubkey, data := range outCoins {
				if acceptedOutputcoins[pubkey] == nil {
					acceptedOutputcoins[pubkey] = make([]privacy.OutputCoin, 0)
				}
				acceptedOutputcoins[pubkey] = append(acceptedOutputcoins[pubkey], data...)
			}
			for pubkey, data := range snDs {
				if snDs[pubkey] == nil {
					snDs[pubkey] = make([]big.Int, 0)
				}
				snDs[pubkey] = append(snDs[pubkey], data...)
			}
			if crossTransaction.TokenPrivacyData != nil && len(crossTransaction.TokenPrivacyData) > 0 {
				for index, tokenPrivacyData := range crossTransaction.TokenPrivacyData {
					subView := NewTxViewPoint(block.Header.ShardID)
					subView.tokenID = &tokenPrivacyData.PropertyID
					subView.privacyCustomTokenMetadata.TokenID = tokenPrivacyData.PropertyID
					subView.privacyCustomTokenMetadata.PropertyName = tokenPrivacyData.PropertyName
					subView.privacyCustomTokenMetadata.PropertySymbol = tokenPrivacyData.PropertySymbol
					subView.privacyCustomTokenMetadata.Amount = tokenPrivacyData.Amount
					subView.privacyCustomTokenMetadata.Mintable = tokenPrivacyData.Mintable
					commitmentsP, outCoinsP, snDsP, err := view.processFetchCrossOutputViewPoint(block.Header.ShardID, db, tokenPrivacyData.OutputCoin, subView.tokenID)
					if err != nil {
						return NewBlockChainError(UnExpectedError, err)
					}
					for pubkey, data := range commitmentsP {
						if subView.mapCommitments[pubkey] == nil {
							subView.mapCommitments[pubkey] = make([][]byte, 0)
						}
						subView.mapCommitments[pubkey] = append(subView.mapCommitments[pubkey], data...)
					}
					for pubkey, data := range outCoinsP {
						if subView.mapOutputCoins[pubkey] == nil {
							subView.mapOutputCoins[pubkey] = make([]privacy.OutputCoin, 0)
						}
						subView.mapOutputCoins[pubkey] = append(subView.mapOutputCoins[pubkey], data...)
					}
					for pubkey, data := range snDsP {
						if subView.mapSnD[pubkey] == nil {
							subView.mapSnD[pubkey] = make([]big.Int, 0)
						}
						subView.mapSnD[pubkey] = append(subView.mapSnD[pubkey], data...)
					}
					view.privacyCustomTokenViewPoint[int32(index)] = subView
				}
			}
		}
	}

	if len(acceptedCommitments) > 0 {
		view.mapCommitments = acceptedCommitments
	}
	if len(acceptedOutputcoins) > 0 {
		view.mapOutputCoins = acceptedOutputcoins
	}
	if len(acceptedSnD) > 0 {
		view.mapSnD = acceptedSnD
		// view.listSnD = acceptedSnD
	}
	return nil
}
