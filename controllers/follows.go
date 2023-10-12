package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moikk-app/helpers"
)

type FollowsController struct{}

// init
func FollowsControllerInit() *FollowsController {
	return &FollowsController{}
}

// @Summary Get User's Followers
// @Description Get all followers by username
// @ID get-all-followers-by-username
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /follows/followers/{username} [get]
// @Tags follows
func (x *FollowsController) GetFollowers(c *gin.Context) {
	targetUsername := c.Param("username")
	followers, err := followsDB.GetUserFollowers(targetUsername)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"followers": followers,
	})
}

// @Summary Get User's Followeds
// @Description Get all followeds by username
// @ID get-all-followeds-by-username
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /follows/followeds/{username} [get]
// @Tags follows
func (x *FollowsController) GetFolloweds(c *gin.Context) {
	targetUsername := c.Param("username")
	followeds, err := followsDB.GetUserFolloweds(targetUsername)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"followers": followeds,
	})
}


// @Summary Follow A User
// @Description Follow a user by username
// @ID follow-user
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /follows/{username} [post]
// @Tags follows
func (x *FollowsController) FollowUser(c *gin.Context) {
	//get user's own id from jwt
	userID, err := helpers.GetUserIDFromJWT(c)
	targetUserID := c.Param("username")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "something went wrong try again later",
		})
		return
	}
	//add the target user as followed
	err = followsDB.FollowUser(userID, targetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//
	c.JSON(200, gin.H{
		"success":        true,
		"message":        "you started to follow user",
		"followedUserID": targetUserID,
	})
}

// @Summary Unfollow A User
// @Description Unfollow a user by username
// @ID unfollow-user
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /follows/{username} [patch]
// @Tags follows
func (x *FollowsController) UnfollowUser(c *gin.Context) {

	//get user's own id from jwt
	userID, err := helpers.GetUserIDFromJWT(c)
	targetUserID := c.Param("username")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "something went wrong try again later",
		})
		return
	}
	//add the target user as followed
	err = followsDB.UnfollowUser(userID, targetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//
	c.JSON(200, gin.H{
		"success":          true,
		"message":          "you succesfully unfollowed the account",
		"unFollowedUserID": targetUserID,
	})
}
