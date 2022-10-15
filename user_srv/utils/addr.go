/**
 * @file: addr.go
 * @time: 2022-10-15 17:24
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package utils

import (
	"net"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
