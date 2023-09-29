package client

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDifference(t *testing.T) {
	queries := []string{"test"}
	filenames := []string{"vm1.log"}
	expectedCounts := map[string]int{
		"192.168.10.20": 5, // 对应 machine01
		"192.168.10.30": 3, // 对应 machine02
		"192.168.10.40": 3, // 对应 machine03
		"192.168.10.50": 1, // 对应 machine04
	}
	machineIPs := []string{"192.168.10.20", "192.168.10.30", "192.168.10.40", "192.168.10.50"}

	totalElapsed := time.Duration(0)
	testRuns := 5

	for run := 1; run <= testRuns; run++ {
		start := time.Now()

		var wg sync.WaitGroup
		results := make(chan struct {
			IP    string
			Count int
		}, len(machineIPs))

		for _, ip := range machineIPs {
			wg.Add(1)
			go func(ip string) {
				defer wg.Done()
				cmd := exec.Command("./client", queries[0], filenames[0], ip)
				var out bytes.Buffer
				cmd.Stdout = &out
				cmd.Run()

				ret, _ := lineCountFromOutput(out.String(), ip)
				results <- struct {
					IP    string
					Count int
				}{IP: ip, Count: ret}
			}(ip)
		}

		wg.Wait()
		close(results)

		fmt.Printf("Run %d:\n", run)
		for result := range results {
			fmt.Printf("Machine %s: Count: %d\n", result.IP, result.Count)
			if result.Count != expectedCounts[result.IP] {
				t.Errorf("Test failed for %s, expected count: %d, actual count: %d", result.IP, expectedCounts[result.IP], result.Count)
			}
		}

		end := time.Now()
		elapsed := end.Sub(start)
		totalElapsed += elapsed

		fmt.Printf("Elapsed time: %v\n", elapsed)
	}

	averageElapsed := totalElapsed / time.Duration(testRuns)
	fmt.Printf("Average elapsed time: %v\n", averageElapsed)
}

func deleteFile(f string) {
	if _, err := os.Stat("./" + f); err == nil {
		err := os.Remove("./" + f)
		_ = err
	}
}
func lineCountFromOutput(output string, ip string) (int, error) {
	count := 0
	scanner := bufio.NewScanner(strings.NewReader(output))
	machinePrefix := ipToMachinePrefix(ip)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, machinePrefix) && strings.Contains(line, ":test vm") {
			count++
		}
	}
	return count, scanner.Err()
}
func ipToMachinePrefix(ip string) string {
	switch ip {
	case "192.168.10.20":
		return "machine01"
	case "192.168.10.30":
		return "machine02"
	case "192.168.10.40":
		return "machine03"
	case "192.168.10.50":
		return "machine04"
	default:
		return ""
	}
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
