package router

import (
	"encoding/json"
	"net/http"

	"MyIOSWebServer/server"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

//ResponseHandler 用于返回处理函数接口
type ResponseHandler func(w http.ResponseWriter, r *http.Request) (bool, interface{})

func (h ResponseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		OK   bool        `json:"ok"`
		Data interface{} `json:"data"`
	}
	ok, ret := h(w, r)
	res := Response{ok, ret}
	byteRes, err := json.Marshal(&res)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(byteRes)
}

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.Handle("/signup", ResponseHandler(server.CreateUserHandler)).Methods("POST")
	mx.Handle("/login", ResponseHandler(server.UserLoginHandler)).Methods("POST")
	mx.Handle("/{email}/createblog", ResponseHandler(server.CreateBlogHandler)).Methods("POST")
	mx.Handle("/blogground", ResponseHandler(server.GetAllBlogPublic)).Methods("GET")
	mx.Handle("/{email}/goodtargetid", ResponseHandler(server.GetAllGoodTargetByEmail)).Methods("GET")
	mx.Handle("/{email}/setDisliked", ResponseHandler(server.DeleteGoodTarget)).Methods("POST")
	mx.Handle("/{email}/setLiked", ResponseHandler(server.SetGoodTarget)).Methods("POST")
	mx.Handle("/{email}/commentb", ResponseHandler(server.CreateCommentB)).Methods("POST")
	mx.Handle("/{email}/commentc", ResponseHandler(server.CreateCommentC)).Methods("POST")
	mx.Handle("/getcomment", ResponseHandler(server.GetAllCommentByBlogID)).Methods("GET")
	mx.Handle("/{email}/uploadPicture", ResponseHandler(server.UploadPictureHandler)).Methods("POST")
	mx.Handle("/{email}/images", ResponseHandler(server.GetAllMyPicture)).Methods("GET")
	mx.Handle("/{email}/images", ResponseHandler(server.DeleteMyPicture)).Methods("DELETE")
	mx.HandleFunc("/blogground/download", server.DownloadPictureHandler).Methods("GET")
}
