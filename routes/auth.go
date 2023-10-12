package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moikk-app/middlewares"
)

func (r MainRoute) AuthRoutes(authR *gin.RouterGroup) {
	auth := authR.Group("/auth")
	{
		auth.POST("/login", middlewares.LogoutMW(), contrl.Auth.Login)
		auth.POST("/register", middlewares.LogoutMW(), contrl.Auth.Register)
		auth.Use(middlewares.LoginMW())
		auth.PATCH("/updateUser", contrl.Auth.UpdateUserInfo)
		auth.POST("/profilePhoto", contrl.Auth.UpdateProfilePhoto)
		auth.POST("/logout", contrl.Auth.Logout)
	}

}
