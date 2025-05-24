package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":5050")

	if err != nil {
		fmt.Println(err)
		listener.Close()
	}

	for {
		conn, errLis := listener.Accept()

		if errLis != nil {
			fmt.Println(errLis)
			conn.Close()
			listener.Close()
		}

		reader := bufio.NewReader(conn)
		firstLine, readErr := reader.ReadString('\n')

		if readErr != nil {
			fmt.Println("Error while reading first line of request")
		}

		splittedFirstLine := strings.Split(firstLine, " ")
		path := splittedFirstLine[1]
		if path == "/firstendpoint" {
			fmt.Println("Connection established")
			content := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 8\r\n\r\nHalloDu!")
			conn.Write(content)
		} else {
			errorHttp404 := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\n404 Not Found")
			conn.Write(errorHttp404)
		}

		conn.Close()

	}
}
