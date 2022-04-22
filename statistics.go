package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

type event struct {
	ClientIPv4      string `json:"clientIPv4,omitempty"`
	SourceType      string `json:"sourceType,omitempty"`
	SourceID        string `json:"sourceID,omitempty"`
	ProjectTitle    string `json:"projectTitle,omitempty"`
	ProjectSize     uint32 `json:"projectSize,omitempty"`
	ProjectHasError bool   `json:"projectHasError,omitempty"`
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

func sendToStatisticServer(e event) error {
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

	var b []byte
	b, err = json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = c.Write(b)
	if err != nil {
		return err
	}

	return nil
}
