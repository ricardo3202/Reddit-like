package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"math"
	"strconv"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票对应的得分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")

	Ctx = context.Background()
)

func CreatePost(postID, communityID int64) error {

	//事件
	pipeline := client.TxPipeline()

	// 帖子时间
	pipeline.ZAdd(Ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(Ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(Ctx, cKey, postID)

	_, err := pipeline.Exec(Ctx)
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票的限制,距离发布时间超过一周就不暴露投票接口了
	// 从redis中读取帖子发布时间
	postTime := client.ZScore(Ctx, getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 步骤2与步骤3 也需要放到一个pipeline中
	// 2. 更新帖子分数
	// 先查询当前用户给当前帖子的投票记录
	oValue := client.ZScore(Ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	if value == oValue {
		return ErrVoteRepeated
	}
	var op float64
	if value > oValue {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(oValue - value) // 计算两次投票的差值

	pipeline := client.TxPipeline()

	pipeline.ZIncrBy(Ctx, getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	// 3. 记录用户为该帖子投票的数据
	if value == 0 { // 说明是取消投票
		pipeline.ZRem(Ctx, getRedisKey(KeyPostVotedZSetPF+postID), postID)
	} else {
		pipeline.ZAdd(Ctx, getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, // 当前用户投的是赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec(Ctx)

	return err
}
