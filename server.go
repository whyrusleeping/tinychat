package main

import (
	"net"
	"fmt"
	"encoding/gob"
)

var out chan string

func HandleConnection(con net.Conn, tosend chan string) {
	defer con.Close()
	enc := gob.NewEncoder(con)
	dec := gob.NewDecoder(con)
	go func() {
		var s string
		for {
			dec.Decode(&s)
			fmt.Println(s)
			out<-s
		}
	}()

	for {
		o := <-tosend
		err := enc.Encode(o)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	list, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}
	add := make(chan chan string)
	out = make(chan string)
	go func() {
		clis := make([]chan string, 0)
		for {
			select {
				case nc := <-add:
					fmt.Println("add client.")
					clis = append(clis, nc)
				case mes := <-out:
					for _,v := range clis {
						v<-mes
					}
			}
		}
	}()

	for {
		con, err := list.Accept()
		if err != nil {
			fmt.Println(err)
		}
		nch := make(chan string)
		add <- nch
		go HandleConnection(con, nch)
	}
}
