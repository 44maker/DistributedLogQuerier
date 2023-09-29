package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
)

type Cache struct {
	data sync.Map
}

func (c *Cache) Get(key string) ([]string, bool) {
	value, ok := c.data.Load(key)
	if ok {
		return value.([]string), true
	}
	return nil, false
}

func (c *Cache) Set(key string, value []string) {
	c.data.Store(key, value)
}

var cache = Cache{}

func printErr(err error, s string) {
	if err != nil {
		fmt.Println("在", s, "发生错误\n", err.Error())
		os.Exit(1)
	}
}

func executeGrep(query string, vm string) ([]string, error) {
	cachedResult, ok := cache.Get(query)
	if ok {
		return cachedResult, nil
	}

	file, err := os.Open(vm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re, err := regexp.Compile(query)
	if err != nil {
		return nil, err
	}

	var matches []string
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			matches = append(matches, fmt.Sprintf("%d:%s", lineNumber, line))
		}
		lineNumber++
	}

	cache.Set(query, matches)
	return matches, scanner.Err()
}

func parseRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	printErr(err, "读取")

	reqArr := strings.Split(string(buf[:reqLen]), " ")

	matches, err := executeGrep(reqArr[0], reqArr[2])
	if err != nil {
		fmt.Println("执行 grep 时发生错误:", err)
		conn.Close()
		return
	}

	out := ""
	for i, match := range matches {
		if i == len(matches)-1 {
			out = out + reqArr[1] + " " + "行 " + match
		} else {
			out = out + reqArr[1] + " " + "行 " + match + "\n"
		}
	}

	conn.Write([]byte(out))
	conn.Close()
}

func getIPAddrAndLogfile() string {
	data, err := ioutil.ReadFile("ip_address")
	if err != nil {
		panic(err)
	}

	ip := string(data[:len(data)])

	if strings.HasSuffix(ip, "\n") {
		ip = ip[:(len(ip) - 1)]
	}
	fmt.Println("当前 VM 的 IP 地址:\n", ip)
	return ip
}

func main() {
	ip := getIPAddrAndLogfile()
	l, err := net.Listen("tcp", ip+":3000")
	printErr(err, "监听")

	defer l.Close()
	fmt.Println("监听端口 3000")

	for {
		conn, err := l.Accept()
		fmt.Println("接受:", conn.RemoteAddr().String())
		printErr(err, "接受")

		go parseRequest(conn)
	}
}
