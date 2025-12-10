package migrate

import (
	"fmt"
	"gorm.io/gorm"
	"time"
	"web_server/db/models"
	"web_server/pkg/password"
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
	db.Where(models.Role{Code: "admin"}).FirstOrCreate(&models.Role{Name: "管理员", Code: "admin"})
	db.Where(models.Role{Code: "user"}).FirstOrCreate(&models.Role{Name: "用户", Code: "user"})

	// categories
	cats := []string{"文学类", "音乐类", "志愿服务", "社会实践", "电竞类", "摄影类", "创业类"}
	catIDs := make([]uint, len(cats))
	for i, name := range cats {
		var c models.ClubCategory
		db.Where(models.ClubCategory{Name: name}).FirstOrCreate(&c)
		catIDs[i] = c.ID
	}

	// demo users (leaders)
	var userRole models.Role
	db.Where(models.Role{Code: "user"}).First(&userRole)
	makePhone := func(i int) string { return fmt.Sprintf("138%08d", 10000+i) }
	hash, _ := password.Hash("123456")
	leaders := make([]models.User, 0, 12)
	for i := 0; i < 12; i++ {
		u := models.User{Account: fmt.Sprintf("leader%02d", i+1), Password: hash, Name: fmt.Sprintf("社长%02d号", i+1), Gender: "M", College: "信息学院", StudentNo: fmt.Sprintf("2025%04d", i+1), Phone: makePhone(i + 1), RoleID: userRole.ID}
		db.Where(models.User{Account: u.Account}).FirstOrCreate(&u)
		leaders = append(leaders, u)
	}

	// clubs
	type ClubSeed struct {
		Name, Intro string
		CatIdx      int
	}
	seeds := []ClubSeed{
		{"文学社", "阅读交流与写作分享，提升文学素养。", 0},
		{"诗词社", "古典诗词鉴赏与创作交流。", 0},
		{"吉他社", "吉他弹唱练习、舞台表演实践。", 1},
		{"合唱团", "多声部合唱训练与演出。", 1},
		{"志愿者协会", "校园及社区志愿服务活动组织。", 2},
		{"红十字会", "应急救护与公益慈善活动。", 2},
		{"社会调研社", "社会议题调研与报告写作。", 3},
		{"实践创新社", "项目实践与校企合作体验。", 3},
		{"电竞俱乐部", "热门电竞项目训练与赛事组织。", 4},
		{"桌游俱乐部", "策略桌游对战与线下联赛。", 4},
		{"摄影协会", "摄影外拍与后期交流。", 5},
		{"创业者联盟", "创业思维训练与路演交流。", 6},
	}
	clubs := make([]models.Club, 0, len(seeds))
	for i, s := range seeds {
		catID := catIDs[s.CatIdx]
		cl := models.Club{Name: s.Name, Logo: "", Intro: s.Intro, Contact: "", CategoryID: catID}
		db.Where(models.Club{Name: cl.Name}).FirstOrCreate(&cl)
		clubs = append(clubs, cl)
		// membership: assign leader
		if i < len(leaders) {
			m := models.Membership{UserID: leaders[i].ID, ClubID: cl.ID, Status: "approved", Role: "leader"}
			db.Where(models.Membership{UserID: m.UserID, ClubID: m.ClubID}).FirstOrCreate(&m)
		}
	}

	// activities (public)
	now := time.Now()
	formats := []string{"讲座", "招新", "比赛"}
	for _, cl := range clubs {
		for j, f := range formats {
			a := models.Activity{Subject: fmt.Sprintf("%s：%s第%d期", f, cl.Name, j+1), Time: now.AddDate(0, 0, 7*(j+1)).Format("2006-01-02 15:04"), Place: "学生活动中心", Target: "全体学生", Scope: "public", ClubID: cl.ID}
			db.Where(models.Activity{Subject: a.Subject, ClubID: cl.ID}).FirstOrCreate(&a)
		}
	}
	return nil
}
