package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"

	"github.com/nyarlabo/go-crypt"
)

func GenerateCrypt() string {
	text := "Iab" //"abbb" //"aaab" //"abb" //"Jaaa" //"Iabb" //"Naab" //text := "aab"
	code := (crypt.Crypt(text, ""))
	return code
}

func main() {

	cipherText := flag.String("cipherText", "#noText", "a string")
	hostName := flag.String("hostName", "127.0.0.1", "a string") //127.0.0.1
	port := flag.String("port", "1217", "a string")              //2600

	flag.Parse()

	fmt.Println("cipherText:", *cipherText)
	fmt.Println("hostName:", *hostName)
	fmt.Println("port:", *port)

	address := *hostName + ":" + *port

	code := *cipherText
	fmt.Println("Text to send: ", code)

	// connect to this socket
	conn, _ := net.Dial("tcp", address)

	code = GenerateCrypt()
	fmt.Println("code:" + code)
	fmt.Fprintf(conn, code+"\n")
	// listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Message from server: " + message)

}
