package main

import (
	"diy_load_balancer/src/tcpip"
	"github.com/songgao/packets/ethernet"
	"log"
	"net"
)

const (
	tapAddr = "10.1.0.10/24"
	tapName = "tap0"
)

var tapIp = []byte{10, 1, 0, 11}

func main() {
	tap := tcpip.NewTap(tapName)
	err := tap.Open()
	log.Println("tapped opened:", tap.Fd)
	if err != nil {
		panic(err)
	}

	err = tap.SetAddress(tapAddr)
	if err != nil {
		log.Println(err)
	}

	err = tap.SetUp()
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Println(err)
	}

	ifrInterface, err := net.InterfaceByName(tapName)
	if err != nil {
		log.Println(err)
	}

	log.Println("MTU:", ifrInterface.MTU)
	tcpip.NetDevAllocate(tapIp, ifrInterface.HardwareAddr, ifrInterface.MTU)

	//Todo : find out how the byte size is calculated
	buf := make([]byte, 1<<16)
	for {
		rn, err := tap.Read(buf)
		if err != nil {
			panic(err)
		}

		frame := make(ethernet.Frame, rn)
		copy(frame, buf[:rn])

		log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		log.Printf("Receiving packets:\n"+
			"Destination -> %s\n"+
			"Source -> %s\n"+
			"Ethernet Type -> %x\n",
			frame.Destination(),
			frame.Source(),
			frame.Ethertype())
		log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		switch frame.Ethertype() {
		case ethernet.ARP:
			log.Print("ARP request")
			err = tcpip.HandleArpRequest(frame.Payload())

			if err != nil {
				log.Println("err:", err)
			}
			log.Println("done handling arp")
		default:
			log.Println("Not supported request")
		}

	}

}
