package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type client struct {
	conn *websocket.Conn
}

func (c *client) echo() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}
		fmt.Printf("rec: %s\n", string(msg))

		c.conn.SetWriteDeadline(time.Now().Add(60 * time.Second))
		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}

		w.Write(msg)

		if err := w.Close(); err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}
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
