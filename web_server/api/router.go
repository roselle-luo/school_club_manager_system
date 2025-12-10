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
	pub.GET("/clubs", controllers.ListClubs)
	pub.GET("/clubs/:clubId", controllers.GetClubDetailPublic)
	pub.GET("/announcements", controllers.ListPublicAnnouncements)
	pub.GET("/activities", controllers.ListPublicActivities)
	pub.GET("/categories", controllers.ListCategories)

	auth := v1.Group("")
	auth.Use(middleware.JWT())
	auth.POST("/upload/image", controllers.UploadImage)
	student := auth.Group("/student")
	student.POST("/clubs/:clubId/apply", controllers.ApplyJoinClub)
	student.POST("/clubs/:clubId/exit", controllers.ExitClub)
	student.GET("/memberships/my", controllers.MyMemberships)
	member := auth.Group("/member")
	leader := auth.Group("/leader")
	leader.GET("/clubs/:clubId/users", controllers.GetClubLeaders)
	leader.GET("/users/:userId/clubs", controllers.GetUserLeaderClubs)
	leader.POST("/clubs/:clubId/members/:userId/role", controllers.SetMemberRoleByLeader)
	leader.GET("/clubs/:clubId/attendance", controllers.ClubAttendance)
	leader.DELETE("/attendance/:id", controllers.DeleteAttendance)
	leader.GET("/clubs/:clubId/memberships", controllers.ListPendingMemberships)
	leader.POST("/clubs/:clubId/memberships/:id/approve", controllers.ApproveMembership)
	leader.POST("/clubs/:clubId/memberships/:id/reject", controllers.RejectMembership)
	leader.GET("/clubs/:clubId/members/users", controllers.ListClubMembers)
	admin := auth.Group("/admin")
	admin.POST("/memberships/:id/role", controllers.UpdateMembershipRole)

	member.POST("/activities/:activityId/signin", controllers.SignIn)
	member.POST("/activities/:activityId/signout", controllers.SignOut)
	member.GET("/attendance/my", controllers.MyAttendance)
	_ = leader
	_ = admin
}
