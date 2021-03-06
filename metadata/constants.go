package metadata

const (
	InvalidMeta = 1

	IssuingRequestMeta     = 24
	IssuingResponseMeta    = 25
	ContractingRequestMeta = 26

	ResponseBaseMeta             = 35
	ShardBlockReward             = 36
	ShardBlockSalaryRequestMeta  = 37
	ShardBlockSalaryResponseMeta = 38
	BeaconSalaryRequestMeta      = 39
	BeaconSalaryResponseMeta     = 40
	ReturnStakingMeta            = 41

	//statking
	ShardStakingMeta  = 63
	BeaconStakingMeta = 64
)

const (
	MaxDivTxsPerBlock = 1000
)

var minerCreatedMetaTypes = []int{
	ShardBlockReward,
	BeaconSalaryResponseMeta,
	IssuingResponseMeta,
	ReturnStakingMeta,
}

// Special rules for shardID: stored as 2nd param of instruction of BeaconBlock
