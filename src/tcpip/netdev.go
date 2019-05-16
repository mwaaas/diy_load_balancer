package tcpip

import (
	"bytes"
	"log"
)

var rawDev *netDev

type netDev struct {
	addr   []byte
	hwaddr []byte
	mtu    int
}

func NetDevAllocate(addr []byte, hwaddr []byte, mtu int) *netDev {
	rawDev = &netDev{addr, hwaddr, mtu}
	return rawDev
}

func netDevGet(dip []byte) *netDev {
	if bytes.Equal(dip, rawDev.addr) {
		return rawDev
	}
	return nil
}

func netdevTransmit(data []byte, dMac []byte, etherType uint16, net *netDev) error {
	frame := encodeEthernetFrame(dMac, net.hwaddr[:],
		etherType, data)
	log.Println("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	log.Printf("sending packet:\n "+
		"Destination -> %s\n"+
		"Source -> %s\n"+
		"Ethernet type: %x\n",
		frame.Destination(),
		frame.Source(),
		frame.Ethertype())
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	//log.Printf("data to send: %x\n", frame)
	_, err := tap.Write(frame)
	if err != nil {
		log.Println("error sending data:", err)
		return err
	}
	return nil
}
