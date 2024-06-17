package internal

import (
	"fmt"
)

// BuildAddress - build address from ip and port
//
//	ip: "127.0.0.1", port: "8080"
//	result: "127.0.0.1:8080"
func BuildAddress(ip, port string) string {
	return fmt.Sprintf("%s:%s", ip, port)
}
