package router

import (
	"net/http"

	"wiki/controllers/article"
	"wiki/controllers/auth"
	"wiki/controllers/welcome"
)

// Routes 操作
//
// TODO
// 问题1： 不能使用同一个路由地址，会出现冲突，建议定义自己的路由方式
//
func Routes() {
	// 添加路由配置
	// Auth login/register
	GET("/auth/login/", auth.Login)
	POST("/auth/login/submit", auth.Login)

	// 文章
	GET("/articles/create/", article.Create)
	GET("/articles/update/", article.Update)
	GET("/articles/delete/", article.Delete)
	GET("/articles/view/", article.View)
	GET("/articles/error", article.Error)
	GET("/articles/", article.Index)
	POST("/articles/save/", article.Save)

	// XML文件操作
	GET("/process/xml/", welcome.ProcessXML)
	GET("/xml/", welcome.XML)

	// 其他
	// GET("/json/to/test/", controllers.JsonToTest)
	// GET("/json/process/", controllers.JsonProcess)
	// GET("/json/to/interface/", controllers.JsonToInterface)
	// GET("/json/", controllers.Json)
	// GET("/upload/", controllers.UploadHandler)
	// GET("/postFile/", controllers.PostFileHandler)

	// 首页
	GET("/", welcome.Index)

	// 静态文件路由
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.HandleFunc("/favicon.ico/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})
}
