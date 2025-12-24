package models

type OperationLog struct {
	BaseModel
	OperatorID   uint   `gorm:"column:operator_id" json:"operator_id"`
	OperatorName string `gorm:"column:operator_name" json:"operator_name"`
	ActionType   string `gorm:"column:action_type" json:"action_type"` // e.g. "APPROVE_CLUB", "MODIFY_ATTENDANCE"
	Content      string `gorm:"column:content" json:"content"`
	ClubID       uint   `gorm:"column:club_id;default:0" json:"club_id"` // 0 if system level
}
