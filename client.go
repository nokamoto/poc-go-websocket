package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type client struct {
	conn *websocket.Conn
}

func (c *client) echo() {
	defer func() {
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}
		fmt.Printf("rec: %v\n", msg)
	}
}

func serve(w http.ResponseWriter, r *http.Request) error {
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	cl := &client{
		conn: conn,
	}

	go cl.echo()

	return nil
}
