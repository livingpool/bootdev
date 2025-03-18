package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// In a separate terminal, run nc -u -l 42069 to listen for UDP packets and log them to the console.

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println(err)
		}
	}
}
