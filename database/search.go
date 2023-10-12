package database

import (
	"context"
	"fmt"

	"github.com/moikk-app/helpers"
	mytypes "github.com/moikk-app/types"
)

type SearchDB struct{}

// init
func SearchDBInit() *SearchDB {
	return &SearchDB{}
}

// search users
func (x *SearchDB) SearchUsers(keyword string, pageI int) ([]*mytypes.GetUser, *mytypes.Pagination, error) {
	//context
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query for pagination
	query := `
	SELECT count(*) FROM users WHERE 
	lower(username) LIKE '%' || $1  || '%' 
	OR 
	lower(name) LIKE '%' || $1 || '%'`
	var totalCount int
	err := conn.QueryRow(ctx, query, keyword).Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}
	//query for main
	currentPage, totalPage := helpers.Pagination(pageI, totalCount, dataPerPage)
	query = `
			SELECT id,name,username,profile_photo FROM users 
			WHERE 
			lower(username) LIKE '%' || $1  || '%' 
			OR 
			lower(name) LIKE '%' || $1 || '%'
			LIMIT $2
			OFFSET $3
	`
	rows, err := conn.Query(ctx, query, keyword, dataPerPage, dataPerPage*currentPage-dataPerPage)
	if err != nil {
		return nil, nil, err
	}
	var users []*mytypes.GetUser
	for rows.Next() {
		var user mytypes.GetUser
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Username,
			&user.ProfilePhoto,
		)
		if err != nil {
			return nil, nil, err
		}
		users = append(users, &user)
	}
	//query for

	//return
	return users, &mytypes.Pagination{
		TotalPage:   totalPage,
		CurrentPage: currentPage,
		TotalCount:  totalCount,
	}, nil

}

// SEARCH POSTS
func (x *SearchDB) SearchPosts(keyword string, page int) ([]*mytypes.GetPostWithUsername, *mytypes.Pagination, error) {
	//context
	ctx, cancel := context.WithTimeout(xctx, dbTimeout)
	defer cancel()
	//query for pag
	query := `
			SELECT COUNT(*) FROM 
				posts 
			WHERE
				content LIKE 
			'%' || $1 || '%'
			`
	var totalCount int
	err := conn.QueryRow(ctx, query, keyword).Scan(&totalCount)
	if err != nil {
		return nil, nil, err
	}
	currentPage, totalPage := helpers.Pagination(page, totalCount, dataPerPage)
	//main
	query = `
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
			LEFT JOIN users ON posts.user_id = users.id
			WHERE
				content LIKE 
					'%' || $1 || '%' 
			ORDER BY posts.created_at DESC
			LIMIT $2 OFFSET $3
			`
	rows, err := conn.Query(ctx, query, keyword, dataPerPage, currentPage*dataPerPage-dataPerPage)
	if err != nil {
		return nil, nil, err
	}
	var posts []*mytypes.GetPostWithUsername
	for rows.Next() {
		post, err := postScannerUsername(rows)
		if err != nil {
			fmt.Println("query err", err)

			return nil, nil, err
		}
		posts = append(posts, post)
	}
	//query for

	//return
	return posts, &mytypes.Pagination{
		CurrentPage: currentPage,
		TotalCount:  totalCount,
		TotalPage:   totalPage,
	}, nil
}
