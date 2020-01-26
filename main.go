package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

const STATIC_URL string = "/public/static/"
const STATIC_ROOT string = "public/static/"

type Context struct {
	Title  string
	Body   string
	Static string
}

func Home(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!", Body: "This is an HTML web page using JS on websocket and served by a GOlang server. Ciao!"}
	render(w, "index", context)
}

func Close(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Close", Body: "This page is all about something simple."}
	render(w, "close", context)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = STATIC_URL

	tmpl_list := []string{"public/templates/base.html",
		fmt.Sprintf("public/templates/%s.html", tmpl)}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}
			if string(msg) == "ping" {
				fmt.Println("client send ping")
				time.Sleep(2 * time.Second)
				err = conn.WriteMessage(msgType, []byte("Server reply pong after client send ping"))
				if err != nil {
					fmt.Println(err)
					return
				}
			} else {
				conn.Close()
				fmt.Println(string(msg))
				fmt.Println("\nShutting down the server...")
				os.Exit(3)
			}
		}
	})
	go open("http://localhost:8080/")
	http.HandleFunc("/", Home)
	http.HandleFunc("/close/", Close)
	http.HandleFunc(STATIC_URL, StaticHandler)

	http.ListenAndServe(":8080", nil)
}
