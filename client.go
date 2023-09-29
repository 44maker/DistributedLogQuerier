package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"
)

type Servers struct {
	Servers []ServerInfo `json:"server_list"`
}

type ServerInfo struct {
	Id       string `json:"id"`
	Hostname string `json:"hostname"`
	Logfile  string `json:"logfile"`
}

func lineCount(filename string) (int, error) {
	count := 0
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		if len(s.Text()) > 0 {
			count++
		}
	}
	return count, s.Err()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("请在命令行参数中输入 grep 命令和端口号！")
		return
	}
	grepCmd := os.Args[1]
	portNum := "3000"
	fileName := os.Args[2]

	jsonFile, err := os.Open("servers.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var serverInfo Servers
	json.Unmarshal(byteValue, &serverInfo)

	var wg sync.WaitGroup
	wg.Add(len(serverInfo.Servers))
	start := time.Now()

	results := make(map[string]int)
	var resultsMutex sync.Mutex

	for i := 0; i < len(serverInfo.Servers); i++ {
		go func(hostname string, logfile string, grepCmd string, portNum string, id string) {
			defer wg.Done()
			conn, err := net.Dial("tcp", hostname+":"+portNum)
			if err != nil {
				fmt.Println(err)
				return
			}
			name := "machine" + id
			fmt.Fprintf(conn, grepCmd+" "+name+" "+fileName)
			f, err := os.Create(logfile)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()

			for {
				message := make([]byte, 5120)
				n1, err := conn.Read(message)
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Println(err)
					return
				}
				text := string(message[:n1])
				fmt.Println(text)
				n2, err := f.WriteString(text)
				if err != nil {
					fmt.Println(err)
					return
				}
				_ = n1
				_ = n2
			}

			lc, _ := lineCount(logfile)
			name = "machine" + id
			fmt.Println(name, lc)

			resultsMutex.Lock()
			results[hostname] = lc
			resultsMutex.Unlock()
		}(serverInfo.Servers[i].Hostname, serverInfo.Servers[i].Logfile, grepCmd, portNum, serverInfo.Servers[i].Id)
	}
	wg.Wait()
	end := time.Now()
	elapsed := end.Sub(start)

	fmt.Println("结果:", results)
	fmt.Println("延迟: ", elapsed)
}