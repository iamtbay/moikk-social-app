package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moikk-app/middlewares"
)

func (r MainRoute) FollowsRoutes(followsR *gin.RouterGroup) {
	follows := followsR.Group("/follows")
	{
		//can add get follows for user
		follows.Use(middlewares.LoginMW())
		follows.POST("/:username", contrl.Follows.FollowUser)
		follows.PATCH("/:username", contrl.Follows.UnfollowUser)
		follows.GET("/followers/:username", contrl.Follows.GetFollowers)
		follows.GET("/followeds/:username", contrl.Follows.GetFolloweds)
	}
}
