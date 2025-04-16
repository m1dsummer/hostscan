package core

import (
	"net"
	"sort"
	"strconv"
	"strings"
)

const PORTMAX = 65535
const PORTMIN = 1

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func HandleIpRange(ipRange string) []string {
	ip, ipNet, err := net.ParseCIDR(ipRange)
	if err != nil {
		return []string{}
	}

	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// Remove network address and broadcast address
	return ips[1 : len(ips)-1]
}

func parsePort(ports string) []int {
	var scanPorts []int
	slices := strings.Split(ports, ",")
	var startStr, endStr string
	for _, port := range slices {
		port = strings.Trim(port, " ")
		if len(port) == 0 {
			continue
		}
		if strings.Contains(port, "-") {
			ranges := strings.Split(port, "-")
			if len(ranges) < 2 {
				continue
			}
			sort.Strings(ranges)
			startStr = ranges[0]
			endStr = ranges[1]
			start, err := strconv.Atoi(startStr)
			if err != nil {
				continue
			}
			end, err := strconv.Atoi(endStr)
			if err != nil {
				continue
			}
			if start < PORTMIN {
				start = PORTMIN
			}
			if end > PORTMAX {
				end = PORTMAX
			}
			for i := start; i <= end; i++ {
				scanPorts = append(scanPorts, i)
			}
		} else {
			targetPort, err := strconv.Atoi(port)
			if err == nil {
				scanPorts = append(scanPorts, targetPort)
			}
		}

	}
	return scanPorts
}
