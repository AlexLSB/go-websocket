package main

import (
	// "io"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"html/template"

	"github.com/Joker/jade"

	"golang.org/x/net/websocket"
)

func echoHandler(ws *websocket.Conn) {

	buf, err := ioutil.ReadFile("mesg.jade")
	if err != nil {
		fmt.Printf("\nReadFile error: %v", err)
		return
	}
	jadeTpl, err := jade.Parse("jade_tp", string(buf))
	if err != nil {
		fmt.Printf("\nParse error: %v", err)
		return
	}
	goTpl, err := template.New("html").Parse(jadeTpl)
	if err != nil {
		fmt.Printf("\nTemplate parse error: %v", err)
		return
	}

	for {
		// Read
		msg := ""
		err = websocket.Message.Receive(ws, &msg)
		if err != nil {
			fmt.Printf("\nExecute error: %v", err)
			return
		}

		var tpl bytes.Buffer
		err = goTpl.Execute(&tpl, msg)
		if err != nil {
			fmt.Printf("\nExecute error: %v", err)
			return
		}

		// Write
		err := websocket.Message.Send(ws, tpl.String())
		if err != nil {
			fmt.Printf("\nExecute error: %v", err)
			return
		}
	}
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHandler))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
