package redis

import (
	"bluebell/models"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	// 2. 确定查询索引的起始点
	start := (page - 1) * size // 第一页从0开始
	end := start + size - 1
	// 3. ZRevRange 按分数以倒序即从大到小查询给出范围元素
	return client.ZRevRange(Ctx, key, start, end).Result()
}

// GetPostIDsInOrder 根据Order 查询每篇帖子的ids,排序后返回帖子id
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis中获取id
	// 1. 根据用户请求中的参数order决定要查询的redis key是time或score
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteDate 根据ids 查询每篇帖子投赞投票的数据
func GetPostVoteDate(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))

	//每查询一个帖子的点赞数量就要对redis缓存发送一次请求,这种做法不太好
	//for _, id := range ids {
	//	Key := getRedisKey(KeyPostVotedZSetPF + id)
	//	v := client.ZCount(Ctx, Key, "1", "1").Val() // 统计va
	//	// lue在min和max之间的值，这里只有1，即统计每篇帖子的赞成票的数量
	//	data = append(data, v)
	//}

	// 使用pipeline进行打包发送，即一次请求
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(Ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(Ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetPostIDsInOrderByComuniy 根据社区查询ids
func GetPostIDsInOrderByComuniy(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用zinterstore把分区的帖子set与按分数排序zset 生成一个新的zset
	// 针对新的zset按照之前的逻辑取数据

	// 社区的key
	//cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	// 利用缓存key减少zinterstore的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID)) // 生成一个新的key，即原key加上社区的id

	// 把查询到的数据返回上一层调用，但是目前这里代码存在问题
	if client.Exists(Ctx, key).Val() < 1 { // Key若不存在
		//
		pipeline := client.Pipeline()
		pipeline.ZInterStore(Ctx, key, &redis.ZStore{
			Aggregate: "MAX",
		})
		//pipeline.ZInterStore(Ctx, key, redis.ZStore{
		//	Aggregate: "MAX",
		//}, cKey, orderKey)
		pipeline.Expire(Ctx, key, 60*time.Second)
		_, err := pipeline.Exec(Ctx)
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)

}
