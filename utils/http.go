package utils

import (
	"crypto/tls"
	"fmt"
	"hostscan/vars"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CheckPort(ip string, port string) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ip, port), time.Duration(*vars.Timeout)*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func GetHttpBody(url, host string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(*vars.Timeout) * time.Second,
	}

	reqest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return ""
	}

	reqest.Host = host

	var ua string
	if *vars.IsRandUA == true {
		ua = RandUA()
	} else {
		ua = fmt.Sprintf("golang-hostscan/%v", vars.Version)
	}

	reqest.Header.Add("User-Agent", ua)

	response, err := client.Do(reqest)
	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		//elog.Error(fmt.Sprintf("DoGet: %s [%s]", url, err))
		return ""
	}

	filterStatusCodes := []int{}
	filters := strings.TrimSpace(*vars.FilterRespStatusCodes)
	if len(filters) > 0 {
		for _, statusCode := range strings.Split(filters, ",") {
			filterStatusCode, err := strconv.Atoi(strings.TrimSpace(statusCode))
			if err != nil {
				continue
			}
			filterStatusCodes = append(filterStatusCodes, filterStatusCode)
		}
		if !containsStatusCode(response.StatusCode, filterStatusCodes) {
			return ""
		}
	}

	bodyByte, _ := io.ReadAll(response.Body)
	body := string(bodyByte)

	return body
}

func containsStatusCode(a int, l []int) bool {
	for _, item := range l {
		if a == item {
			return true
		}
	}

	return false
}
