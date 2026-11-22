package redis

import (
	"RedBlog/models"
	"context"

	"github.com/redis/go-redis/v9"
)

// GetPostIDsInOrder 从redis中以时间/分数顺序查询帖子的id
func GetPostIDsInOrder(ctx context.Context, p *models.ParamPostList) ([]string, error) {
	//根据用户请求中携带的order参数确定要查询的redis key
	//默认order为时间排序
	key := getRedisKey(keyPostTimeZset)
	if p.Order == KeyPostScoreZset {
		key = getRedisKey(KeyPostScoreZset)
	}
	// 2. 确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	// 3. ZREVRANGE 查询（旧）
	return client.ZRevRange(ctx, key, start, end).Result()
	// ZRangeArgs 按分数从大到小的顺序查询指定数量的额元素（新）
	//res, err := client.ZRangeArgs(ctx, redis.ZRangeArgs{
	//	Key:   key,
	//	Start: start,
	//	Stop:  end,
	//	Rev:   true, // 开启倒序，等价 ZREVRANGE
	//}).Result()

}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ctx context.Context, ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZsetPrefix + id)
	//	// 查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZsetPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
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
