package api

import (
    "github.com/gin-gonic/gin"
    "web_server/internal/middleware"
    "web_server/internal/controllers"
)

func Register(r *gin.Engine) {
    r.Use(middleware.CORS())
    v1 := r.Group("/api/v1")
    v1.GET("/health", controllers.Health)

    pub := v1.Group("/public")
    pub.POST("/login", controllers.Login)

    auth := v1.Group("")
    auth.Use(middleware.JWT())
    auth.POST("/upload/image", controllers.UploadImage)
    student := auth.Group("/student")
    member := auth.Group("/member")
    leader := auth.Group("/leader")
    leader.GET("/clubs/:clubId/users", controllers.GetClubLeaders)
    leader.GET("/users/:userId/clubs", controllers.GetUserLeaderClubs)
    leader.POST("/clubs/:clubId/members/:userId/role", controllers.SetMemberRoleByLeader)
    admin := auth.Group("/admin")
    admin.POST("/memberships/:id/role", controllers.UpdateMembershipRole)

    _ = student
    _ = member
    _ = leader
    _ = admin
}
