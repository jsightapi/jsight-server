package main

import (
	"net"
	"net/http"
	"strings"
)

func sendToStatisticServer(b []byte) error {
	a, err := net.ResolveUDPAddr("udp4", "stat.jsight.io:1053")
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp4", nil, a)
	if err != nil {
		return err
	}

	defer func() {
		_ = c.Close()
	}()

	_, err = c.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip
	}

	ips := r.Header.Get("X-Forwarder-For")
	splitIps := strings.Split(ips, ",")
	for _, ip = range splitIps {
		netIP = net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return ""
}
