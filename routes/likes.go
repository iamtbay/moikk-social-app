package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moikk-app/middlewares"
)

func (r MainRoute) LikesRoutes(likesR *gin.RouterGroup) {

	likes := likesR.Group("/likes")
	{
		likes.GET("/:id", contrl.Likes.GetLikesByPostID)
		likes.GET("/user/:username", contrl.Likes.GetUserLikes)
		likes.Use(middlewares.LoginMW())
		likes.POST("/:id", contrl.Likes.LikePost)
		likes.PATCH("/:id", contrl.Likes.UnlikePost)
	}
}
