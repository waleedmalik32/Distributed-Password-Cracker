package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/muaazbinsaeed/project/cpu"
	"github.com/nyarlabo/go-crypt"
)

/*
abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
*/
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var stop bool = false
var answer string = "$$"
var external_found bool = false

func FirstJob(maxLength int, code string) {
	runtime.GOMAXPROCS(1)
	var wg sync.WaitGroup
	wg.Add(1)
	go FirstCombo(maxLength, code, &wg)
	wg.Wait()
	stop = true
	//external_found = false
}
func SecondJob(maxLength int, code string, job int) {
	cores := cpu.NoOfProc() // //1 //2 //3 //4
	if cores > 4 {
		cores = 4
	}
	combo := 4 //len(letters)
	divisions := combo / (cores)
	runtime.GOMAXPROCS(cores)
	var wg sync.WaitGroup
	wg.Add(cores)

	fmt.Println("Max Cores:", cores)
	fmt.Println("div: ", divisions)
	fmt.Println("Code: " + code)
	start, end := genJobIndex(job)
	if cores == 1 {
		go SecondCombo(maxLength, code, &wg, start, end)
	} else {
		_, real_end := genJobIndex(job)
		end = start + divisions - 1
		for i := 0; i < cores; i++ {
			if i == cores-1 {
				end = real_end
			}
			go SecondCombo(maxLength, code, &wg, start, end)
			start += divisions
			end += divisions
		}
	}
	wg.Wait()
	stop = true
	//external_found = false
}

/*func main() {
	maxLength := 4
	code := GenerateCrypt()
	job := 13
	if job == 0 {
		FirstJob(maxLength, code)
	} else {
		SecondJob(maxLength, code, job)
	}
}*/
func GenerateCrypt() string {
	text := "abb" //"Jaaa" //"Iabb" //"Naab" //text := "aab"
	code := (crypt.Crypt(text, ""))
	return code
}
func FirstCombo(maxLength int, code string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i < maxLength; i++ {
		if stop == true {
			return
		}
		generateCombo("", i, code)
	}
}
func SecondCombo(maxLength int, code string, wg *sync.WaitGroup, start int, end int) {
	defer wg.Done()
	fmt.Println("Second", string(letters[start]), string(letters[end]))
	for i := start; i <= end; i++ {
		if stop == true {
			return
		}
		generateCombo(string(letters[i]), maxLength-1, code)
	}
}
func generateCombo(prefix string, k int, code string) {
	if stop == true {
		return
	}
	if k == 0 {
		fmt.Println(prefix)
		if code == (crypt.Crypt(prefix, "")) {
			fmt.Println("Found:" + prefix)
			stop = true
			answer = prefix
		}
		return
	}
	for i := 0; i < len(letters); i++ {
		newPrefix := prefix + string(letters[i])
		generateCombo(newPrefix, k-1, code)
	}
}
func genJobIndex(job int) (int, int) {
	start := 0
	end := 0
	if job == 1 {
		start = 1
		end = 4
	} else if job == 2 {
		start = 5
		end = 8
	} else if job == 3 {
		start = 8
		end = 12
	} else if job == 4 {
		start = 13
		end = 16
	} else if job == 5 {
		start = 17
		end = 20
	} else if job == 6 {
		start = 21
		end = 24
	} else if job == 7 {
		start = 25
		end = 28
	} else if job == 8 {
		start = 29
		end = 32
	} else if job == 9 {
		start = 33
		end = 36
	} else if job == 10 {
		start = 37
		end = 40
	} else if job == 11 {
		start = 41
		end = 44
	} else if job == 12 {
		start = 45
		end = 48
	} else if job == 13 {
		start = 49
		end = 52
	}
	start--
	end--
	return start, end
}

func main() {

	cipherText := flag.String("cipherText", "#noText", "a string")
	hostName := flag.String("hostName", "127.0.0.1", "a string") //127.0.0.1
	port := flag.String("port", "1216", "a string")              //2600

	flag.Parse()

	fmt.Println("cipherText:", *cipherText)
	fmt.Println("hostName:", *hostName)
	fmt.Println("port:", *port)

	address := *hostName + ":" + *port
	text1 := *cipherText
	fmt.Println("Text to send: ", text1)
	// connect to this socket
	conn, _ := net.Dial("tcp", address)
	// send to socket

	maxLength := 3 //4
	//code := GenerateCrypt()
	//job := 13

	// listen for reply
	//conn.Write([]byte("$$" + "\n"))
	//fmt.Fprintf(conn, "$$"+"\n")

	for {

		msg_send := ""
		if answer != "$$" {
			msg_send = answer
			time.Sleep(1000 * time.Millisecond)
			answer = "$$"
			stop = false
		} else if answer == "$$" && stop == true && external_found == false {
			msg_send = "#"
			stop = false
		} else {
			msg_send = "$$"
		}
		external_found = false
		fmt.Println("sending msg", msg_send)
		fmt.Fprintf(conn, msg_send+"\n")
		fmt.Println("sending start")
		//fmt.Fprintf(conn, "$$")
		msg, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println("Message from server: " + msg)

		message := string(msg[:len(msg)-1])
		fmt.Println("Message from server: " + message)

		//message := msg
		if message == "found" {
			stop = true

			time.Sleep(1000 * time.Millisecond)
			//stop = false
			external_found = true

		} else if message != "$$" {
			msg := string(message)
			i := strings.Index(msg, "#")
			fmt.Println("index: ", string(i))
			//job, _ := strconv.Atoi(msg[:i-1])
			job, _ := strconv.Atoi(msg[:i])
			code := msg[i+1:]
			fmt.Println("job: ", job, " code: ", code)
			stop = false
			var wg sync.WaitGroup
			if job == 0 {
				runtime.GOMAXPROCS(1)

				wg.Add(1)
				go FirstCombo(maxLength, code, &wg)

			} else {
				cores := cpu.NoOfProc() // //1 //2 //3 //4
				if cores > 4 {
					cores = 4
				}
				combo := 4 //len(letters)
				divisions := combo / (cores)
				runtime.GOMAXPROCS(cores)
				//var wg sync.WaitGroup
				wg.Add(cores)

				fmt.Println("Max Cores:", cores)
				fmt.Println("div: ", divisions)
				fmt.Println("Code: " + code)
				start, end := genJobIndex(job)
				if cores == 1 {
					go SecondCombo(maxLength, code, &wg, start, end)
				} else {
					_, real_end := genJobIndex(job)
					end = start + divisions - 1
					for i := 0; i < cores; i++ {
						if i == cores-1 {
							end = real_end
						}
						go SecondCombo(maxLength, code, &wg, start, end)
						start += divisions
						end += divisions
					}
				}

			}
			wg.Wait()
			stop = true
		}

	}

}
