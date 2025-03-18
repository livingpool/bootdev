package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("connection accepted:", conn.RemoteAddr())

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}

		fmt.Println("connection and channel closed.")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	resCh := make(chan string)

	go func() {
		defer func() {
			f.Close()
			close(resCh)
		}()

		line := ""
		for {
			data := make([]byte, 8)

			n, err := f.Read(data)
			if err != nil && err != io.EOF {
				panic(err)
			}

			if n == 0 {
				if len(line) > 0 {
					resCh <- line
				}
				break
			}

			dataStr := string(data)
			if strings.Contains(dataStr, "\n") {
				tokens := strings.Split(dataStr, "\n")
				line += tokens[0]
				resCh <- line
				line = tokens[1]
			} else {
				line += dataStr
			}
		}
	}()

	return resCh
}
