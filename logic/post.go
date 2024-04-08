package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
)

// CreatePost 根据用户输入的内容创建帖子内容
func CreatePost(p *models.Post) (err error) {
	// 1. 生成postID
	p.ID = snowflake.GenID()
	// 2. 保存到数据库中（新增）
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	// 3. 保存包redis缓存中
	err = redis.CreatePost(p.ID, p.CommunityID)
	if err != nil {
		return err
	}
	return
	// 3. 返回
}

// GetPostByID 通过帖子id查找单篇帖子
func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {

	// 查询并组合我们接口需要的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(post.AuthorID) failed", zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

// GetPostList 查询所有帖子即帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(post.AuthorID) failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}

// GetPostListByScore 查找帖子列表，根据时间或者分数降序查找
func GetPostListByScore(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis缓存查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 根据拿到的id去MySQL数据库查询帖子详细信息,并且，是已经排好序的顺序
	// 这里是不是重复了，redis拿出来就已经排好序了吧，查mysql出来还要再查一遍吗
	return postListInfoPackage(ids)
}

// GetPostListByCommunityID 根据社区id去查询帖子详细信息
func GetPostListByCommunityID(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis缓存查询id列表
	ids, err := redis.GetPostIDsInOrderByComuniy(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0")
		return
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids", ids))

	return postListInfoPackage(ids)
}

// GetPostListBySpecial 将两个查询接口合并
func GetPostListBySpecial(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostListByScore(p)
	} else {
		// 相比较这里多了一个社区的范畴
		data, err = GetPostListByCommunityID(p)
	}
	if err != nil {
		zap.L().Error("GetPostListBySpecial failed", zap.Error(err))
		return nil, err
	}
	return
}

// PostListInfoPackage 封装打包帖子列表信息
func postListInfoPackage(ids []string) (data []*models.ApiPostDetail, err error) {
	// 根据拿到的id去MySQL数据库查询帖子详细信息,并且，是已经排好序的顺序
	posts, err := mysql.GetPostListByIDs(ids)
	zap.L().Debug("GetCommunityPostList", zap.Any("posts", posts))
	if err != nil {
		return
	}
	//提前查好每篇帖子的投票数
	voteData, err := redis.GetPostVoteDate(ids)
	if err != nil {
		return nil, err
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))

	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(post.AuthorID) failed", zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed", zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)
	}
	return
}
