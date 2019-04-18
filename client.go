package main

import (
	speech "cloud.google.com/go/speech/apiv1"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	"io"
	"net/http"
)

type client struct {
	conn *websocket.Conn
	sp   *speech.Client
}

func (c *client) echo() {
	ctx := context.Background()

	stream, err := c.sp.StreamingRecognize(ctx)
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}

	fmt.Printf("stream %v\n", stream)

	defer func() {
		c.conn.Close()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				fmt.Println("eof")
				break
			}

			if err != nil {
				fmt.Printf("err: %v\n", err)
				break
			}

			if err := res.Error; err != nil {
				// Workaround while the API doesn't give a more informative error.
				if err.Code == 3 || err.Code == 11 {
					fmt.Println("WARNING: Speech recognition request exceeded limit of 60 seconds.")
				}
				fmt.Printf("Could not recognize: %v\n", err)
				break
			}

			for _, result := range res.Results {
				fmt.Printf("Result: %+v\n", result)
			}
		}
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("err: %v\n", err)
			break
		}
		// fmt.Printf("rec: %d bytes\n", len(msg))

		if err := stream.Send(&speechpb.StreamingRecognizeRequest{
			StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
				AudioContent: msg,
			},
		}); err != nil {
			fmt.Printf("Could not send audio: %v\n", err)
		} else {
			fmt.Println("ok")
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

	ctx := context.Background()

	sp, err := speech.NewClient(ctx)
	if err != nil {
		return err
	}

	cl := &client{
		conn: conn,
		sp:   sp,
	}

	go cl.echo()

	return nil
}
