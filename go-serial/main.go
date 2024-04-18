package main

import (
	"github.com/jacobsa/go-serial/serial"
	"log"
)

func main() {
	// 配置串口参数
	options := serial.OpenOptions{
		//PortName:        "/dev/ttyS0",
		PortName:        "COM5",
		BaudRate:        9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// 打开串口
	port, err := serial.Open(options)
	if err != nil {
		log.Fatal(err)
	}

	// 关闭串口
	defer port.Close()

	// 发送数据
	sendData := []byte("Hello, Serial!")
	n, err := port.Write(sendData)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sent %d bytes: %s", n, sendData)
	go func() {
		// 接收数据
		buf := make([]byte, 1024)
		n, err = port.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Received %d bytes: %s", n, buf)
	}()
	for {

	}
}
