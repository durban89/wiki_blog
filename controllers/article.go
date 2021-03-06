package controllers

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/durban89/wiki/helpers"
	"github.com/durban89/wiki/models"
)

// ArticleSave 存储文章
func ArticleSave(w http.ResponseWriter, r *http.Request, title string) {
	r.ParseForm()

	body := r.FormValue("body")
	author := r.FormValue("author")
	category := r.FormValue("category")
	ctime := r.FormValue("ctime")
	openStatus := r.FormValue("openStatus")
	token := r.FormValue("token")
	channel := r.Form["channel"]
	method := r.Method

	fmt.Println("author ", author)
	fmt.Println("category ", category)
	fmt.Println("ctime ", ctime)
	fmt.Println("method: ", method)
	fmt.Println("openStatus: ", openStatus)
	fmt.Println("channels: ", channel)
	fmt.Println("token: ", token)

	if len(r.Form.Get("author")) == 0 {
		fmt.Println("author is empty")
	}

	slice := []string{"php", "java", "golang"}
	if !helpers.ValidateInArray(slice, category) {
		fmt.Println("category not in slice")
	}

	chennelSlice := []string{"news", "technology", "other"}
	a := helpers.ValidateSliceIntersection(channel, chennelSlice)
	if len(a) == 0 {
		fmt.Println("channel is empty")
	}

	// XSS Example
	// fmt.Println("author: ", template.HTMLEscapeString(r.Form.Get("author"))) // print at server side
	// fmt.Println("author: ", template.HTMLEscapeString(r.Form.Get("author")))
	// template.HTMLEscape(w, []byte(r.Form.Get("author"))) // responded to clients
	// return

	// 重复提交 Example
	// if token != "" {
	// 	// check token validity
	// 	fmt.Println("To Validate Token")
	// } else {
	// 	// give error if no token
	// 	fmt.Println("Token is empty")
	// }

	p := &helpers.Page{
		Title: title,
		Body:  []byte(body),
	}

	// 开放状态验证
	openStatusSlice := []string{"1", "2"}

	if !helpers.ValidateInArray(openStatusSlice, openStatus) {
		fmt.Println("openStatus not in openStatusSlice")
	}

	err := p.Save()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// ArticleEdit 编辑文章
func ArticleEdit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	where := models.WhereValues{
		"autokid": models.WhereCondition{
			Operator: "=",
			Value:    id,
		},
	}

	update := map[string]string{}

	var blogModel models.BlogModel

	if strings.ToLower(r.Method) == "get" {
		var autokid int64
		var title string
		selectField := models.SelectValues{
			"autokid": &autokid,
			"title":   &title,
		}

		err := blogModel.QueryOne(selectField, where)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		p := helpers.Page{
			Title: title,
			ID:    autokid,
		}

		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		p.Token = token
		helpers.RenderTemplate(w, "edit", &p)
	} else if strings.ToLower(r.Method) == "post" {
		title := r.FormValue("title")
		if title == "" {
			http.Redirect(w, r, fmt.Sprintf("/edit?id=%s", id), http.StatusFound)
			return
		}

		update["title"] = title

		var blogInstance models.BlogModel
		_, err := blogInstance.Update(update, where)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/edit?id=%s", id), http.StatusFound)
		return
	}

}
