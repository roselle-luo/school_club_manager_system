package authz

import (
	"web_server/db/models"
	"web_server/internal/store"
)

func IsAdmin(u *models.User) bool { return u.Role.Code == "admin" }

func IsClubLeader(userID uint, clubID uint) bool {
	var m models.Membership
	if err := store.DB().Where("user_id = ? AND club_id = ? AND role IN ?", userID, clubID, []string{"leader", "advisor"}).First(&m).Error; err != nil {
		return false
	}
	return true
}

func IsClubMember(userID uint, clubID uint) bool {
	var m models.Membership
	if err := store.DB().Where("user_id = ? AND club_id = ?", userID, clubID).First(&m).Error; err != nil {
		return false
	}
	return true
}
