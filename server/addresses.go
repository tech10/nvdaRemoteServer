package server

import (
	"errors"
	"net"
	"strconv"
)

type addressList []string

func (a *addressList) String() string {
	l := ""
	if len(*a) == 0 {
		return ""
	}
	for _, v := range *a {
		if v == "" {
			continue
		}
		l += v + "\n"
	}
	return l
}

func (a *addressList) Set(v string) error {
	err := address_valid(v)
	if err != nil {
		return err
	}
	*a = append(*a, v)
	return nil
}

func address_valid(v string) error {
	var ip string
	var portstr string
	var port int
	var err error
	ip, portstr, err = net.SplitHostPort(v)
	if err != nil {
		return err
	}
	if ip != "" {
		ipcheck := net.ParseIP(ip)
		if ipcheck == nil {
			return errors.New("The ip address " + ip + " is invalid.")
		}
	}
	port, err = strconv.Atoi(portstr)
	if err != nil {
		return err
	}
	if port < 1 || port > 65536 {
		return errors.New("Invalid port range, " + portstr + " is not a valid port number.")
	}
	return nil
}
