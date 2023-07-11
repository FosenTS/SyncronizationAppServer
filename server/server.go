package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var Clients = make(map[*Client]bool)

func StartWebsocket() {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{
			conn: conn,
			send: make(chan []byte),
		}
		Clients[client] = true
		log.Println("Клиент подключен!")
		if err := client.conn.WriteMessage(websocket.TextMessage, []byte("Directory name: "+os.Args[1])); err != nil {
			log.Println(err)
			return
		}
	})
	http.ListenAndServe(":8080", nil)
}

func SendAllClientMessage(clients map[*Client]bool, pathFile string) {
	file, err := ioutil.ReadFile(pathFile)
	if err != nil {
		panic(err)
	}
	for c, _ := range clients {
		if err := c.conn.WriteMessage(websocket.TextMessage, []byte("FileName:"+pathFile+"FileBit:"+string(file))); err != nil {
			log.Println(err)
			delete(clients, c)
		}
	}
}
