package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "this is api.microcodor.com app api doc!\n")
	//url := "http://microcodor.com/"
	//http.Redirect(w, r, url, http.StatusFound)
	fmt.Fprintln(w, "GET /getwpuser/{userId}:获取单个用户信息")
	fmt.Fprintln(w, "GET /getwpterms:获取所有分类信息")
	fmt.Fprintln(w, "GET /getwppost/{postId}:获取单篇文章")
	fmt.Fprintln(w, "GET /getsimpleposts/{termId}/{postId}/{num}:简版文章分类列表")
}

/*
接口一：获取APP文章分类
*/
func GetWpterms(w http.ResponseWriter, r *http.Request) {
	QueryWpterms(w)
}

func QueryWpterms(w http.ResponseWriter) {
	seelog.Error("QueryWpterms")
	wptermsbean := new(WpTermsBean)
	wptermsbean.Common.Code = 0
	wptermsbean.Common.Msg = "数据异常"

	rows, err := db.Query("SELECT * FROM wp_terms WHERE slug like 'wx_%'")
	if err != nil {
		panic(err.Error())
		seelog.Error(err.Error())
		wptermsbean.Common.Msg = err.Error()
	}
	defer rows.Close()

	for rows.Next() {
		var wpterm Wpterm

		err = rows.Scan(&wpterm.Termid, &wpterm.Name, &wpterm.Slug, &wpterm.Termgroup)
		if err != nil {
			seelog.Error(err.Error())
			wptermsbean.Common.Msg = err.Error()
		}
		wptermsbean.WpTerms = append(wptermsbean.WpTerms, wpterm)

	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if len(wptermsbean.WpTerms) > 0 {
		wptermsbean.Common.Code = 1
		wptermsbean.Common.Msg = "查询分类成功"
	}
	if err := json.NewEncoder(w).Encode(wptermsbean); err != nil {
		//panic(err)
		seelog.Error(err.Error())
	}
}

/*
接口二：获取单个用户
*/
func GetWpuser(w http.ResponseWriter, r *http.Request) {
	wpuserbean := new(WpuserBean)
	wpuserbean.Common.Code = 0
	wpuserbean.Common.Msg = "数据异常"

	vars := mux.Vars(r)
	var userId int
	var err error
	if userId, err = strconv.Atoi(vars["userId"]); err != nil {
		seelog.Error(err.Error())
	}
	//需要完善
	wpuserbean.User, err = QueryWpuser(userId)
	if err == nil {
		wpuserbean.Common.Code = 1
		wpuserbean.Common.Msg = "查询用户成功"
	} else {
		if err == sql.ErrNoRows {
			log.Print("没有结果")
			wpuserbean.Common.Msg = "未查到相关用户"
		} else {
			seelog.Error(err.Error())
			wpuserbean.Common.Msg = err.Error()
		}

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(wpuserbean); err != nil {
		seelog.Error(err.Error())
	}
}

func QueryWpuser(userId int) (Wpuser, error) {
	var wpuser Wpuser
	stmt, _ := db.Prepare("select ID, user_login, user_nicename, " +
		"user_email, user_url, user_registered, user_activation_key, user_status, display_name from wp_users " +
		"where id = ?")
	rows := stmt.QueryRow(userId)

	err := rows.Scan(&wpuser.Id, &wpuser.LoginName,
		&wpuser.NickName, &wpuser.Email, &wpuser.Url, &wpuser.RegisteredTime, &wpuser.ActivationKey, &wpuser.UserStatus, &wpuser.DisplayName)
	return wpuser, err
}

/*
接口三：获取单篇文章
*/
func GetWppost(w http.ResponseWriter, r *http.Request) {
	wppostbean := new(WppostBean)
	wppostbean.Common.Code = 0
	wppostbean.Common.Msg = "数据异常"

	vars := mux.Vars(r)
	var postId int
	var err error
	if postId, err = strconv.Atoi(vars["postId"]); err != nil {
		seelog.Error(err.Error())
	} else {
		wppostbean.Wppost, err = QuerytWppost(postId)
		if err == nil {
			wppostbean.Common.Code = 1
			wppostbean.Common.Msg = "查询单篇文章成功"
		} else {
			wppostbean.Common.Code = 0
			wppostbean.Common.Msg = err.Error()
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(wppostbean); err != nil {
		seelog.Error(err.Error())
	}

}
func QuerytWppost(postId int) (Wppost, error) {

	stmt, _ := db.Prepare("select posts.ID, termships.term_taxonomy_id, posts.post_author, posts.post_date, posts.post_title, " +
		"posts.post_content, posts.comment_status, posts.guid, posts.comment_count postmeta.meta_value as views_count from wp_posts posts " +
		"inner join wp_term_relationships termships on termships.object_id=posts.ID " +
		"inner join wp_postmeta postmeta on posts.ID = postmeta.post_id " +
		"where ID = ?")
	rows, err := stmt.Query(postId)
	if err != nil {
		seelog.Error(err.Error())
	}
	defer rows.Close()

	var wppost Wppost
	for rows.Next() {
		err = rows.Scan(&wppost.Id, &wppost.TermId, &wppost.PostAuthor,
			&wppost.PostDate, &wppost.PostTitle, &wppost.PostContent,
			&wppost.CommentStatus, &wppost.PostUrl, &wppost.CommentCount, &wppost.ViewsCount)
		if err != nil {
			seelog.Error(err.Error())
		}
	}
	//需要完善
	wppost.User, err = QueryWpuser(wppost.PostAuthor)
	if err != nil {
		seelog.Error(err.Error())

	}
	return wppost, err
}

/*
接口四：简版文章分类列表
*/
func GetSimplePosts(w http.ResponseWriter, r *http.Request) {
	simplewppostsbean := new(SimpleWppostsBean)
	simplewppostsbean.Common.Code = 0
	simplewppostsbean.Common.Msg = "数据异常"

	vars := mux.Vars(r)
	var termId int
	var postId int
	var num int
	var err error
	if termId, err = strconv.Atoi(vars["termId"]); err != nil {
		seelog.Error(err.Error())
	}
	if postId, err = strconv.Atoi(vars["postId"]); err != nil {
		seelog.Error(err.Error())
	}
	if num, err = strconv.Atoi(vars["num"]); err != nil {
		seelog.Error(err.Error())
	}
	if termId == 0 || postId == 0 || num == 0 {
		simplewppostsbean.Common.Msg = "请求参数错误"
	} else {
		simplewppostsbean.Wpposts, err = QuerySimplePosts(termId, postId, num)
		if err == nil {
			simplewppostsbean.Common.Code = 1
			simplewppostsbean.Common.Msg = "查询简版文章列表成功"
		} else {
			simplewppostsbean.Common.Code = 0
			simplewppostsbean.Common.Msg = err.Error()
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if len(simplewppostsbean.Wpposts) <= 0 {
		simplewppostsbean.Common.Code = 0
		simplewppostsbean.Common.Msg = "暂无新的文章列表"
		simplewppostsbean.Wpposts = make([]Wppost, 0)
	}
	if err := json.NewEncoder(w).Encode(simplewppostsbean); err != nil {
		seelog.Error(err.Error())
	}
}

func QuerySimplePosts(termId int, postId int, num int) ([]Wppost, error) {
	stmt, _ := db.Prepare("select * from (select posts.ID as post_id,termships.term_taxonomy_id as term_id,posts.post_author as user_id," +
		"posts.post_title as post_title,posts.guid as post_url,posts.post_date as post_date," +
		"posts.comment_count as comment_count,users.user_nicename as user_nicename postmeta.meta_value as views_count from db_wordpress.wp_posts posts " +
		"inner join db_wordpress.wp_users users on posts.post_author=users.ID " +
		"inner join db_wordpress.wp_term_relationships termships on termships.object_id=posts.ID " +
		"inner join wp_postmeta postmeta on posts.ID = postmeta.post_id " +
		"where posts.post_status='publish' and termships.term_taxonomy_id=? and posts.ID>?  LIMIT ?) as art ORDER BY art.post_id DESC")
	rows, err := stmt.Query(termId, postId, num)
	if err != nil {
		seelog.Error(err.Error())
	}
	defer rows.Close()

	var wpposts []Wppost
	for rows.Next() {
		var wppost Wppost
		err = rows.Scan(&wppost.Id, &wppost.TermId, &wppost.PostAuthor,
			&wppost.PostTitle, &wppost.PostUrl, &wppost.PostDate, &wppost.CommentCount, &wppost.User.NickName, &wppost.ViewsCount)
		if err != nil {
			seelog.Error(err.Error())
		} else {
			wppost.User.Id = wppost.PostAuthor
			wpposts = append(wpposts, wppost)
		}

	}

	return wpposts, err
}

func TodoJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}
func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todoId int
	var err error
	if todoId, err = strconv.Atoi(vars["todoId"]); err != nil {
		panic(err)
	}
	todo := RepoFindTodo(todoId)
	if todo.Id > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}' http://localhost:8080/todos

*/
func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
