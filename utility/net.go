package utility

import (
	"net"
	"os"

	"github.com/hashicorp/mdns"
)

func BroadcastServer(port int) *mdns.Server {
	host, _ := os.Hostname()
	//port := listener.Addr().(*net.TCPAddr).Port
	info := []string{"Remote Camera Service"}
	service, _ := mdns.NewMDNSService(host, "_goshot._tcp", "", "", port, getLocalIPS(), info)

	// Create the mDNS server, defer shutdown
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	return server
	//defer server.Shutdown()
}

// GetLocalIP returns the non loopback local IP of the host
func getLocalIPS() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	var ips []net.IP
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.String() != "127.0.1.1" {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips
}
