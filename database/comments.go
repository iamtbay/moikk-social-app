package database

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/moikk-app/helpers"
	mytypes "github.com/moikk-app/types"
)

type CommentsDB struct{}

// INIT
func CommentsDBInit() *CommentsDB {
	return &CommentsDB{}
}

var commentPerPage = 10

// NEW COMMENT
func (x *CommentsDB) NewComment(postID, userID string, commentInfo *mytypes.NewComment) (*mytypes.GetComment, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	var newComment mytypes.GetComment
	query := `
			INSERT INTO 
			comments(
				post_id, 
				user_id, 
				content
				) 
			VALUES ($1, $2, $3) 
			RETURNING *`
	err := conn.QueryRow(ctx, query,
		postID,
		userID,
		&commentInfo.Content,
	).Scan(
		&newComment.ID,
		&newComment.PostID,
		&newComment.UpperCommentID,
		&newComment.UserID,
		&newComment.Content,
		&newComment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newComment, nil
}

// UPDATE COMMENT
func (x *CommentsDB) UpdateComment(postID, commentID, userID string, updateInfo *mytypes.NewComment) ([]*mytypes.GetComment, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	query := `
			UPDATE comments 
			SET
				content=$1
			WHERE
				post_id = $2 AND 
				id=$3 AND
				user_id = $4
			RETURNING *
	`
	rows, err := conn.Query(ctx, query, updateInfo.Content, postID, commentID, userID)
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	//scan
	var updatedInfo []*mytypes.GetComment
	for rows.Next() {
		singleUpdate, err := commentsScanner(rows)
		if err != nil {
			return nil, err
		}
		updatedInfo = append(updatedInfo, singleUpdate)
	}
	//return
	return updatedInfo, nil
}

// DELETE COMMENT
func (x *CommentsDB) DeleteComment(postID, commentID, userID string) error {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	query := `
			DELETE FROM 
				comments 
			WHERE
				id=$1 AND
				post_id=$2 AND
				user_id=$3
			`
	_, err := conn.Exec(ctx, query, commentID, postID, userID)
	if err != nil {
		return errors.New("something went wrong")
	}
	//return
	return nil
}

// GET COMMENTS
func (x *CommentsDB) GetComments(postID, pageStr string) ([]*mytypes.GetCommentWithUsername, *mytypes.Pagination, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//pagination
	var totalCount int
	page, _ := strconv.Atoi(pageStr)
	countQ := `SELECT COUNT(*) from comments WHERE post_id=$1`
	conn.QueryRow(ctx, countQ, postID).Scan(&totalCount)
	currentPage, totalPage := helpers.Pagination(page, totalCount, commentPerPage)
	//todo
	//query
	query := `
			SELECT
			comments.id,
			comments.post_id,
			comments.upper_comment_id,
			comments.user_id,
			users.username,
			comments.content,
			comments.created_at
			FROM comments 
			LEFT JOIN users ON comments.user_id=users.id
			WHERE comments.post_id=$1 LIMIT $2 OFFSET $3`
	rows, err := conn.Query(ctx, query, postID, commentPerPage, (currentPage*commentPerPage)-commentPerPage)
	if err != nil {
		return nil, nil, errors.New("something went wrong")
	}
	//scan
	var comments []*mytypes.GetCommentWithUsername
	for rows.Next() {
		comment, err := commentsScannerWithUsername(rows)
		if err != nil {
			return nil, nil, err
		}
		comments = append(comments, comment)
	}
	//return
	return comments, &mytypes.Pagination{
		CurrentPage: currentPage,
		TotalPage:   totalPage,
		TotalCount:  totalCount,
	}, nil
}

// GET SIGNLE COMMENT
func (x *CommentsDB) GetSingleComment(commentID string) ([]*mytypes.GetCommentWithUsername, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	var comments []*mytypes.GetCommentWithUsername
	query := `SELECT
				comments.id,
				comments.post_id,
				comments.upper_comment_id,
				comments.user_id,
				users.username,
				comments.content,
				comments.created_at
			FROM comments 
			LEFT JOIN users ON comments.user_id=users.id
			WHERE comments.id=$1`
	rows, err := conn.Query(ctx, query, commentID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(notFoundMsg)
		}
		return nil, err
	}
	//scan
	for rows.Next() {
		comment, err := commentsScannerWithUsername(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	//return
	return comments, nil
}

// todo SCANNER
func commentsScanner(rows pgx.Rows) (*mytypes.GetComment, error) {
	var singleComment mytypes.GetComment
	err := rows.Scan(
		&singleComment.ID,
		&singleComment.PostID,
		&singleComment.UpperCommentID,
		&singleComment.UserID,
		&singleComment.Content,
		&singleComment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &singleComment, nil
}
func commentsScannerWithUsername(rows pgx.Rows) (*mytypes.GetCommentWithUsername, error) {
	var singleComment mytypes.GetCommentWithUsername
	err := rows.Scan(
		&singleComment.ID,
		&singleComment.PostID,
		&singleComment.UpperCommentID,
		&singleComment.UserID,
		&singleComment.Username,
		&singleComment.Content,
		&singleComment.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &singleComment, nil
}
