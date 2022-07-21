package main

import (
	"complexhttp/tracer"
	"flag"
	"github.com/joho/godotenv"
	"github.com/stretchr/objx"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatarAvatar,
}

type templateHandler struct {
	once     sync.Once
	fileName string
	templ    *template.Template
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Cannot load .env file")
	}
}

func (t *templateHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	//t.once.Do(func() {
	t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	//})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	_ = t.templ.Execute(writer, data)
}

func main() {
	addr := flag.String("addr", ":9090", "Address of application")
	flag.Parse()
	InitAuthProvider()
	r := newRoom()
	r.tracer = tracer.New(os.Stdout)
	go r.run()
	http.Handle("/avatars/", http.StripPrefix("/avatars",
		http.FileServer(http.Dir("avatars"))))
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/upload", &templateHandler{fileName: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/login", &templateHandler{fileName: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logout)
	http.Handle("/room", r)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln("Listen and serve ", err)
	}
}

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	userData := objx.MustFromBase64(cookie.Value)
	userId := userData.Get("userId").MustStr()
	avatarFile, header, err := r.FormFile("avatarFile")
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(avatarFile)
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("avatars", userId+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = io.WriteString(w, "Uploaded successfully")
}
