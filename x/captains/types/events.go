package types

const (
	EventTypeCreateNode = "create_node"

	AttributeValueCategory = ModuleName

	AttributeKeyNodeID     = "node_id"
	AttributeKeyDivisionID = "division_id"
	AttributeKeyReceiver   = "receiver"

	EventTypeUpdatePowerOnPeriod = "update_power_on_period"
	AttributeKeyOldPowerOnPeriod = "old_power_on_period"
	AttributeKeyNewPowerOnPeriod = "new_power_on_period"

	EventTypeReceiveExperience = "receive_experience"
	AttributeKeyExperience     = "experience"

	EventTypeUpdateUserExperience = "update_user_experience"
	AttributeKeyOwner             = "owner"
	AttributeKeyOldExperience     = "old_experience"
	AttributeKeyNewExperience     = "new_experience"

	EventTypeAddCaller = "add_caller"
	AttributeCaller    = "caller"

	EventTypeRemoveCaller = "remove_caller"

	EventTypeUpdateSaleLevel = "update_sale_level"
	AttributeKeyOldSaleLevel = "old_sale_level"
	AttributeKeyNewSaleLevel = "new_sale_level"
)
