package component

const (
	VoteProposalIns              = 100 + iota
	NewDCBConstitutionIns        //1
	NewGOVConstitutionIns        //2
	UpdateDCBConstitutionInsType //3
	UpdateGOVConstitutionInsType //4
	VoteBoardIns                 //5
	SubmitProposalIns            //6

	AcceptDCBProposalInsType //7
	AcceptDCBBoardInsType    //8
	AcceptGOVProposalInsType //9
	AcceptGOVBoardInsType    //10

	RewardDCBProposalSubmitterIns       //11
	RewardGOVProposalSubmitterIns       //12
	ShareRewardOldDCBBoardSupportterIns //13
	ShareRewardOldGOVBoardSupportterIns //14
	SendBackTokenVoteBoardFailIns       //15

	ConfirmBuySellRequestMeta      //16
	ConfirmBuyBackRequestMeta      //17
	RewardDCBProposalVoterIns      //18
	RewardGOVProposalVoterIns      //19
	KeepOldDCBProposalIns          //20
	KeepOldGOVProposalIns          //21
	SendBackTokenToOldSupporterIns //22
)

const (
	AllShards  = -1
	BeaconOnly = -2
)
