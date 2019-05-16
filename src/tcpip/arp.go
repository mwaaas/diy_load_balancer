package tcpip

import (
	"encoding/binary"
	"fmt"
	"log"
)

// https://tools.ietf.org/html/rfc826
/*
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|      hwtype(2 bytes)        |       protype(2 bytes)          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|hwsize(1bytes) |psize(1bytes)|        opcode(2 bytes)          |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ 8bytes
|       smac(6 bytes)         |				sip(4 bytes)        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|       dmac(6 bytes)         |				dip(4 bytes)        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+ 28 bytes
*/

// 8 bytes

type ArpRequest []byte

const (
	ARP_ETHERNET = 0x0001
	ARP_IPV4     = 0x0800
	ARP_REQUEST  = 0x0001
	ARP_REPLY    = 0x0002
)

func (arp ArpRequest) hwType() uint16 {
	return binary.BigEndian.Uint16(arp[0:2])
}

func (arp ArpRequest) proType() uint16 {
	return binary.BigEndian.Uint16(arp[2:4])
}

func (arp ArpRequest) hwSize() uint8 {
	return arp[4]
}

func (arp ArpRequest) proSize() uint8 {
	return arp[5]
}

func (arp ArpRequest) opCode() uint16 {
	return binary.BigEndian.Uint16(arp[6:8])
}

func (arp ArpRequest) sMac() []byte {
	return arp[8:14]
}

func (arp ArpRequest) sIp() []byte {
	return arp[14:18]
}

func (arp ArpRequest) dMac() []byte {
	return arp[18:24]
}

func (arp ArpRequest) dIp() []byte {
	return arp[24:28]
}

func (arp ArpRequest) string() string {
	return fmt.Sprintf("Arp payload\n"+
		"hwType -> %x\n"+
		"proType -> %x\n"+
		"hwSize -> %x\n"+
		"proSize -> %x\n"+
		"opCode -> %x\n"+
		"sMac -> %x\n"+
		"sIp -> %x\n"+
		"dMac -> %x\n"+
		"dIp -> %x\n", arp.hwType(), arp.proType(),
		arp.hwSize(), arp.proSize(), arp.opCode(),
		arp.sMac(), arp.sIp(), arp.dMac(), arp.dIp())
}

func arpReply(arp ArpRequest, net *netDev) error {
	// set opcode
	binary.BigEndian.PutUint16(arp[6:8], ARP_REPLY)

	// set destination mac and ip as the arp src and mac and ip
	copy(arp[18:24], arp.sMac())
	copy(arp[24:28], arp.sIp())

	// set the source mac and ip
	copy(arp[8:14], net.hwaddr)
	copy(arp[14:18], net.addr)

	log.Println()
	log.Println("arp reply:\n", arp.string())
	log.Println()
	return netdevTransmit(arp, arp.sMac(), 0x806, net)
}

func HandleArpRequest(arp ArpRequest) error {
	if arp.hwType() != ARP_ETHERNET {
		log.Printf("ARP: Unsupported HW type\n")
		return nil
	}

	if arp.proType() != ARP_IPV4 {
		log.Printf("ARP: Unsupported protocol\n")
		return nil
	}

	log.Println()
	log.Println("arp request:\n", arp.string())
	log.Println()

	net := netDevGet(arp.dIp())

	if net == nil {
		log.Printf("ARP was not for us\n")
		return nil
	}

	switch arp.opCode() {
	case ARP_REQUEST:
		err := arpReply(arp, net)
		return err
	default:
		log.Printf("ARP: Opcode not supported %x\n", arp.opCode())

	}
	return nil
}
