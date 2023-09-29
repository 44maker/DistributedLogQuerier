package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getIPAddr() string {
	data, err := ioutil.ReadFile("ip_address")
	if err != nil {
		panic(err)
	}

	ip := string(data[:len(data)])

	if strings.HasSuffix(ip, "\n") {
		ip = ip[:(len(ip) - 1)]
	}
	fmt.Println("ip address of current VM:\n", ip)
	return ip
}

func randomStr(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	check(err)
	return string(b[:length]) + "\n"
}

func generateContent(machineNum string) bytes.Buffer {
	var content bytes.Buffer
	content.WriteString("This is an apple.\nToday is a sunny day.\nI love pizza.\n")

	remainingBytes := 60*1024*1024 - content.Len()
	lineLength := 200
	numLines := remainingBytes / (lineLength + 1)

	for i := 0; i < numLines; i++ {
		content.WriteString(randomStr(lineLength))
	}

	return content
}

func main() {
	ip := getIPAddr()
	fmt.Println("current ip address:", ip)
	machineNum := ip[strings.LastIndex(ip, ".")+1:]
	fmt.Println("machine number:", machineNum)

	content := generateContent(machineNum)

	f, err := os.Create("./vm1.log")
	check(err)
	defer f.Close()

	n, err := f.Write(content.Bytes())
	check(err)
	fmt.Printf("Wrote %d bytes\n", n)
}