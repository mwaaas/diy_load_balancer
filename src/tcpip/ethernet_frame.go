package tcpip

import (
	"encoding/binary"
	"github.com/songgao/packets/ethernet"
	"net"
)

/*
https://en.wikipedia.org/wiki/Ethernet_frame#Ethernet_II
 dmac	smac 	 type	 数据	 		    FCS
 6bytes	6bytes	2bytes	 46～1500bytes	 4bytes
*/
type EthernetFrame []byte

type Tagging byte

// Const values for different taggings
const (
	NotTagged    Tagging = 0
	Tagged       Tagging = 4
	DoubleTagged Tagging = 8
)

func (f EthernetFrame) Destination() net.HardwareAddr {
	return net.HardwareAddr(f[0:6])
}

func (f EthernetFrame) Source() net.HardwareAddr {
	return net.HardwareAddr(f[6:12])
}

func (f EthernetFrame) Ethertype() ethernet.Ethertype {
	return ethernet.Ethertype{f[12], f[13]}
}

func (f EthernetFrame) Payload() []byte {
	return f[14:]
}

func encodeEthernetFrame(dMac []byte, sMac []byte,
	ethernetType uint16, payload []byte) EthernetFrame {
	ethernetFrame := make(EthernetFrame, 1<<16)

	// header
	copy(ethernetFrame[0:6], dMac)
	copy(ethernetFrame[6:12], sMac)
	binary.BigEndian.PutUint16(ethernetFrame[12:14], ethernetType)

	// payload
	copy(ethernetFrame[14:], payload)
	return ethernetFrame
}
