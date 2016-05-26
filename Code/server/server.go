package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type InputType struct {
	conn net.Conn
}

var mutex = &sync.Mutex{}
var code string = "$$"

//var found bool = false
var answer string = "$$"

// only needed below for sample processing

func main() {
	clientPort := ":1217"
	slavePort := ":1216"

	go func() {
		server(clientPort)
	}()
	server(slavePort)

}
func server(runAtPort string) {
	fmt.Println("Launching server...")

	inputClients := []InputType{}

	tcpAddr, err := net.ResolveTCPAddr("tcp4", runAtPort)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	fmt.Println("Server Running..." + runAtPort)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		if runAtPort == ":1217" {
			in := InputType{conn}
			inputClients = append(inputClients, in)
			handleClient(inputClients[0].conn)

			inputClients = removeTop(inputClients)
		} else {
			go handleSlave(conn)
		}
	}
}
func handleClient(conn net.Conn) {
	defer conn.Close()
	answer = "$$"
	msg, _ := bufio.NewReader(conn).ReadString('\n')
	code = string(msg[:len(msg)-1])
	fmt.Print("Message from server:|" + code + "|")

	for {
		if answer != "$$" {
			fmt.Print("Sending Clinent ans:|" + answer + "|")
			conn.Write([]byte(answer))
			code = "$$"

			time.Sleep(2000 * time.Millisecond)
			taskDeallocateAll()
			slaveDeAllotAll()
			//time.Sleep(2000 * time.Millisecond)

			//taskDeallocateAll()
			//slaveDeAllotAll()
			break
		}
	}

	answer = "$$"

}

/*
func handleSlave(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte(code))
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from slave: " + message)
}*/

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func removeTop(FILO []InputType) []InputType {
	return FILO[1:len(FILO)]
}

//////////////////////////////////////////Slave//////////////////////////
var slaveid int = 0
var slave []slavestruct
var task = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func handleSlave(conn net.Conn) {

	slaveadd(conn)
	log.Println("new slave", conn, " len:", len(slave))

	nextjob := true

	msg := "$$"
	job := 0
	for {

		for {
			if code != "$$" {
				//time.Sleep(1000 * time.Millisecond)
				break
			}
		}
		if nextjob == true {
			taskDeallocate(conn)
			job = taskAllocate(conn)
			taskshow()
			nextjob = false
			msg = strconv.Itoa(job) + "#" + code
		} else {
			msg = "$$"
		}
		//fmt.Println("listing to slave")

		msgg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			conn.Close()
			break
		}
		n := string(msgg[:len(msgg)-1])

		//fmt.Println("Message from slave: " + n)
		if err != nil {
			conn.Close()
			break
		} else if string(n[0]) == "#" {
			fmt.Print("found: " + n[1:])
			taskNotFound(conn)
			nextjob = true
			time.Sleep(1000 * time.Millisecond)
		} else if string(n[0]) != "#" && n != "$$" {
			fmt.Print("found: " + n)
			nextjob = true
			answer = n
			/*
				time.Sleep(1000 * time.Millisecond)
				taskDeallocateAll()
				slaveDeAllotAll()*/
		}
		if answer != "$$" {
			msg = "found"
			/*time.Sleep(1000 * time.Millisecond)
			taskDeallocateAll()
			slaveDeAllotAll()*/

		}
		if n != "$$" {
			fmt.Println(conn, "Message from slave: "+n)
		}
		if msg != "$$" {
			fmt.Println(conn, "Message To slave:"+msg)
		}
		_, err = fmt.Fprintf(conn, msg+"\n")
		if err != nil {
			conn.Close()
			break
		}
		/*_, err = conn.Write([]byte(msg + "\n"))
		if err != nil {
			conn.Close()
			break
		}*/

	}
	log.Printf("Connection from %v closed.", conn.RemoteAddr())
	slave_index := slaveIndex(conn)
	taskDeallocate(conn)
	slavedelete(slave_index)
	taskshow()
	fmt.Println("Connection from slave.", conn, "at", slave_index, "now len:", len(slave))
}
func slaveadd(conn net.Conn) {
	slaveid++
	slave = append(slave, slavestruct{conn, -1, slaveid})
}
func slavedelete(i int) {
	slave = append(slave[:i], slave[i+1:]...)
}
func slaveIndex(c net.Conn) int {
	for i := 0; i < len(slave); i++ {
		if slave[i].conn == c {
			return i
		}
	}
	return -1
}
func slaveDeAllotAll() {
	for i := 0; i < len(slave); i++ {
		slave[i].assig = -1
	}
}
func slaveNo(c net.Conn) int {
	for i := 0; i < len(slave); i++ {
		if slave[i].conn == c {
			return slave[i].id
		}
	}
	return -1
}

type slavestruct struct {
	conn  net.Conn
	assig int
	id    int
}

func taskAllocate(conn net.Conn) int {
	slave_no := slaveIndex(conn)
	for i := 0; i < len(task); i++ {

		if task[i] == 0 {

			task[i] = slave[slave_no].id
			slave[slave_no].assig = i

			return i
		}
	}
	return -1
}
func taskDeallocate(conn net.Conn) {
	slave_no := slaveIndex(conn)
	for i := 0; i < len(task); i++ {
		if task[i] == slave[slave_no].id {

			task[i] = 0

			slave[slave_no].assig = 0
			return
		}
	}
}
func taskDeallocateAll() {
	for i := 0; i < len(task); i++ {

		task[i] = 0

	}
}
func taskshow() {
	for i := 0; i < len(task); i++ {
		fmt.Println("Index:", i, " - Slave :", task[i])
	}
}
func taskNotFound(conn net.Conn) {
	slave_no := slaveIndex(conn)
	for i := 0; i < len(task); i++ {
		if task[i] == slave[slave_no].id {

			task[i] = -1

			slave[slave_no].assig = -1
			return
		}
	}
}
