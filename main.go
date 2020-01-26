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

// StaticURL : absolute
const StaticURL string = "/public/static/"

// StaticRoot : landing
const StaticRoot string = "public/static/"

// Context : struct
type Context struct {
	Title  string
	Body   string
	Static string
}

// Home : homepage
func Home(w http.ResponseWriter, req *http.Request) {
	Context := Context{Title: "Welcome!", Body: "This is an HTML web page using JS on websocket and served by a GOlang server. Ciao!"}
	render(w, "index", Context)
}

// Close : bye
func Close(w http.ResponseWriter, req *http.Request) {
	Context := Context{Title: "Close", Body: "This page is all about something simple."}
	render(w, "close", Context)
}

func render(w http.ResponseWriter, tmpl string, Context Context) {
	Context.Static = StaticURL

	templateList := []string{"public/templates/base.html",
		fmt.Sprintf("public/templates/%s.html", tmpl)}
	t, err := template.ParseFiles(templateList...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, Context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

// StaticHandler : handle
func StaticHandler(w http.ResponseWriter, req *http.Request) {
	staticFile := req.URL.Path[len(StaticURL):]
	if len(staticFile) != 0 {
		f, err := http.Dir(StaticRoot).Open(staticFile)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, staticFile, time.Now(), content)
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
	http.HandleFunc(StaticURL, StaticHandler)

	http.ListenAndServe(":8080", nil)
}
