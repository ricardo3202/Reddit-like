package mysql

import (
	"bluebell/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `INSERT INTO post (post_id, title, content, author_id, community_id) values (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

// GetPostByID 根据id查询单个帖子详情数据
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time FROM post WHERE post_id = ?`
	if err := db.Get(post, sqlStr, pid); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("wrong postID")
			err = ErrorInvalidID
		}
	}
	return
}

// GetPostList 根据给出的索引查询帖子列表函数
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	posts = make([]*models.Post, 0, 2)
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time FROM post ORDER BY create_time  DESC LIMIT ?,?`
	if err := db.Select(&posts, sqlStr, (page-1)*size, size); err != nil { // db.Select查询多行并返回
		if err == sql.ErrNoRows { //在数据库中什么都没查到
			zap.L().Warn("there is no post in db")
			err = nil
		}
	}
	return
}

// GetPostListByIDs 根据给定的帖子id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	postList = make([]*models.Post, len(ids))
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time FROM post WHERE post_id IN (?) 
                                                                               order by FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
