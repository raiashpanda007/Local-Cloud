package service_networks

import "net"

func getRandomPort() (net.Listener, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	return listener, nil
}
