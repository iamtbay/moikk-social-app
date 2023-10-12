package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moikk-app/helpers"
	mytypes "github.com/moikk-app/types"
)

type CommentsController struct{}

// COMMENT INIT
func CommentsControllerInit() *CommentsController {
	return &CommentsController{}
}

// @Summary New Comment
// @Description New Comment to a Post ID
// @ID newcomment
// @Accept json
// @Produce json
// @Param commentInfo body mytypes.NewComment true "New Comment data"
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /comments/new/{id} [post]
// @Tags comments
func (x *CommentsController) NewComment(c *gin.Context) {
	//get userid from jwt and postid from url
	postID := c.Param("id")
	userID, _ := helpers.GetUserIDFromJWT(c)
	//get user input
	var newCommentInfo *mytypes.NewComment
	err := c.BindJSON(&newCommentInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	//send it to make db ops
	newComment, err := commentsDB.NewComment(postID, userID, newCommentInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//return values
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Succesfully commented",
		"data":    newComment,
	})
}

// @Summary Update Comment
// @Description Update a Comment by ID
// @ID updatecomment
// @Accept json
// @Produce json
// @Param commentInfo body mytypes.NewComment true "Update Comment data"
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /comments/update/{id}/{commentID} [patch]
// @Tags comments
func (x *CommentsController) UpdateComment(c *gin.Context) {
	//get userid from jwt and postid&commentid from url
	postID := c.Param("id")
	commentID := c.Param("commentID")
	userID, _ := helpers.GetUserIDFromJWT(c)
	//get user inputs
	var updateInfo *mytypes.NewComment
	err := c.BindJSON(&updateInfo)
	if err != nil {
		return
	}
	//send it to make db ops
	updatedInfo, err := commentsDB.UpdateComment(postID, commentID, userID, updateInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//return values
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Comment succesfully updated",
		"data":    updatedInfo,
	})
}

// @Summary Delete a Comment
// @Description Delete a Comment by Comment ID
// @ID deletecomment
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /comments/delete/{id}/{commentID} [delete]
// @Tags comments
func (x *CommentsController) DeleteComment(c *gin.Context) {
	//get commentid and postid
	postID := c.Param("id")
	commentID := c.Param("commentID")
	//get userid
	userID, _ := helpers.GetUserIDFromJWT(c)
	//send it to db
	err := commentsDB.DeleteComment(postID, commentID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//return message
	c.JSON(200, gin.H{
		"success": true,
		"message": "deleted succesfully",
	})
}

// @Summary Get All Comments By PostID
// @Description Get all comments relavant to post id
// @ID get-all-comments-by-postID
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /comments/{id} [get]
// @Tags comments
func (x *CommentsController) GetAllCommentsByPostID(c *gin.Context) {
	//get post id
	postID := c.Param("id")
	page := c.Query("page")
	//send it to db
	comments, pages, err := commentsDB.GetComments(postID, page)
	//check count of coming data
	var msg = "Comments got succesfully"
	if len(comments) < 1 {
		msg = "There is no comments in this post"
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//return items
	c.JSON(200, gin.H{
		"success": true,
		"message": msg,
		"pages":   pages,
		"data":    comments,
	})
}

// @Summary Get Single Comment By ID
// @Description Get single comment by commentID
// @ID get-a-comment-by-id
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 400 {string} Bad Request
// @Router /comments/single/{id} [get]
// @Tags comments
func (x *CommentsController) GetSingleCommentByID(c *gin.Context) {
	//get comment id from url
	commentID := c.Param("id")
	//send it to db
	msg := "Comment got succesfully"
	comment, err := commentsDB.GetSingleComment(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(comment) < 1 {
		msg = "No comments to show."
	}
	//return
	c.JSON(200, gin.H{
		"success": true,
		"message": msg,
		"data":    comment,
	})

}
