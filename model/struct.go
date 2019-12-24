package model

//User 定义在数据库中保存的user数据结构
type User struct {
	Username string `gorm:"username"`
	Password string `gorm:"password"`
	Email    string `gorm:"email;PRIMARY_KEY"`
}

//Blog 保存发布内容信息
type Blog struct {
	ID           string   `json:"id" gorm:"id;PRIMARY_KEY"`
	CreatorEmail string   `json:"creatoremail" gorm:"creator_email"`
	CreateTime   string   `json:"createtime" gorm:"create_time"`
	IsPublic     string   `json:"ispublic" gorm:"is_public"`
	Content      string   `json:"content" gorm:"content"`
	PictureName  string `json:"picturename" gorm:"picture_name"`
	GoodCount    int      `json:"goodcount" gorm:"good_count"`
}

//Good 保存所有点赞信息
type Good struct {
	ID            string `json:"id" gorm:"id;PRIMARY_KEY"`
	FromUserEmail string `json:"from_user_email" gorm:"from_user_email"`
	TargetBlogID  string `json:"target_blog_id" gorm:"target_blog_id"`
}

//CommentItem 保存评论
type CommentItem struct {
	ID               string `json:"id" gorm:"id;PRIMARY_KEY"`
	FromUserEmail    string `json:"from_user_email" gorm:"from_user_email"`
	TargetID         string `json:"target_blog_id" gorm:"target_id"`
	Content          string `json:"content" gorm"content"`
	TargetBlogID     string `json:"at_blog_id" gorm:"target_blog_id"`
	TargetCommentcID string `json:"target_commentc_id" gorm:"target_commentc_id"`
}

//ReturnCommentItem 保存所有子评论
type ReturnCommentItem struct {
	ID               string `json:"id"`
	FromUsername     string `json:"from_user_name"`
	TargetID         string `json:"target_blog_id"`
	Content          string `json:"content"`
	TargetBlogID     string `json:"at_blog_id"`
	TargetCommentcID string `json:"target_commentc_id"`
}

//CommentResponse 返回当前的所有评论
type CommentResponse struct {
	ID           string              `json:"id"`
	FromUsername string              `json:"from_user_name"` //这里有问题
	TargetBlogID string              `json:"target_blog_id"`
	Content      string              `json:"content" gorm"content"`
	SubComments  []ReturnCommentItem `json:"sub_comments" `
}
