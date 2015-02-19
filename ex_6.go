package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

var currentNumber int = 0

func main() {
	data_ch := make(chan int, 1)
	send_ch := make(chan int, 1)
	wait_ch := make(chan bool)
	go server(data_ch)
	i := <-data_ch
	if i < 0 {
		initialization(send_ch)
		<-wait_ch
	} else {
		go reader(data_ch, send_ch)
		currentNumber = i
		<-wait_ch
	}

}
func reader(data_ch chan int, send_ch chan int) {
	for {
		num := <-data_ch
		if num < 0 {
			initialization(send_ch)
			break
		} else {
			currentNumber = num
		}
	}
}

func initialization(send_ch chan int) {
	go transmit(send_ch)
	go counter(send_ch, currentNumber)
	newProcess()
}

func newProcess() {
	cmd := exec.Command("start", "ny", "call", "ex6.go")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(out))

}
func counter(send_ch chan int, currentNumber int) {
	i := currentNumber
	for {
		send_ch <- i
		fmt.Print(i)
		i++
		time.Sleep(250 * time.Millisecond)
	}
}

const (
	master_hostname = "localhost"
	master_port     = ":8080"
)

func transmit(send_ch chan int) {
	conn, err := net.Dial("udp", master_hostname+master_port)

	if err != nil {
		panic(err)
	}
	for {
		i := <-send_ch
		payloadString := strconv.Itoa(i)
		conn.Write([]byte(payloadString + "/"))
	}

}

func server(data chan<- int) {
	laddr, err := net.ResolveUDPAddr("udp", master_port)
	if err != nil {
		panic(err)
	}
	buffer := make([]byte, 1024)
	newData := make(chan int)
	for {
		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			panic(err)
		}
		go func(newData chan int, conn *net.UDPConn) {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println(n)
				panic(err)
			}
			s := ""
			for i := 0; i < n; i++ {
				if string(buffer[i]) == "/" {
					s = string(buffer[:i-1])
				}
			}
			i, err := strconv.Atoi(s)
			newData <- i
		}(newData, conn)
		for {
			select {
			case i := <-newData:
				data <- i
			case <-time.After(1 * time.Second):
				data <- -1
				conn.Close()
				break
			}
		}
	}
}
