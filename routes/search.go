package routes

import "github.com/gin-gonic/gin"



func (r MainRoute) SearchRoutes(searchR *gin.RouterGroup) {

	search := searchR.Group("/search")
	{
		search.GET("/posts", contrl.Search.SearchPosts)
		search.GET("/users", contrl.Search.SearchUsers)
	}

}
