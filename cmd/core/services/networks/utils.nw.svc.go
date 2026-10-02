package service_networks

import "net"

func getRandomPort() (net.Listener, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	return listener, nil
}
