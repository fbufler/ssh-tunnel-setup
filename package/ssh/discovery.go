package ssh

import (
	"fmt"
	"net"
	"time"
)

// DiscoverRemote discovers if a remote host is reachable on the provided port
func DiscoverRemote(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	timeout := 5 * time.Second

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
