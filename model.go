package main

/* 通用返回数据结构 */
type CommonData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

/* 定义分类对象结构体 */
type Wpterm struct {
	Termid    int    `json:"term_id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Termgroup int    `json:"term_group"`
}

/* 定义用户信息结构体  */
/*
ID: "1",
user_login: "wxdroid",
user_pass: "$P$BpruU3CXAbCSywyT7HNKR7YVgi4GjN/",
user_nicename: "wxdroid",
user_email: "jinchun8023@163.com",
user_url: "",
user_registered: "2016-03-29 18:05:40",
user_activation_key: "",
user_status: "0",
display_name: "wxdroid"
*/
type Wpuser struct {
	Id             int    `json:"id"`
	LoginName      string `json:"user_login"`
	NickName       string `json:"user_nicename"`
	Email          string `json:"user_email"`
	Url            string `json:"user_url"`
	RegisteredTime string `json:"user_registered"`
	ActivationKey  string `json:"user_activation_key"`
	UserStatus     int    `json:"user_status"`
	DisplayName    string `json:"display_name"`
}

/*
 Name string `json:"msg_name"`       // 对应JSON的msg_name
    Body string `json:"body,omitempty"` // 如果为空置则忽略字段
    Time int64  `json:"-"`              // 直接忽略字段
*/

/* 定义文章结构体  */
type Wppost struct {
	Id            int    `json:"id"`
	TermId        int    `json:term_id`
	PostAuthor    int    `json:"post_author"`
	PostDate      string `json:"post_date"`
	PostTitle     string `json:"post_title"`
	PostContent   string `json:"post_content, omitempty"`
	CommentStatus string `json:"comment_status"`
	PostUrl       string `json:"guid"`
	CommentCount  int    `json:"comment_count"`
	ViewsCount    int    `json:"views_count"`
	User          Wpuser `json:"user, omitempty"`
}

/* 评论结构体 */
type WpComment struct {
	CommentId          int    `json:"comment_id"`
	PostId             int    `json:"post_id"`
	CommentAuthor      string `json:"comment_author"`
	CommentAuthorEmail string `json:"comment_author_email"`
	CommentAuthorIp    string `json:"comment_author_IP"`
	CommentDate        string `json:"comment_date"`
	CommentContent     string `json:"comment_content"`
	CommentAgent       string `json:"comment_agent"`
	UserId             int    `json:"user_id"`
}

/*
JSON返回封装结构体
*/
type WpTermsBean struct {
	Common  CommonData `json:"common"`
	WpTerms []Wpterm   `json:"data, omitempty"`
}
type WpuserBean struct {
	Common CommonData `json:"common"`
	User   Wpuser     `json:"data, omitempty"`
}
type WppostBean struct {
	Common CommonData `json:"common"`
	Wppost Wppost     `json:"data, omitempty"`
}
type SimpleWppostsBean struct {
	Common  CommonData `json:"common"`
	Wpposts []Wppost   `json:"data, omitempty"`
}
