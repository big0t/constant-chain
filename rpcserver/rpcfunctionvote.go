package rpcserver

import (
	"fmt"

	"github.com/constant-money/constant-chain/database/lvdb"

	"github.com/constant-money/constant-chain/common"
	"github.com/constant-money/constant-chain/metadata"
	"github.com/pkg/errors"
)

func iPlusPlus(x *int) int {
	*x += 1
	return *x - 1
}

// ============================== VOTE PROPOSAL

func (rpcServer RpcServer) handleCreateRawVoteProposalTransaction(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	//VoteProposal - Step 2: Create Raw vote proposal transaction
	params = setBuildRawBurnTransactionParams(params, FeeVote)
	return rpcServer.createRawTxWithMetadata(
		params,
		closeChan,
		metadata.NewVoteProposalMetadataFromRPC,
	)
}

func (rpcServer RpcServer) handleCreateAndSendVoteProposalTransaction(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	//VoteProposal - Step 1: Client call rpc function to create vote proposal transaction
	return rpcServer.createAndSendTxWithMetadata(
		params,
		closeChan,
		RpcServer.handleCreateRawVoteProposalTransaction,
		RpcServer.handleSendRawTransaction,
	)
}

func (rpcServer RpcServer) handleGetDCBBoardIndex(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	return rpcServer.config.BlockChain.BestState.Beacon.StabilityInfo.DCBGovernor.BoardIndex, nil
}
func (rpcServer RpcServer) handleGetGOVBoardIndex(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	return rpcServer.config.BlockChain.BestState.Beacon.StabilityInfo.GOVGovernor.BoardIndex, nil
}

func setBuildRawBurnTransactionParams(params interface{}, fee float64) interface{} {
	arrayParams := common.InterfaceSlice(params)
	x := make(map[string]interface{})
	x[common.BurningAddress] = fee
	arrayParams[1] = x
	return arrayParams
}

func (rpcServer RpcServer) handleGetProposalTxIDbyConstitutionID(params interface{}, closeChan <-chan struct{}) (interface{}, *RPCError) {
	arrayParams := common.InterfaceSlice(params)
	if len(arrayParams) != 2 {
		return nil, NewRPCError(ErrRPCInvalidParams, errors.New("Params for this rpc function is [ <BoardType>, <Constitution Index>]"))
	}
	boardType := common.NewBoardTypeFromString(arrayParams[0].(string))
	if boardType == 255 {
		return nil, NewRPCError(ErrRPCInvalidParams, errors.New("Wrong board type!"))
	}
	constitutionIndex := uint32(arrayParams[1].(float64))
	if constitutionIndex > rpcServer.config.BlockChain.GetConstitution(boardType).GetConstitutionIndex() {
		return nil, NewRPCError(ErrRPCInvalidRequest, errors.New("Can not find any data of this constitution"))
	}
	fmt.Printf("[ndh] - - - Request params: ", boardType, constitutionIndex)
	db := rpcServer.config.BlockChain.GetDatabase()
	gg := lvdb.ViewDetailDBByPrefix(db, lvdb.SubmitProposalPrefix)
	for _, value := range gg {
		boardTypee, constitutionIndexx, proposalTxIDD, _ := lvdb.ParseKeySubmitProposal(value)
		fmt.Printf("[ndh] - - - - - - - %+v %+v %+v\n", boardTypee, constitutionIndexx, proposalTxIDD)
	}
	proposalTxBytes, err := db.GetProposalTXIDByConstitutionIndex(boardType, constitutionIndex)
	if err != nil {
		return nil, NewRPCError(ErrUnexpected, err)
	}
	proposalTxHash, err := common.NewHash(proposalTxBytes)
	if err != nil {
		return nil, NewRPCError(ErrRPCParse, err)
	}
	return proposalTxHash.String(), nil
}
