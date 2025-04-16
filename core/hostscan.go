package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hostscan/elog"
	"hostscan/models"
	"hostscan/utils"
	"hostscan/vars"
	"io"
	"os"
	"strings"
	"sync"
)

func Scan() error {
	taskWg := sync.WaitGroup{}
	threadWg := sync.WaitGroup{}

	var ipList []string
	var hostList []string
	portList := []string{"80", "443"}
	schemeList := []string{"http", "https"}

	taskChan := make(chan Task, *vars.Thread)
	resultChan := make(chan models.Result)

	// standalone ip
	if strings.Contains(*vars.Ip, "/") {
		ipList = HandleIpRange(*vars.Ip)
	} else {
		ipList = append(ipList, *vars.Ip)
	}

	// standalone host
	if len(*vars.Host) > 0 {
		hostList = append(hostList, *vars.Host)
	}

	// ip file
	if len(*vars.IpFile) > 0 {
		ipF, err := os.Open(*vars.IpFile)
		defer ipF.Close()
		if err != nil {
			return err
		}
		ipBuf := bufio.NewReader(ipF)

		for {
			ip, err := ipBuf.ReadString('\n')
			ip = strings.TrimSpace(ip)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			ipList = append(ipList, ip)
		}
	}

	// host file
	if len(*vars.HostFile) > 0 {
		hostF, err := os.Open(*vars.HostFile)
		defer hostF.Close()
		if err != nil {
			return err
		}
		hostBuf := bufio.NewReader(hostF)

		for {
			host, err := hostBuf.ReadString(10)
			host = strings.TrimSpace(host)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			hostList = append(hostList, host)
		}
	}

	// port list
	if len(*vars.Iports) > 0 {
		portList = strings.Split(*vars.Iports, ",")
	}

	// 创建vars.ThreadNum个协程
	for i := 0; i < *vars.Thread; i++ {
		go goScan(taskChan, resultChan, &taskWg, &threadWg)
		threadWg.Add(1)
	}

	go func(wg *sync.WaitGroup) {
		wg.Add(1)

		defer wg.Done()

		fp, err := os.OpenFile(*vars.OutFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			elog.Warn(fmt.Sprintf("Open file[%s]: %s", *vars.OutFile, err))
			return
		}
		defer fp.Close()

		for {
			select {
			case result, ok := <-resultChan:
				if !ok {
					return
				}
				elog.Info(fmt.Sprintf("Uri: %s, Hostname: %s", result.Uri, result.Host))
				t, _ := json.Marshal(result)
				_, _ = fp.WriteString(string(t) + "\n")
			}
		}
	}(&threadWg)

	for _, ip := range ipList {
		for _, port := range portList {
			// check connection first, if <ip:port> is not reachable, skip
			if !utils.CheckPort(ip, port) {
				if *vars.Verbose {
					elog.Warn(fmt.Sprintf("IP: %s, Port: %s --> Connection Failed", ip, port))
				}
				continue
			}

			for _, host := range hostList {
				for _, scheme := range schemeList {
					if scheme == "http" && port == "443" {
						continue
					}
					if scheme == "https" && port == "80" {
						continue
					}
					taskChan <- Task{
						Uri:  fmt.Sprintf("%s://%s:%s", scheme, ip, port),
						Host: host,
					}
					taskWg.Add(1)
				}
			}
		}
	}

	// wait for all task to done
	taskWg.Wait()
	close(taskChan)
	close(resultChan)

	// wait all threads to done
	threadWg.Wait()

	return nil
}
