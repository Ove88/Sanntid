package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

var currentNumber int = 0

func main() {
	data_ch := make(chan int, 1)
	send_ch := make(chan int, 1)
	wait_ch := make(chan bool)
	go server(data)
	i := <-data_ch
	if i < 0 {
		initialization()
		<-wait_ch
	} else {
		go reader(data_ch)
		currentNumber = i
		<-wait_ch
	}

}
func reader(data_ch chan int) {
	for {
		num := <-data_ch
		if num < 0 {
			initialization(num)
			break
		} else {
			currentNumber = num
		}
	}
}

func initialization() {
	go transmit(send_ch)
	go counter(send_ch, currentNumber)
}

func counter(send_ch chan int, currentNumber int) {
	i := currentNumber
	for {
		send_ch <- i
		fmt.Print(i)
		i++
		time.Sleep(250 * time.Millisecond)
	}
	//inkrementer
	//send tall for utskriving(printNumbers)
	//oppdater backup
	//vent
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
		conn.Write([]byte(payloadString + "/0"))
	}

}

func server(data chan<- int) {
	laddr, err := net.ResolveUDPAddr("udp", master_port)
	if err != nil {
		panic(err)
	}
	buffer := make([]byte, 1024)
	for {
		conn, err := net.ListenUDP("udp", laddr)
		if err != nil {
			panic(err)
		}
		for {
			select {
			case n, err := conn.Read(buffer):
				if err != nil {
					fmt.Println(n)
					panic(err)
				}
				s := string(buffer[:n])
				i, err := strconv.Atoi(s)
				data <- i
			case time.After(1 * time.Second):
				data <- -1
				conn.Close()
				break
			}
		}
	}
}
