package types

const (
	EventTypeCreateNode              = "create_node"
	EventTypeCommitReport            = "commit_report"
	EventTypeAddAuthorizedMembers    = "add_authorized_members"
	EventTypeRemoveAuthorizedMembers = "remove_authorized_members"
	EventTypeUpdateSaleLevel         = "update_sale_level"
	EventTypeCommitComputingPower    = "commit_computing_power"
	EventTypeClaimComputingPower     = "claim_computing_power"
	EventTypeEpochPhase              = "epoch_phase"
	EventTypeBeginBlock              = "begin_block"
	EventTypeEndBlock                = "end_block"

	AttributeKeyEpochID              = "epoch_id"
	AttributeKeyNodeID               = "node_id"
	AttributeKeyDivisionID           = "division_id"
	AttributeKeyReceiver             = "receiver"
	AttributeKeyOwner                = "owner"
	AttributeKeyReportType           = "report_type"
	AttributeKeyAuthorizedMember     = "authorized_member"
	AttributeKeyComputingPower       = "computing_power"
	AttributeKeyComputingPowerBefore = "computing_power_before"
	AttributeKeyComputingPowerAfter  = "computing_power_after"
	AttributeKeySaleLevelBefore      = "sale_level_before"
	AttributeKeySaleLevelAfter       = "sale_level_after"

	AttributeValueCategory = ModuleName
)
