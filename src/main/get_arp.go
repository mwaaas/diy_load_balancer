// Command arpc provides a simple ARP client which can be used to retrieve
// hardware addresses of other machines in a LAN using their IPv4 address.
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mdlayher/arp"
)

var (
	// durFlag is used to set a timeout for an ARP request
	durFlag = flag.Duration("d", 1*time.Second, "timeout for ARP request")

	// ifaceFlag is used to set a network interface for ARP requests
	ifaceFlag = flag.String("i", "tap0", "network interface to use for ARP request")

	// ipFlag is used to set an IPv4 address destination for an ARP request
	ipFlag = flag.String("ip", "", "IPv4 address destination for ARP request")
)

func main() {
	flag.Parse()

	// Ensure valid network interface
	ifi, err := net.InterfaceByName(*ifaceFlag)
	if err != nil {
		log.Fatal("invalid network interface:", err)
	}
	log.Println("ifi:", ifi.Name, "address:", ifi.HardwareAddr)
	// Set up ARP client with socket
	c, err := arp.Dial(ifi)
	if err != nil {
		log.Fatal("error setting up socket:", err)
	}
	defer c.Close()

	// Set request deadline from flag
	if err := c.SetDeadline(time.Now().Add(*durFlag)); err != nil {
		log.Fatal("error setting flags:", err)
	}

	// Request hardware address for IP address
	ip := net.ParseIP(*ipFlag).To4()
	log.Println("ip to request for:", ip)
	mac, err := c.Resolve(ip)
	if err != nil {
		log.Fatal("error parsing ip:", err)
	}

	fmt.Printf("end %s -> %s", ip, mac)
}
