package server

import (
	"encoding/json"

	"MyIOSWebServer/lib"
	"MyIOSWebServer/model"
	noti "MyIOSWebServer/notification"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/garyburd/redigo/redis"

	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//DeletePictureRequest 用户删除云空间图片请求
type DeletePictureRequest struct {
	PictureName string `json:"picture_name"`
}

//CreateCommentRequest 用户创建评论请求
type CreateCommentRequest struct {
	ID      string `json:"id"`
	Content string `json"content"`
}

//IDRequest 解析存在ID的request
type IDRequest struct {
	ID string `json:"id"`
}

//DownloadPictureRequest 处理用户的下载请求
type DownloadPictureRequest struct {
	PictureName string `json:"picture_name"`
}

//ErrorReturnType 返回错误
type ErrorReturnType struct {
	OK   bool   `json:"ok"`
	Data string `json:"data"`
}

//CreateUserRequest 用于创建网站用户请求
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

//DeleteBlogRequest 用于用户删除自己的博客
type DeleteBlogRequest struct {
	ID string `json:"id"`
}

//TokenResponse 返回token
type TokenResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

//UserLogInRequest 用于用户登录请求
type UserLogInRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

//CreateBlogRequest 用户创建博客请求
type CreateBlogRequest struct {
	IsPublic    string `json:"ispublic"`
	Content     string `json:"content"`
	PictureName string `json:"picture_name"`
}

//PublicBlogsResponse 返回用户所有可见博客
type PublicBlogsResponse struct {
	ID          string `json:"id"`
	CreateTime  string `json:"create_time"`
	Content     string `json:"content"`
	CreatorName string `json:"creator_name"`
	PictureName string `json:"picture_name"`
	GoodCount   int    `json:"good_count"`
}

//GetAllCommentByBlogID 返回一个博客所有评论
func GetAllCommentByBlogID(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	idrequest := IDRequest{}
	if err := json.Unmarshal(body, &idrequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}
	commentResponseList := []model.CommentResponse{}
	commentBIDList, err := dbServer.GetAllCommentBIDFromBlogID(idrequest.ID)
	if err == gorm.ErrRecordNotFound {
		return true, commentResponseList
	}
	if err != nil {
		return false, "查询一级评论错误"
	}
	for _, v := range commentBIDList {
		commentBBody, err := dbServer.GetCommentFromCommentID(v)
		if err != nil {
			return false, "获取一级评论体失败"
		}
		commentResponse := model.CommentResponse{
			ID:           commentBBody.ID,
			FromUsername: commentBBody.FromUserEmail,
			TargetBlogID: commentBBody.TargetBlogID,
			Content:      commentBBody.Content,
		}
		commentCBody, err := dbServer.GetAllCommentCBodyFromCommentBID(v)
		if err != nil && err != gorm.ErrRecordNotFound {
			return false, "查询二级ID错误"
		}
		returnItems := []model.ReturnCommentItem{}
		for _, v := range commentCBody {
			user, err := dbServer.GetUserFromEmail(v.FromUserEmail)
			if err != nil {
				return false, "获取二级评论信息错误"
			}
			returnItems = append(returnItems, model.ReturnCommentItem{
				ID:               v.ID,
				FromUsername:     user.Username,
				TargetID:         v.TargetID,
				Content:          v.Content,
				TargetBlogID:     v.TargetBlogID,
				TargetCommentcID: v.TargetCommentcID,
			})
		}
		commentResponse.SubComments = returnItems
		commentResponseList = append(commentResponseList, commentResponse)
	}
	return true, commentResponseList
}

//CreateCommentC 创建二级评论
func CreateCommentC(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	useremail := mux.Vars(req)["email"]
	// Check token
	ok, Tuseremail := lib.GetUserEmailFromToken(req.Header.Get("token"), lib.SignKey)
	if !ok {
		return false, "身份验证失败"
	}
	if Tuseremail != useremail {
		return false, "身份验证失败"
	}
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		log.Println("用户不存在")
		return false, "用户不存在"
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	createcommentRequest := CreateCommentRequest{}
	if err := json.Unmarshal(body, &createcommentRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}
	var targetUserEmail string
	targetBlogID, targetBlogCID, targetUserEmail, err := dbServer.GetBlogIDAndCIDAndEmailFromCommentID(createcommentRequest.ID)
	if err != nil {
		return false, "获取评论对象对应的博客ID失败"
	}
	if targetBlogCID == "self" {
		targetBlogCID = createcommentRequest.ID
	}
	if err := dbServer.NewCommentToComment(createcommentRequest.Content, createcommentRequest.ID, targetBlogID, targetBlogCID, useremail); err != nil {
		return false, "创建二级评论失败"
	}

	//发送通知
	if err := hubServer.SendMessage(targetUserEmail, []byte(useremail+" comments your comment")); err != nil {
		log.Print(err)
	}
	return true, nil
}

//CreateCommentB 创建一级评论
func CreateCommentB(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	useremail := mux.Vars(req)["email"]
	// Check token
	ok, Tuseremail := lib.GetUserEmailFromToken(req.Header.Get("token"), lib.SignKey)
	if !ok {
		return false, "身份验证失败"
	}
	if Tuseremail != useremail {
		return false, "身份验证失败"
	}
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		log.Println("用户不存在")
		return false, "用户不存在"
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	createcommentRequest := CreateCommentRequest{}
	if err := json.Unmarshal(body, &createcommentRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}

	if err := dbServer.NewCommentToBlog(createcommentRequest.Content, createcommentRequest.ID, useremail); err != nil {
		return false, "创建一级评论失败"
	}

	var targetUserEmail string
	if targetUserEmail, err = dbServer.GetBlogCreatorFromBlogID(createcommentRequest.ID); err != nil {
		return false, "获取评论对象email失败"
	}

	//发送通知
	if err := hubServer.SendMessage(targetUserEmail, []byte(useremail+" comments your blog")); err != nil {
		log.Print(err)
	}
	return true, nil
}

//GetAllBlogPublic 获取当前所有public博客
func GetAllBlogPublic(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	allBlogData, err := dbServer.GetAllPublicBlog()
	if err != nil {
		log.Println(err)
		return false, "获取所有公开博客失败"
	}
	var publicBlogResponseList []PublicBlogsResponse

	for _, v := range allBlogData {
		user, err := dbServer.GetUserFromEmail(v.CreatorEmail)
		if err != nil {
			return false, "获取所有公开博客失败"
		}
		publicBlogResponseList = append(publicBlogResponseList, PublicBlogsResponse{
			ID:          v.ID,
			CreateTime:  v.CreateTime,
			Content:     v.Content,
			CreatorName: user.Username,
			GoodCount:   v.GoodCount,
			PictureName: v.PictureName,
		})
	}
	return true, publicBlogResponseList
}

//DeleteGoodTarget 删除用户点赞
func DeleteGoodTarget(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	useremail := mux.Vars(req)["email"]
	// // Check token
	// ok, Tuseremail := lib.GetUserEmailFromToken(req.Header.Get("token"), lib.SignKey)
	// if !ok {
	// 	return false, "身份验证失败"
	// }
	// if Tuseremail != useremail {
	// 	return false, "身份验证失败"
	// }
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		log.Println("用户不存在")
		return false, "用户不存在"
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	idrequest := IDRequest{}
	if err := json.Unmarshal(body, &idrequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}
	if err := dbServer.DeleteGoodTargetFromGoodTableAndBlog(idrequest.ID); err != nil {
		log.Println(err)
		return false, "数据库中删除数据失败"
	}
	return true, nil
}

//SetGoodTarget 用户点赞=》更新good库和blog库
func SetGoodTarget(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	useremail := mux.Vars(req)["email"]
	// Check token
	ok, Tuseremail := lib.GetUserEmailFromToken(req.Header.Get("token"), lib.SignKey)
	if !ok {
		return false, "身份验证失败"
	}
	if Tuseremail != useremail {
		return false, "身份验证失败"
	}
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		log.Println("用户不存在")
		return false, "用户不存在"
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	idrequest := IDRequest{}
	var targetCreatorEmail string
	if err := json.Unmarshal(body, &idrequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}
	if err, targetCreatorEmail = dbServer.SetGoodTargetToGoodTableAndBlog(useremail, idrequest.ID); err != nil {
		log.Println(err)
		return false, "数据库中添加数据失败"
	}
	//发送通知
	if err := hubServer.SendMessage(targetCreatorEmail, []byte(useremail+" likes your blog")); err != nil {
		log.Println(err)
	}
	return true, nil
}

//GetAllGoodTargetByEmail 返回当前用户所有点赞内容
func GetAllGoodTargetByEmail(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	useremail := mux.Vars(req)["email"]
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		log.Println("用户不存在")
		return false, "用户不存在"
	}
	var goodList []model.Good
	var err error
	if goodList, err = dbServer.GetAllGoodTargetID(useremail); !(err == gorm.ErrRecordNotFound || err == nil) {
		log.Panicln(err)
		return false, "获取所有点赞id失败"
	}
	return true, goodList
}

//CreateBlogHandler 提供创建博客服务
func CreateBlogHandler(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	vars := mux.Vars(req)
	useremail := vars["email"]
	// Check token
	ok, Tuseremail := lib.GetUserEmailFromToken(req.Header.Get("token"), lib.SignKey)
	if !ok {
		return false, "身份验证失败"
	}
	if Tuseremail != useremail {
		return false, "身份验证失败"
	}

	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		return false, "用户不存在"
	}

	crateBlogRequest := CreateBlogRequest{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	if err := json.Unmarshal(body, &crateBlogRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}

	if crateBlogRequest.Content == "" {
		log.Print("博客内容不能为空")
		return false, "博客内容不能为空"
	}
	blog := model.Blog{
		CreatorEmail: useremail,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		ID:           lib.GetUniqueID(),
		IsPublic:     crateBlogRequest.IsPublic,
		Content:      crateBlogRequest.Content,
		PictureName:  crateBlogRequest.PictureName,
		GoodCount:    0,
	}

	if ok, err := dbServer.UserCreateBlog(blog); ok != true || err != nil {
		log.Println(err)
		return false, "创建博客失败"
	}
	return true, ""
}

//CreateUserHandler 提供创建用户服务
func CreateUserHandler(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	createUserRequest := CreateUserRequest{}
	if err := json.Unmarshal(body, &createUserRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}

	// check if given information is valid
	if ok := lib.CheckEmail(createUserRequest.Email); !ok {
		return false, "无效的邮箱地址"
	}

	user := model.User{
		Password: createUserRequest.Password,
		Email:    createUserRequest.Email,
		Username: createUserRequest.Username,
	}
	if ok, err := dbServer.AddNewSignUpUser(user); ok != true || err != nil {
		log.Println(err)
		return false, "该邮箱已存在"
	}
	return true, ""
}

//UserLoginHandler 提供用户登录服务
func UserLoginHandler(w http.ResponseWriter, req *http.Request) (bool, interface{}) {
	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Println(err)
		return false, "无法读取用户的请求"
	}
	userLogInRequest := UserLogInRequest{}
	if err := json.Unmarshal(body, &userLogInRequest); err != nil {
		log.Println(err)
		return false, "无效的json信息"
	}

	// check if given information is valid
	if ok := lib.CheckEmail(userLogInRequest.Email); !ok {
		return false, "无效的邮箱地址"
	}

	user, err := dbServer.GetUserFromEmail(userLogInRequest.Email)
	if err == gorm.ErrRecordNotFound {
		return false, "该用户不存在"
	}
	if err != nil {
		return false, "服务端查询错误"
	}
	if user.Password != userLogInRequest.Password {
		return false, "密码错误"
	}

	//Token 生命周期起始点
	token := req.Header.Get("token")
	if ok, _ := lib.CheckToken(token); ok != true {
		//需要重建token
		newToken, err := lib.GenerateToken(userLogInRequest.Email)
		if err != nil {
			return false, "生成Token失败"
		}
		return true, TokenResponse{
			Token:    newToken,
			Username: user.Username,
		}

	}
	//不需要重建Token 使用旧token
	return true, TokenResponse{
		Token:    token,
		Username: user.Username,
	}
}

//DeleteMyPicture 删除用户图片
func DeleteMyPicture(w http.ResponseWriter, r *http.Request) (bool, interface{}) {
	useremail := mux.Vars(r)["email"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, "请求格式错误"
	}
	//Check token
	ok, Tuseremail := lib.GetUserEmailFromToken(r.Header.Get("token"), lib.SignKey)
	if !ok {
		return false, "身份验证失败"
	}
	if Tuseremail != useremail {
		return false, "身份验证失败"
	}
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		return false, "用户不存在"
	}

	deletePictureRequest := DeletePictureRequest{}
	if err := json.Unmarshal(body, &deletePictureRequest); err != nil {
		log.Println("json格式错误")
		return false, "json格式错误"
	}
	_, err = conn.Do("lrem", useremail, "0", deletePictureRequest.PictureName)
	if err != nil {
		return false, "redis系统错误"
	}
	return true, ""
}

//GetAllMyPicture 获取用户云存储内所有图片名称
func GetAllMyPicture(w http.ResponseWriter, r *http.Request) (bool, interface{}) {
	useremail := mux.Vars(r)["email"]
	//Check token
	ok, Tuseremail := lib.GetUserEmailFromToken(r.Header.Get("token"), lib.SignKey)
	if !ok {
		return false, "身份验证失败"
	}
	if Tuseremail != useremail {
		return false, "身份验证失败"
	}
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		return false, "用户不存在"
	}
	values, err := redis.Values(conn.Do("lrange", useremail, "0", "-1"))
	if err != nil {
		return false, "redis系统错误"
	}
	var allPictureNames []string
	for _, v := range values {
		allPictureNames = append(allPictureNames, string(v.([]byte)))
	}
	return true, allPictureNames
}

//UploadPictureHandler 处理上传文件请求
func UploadPictureHandler(w http.ResponseWriter, r *http.Request) (bool, interface{}) {
	useremail := mux.Vars(r)["email"]
	//Check token
	ok, Tuseremail := lib.GetUserEmailFromToken(r.Header.Get("token"), lib.SignKey)
	if !ok {
		return false, "身份验证失败"
	}
	if Tuseremail != useremail {
		return false, "身份验证失败"
	}
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		log.Println("用户不存在")
		return false, "用户不存在"
	}
	file, header, err := r.FormFile("picture")
	if err != nil {
		log.Println(err)
		return false, "invalid input"
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return false, "invalid input"
	}
	// get file extension
	var picExt string
	if !strings.Contains(header.Filename, ".") {
		picExt = ""
	} else {
		name := strings.Split(header.Filename, ".")
		picExt = "." + name[len(name)-1]
	}
	// save file
	filename := lib.GetMD5(string(content)) + picExt
	path := filepath.Join("./src", filename)
	if err = ioutil.WriteFile(path, content, 0644); err != nil {
		log.Println(err)
		return false, "服务错误"
	}
	conn.Do("lpush", useremail, filename)
	return true, filename
}

//DownloadPictureHandler 处理下载图片请求
func DownloadPictureHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pictureName := r.FormValue("picturename")
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("picture", pictureName)
	if err != nil {
		w.WriteHeader(500)
		log.Println("error in creating form file")
		return
	}

	path := filepath.Join("./src", pictureName)
	fileHandler, err := os.Open(path)
	if err != nil {
		log.Println("error in opening file" + path)
		w.WriteHeader(500)
		return
	}
	defer fileHandler.Close()

	if _, err = io.Copy(fileWriter, fileHandler); err != nil {
		w.WriteHeader(500)
		log.Println("error in copying file")
		return
	}

	w.WriteHeader(200)
	w.Write(bodyBuf.Bytes())
}

//ServeWebSocket 提供通知socket连接服务
func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	useremail := r.FormValue("email")
	if _, err := dbServer.GetUserFromEmail(useremail); err == gorm.ErrRecordNotFound {
		log.Println("用户不存在")
		w.WriteHeader(404)
		return
	}
	//创建一个用户
	client, err := noti.NewClientInstance(w, r, useremail, hubServer)
	if err != nil {
		log.Println("创建用户失败")
		w.WriteHeader(500)
		return
	}
	//将用户加入hub中
	hubServer.LogInChan <- client
	//开启websocket传递
	go client.SendNoti()
}

//接下来要做的，使用hubServer给特定用户发送通知。
