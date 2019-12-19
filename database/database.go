package database

import (
	"MyIOSWebServer/lib"
	"MyIOSWebServer/model"
	"fmt"
)

//DBServiceInterface 定义数据库操作接口
type DBServiceInterface interface {
	GetUserFromEmail(email string) (model.User, error)
	AddNewSignUpUser(user model.User) (bool, error)
	UserCreateBlog(blog model.Blog) (bool, error)
	GetAllPublicBlog() ([]model.Blog, error)
	GetAllGoodTargetID(email string) ([]model.Good, error)
	SetGoodTargetToGoodTableAndBlog(email string, id string) (error, string)
	DeleteGoodTargetFromGoodTableAndBlog(id string) error
	NewCommentToComment(content string, targetCommentID string, targetBlogID, targetCommentCID string, useremail string) error
	NewCommentToBlog(content string, targetID string, useremail string) error
	GetCommentFromCommentID(commentid string) (model.CommentItem, error)
	GetAllCommentCBodyFromCommentBID(commentBid string) ([]model.CommentItem, error)
	GetAllCommentBIDFromBlogID(blogid string) ([]string, error)
	GetBlogIDAndCIDAndEmailFromCommentID(commentid string) (string, string, string, error)
	GetBlogCreatorFromBlogID(blogID string) (string, error)
}

//DBService 接口的实现
type DBService struct{}

//GetBlogCreatorFromBlogID 从blogID获取blogcreator的email
func (dbservice *DBService) GetBlogCreatorFromBlogID(BlogID string) (string, error) {
	tempBlog := model.Blog{
		ID: BlogID,
	}
	if err := db.Table("blog").First(&tempBlog).Error; err != nil {
		return "", err
	}
	return tempBlog.CreatorEmail, nil
}

//GetCommentFromCommentID 通过ID获取评论体
func (dbservice *DBService) GetCommentFromCommentID(commentid string) (model.CommentItem, error) {
	needComment := model.CommentItem{}
	if err := db.Table("comment").Where("id = ?", commentid).First(&needComment).Error; err != nil {
		return needComment, err
	}
	return needComment, nil
}

//GetAllCommentCBodyFromCommentBID 通过一级博客ID获得所有二级博客评论体
func (dbservice *DBService) GetAllCommentCBodyFromCommentBID(commentBid string) ([]model.CommentItem, error) {
	allCommentC := []model.CommentItem{}
	if err := db.Table("comment").Where("target_id = ? or target_commentc_id = ?", commentBid, commentBid).Find(&allCommentC).Error; err != nil {
		return allCommentC, err
	}

	return allCommentC, nil
}

//GetAllCommentBIDFromBlogID 通过博客ID获取所有一级评论ID
func (dbservice *DBService) GetAllCommentBIDFromBlogID(blogid string) ([]string, error) {
	allCommentB := []model.CommentItem{}
	if err := db.Table("comment").Where("target_id = ?", blogid).Find(&allCommentB).Error; err != nil {
		return []string{}, err
	}
	var idList []string
	for _, v := range allCommentB {
		idList = append(idList, v.ID)
	}
	return idList, nil
}

//GetBlogIDAndCIDAndEmailFromCommentID 从评论id获取博客id和commentcid和被评论对象的EMail
func (dbservice *DBService) GetBlogIDAndCIDAndEmailFromCommentID(commentid string) (string, string, string, error) {
	tempComment := model.CommentItem{
		ID: commentid,
	}
	if err := db.Table("comment").First(&tempComment).Error; err != nil {
		return "", "", "", err
	}
	return tempComment.TargetBlogID, tempComment.TargetCommentcID, tempComment.FromUserEmail, nil
}

//NewCommentToComment 为用户新建二级评论
func (dbservice *DBService) NewCommentToComment(content string, targetCommentID string, targetBlogID, targetCommentCID string, useremail string) error {
	newCommentID := lib.GetUniqueID()
	newComment := model.CommentItem{
		ID:               newCommentID,
		FromUserEmail:    useremail,
		TargetID:         targetCommentID,
		TargetBlogID:     targetBlogID,
		Content:          content,
		TargetCommentcID: targetCommentCID,
	}
	if err := db.Table("comment").Create(&newComment).Error; err != nil {
		return err
	}
	return nil
}

//NewCommentToBlog 为用户新建一级评论
func (dbservice *DBService) NewCommentToBlog(content string, targetID string, useremail string) error {
	newCommentID := lib.GetUniqueID()
	newComment := model.CommentItem{
		ID:               newCommentID,
		FromUserEmail:    useremail,
		TargetID:         targetID,
		TargetBlogID:     targetID,
		Content:          content,
		TargetCommentcID: "self",
	}
	if err := db.Table("comment").Create(&newComment).Error; err != nil {
		return err
	}
	return nil
}

//SetGoodTargetToGoodTableAndBlog 将点赞情况插入good数据库
func (dbservice *DBService) SetGoodTargetToGoodTableAndBlog(email string, id string) (error, string) {
	newGoodID := lib.GetUniqueID()
	tx := db.Begin()
	defer tx.Commit()
	newgood := model.Good{
		FromUserEmail: email,
		TargetBlogID:  id,
		ID:            newGoodID,
	}
	if err := tx.Table("good").Create(&newgood).Error; err != nil {
		tx.Rollback()
		return err, ""
	}
	targetblog := model.Blog{
		ID: id,
	}
	if err := tx.Table("blog").Find(&targetblog).Error; err != nil {
		tx.Rollback()
		return err, ""
	}
	targetblog.GoodCount++
	if err := tx.Table("blog").Save(&targetblog).Error; err != nil {
		tx.Rollback()
		return err, ""
	}

	return nil, targetblog.CreatorEmail
}

//DeleteGoodTargetFromGoodTableAndBlog 将点赞情况从good数据库删除
func (dbservice *DBService) DeleteGoodTargetFromGoodTableAndBlog(id string) error {
	targetgood := model.Good{
		ID: id,
	}
	tx := db.Begin()
	defer tx.Commit()
	if err := tx.Table("good").Find(&targetgood).Error; err != nil {
		tx.Rollback()
		return err
	}

	//从blog表减一
	targetblog := model.Blog{
		ID: targetgood.TargetBlogID,
	}
	if err := tx.Table("blog").Find(&targetblog).Error; err != nil {
		tx.Rollback()
		return err
	}
	targetblog.GoodCount--
	if err := tx.Table("blog").Save(&targetblog).Error; err != nil {
		tx.Rollback()
		return err
	}

	//从good表彻底删除
	if err := tx.Table("good").Delete(&targetgood).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

//GetAllGoodTargetID 获取当前用户所有点赞信息
func (dbservice *DBService) GetAllGoodTargetID(email string) ([]model.Good, error) {
	targetgoodlist := []model.Good{}
	if err := db.Table("good").Where("from_user_email = ?", email).Find(&targetgoodlist).Error; err != nil {
		return targetgoodlist, err
	}
	return targetgoodlist, nil
}

//GetAllPublicBlog 获取所有用户的public博客
func (dbservice *DBService) GetAllPublicBlog() ([]model.Blog, error) {
	var bloglist []model.Blog
	if err := db.Table("blog").Order("create_time").Where("is_public = 1").Find(&bloglist).Error; err != nil {
		return bloglist, err
	}
	for i, j := 0, len(bloglist)-1; i < j; i, j = i+1, j-1 {
		bloglist[i], bloglist[j] = bloglist[j], bloglist[i]
	}
	return bloglist, nil
}

//UserCreateBlog 用户新建博客
func (dbservice *DBService) UserCreateBlog(blog model.Blog) (bool, error) {
	fmt.Println(blog)
	if err := db.Table("blog").Create(&blog).Error; err != nil {
		return false, err
	}
	return true, nil
}

//GetUserFromEmail 根据邮箱获取用户信息
func (dbservice *DBService) GetUserFromEmail(email string) (model.User, error) {
	user := model.User{
		Email: email,
	}
	if err := db.Table("user").First(&user).Error; err != nil {
		fmt.Println(err)
		return user, err
	}
	return user, nil
}

//AddNewSignUpUser 添加博客网站注册用户
func (dbservice *DBService) AddNewSignUpUser(user model.User) (bool, error) {
	if err := db.Table("user").Create(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}

// //DeleteBlogByID 删除用户自己的博客
// func (dbservice *DBService) DeleteBlogByID(id string, useremail string) error {
// 	//删除博客db内部的信息
// 	db, err := gorm.Open("")
// 	if err != nil {
// 		fmt.Println("open failed")
// 		return err
// 	} else {
// 		fmt.Println("open succeed!")
// 	}
// 	defer db.Close()
// 	var blogTitle string
// 	if err := db.Update(func(tx *bolt.Tx) error {
// 		blog := tx.Bucket([]byte(id))
// 		if blog == nil {
// 			return errors.New("不存在当前id")
// 		}
// 		if string(blog.Get([]byte("creatoremail"))) != useremail {
// 			return errors.New("用户不能删除别人的博客")
// 		}
// 		blogTitle = string(blog.Get([]byte("title")))
// 		if err := tx.DeleteBucket([]byte(id)); err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return err
// 	}

// 	//删除用户db内部的信息
// 	dbUser, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
// 	if err != nil {
// 		fmt.Println("open failed")
// 		return err
// 	} else {
// 		fmt.Println("open succeed!")
// 	}
// 	defer dbUser.Close()
// 	if err := dbUser.Update(func(tx *bolt.Tx) error {
// 		blog := tx.Bucket([]byte(useremail))
// 		if blog == nil {
// 			return errors.New("不存在当前用户")
// 		}
// 		blogListBucket := blog.Bucket([]byte("bloglist"))
// 		if string(blogListBucket.Get([]byte(blogTitle))) != id {
// 			return errors.New("系统错误：当前用户不存在此博客")
// 		}
// 		if err := blogListBucket.Delete([]byte(blogTitle)); err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return err
// 	}
// 	return nil
// }

// //GetAllBlogData 从用户名ID的list获取所有博客内容
// func (dbservice *DBService) GetAllBlogData(idlist []string) (bool, []model.Blog) {
// 	db, err := bolt.Open("blog.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
// 	var bloglist []model.Blog
// 	if err != nil {
// 		fmt.Println("open failed")
// 		return false, bloglist
// 	} else {
// 		fmt.Println("open succeed!")
// 	}

// 	defer db.Close()
// 	if err := db.View(func(tx *bolt.Tx) error {
// 		for _, v := range idlist {
// 			if v == "" {
// 				continue
// 			}
// 			blog := tx.Bucket([]byte(v))
// 			if blog == nil {
// 				return errors.New("不存在当前id")
// 			}
// 			tempBlogData := model.Blog{
// 				CreatorEmail: string(blog.Get([]byte("creatoremail"))),
// 				Title:        string(blog.Get([]byte("title"))),
// 				CreateTime:   string(blog.Get([]byte("createtime"))),
// 				Tag:          string(blog.Get([]byte("tag"))),
// 				ID:           string(blog.Get([]byte("ID"))),
// 				IsPublic:     string(blog.Get([]byte("ispublic"))),
// 				Content:      string(blog.Get([]byte("content"))),
// 			}
// 			bloglist = append(bloglist, tempBlogData)
// 		}
// 		return nil
// 	}); err != nil {
// 		return false, bloglist
// 	}
// 	return true, bloglist
// }

// //GetAllBlogIDFromUserEmail 为用户增加从文件名到ID的映射
// func (dbservice *DBService) GetAllBlogIDFromUserEmail(email string) (bool, []string) {
// 	//查找对应用户
// 	db, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
// 	if err != nil {
// 		fmt.Println("open failed", err)
// 		return false, nil
// 	} else {
// 		fmt.Println("open succeed!")
// 	}
// 	var resultList []string
// 	defer db.Close()
// 	if err := db.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket([]byte(email))
// 		if b == nil {
// 			return errors.New("不存在当前用户")
// 		}
// 		bList := b.Bucket([]byte("bloglist"))
// 		bList.ForEach(func(k, v []byte) error {
// 			//K 、 V 为从博客名 -> 博客ID的映射
// 			resultList = append(resultList, string(v))
// 			return nil
// 		})
// 		return nil
// 	}); err != nil {
// 		return false, resultList
// 	}
// 	return true, resultList
// }

// //AddNewBlogToUser 为用户增加从文件名到ID的映射
// func (dbservice *DBService) AddNewBlogToUser(blog model.Blog) bool {
// 	//查找对应用户
// 	db, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
// 	if err != nil {
// 		fmt.Println("open failed", err)
// 		return false
// 	} else {
// 		fmt.Println("open succeed!")
// 	}
// 	defer db.Close()
// 	if err := db.Update(func(tx *bolt.Tx) error {
// 		blogerOwner := tx.Bucket([]byte(blog.CreatorEmail))
// 		if blogerOwner == nil {
// 			return errors.New("不存在此创建者")
// 		}
// 		blogList := blogerOwner.Bucket([]byte("bloglist"))
// 		if blogList == nil {
// 			return errors.New("该用户不存在博客列表")
// 		}
// 		if err := blogList.Put([]byte(blog.Title), []byte(blog.ID)); err != nil {
// 			return err
// 		}
// 		return nil
// 	}); err != nil {
// 		return false
// 	}
// 	return true
// }

// //GetUserAllBlogName 从
// func (dbservice *DBService) GetUserAllBlogName(user model.User) []string {
// 	var blogList []string
// 	return blogList
// }

// func worker(user model.User) {
// 	db, err := bolt.Open("kes.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
// 	if err != nil {
// 		fmt.Println("open failed")
// 		return
// 	} else {
// 		fmt.Println("open succeed!")
// 	}
// 	time.Sleep(1 * time.Second)
// 	db.Close()
// 	return
// }
