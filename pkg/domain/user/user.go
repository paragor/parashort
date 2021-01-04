package user

import "net"

type UserInfo struct {
	LastIp net.IP
}

func NewUserInfo(lastIp net.IP) *UserInfo {
	return &UserInfo{LastIp: lastIp}
}
