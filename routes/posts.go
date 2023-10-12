package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moikk-app/middlewares"
)

func (R MainRoute) PostsRoutes(postsR *gin.RouterGroup) {
	posts := postsR.Group("/posts")
	{
		//get all posts
		posts.GET("/", contrl.Posts.GetAllPosts)
		//get single post
		posts.GET("/single/:id", contrl.Posts.GetAPost)
		//get user's posts
		posts.GET("/user/:username", contrl.Posts.GetAllPostsByUserID)
		//
		posts.Use(middlewares.LoginMW())
		
		posts.POST("/upload",contrl.Posts.UploadFile)
		posts.POST("/new", contrl.Posts.CreatePost)
		//CAN URL CHANGE? LIKE /:id only.
		posts.PATCH("/:id", contrl.Posts.UpdatePost)
		posts.DELETE("/:id", contrl.Posts.DeletePost)
	}
}
