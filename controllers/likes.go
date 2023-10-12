package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/moikk-app/helpers"
)

type LikesController struct{}

// init
func LikesControllerInit() *LikesController {
	return &LikesController{}
}

// @Summary Like a Post
// @Description Like a Post by postID
// @ID like-post-by-id
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /likes/{id} [post]
// @Tags likes
func (x *LikesController) LikePost(c *gin.Context) {
	//get postID from url and user id from jwt
	postID := c.Param("id")
	userID, _ := helpers.GetUserIDFromJWT(c)
	//send db
	err := likesDB.LikePost(postID, userID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//return
	c.JSON(200, gin.H{
		"success": true,
		"message": "liked succesfully",
	})
}

// @Summary Unlike a Post
// @Description Unlike a Post by postID
// @ID unlike-post-by-id
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /likes/{id} [patch]
// @Tags likes
func (x *LikesController) UnlikePost(c *gin.Context) {
	//get postID from url and user id from jwt
	postID := c.Param("id")
	userID, _ := helpers.GetUserIDFromJWT(c)
	//send db
	err := likesDB.UnlikePost(postID, userID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//return
	c.JSON(200, gin.H{
		"success": true,
		"message": "unliked succesful",
	})
}

// @Summary Get Likes By PostID
// @Description Get All Likes by PostID
// @ID get-all-likes-by-id
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /likes/{id} [get]
// @Tags likes
func (x *LikesController) GetLikesByPostID(c *gin.Context) {
	//get postID from url and user id from jwt
	postID := c.Param("id")
	//send db
	data, err := likesDB.GetLikesByPostID(postID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//check count of coming data
	var msg string = "Succesfully got likes"
	if len(data) < 1 {
		msg = "anyone liked this post yet."
	}
	//return
	c.JSON(200, gin.H{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

// @Summary Get Likes By Username
// @Description Get User's Likes 
// @ID get-all-likes-by-username
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /likes/user/{username} [get]
// @Tags likes
func (x *LikesController) GetUserLikes(c *gin.Context) {
	//get user id from jwt
	userID := c.Param("username")
	//send db
	data, err := likesDB.GetUserLikes(userID)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//check count of coming data
	var msg string = "Succesfully got user's likes"
	if len(data) < 1 {
		msg = "User doesn't like any post yet."
	}
	//return
	c.JSON(200, gin.H{
		"success": true,
		"message": msg,
		"data":    data,
	})
}
