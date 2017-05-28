package utility

import (
	"github.com/hashicorp/mdns"
	"net"
	"os"
)

func BroadcastServer() *mdns.Server {
	host, _ := os.Hostname()
	/*listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}*/
	//port := listener.Addr().(*net.TCPAddr).Port
	info := []string{"Remote Camera Service"}
	service, _ := mdns.NewMDNSService(host, "_goshot._tcp", "", "", 8080, getLocalIPS(), info)

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
