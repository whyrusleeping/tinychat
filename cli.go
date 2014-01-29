package main

import (
	"net"
	"encoding/gob"
	"fmt"
	"bufio"
	"os"
)


func main() {
	con, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

	scan := bufio.NewScanner(os.Stdin)
	enc := gob.NewEncoder(con)
	dec := gob.NewDecoder(con)
	go func() {
		var s string
		for {
			err := dec.Decode(&s)
			if err != nil {
				panic(err)
			}
			fmt.Println(s)
		}
	}()
	for scan.Scan() {
		err := enc.Encode(scan.Text())
		if err != nil {
			panic(err)
		}
	}
}
