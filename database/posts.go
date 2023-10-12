package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/moikk-app/helpers"
	mytypes "github.com/moikk-app/types"
)

type PostsDB struct{}

// init
func PostsDBInit() *PostsDB {
	return &PostsDB{}
}

// CREATE A POST
func (x *PostsDB) CreatePost(postInfos *mytypes.NewPost) (*mytypes.GetPost, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	query := `
			INSERT INTO 
				posts(user_id,title,content,files,updated_at)
			VALUES($1,$2,$3,$4,$5)
			RETURNING *
				`
	var newPost mytypes.GetPost
	var t time.Time
	err := conn.QueryRow(ctx, query,
		postInfos.UserID,
		postInfos.Title,
		postInfos.Content,
		pq.Array(postInfos.Files),
		t,
	).Scan(
		&newPost.ID,
		&newPost.UserID,
		&newPost.Title,
		&newPost.Content,
		&newPost.Files,
		&newPost.CreatedAt,
		&newPost.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newPost, nil
}

//UPDATE A POST

func (x *PostsDB) UpdatePost(postID, userID string, updatePostInfos *mytypes.UpdatePost) ([]*mytypes.GetPost, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()

	//get files from db first, if it isn't equal to current files
	//delete it
	//don't allow to add new photos here.
	var currentFiles []string
	query := `SELECT files FROM posts WHERE id=$1 AND user_id=$2`
	err := conn.QueryRow(ctx, query, postID, userID).Scan(&currentFiles)
	if err != nil {
		return nil, err
	}
	//TODO
	//DELETE FILE WHICH USER DELETED!
	existingFiles := make(map[string]bool)
	for _, newFile := range updatePostInfos.Files {
		existingFiles[newFile] = true
	}
	for _, currFile := range currentFiles {
		if !existingFiles[currFile] {
			file := strings.Split(currFile, "/")
			os.Remove(fmt.Sprintf("./public/posts/%v", file[len(file)-1]))
			//CAN ERROR HANDLE HERE.
		}
	}
	//TODO
	//QUERY FOR UPDATE
	var updatedPost []*mytypes.GetPost
	query = `UPDATE posts 
		SET
			title=$1,
			content=$2,
			files=$3,
			updated_at=$4
		WHERE
			user_id=$5
		AND
			id=$6
		RETURNING *
		`
	rows, err := conn.Query(ctx, query,
		updatePostInfos.Title,
		updatePostInfos.Content,
		pq.Array(updatePostInfos.Files),
		time.Now(),
		userID,
		postID,
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		post, err := postScanner(rows)
		if err != nil {
			return nil, err
		}
		updatedPost = append(updatedPost, post)
	}
	return updatedPost, nil
}

// DELETE A POST
func (x *PostsDB) DeletePost(userID, postID string) error {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	query := `
			UPDATE posts
			SET 
				title=$1,
				content=$2,
				files=$3,
				updated_at=$4
			WHERE
				user_id=$5
				id=$6

			`
	err := conn.QueryRow(ctx, query, "", "", pq.Array(""), time.Now(), userID, postID).Scan()
	if err != nil {
		return err
	}
	//return
	return nil
}

// GET A POST
func (x *PostsDB) GetAPost(postID string) ([]*mytypes.GetPostWithUsername, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	var posts []*mytypes.GetPostWithUsername
	query := `
			SELECT 
				posts.id,
				posts.user_id,
				users.username,
				posts.title,
				posts.content,
				posts.files,
				posts.created_at,
				posts.updated_at
			FROM posts
			LEFT JOIN users ON posts.user_id=users.id
			WHERE posts.id=$1	
			`
	rows, err := conn.Query(ctx, query, postID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(notFoundMsg)
		}
		return nil, err
	}
	//scan
	for rows.Next() {
		post, err := postScannerUsername(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	//return
	return posts, nil
}

// GET ALL POSTS
func (x *PostsDB) GetAllPostsByUserID(page int, username string) ([]*mytypes.GetPostWithUsername, *mytypes.Pagination, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()

	//find user id
	var userID string
	err := conn.QueryRow(ctx, "select id from users where username=$1", username).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, errors.New(notFoundMsg)
		}
		return nil, nil, err
	}
	// make page setter operations here
	var totalCount int
	conn.QueryRow(ctx, "SELECT COUNT(*) FROM posts WHERE user_id=$1", userID).Scan(&totalCount)
	//
	currentPage, totalPage := helpers.Pagination(page, totalCount, dataPerPage)

	//todo
	//query
	query := `SELECT
				posts.id,
				posts.user_id,
				users.username,
				posts.title,
				posts.content,
				posts.files,
				posts.created_at,
				posts.updated_at
			FROM posts 
			LEFT JOIN users
			ON posts.user_id=users.id
			WHERE 
				posts.user_id=$1 
			LIMIT $2 OFFSET $3`
	rows, err := conn.Query(ctx, query, userID, dataPerPage, (currentPage*dataPerPage)-dataPerPage)
	if err != nil {
		return nil, nil, err
	}
	//scan
	var allPosts []*mytypes.GetPostWithUsername
	for rows.Next() {
		singlePost, err := postScannerUsername(rows)
		if err != nil {
			return nil, nil, err
		}
		allPosts = append(allPosts, singlePost)
	}
	//return
	return allPosts, &mytypes.Pagination{
		TotalPage:   int(totalPage),
		CurrentPage: currentPage,
		TotalCount:  totalCount,
	}, nil
}

// GET ALL POSTS BY USERID
func (x *PostsDB) GetAllPosts(userID string, pageI int) ([]*mytypes.GetPostWithUsername, *mytypes.Pagination, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//count q for pagination
	query := `SELECT COUNT(*) FROM posts`
	//pagination

	//query1 for followed users
	query = `SELECT followed_user_id FROM follows where follower_user_id=$1 `
	rows, err := conn.Query(ctx, query, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil, errors.New(notFoundMsg)
		}
		return nil, nil, err
	}
	var followedUsers []any
	for rows.Next() {
		var singleUser string
		err := rows.Scan(&singleUser)
		if err != nil {
			return nil, nil, err
		}
		followedUsers = append(followedUsers, singleUser)
	}

	//mapping for userids
	followedUserMap := make([]string, len(followedUsers))
	for i := range followedUsers {
		followedUserMap[i] = fmt.Sprintf("$%d", i+1)
	}
	if len(followedUsers) < 1 {
		return nil, nil, errors.New("user doesn't following anyone")
	}
	//pagination
	var totalCount int
	query2 := `
		select count(*) from posts WHERE 
		user_id IN (` + strings.Join(followedUserMap, ",") + `) `
	err = conn.QueryRow(ctx, query2, followedUsers...).Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}
	currentPage, totalPage := helpers.Pagination(pageI, totalCount, dataPerPage)
	//query2 FOR FOLLOWED USER'S POSTS
	query2 = `
				SELECT 
					posts.id,
					posts.user_id,
					users.username,
					posts.title,
					posts.content,
					posts.files,
					posts.created_at,
					posts.updated_at
				FROM
				posts 
				LEFT JOIN users ON posts.user_id=users.id
				WHERE 
				posts.user_id IN (` + strings.Join(followedUserMap, ",") + `) 
				LIMIT (` + fmt.Sprint(dataPerPage) + `)
				OFFSET (` + fmt.Sprint(dataPerPage*currentPage-dataPerPage) + `)
	`
	rows2, err := conn.Query(ctx, query2, followedUsers...)
	var posts []*mytypes.GetPostWithUsername
	if err != nil {
		return nil, nil, err
	}
	//scan
	for rows2.Next() {
		post, err := postScannerUsername(rows2)
		if err != nil {
			return nil, nil, err
		}
		posts = append(posts, post)
	}
	//return

	return posts, &mytypes.Pagination{
		CurrentPage: currentPage,
		TotalPage:   totalPage,
		TotalCount:  totalCount,
	}, nil
}

//rows scan func

func postScanner(rows pgx.Rows) (*mytypes.GetPost, error) {
	var singlePost mytypes.GetPost
	err := rows.Scan(
		&singlePost.ID,
		&singlePost.UserID,
		&singlePost.Title,
		&singlePost.Content,
		&singlePost.Files,
		&singlePost.CreatedAt,
		&singlePost.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &singlePost, nil
}
// W USERNAME
func postScannerUsername(rows pgx.Rows) (*mytypes.GetPostWithUsername, error) {
	var singlePost mytypes.GetPostWithUsername
	err := rows.Scan(
		&singlePost.ID,
		&singlePost.UserID,
		&singlePost.Username,
		&singlePost.Title,
		&singlePost.Content,
		&singlePost.Files,
		&singlePost.CreatedAt,
		&singlePost.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &singlePost, nil
}
