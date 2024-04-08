package redis

// redis key

// redis key注意使用命名空间的方式，方便查询和拆分

const (
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"  // zset;时间作为排序规则  key是postid,值是时间
	KeyPostScoreZSet   = "post:score" // zset;得分作为排序规则  key是postid,值是分数(通过计算)
	KeyPostVotedZSetPF = "post:voted" // zset;记录用户给对应帖子投票的类型是（-1,0,1）
	
	KeyCommunitySetPF = "community:" // set; 保存每个分区下的帖子id
)

// 给redis key 加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
