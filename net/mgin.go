package net

import (
	"net"
	"time"
)

//WaitAddresses
//wait port avaliable
//addresses: ["ip:port","ip:port"]
func WaitAddresses(callResult func(address string, isSuccess bool), addresses ...string) {
	isSuccess := false
	for !isSuccess {
		isSuccess = true
		for _, address := range addresses {
			con, err := net.DialTimeout("tcp", address, 10*time.Second)
			if err == nil {
				con.Close()
				callResult(address, true)
			} else {
				callResult(address, false)
				isSuccess = false
				time.Sleep(5 * time.Second)
				break
			}
		}
	}
}
