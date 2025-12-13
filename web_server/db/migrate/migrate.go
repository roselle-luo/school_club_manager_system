package migrate

import (
	"fmt"
	"time"
	"web_server/db/models"
	"web_server/pkg/password"

	"gorm.io/gorm"
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
		&models.ActivityParticipant{},
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
			start := now.AddDate(0, 0, 7*(j+1))
			end := start.Add(2 * time.Hour)
			publish := now
			a := models.Activity{
				Subject:         fmt.Sprintf("%s：%s第%d期", f, cl.Name, j+1),
				Time:            start.Format("2006-01-02 15:04"),
				Place:           "学生活动中心",
				Target:          "全体学生",
				Scope:           "public",
				ClubID:          cl.ID,
				Content:         fmt.Sprintf("%s活动内容：欢迎参与%s的主题活动，现场有精彩互动与分享。", f, cl.Name),
				StartAt:         &start,
				EndAt:           &end,
				MaxParticipants: 0,
				PublishAt:       &publish,
			}
			db.Where(models.Activity{Subject: a.Subject, ClubID: cl.ID}).FirstOrCreate(&a)
		}
	}
	// announcements (public)
	for _, cl := range clubs {
		an1 := models.Announcement{Title: cl.Name + "招新啦", Content: "欢迎加入" + cl.Name + "，一起成长与交流。", Scope: "public", ClubID: cl.ID}
		an2 := models.Announcement{Title: cl.Name + "社团会议通知", Content: "本周将召开例会，请关注具体时间与地点。", Scope: "public", ClubID: cl.ID}
		db.Where(models.Announcement{Title: an1.Title, ClubID: cl.ID}).FirstOrCreate(&an1)
		db.Where(models.Announcement{Title: an2.Title, ClubID: cl.ID}).FirstOrCreate(&an2)
	}
	return nil
}

func BackfillActivities(db *gorm.DB) error {
	var acts []models.Activity
	if err := db.Model(&models.Activity{}).Preload("Club").Find(&acts).Error; err != nil {
		return err
	}
	for i := range acts {
		a := &acts[i]
		changed := false
		if a.PublishAt == nil {
			p := a.CreatedAt
			a.PublishAt = &p
			changed = true
		}
		if a.Content == "" {
			clubName := ""
			if a.Club.ID != 0 && a.Club.Name != "" {
				clubName = a.Club.Name
			}
			if clubName != "" {
				a.Content = fmt.Sprintf("%s活动详情：由%s组织，欢迎参与交流与互动。", a.Subject, clubName)
			} else {
				a.Content = fmt.Sprintf("%s活动详情：欢迎参与交流与互动。", a.Subject)
			}
			changed = true
		}
		parseTime := func(s string) (time.Time, bool) {
			layouts := []string{"2006-01-02 15:04", time.RFC3339, "2006-01-02 15:04:05"}
			for _, l := range layouts {
				if t, err := time.ParseInLocation(l, s, time.Local); err == nil {
					return t, true
				}
			}
			return time.Time{}, false
		}
		if a.StartAt == nil {
			if a.Time != "" {
				if st, ok := parseTime(a.Time); ok {
					a.StartAt = &st
				}
			}
			if a.StartAt == nil {
				t := a.CreatedAt
				a.StartAt = &t
			}
			changed = true
		}
		if a.EndAt == nil {
			if a.StartAt != nil {
				et := a.StartAt.Add(2 * time.Hour)
				a.EndAt = &et
			} else {
				t := a.CreatedAt.Add(2 * time.Hour)
				a.EndAt = &t
			}
			changed = true
		}
		if a.MaxParticipants < 0 {
			a.MaxParticipants = 0
			changed = true
		}
		if changed {
			_ = db.Save(a).Error
		}
	}
	return nil
}

func MigrateAttendanceActivityNullable(db *gorm.DB) error {
	// 将 attendances.activity_id 改为可空，满足社团级打卡不关联活动的场景
	// 保留外键约束，NULL 值不触发外键检查
	sql := "ALTER TABLE `attendances` MODIFY COLUMN `activity_id` BIGINT UNSIGNED NULL"
	if err := db.Exec(sql).Error; err != nil {
		return err
	}
	return nil
}
