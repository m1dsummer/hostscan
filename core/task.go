package core

import (
	"fmt"
	"hostscan/elog"
	"hostscan/models"
	"hostscan/utils"
	"hostscan/vars"
	"regexp"
	"strings"
	"sync"
)

type Task struct {
	Uri  string
	Host string
}

func getTitle(body string) string {
	re := regexp.MustCompile(`<title>([\s\S]*?)</title>`)
	match := re.FindStringSubmatch(body)
	if match != nil && len(match) > 1 {
		return strings.TrimSpace(match[1])
	} else {
		return ""
	}
}

func goScan(taskChan chan Task, resultChan chan models.Result, taskWg *sync.WaitGroup, threadWg *sync.WaitGroup) {
	defer threadWg.Done()

	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return
			} else {
				body := utils.GetHttpBody(task.Uri, task.Host)
				if len(body) == 0 {
					if *vars.Verbose {
						elog.Info(fmt.Sprintf("Uri: %s, Hostname: %s --> No Response Body", task.Uri, task.Host))
					}
					taskWg.Done()
					continue
				}
				title := getTitle(body)

				if title != "" {
					shouldSkip := false
					for _, blackTitle := range vars.BlackListTitles {
						if strings.Contains(title, blackTitle) {
							shouldSkip = true
							break
						}
					}
					if shouldSkip {
						taskWg.Done()
						continue
					}
				}

				result := models.Result{
					Uri:   task.Uri,
					Host:  task.Host,
					Title: title,
				}
				resultChan <- result
				taskWg.Done()
			}
		}
	}
}
