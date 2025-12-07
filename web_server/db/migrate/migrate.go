package migrate

import (
    "gorm.io/gorm"
    "web_server/db/models"
)

func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Role{},
        &models.ClubCategory{},
        &models.Club{},
        &models.Membership{},
        &models.Announcement{},
        &models.Activity{},
        &models.Attendance{},
        &models.Achievement{},
    )
}

func Seed(db *gorm.DB) error {
    var r models.Role
    db.Where(models.Role{Code: "admin"}).FirstOrCreate(&models.Role{Name: "管理员", Code: "admin"})
    db.Where(models.Role{Code: "user"}).FirstOrCreate(&models.Role{Name: "用户", Code: "user"})
    _ = r
    return nil
}
