package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moikk-app/middlewares"
)

func (r MainRoute) CommentsRoutes(commentR *gin.RouterGroup) {

	comments := commentR.Group("/comments")
	{
		comments.GET("/:id", contrl.Comments.GetAllCommentsByPostID)
		comments.GET("/single/:id", contrl.Comments.GetSingleCommentByID)
		comments.Use(middlewares.LoginMW())
		comments.POST("/new/:id", contrl.Comments.NewComment)
		comments.PATCH("/update/:id/:commentID", contrl.Comments.UpdateComment)
		comments.DELETE("/delete/:id/:commentID", contrl.Comments.DeleteComment)
	}

}
