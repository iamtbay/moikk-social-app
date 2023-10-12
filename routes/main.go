package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/moikk-app/controllers"
	_ "github.com/moikk-app/controllers/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var contrl = controllers.ControllersInit()

type MainRoute struct {
	router *gin.Engine
}

func Routes() *gin.Engine {

	r := MainRoute{
		router: gin.Default(),
	}
	r.router.MaxMultipartMemory = 8 << 20 //8mb
	//
	r.router.Static("/files","./public")
	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.router.Group("api/v1")
	r.SearchRoutes(v1)
	//auth
	r.AuthRoutes(v1)
	//posts
	r.PostsRoutes(v1)
	//comments
	r.CommentsRoutes(v1)
	//likes
	r.LikesRoutes(v1)
	//follows
	r.FollowsRoutes(v1)

	return r.router

}
