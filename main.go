package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":9090", "server address")

	flag.Parse()

	fmt.Printf("start %s\n", *addr)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		err := serve(w, r)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
