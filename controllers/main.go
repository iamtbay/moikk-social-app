package controllers

import (
	"github.com/moikk-app/database"
	//swagger
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	 
)

type Controllers struct {
	Auth     *AuthController
	Comments *CommentsController
	Follows  *FollowsController
	Likes    *LikesController
	Posts    *PostsController
	Search   *SearchStruct
}

var authDB = database.AuthDBInit()
var followsDB = database.FollowsDBInit()
var postsDB = database.PostsDBInit()
var commentsDB = database.CommentsDBInit()
var likesDB = database.LikesDBInit()
var searchDB = database.SearchDBInit()

// @title Moikk-Social-App
// @version 1.0
// @description This is a sample server MoikkSocialApp.
// @BasePath api/v1
func ControllersInit() *Controllers {
	return &Controllers{
		Auth:     AuthControllerInit(),
		Posts:    PostControllerinit(),
		Comments: CommentsControllerInit(),
		Follows:  FollowsControllerInit(),
		Likes:    LikesControllerInit(),
		Search:   SearchControllerInit(),
	}
}

