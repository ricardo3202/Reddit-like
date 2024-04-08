package mysql

import (
	"bluebell/models"
	"database/sql"
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "SELECT community_id, community_name FROM community "
	if err := db.Select(&communityList, sqlStr); err != nil { // db.Select查询多行并返回
		if err == sql.ErrNoRows { //在数据库中什么都没查到
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `SELECT community_id, community_name,introduction,create_time FROM community where community_id = ?`
	if err = db.Get(community, sqlStr, id); err != nil { // db.Select查询一行并返回
		if err == sql.ErrNoRows {
			zap.L().Warn("wrong communityID")
			err = ErrorInvalidID
		}
	}
	return
}
