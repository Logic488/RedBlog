package redis

// redis key
// redis key注意使用命名空间的方式，方便业务查询和拆分
const (
	Prefix                 = "RedBlog:"
	keyPostTimeZset        = "post:time"  // zset;帖子及发帖时间
	KeyPostScoreZset       = "post:score" // zset;帖子及投票分数
	KeyPostVotedZsetPrefix = "post:voted" // zset;记录用户及投票类型;参数是post_id
	KeyCommunitySetPrefix  = "community:" // set;保存每个分区下帖子的id
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
