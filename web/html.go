package web

import (
	"embed"
	"html/template"
	"io"

	"github.com/abayomipopoola/reddit"
)

//go:embed *
var files embed.FS

var (
	home         = parse("html/home.html")
	posts         = parse("html/posts.html")
	postCreate   = parse("html/post/create.html")
	threads       = parse("html/threads.html")
	threadList      = parse("html/thread/lists.html")
	threadCreate = parse("html/thread/create.html")
	userLogin    = parse("html/user/login.html")
	userRegister = parse("html/user/register.html")
)

type HomeParams struct {
	SessionData
	Posts []reddit.Post
}

func Home(w io.Writer, p HomeParams) error {
	return home.Execute(w, p)
}

type PostParams struct {
	SessionData
	CSRF     template.HTML
	Thread   reddit.Thread
	Post     reddit.Post
	Comments []reddit.Comment
}

func Posts(w io.Writer, p PostParams) error {
	return posts.Execute(w, p)
}

type PostCreateParams struct {
	SessionData
	CSRF   template.HTML
	Thread reddit.Thread
}

func PostCreate(w io.Writer, p PostCreateParams) error {
	return postCreate.Execute(w, p)
}

type ThreadParams struct {
	SessionData
	CSRF   template.HTML
	Thread reddit.Thread
	Posts  []reddit.Post
}

func Threads(w io.Writer, p ThreadParams) error {
	return threads.Execute(w, p)
}

type ThreadCreateParams struct {
	SessionData
	CSRF template.HTML
}

func ThreadCreate(w io.Writer, p ThreadCreateParams) error {
	return threadCreate.Execute(w, p)
}

type ThreadListParams struct {
	SessionData
	Threads []reddit.Thread
}

func ThreadList(w io.Writer, p ThreadListParams) error {
	return threadList.Execute(w, p)
}

type UserParams struct {
	SessionData
	CSRF template.HTML
}

func UserLogin(w io.Writer, p UserParams) error {
	return userLogin.Execute(w, p)
}

func UserRegister(w io.Writer, p UserParams) error {
	return userRegister.Execute(w, p)
}

func parse(file string) *template.Template {
	return template.Must(
		template.New("layout.html").ParseFS(files, "html/layout.html", file))
}
