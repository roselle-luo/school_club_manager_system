package models

import "time"

type Announcement struct {
	BaseModel
	Title   string `gorm:"size:128;not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`
	Scope   string `gorm:"size:16" json:"scope"`
	ClubID  uint   `gorm:"index" json:"club_id"`
	Club    Club   `json:"club"`
}

type Activity struct {
	BaseModel
	Subject         string     `gorm:"size:128;not null" json:"subject"`
	Time            string     `gorm:"size:64" json:"time"`
	Place           string     `gorm:"size:128" json:"place"`
	Target          string     `gorm:"size:64" json:"target"`
	Scope           string     `gorm:"size:16" json:"scope"`
	ClubID          uint       `gorm:"index" json:"club_id"`
	Club            Club       `json:"club"`
	Content         string     `gorm:"type:text" json:"content"`
	StartAt         *time.Time `json:"start_at"`
	EndAt           *time.Time `json:"end_at"`
	MaxParticipants int        `json:"max_participants"`
	PublishAt       *time.Time `json:"publish_at"`
}

type Attendance struct {
	BaseModel
	UserID          uint       `gorm:"index" json:"user_id"`
	ActivityID      *uint      `gorm:"index" json:"activity_id"`
	ClubID          uint       `gorm:"index" json:"club_id"`
	Type            string     `gorm:"size:16" json:"type"`
	SigninAt        *time.Time `json:"signin_at"`
	SignoutAt       *time.Time `json:"signout_at"`
	DurationMinutes int        `json:"duration_minutes"`
	DurationHours   float64    `gorm:"type:decimal(8,2)" json:"duration_hours"`
	Activity        Activity   `json:"activity"`
}

type Achievement struct {
	BaseModel
	Name    string `gorm:"size:128;not null" json:"name"`
	Desc    string `gorm:"type:text" json:"desc"`
	Type    string `gorm:"size:32" json:"type"`
	CertImg string `gorm:"size:255" json:"cert_img"`
	ClubID  uint   `gorm:"index" json:"club_id"`
	Status  string `gorm:"size:16" json:"status"`
}

type ActivityParticipant struct {
	BaseModel
	UserID     uint   `gorm:"index" json:"user_id"`
	ActivityID uint   `gorm:"index" json:"activity_id"`
	ClubID     uint   `gorm:"index" json:"club_id"`
	Status     string `gorm:"size:16" json:"status"`
}
