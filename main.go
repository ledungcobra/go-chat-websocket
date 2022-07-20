package main

import (
	"complexhttp/tracer"
	"flag"
	"github.com/joho/godotenv"
	"github.com/stretchr/objx"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

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
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.fileName)))
	})
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
	http.Handle("/chat", MustAuth(&templateHandler{fileName: "chat.html"}))
	http.Handle("/login", &templateHandler{fileName: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logout)
	http.Handle("/room", r)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln("Listen and serve ", err)
	}
}
