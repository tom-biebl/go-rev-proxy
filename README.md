# go-rev-proxy
Learning Project. Creating an Reverse Proxy in Go from scratch

# Learning Notes
## How to create a listener?
- I need a Listener on Port 5050
-> `net` Library supports TCP or Unix Sockets
--> For a Webserver in my case I need an TCP socket.

### Whats the difference between a TCP- and an Unix-Socket?
#### TCP
- Uses the TCP/IP-Protocoll
- Communicates via ip and ports
- can be used local and over a network
- works on any os

#### Unix
- works only local
- No IP/Port, uses a path in the filesystem eg: /tmp/socket.sock
- very performant, because theres no overhead because no network stack is used
- typicall for communication on the same server

**Solution for Setting up an Listener:**

```go
package main

import "net"

func main() {
    listener, err := net.Listener()

    if err != nil {
        // Errorhandling
    }

    for { // Infiniteloop cause we want to listen till server is shut down
        con, errList := listener.Accept()

        if errList != nil {
            // Errorhandling
        }

        if con != nil {
            // Do something with the connection
        }
    }
}
```

##  How to send an response
**Important to know:** HTTP is text with a specific format.
It follows the following rules:

Statusline
Header1: Value
Header2: Value
-- Empty line for body / header separation
Body

**Solution:**
```go
package main

import (
	"fmt"
	"net"
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

		if conn != nil {
			fmt.Println("Connection established")
			content := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 8\r\n\r\nHalloDu!") // conn.Write need bytes as parameter
			conn.Write(content) // writes the content to the connection (Important dont forget content-length in this case)
		}
	}
}

```

## How to setup an Endpoint
> 📌 **Merke:**  
> Die angeforderte URI (z. B. `/firstendpoint`) steht **immer in der ersten Zeile** eines HTTP-Requests.  
> Diese Zeile folgt dem Format: `METHOD /path HTTP/version` – z. B. `GET /firstendpoint HTTP/1.1`.  
> Um Endpunkte zu unterscheiden, genügt es, **nur die erste Zeile** der Anfrage zu lesen und auszuwerten.

#### Solution
```go
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

		reader := bufio.NewReader(conn) // Important NewReader not Reader
		firstLine, readErr := reader.ReadString('\n') // Read Firstline

		if readErr != nil {
			fmt.Println("Error while reading first line of request")
		}

		splittedFirstLine := strings.Split(firstLine, " ") // split line by whitespace
		path := splittedFirstLine[1] // uri is the second string in the first line

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

```