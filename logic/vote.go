package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"go.uber.org/zap"
	"strconv"
)

// 投票功能：
// 1. 用户投票  数据校验  谁给谁

//这里用了简化版的投票分数 即投一票加分432分 86400秒/200  需要两百张赞成票可以给帖子续一天,就是让时间戳增大了一天

// VoteForPost 为帖子投票
func VoteForPost(userID int64, p *models.ParamVoteData) error {
	zap.L().Debug("VoteForPost", zap.Int64("userId", userID),
		zap.String("postID", p.PostID), zap.Int8("direct", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))

}
