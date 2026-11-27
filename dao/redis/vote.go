package redis

import (
	"context"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	OneWeekInSeconds = 7 * 24 * 60 * 60
	scorePerVote     = 432
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeat     = errors.New("不允许重复投票")
)

func VoteForPost(ctx context.Context, userID, postID string, direction float64) error {
	//1.判断投票限制
	//去redis取帖子发布时间
	postTime, err := client.ZScore(ctx, getRedisKey(keyPostTimeZset), postID).Result()
	if err != nil {
		return err
	}
	if float64(time.Now().Unix())-postTime > OneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	//2.更新帖子分数
	//先查询当前用户给当前帖子之前的投票记录
	ov := client.ZScore(ctx, getRedisKey(KeyPostVotedZsetPrefix+postID), userID).Val()
	//如果当前的投票和上一次投票选项一样，就提示不允许重复投票
	if ov == direction {
		return ErrVoteRepeat
	}
	var dir float64
	if direction > ov {
		dir = 1
	} else {
		dir = -1
	}
	diff := math.Abs(ov - direction) // 计算两次投票的差值

	//第2步和第3步需要放到同一个pipeline事务中执行
	pipeline := client.TxPipeline()

	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZset), dir*diff*scorePerVote, postID)
	//3.记录用户为该帖子投票的数据
	if direction == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZsetPrefix+postID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZsetPrefix+postID), redis.Z{
			Score:  direction, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)
	return err
}

func CreatePost(ctx context.Context, postID int64, communityID int64) error {
	//帖子时间和帖子分数的更新写入操作要放到同一个事务中执行
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(ctx, getRedisKey(keyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//帖子分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//把帖子id加到社区的set中
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(communityID)))
	pipeline.SAdd(ctx, cKey, postID)
	_, err := pipeline.Exec(ctx)

	return err
}
