package models

type Role struct {
	BaseModel
	Name string `gorm:"size:32;uniqueIndex;not null" json:"name"`
	Code string `gorm:"size:32;uniqueIndex;not null" json:"code"`
}

type User struct {
	BaseModel
	Account   string `gorm:"size:64;uniqueIndex;not null" json:"account"`
	Password  string `gorm:"size:128;not null" json:"-"`
	Name      string `gorm:"size:64" json:"name"`
	Gender    string `gorm:"size:8" json:"gender"`
	College   string `gorm:"size:64" json:"college"`
	StudentNo string `gorm:"size:32" json:"student_no"`
	Phone     string `gorm:"size:20" json:"phone"`
	RoleID    uint   `gorm:"index" json:"role_id"`
	Role      Role   `json:"role"`
}
