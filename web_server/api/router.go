package api

import (
	"web_server/internal/controllers"
	"web_server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	r.Use(middleware.CORS())
	v1 := r.Group("/api/v1")
	v1.GET("/health", controllers.Health)

	pub := v1.Group("/public")
	pub.POST("/login", controllers.Login)
	pub.POST("/register", controllers.Register)
	pub.POST("/clubs/register", controllers.RegisterClub)
	pub.POST("/upload/image", controllers.UploadImage)
	pub.GET("/clubs", controllers.ListClubs)
	pub.GET("/clubs/:clubId", controllers.GetClubDetailPublic)
	pub.GET("/announcements", controllers.ListPublicAnnouncements)
	pub.GET("/activities", controllers.ListPublicActivities)
	pub.GET("/activities/:activityId", controllers.GetPublicActivityDetail)
	pub.GET("/categories", controllers.ListCategories)

	auth := v1.Group("")
	auth.Use(middleware.JWT())
	auth.POST("/upload/image", controllers.UploadImage)
	student := auth.Group("/student")
	student.POST("/clubs/:clubId/apply", controllers.ApplyJoinClub)
	student.POST("/clubs/:clubId/exit", controllers.ExitClub)
	student.GET("/memberships/my", controllers.MyMemberships)
	student.GET("/me", controllers.MyProfile)
	student.PUT("/me", controllers.UpdateMyProfile)
	student.PUT("/password", controllers.ChangePassword)
	member := auth.Group("/member")
	leader := auth.Group("/leader")
	leader.GET("/clubs/:clubId/users", controllers.GetClubLeaders)
	leader.GET("/users/:userId/clubs", controllers.GetUserLeaderClubs)
	leader.POST("/clubs/:clubId/members/:userId/role", controllers.SetMemberRoleByLeader)
	leader.DELETE("/clubs/:clubId/members/:userId", controllers.KickMember)
	leader.GET("/clubs/:clubId/attendance", controllers.ClubAttendance)
	leader.DELETE("/attendance/:id", controllers.DeleteAttendance)
	leader.GET("/clubs/:clubId/memberships", controllers.ListPendingMemberships)
	leader.POST("/clubs/:clubId/memberships/:id/approve", controllers.ApproveMembership)
	leader.POST("/clubs/:clubId/memberships/:id/reject", controllers.RejectMembership)
	leader.GET("/clubs/:clubId/members/users", controllers.ListClubMembers)

	leader.GET("/clubs/:clubId/announcements", controllers.ListClubAnnouncements)
	leader.POST("/clubs/:clubId/announcements", controllers.CreateAnnouncement)
	leader.PUT("/clubs/:clubId/announcements/:id", controllers.UpdateAnnouncement)
	leader.DELETE("/clubs/:clubId/announcements/:id", controllers.DeleteAnnouncement)

	leader.GET("/clubs/:clubId/logs", controllers.ListOperationLogs)

	// 考勤管理相关接口
	leader.GET("/attendance/list", controllers.ListManagedAttendance)
	leader.POST("/attendance/:id/signout", controllers.ForceSignOut)

	admin := auth.Group("/admin")
	admin.DELETE("/clubs/:clubId", controllers.DissolveClub)
	admin.POST("/memberships/:id/role", controllers.UpdateMembershipRole)
	admin.GET("/attendance", controllers.ListManagedAttendance) // 保留原有路由，指向新控制器
	admin.GET("/clubs/audit", controllers.ListPendingClubs)
	admin.POST("/clubs/:id/audit", controllers.AuditClub)

	member.POST("/activities/:activityId/signin", controllers.SignIn)
	member.POST("/activities/:activityId/signout", controllers.SignOut)
	member.POST("/activities/:activityId/register", controllers.RegisterActivity)
	member.GET("/activities/:activityId/register", controllers.GetRegisterStatus)
	member.DELETE("/activities/:activityId/register", controllers.CancelRegisterActivity)
	member.POST("/clubs/:clubId/signin", controllers.ClubSignIn)
	member.POST("/clubs/:clubId/signout", controllers.ClubSignOut)
	member.GET("/attendance/my", controllers.MyAttendance)
	_ = leader
	_ = admin
}
