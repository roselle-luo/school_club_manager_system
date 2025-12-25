package models

type ClubCategory struct {
	BaseModel
	Name string `gorm:"size:64;uniqueIndex;not null" json:"name"`
}

type Club struct {
	BaseModel
	Name       string       `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Logo       string       `gorm:"size:255" json:"logo"`
	Intro      string       `gorm:"type:text" json:"intro"`
	Contact    string       `gorm:"size:64" json:"contact"`
	CategoryID uint         `gorm:"index" json:"category_id"`
	Category   ClubCategory `json:"category"`
	Status     string       `gorm:"size:32;default:'pending';index" json:"status"` // pending, approved, rejected
}

type Membership struct {
	BaseModel
	UserID uint   `gorm:"index;uniqueIndex:ux_user_club" json:"user_id"`
	ClubID uint   `gorm:"index;uniqueIndex:ux_user_club" json:"club_id"`
	Status string `gorm:"size:16" json:"status"`
	Role   string `gorm:"size:32;index" json:"role"`
	User   User   `json:"user"`
	Club   Club   `json:"club"`
}
