package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	mytypes "github.com/moikk-app/types"
)

type LikesDB struct{}

//init

func LikesDBInit() *LikesDB {
	return &LikesDB{}
}

// LIKE POST
func (x *LikesDB) LikePost(postID, userID string) error {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//
	query := `SELECT EXISTS ( SELECT 1 FROM likes WHERE post_id=$1 AND user_id=$2)`
	var isExist bool
	conn.QueryRow(ctx, query, postID, userID).Scan(&isExist)
	if isExist {
		return errors.New("user already liked this post")
	}
	query = `INSERT INTO likes(post_id,user_id) VALUES($1,$2)`
	_, err := conn.Exec(ctx, query, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

// UNLIKE POST
func (x *LikesDB) UnlikePost(postID, userID string) error {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	query := `SELECT EXISTS ( SELECT 1 FROM likes WHERE post_id=$1 AND user_id=$2)`
	var isExist bool
	conn.QueryRow(ctx, query, postID, userID).Scan(&isExist)
	if !isExist {
		return errors.New("user didn't like this post")
	}
	query = `DELETE FROM likes WHERE post_id=$1 AND user_id=$2`
	_, err := conn.Exec(ctx, query, postID, userID)
	if err != nil {
		return err
	}
	return nil
}

// GET LIKES BY POST ID
func (x *LikesDB) GetLikesByPostID(postID string) ([]*mytypes.GetLikes, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	query := `
		SELECT 
		likes.id,
		likes.post_id,
		likes.user_id,
		users.username,
		likes.created_at
		FROM likes 
		LEFT JOIN users ON likes.user_id=users.id
		WHERE post_id=$1`
	rows, err := conn.Query(ctx, query, postID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(notFoundMsg)
		}
		return nil, err
	}
	//scan
	var allLikes []*mytypes.GetLikes
	for rows.Next() {
		like, err := likesScanner(rows)
		if err != nil {
			return nil, err
		}
		allLikes = append(allLikes, like)
	}
	//
	return allLikes, nil
}

// GET USER'S LIKES
func (x *LikesDB) GetUserLikes(username string) ([]*mytypes.GetLikes, error) {
	//ctx
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query
	//find the userid
	var userID string
	query := `select id from users where username=$1`
	err := conn.QueryRow(ctx, query, username).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(notFoundMsg)
		}
		return nil, err
	}
	//query for likes
	query = `
		SELECT 
			likes.id,
			likes.post_id,
			likes.user_id,
			users.username,
			likes.created_at
		FROM likes 
			LEFT JOIN users ON likes.user_id=users.id
		WHERE user_id=$1`
	rows, err := conn.Query(ctx, query, userID)
	if err != nil {

		return nil, err
	}
	//scan
	var allLikes []*mytypes.GetLikes
	for rows.Next() {
		like, err := likesScanner(rows)
		if err != nil {
			return nil, err
		}
		allLikes = append(allLikes, like)
	}
	//
	return allLikes, nil
}

// likeScanner
func likesScanner(rows pgx.Row) (*mytypes.GetLikes, error) {
	var singleLike mytypes.GetLikes
	err := rows.Scan(
		&singleLike.ID,
		&singleLike.PostID,
		&singleLike.UserID,
		&singleLike.Username,
		&singleLike.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &singleLike, nil
}
