package xFacility

import (
	"net"
)

type NetworkCardInfo struct {
	Interface string `json:"interface"`
	Ip        string `json:"ip"`
	MacAddr   string `json:"macAddr"`
}

// GetNetworkCardInfoList 获取网卡信息列表
func GetNetworkCardInfoList() ([]*NetworkCardInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var network []*NetworkCardInfo
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
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
				continue
			}
			network = append(network, &NetworkCardInfo{
				Interface: iface.Name,
				Ip:        ip.String(),
				MacAddr:   iface.HardwareAddr.String(),
			})
		}
	}
	return network, nil
}
