package utils

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// Createfile creates the "config" in the script folder
func Createfile(s string) {
	_, err := os.Create(s)
	if err != nil {
		fmt.Printf("error creating config file: %v", err)
		return
	}
}

// Openfile opens the "config" file and returns its file descriptor
func Openfile(s string) *os.File {
	f, err := os.OpenFile(s, os.O_RDWR, 0)

	if err != nil {
		fmt.Printf("error opening config file: %v", err)
		return nil
	}
	return f
}

// GetMyIP returns your external IP and gives an error if it does not find it.
func GetMyIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("Can't find external IP")
}
