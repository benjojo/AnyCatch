package main

import (
	"github.com/akrennmair/gopcap"
	"log"
	"net"
)

const (
	TYPE_IP  = 0x0800
	TYPE_ARP = 0x0806
	TYPE_IP6 = 0x86DD

	IP_ICMP = 1
	IP_INIP = 4
	IP_TCP  = 6
	IP_UDP  = 17
)

var lastIPs []string
var ipptr int

func StartListeningForPings(device, anycastIP string, snaplen int) {
	lastIPs = make([]string, 255)
	ipptr = 0

	var incomingIP net.IP = net.ParseIP(anycastIP)

	if incomingIP == nil || anycastIP == "1.2.3.4" {
		log.Fatal("Incorrect Anycast IP given")
	}

	if device == "" {
		devs, err := pcap.Findalldevs()
		if err != nil {
			log.Fatal("tcpdump: couldn't find any devices: %s\n", err)
		}
		if 0 == len(devs) {
			log.Fatal("tcpdump: Device error, RTFM please")
		}
		device = devs[0].Name
	}

	h, err := pcap.Openlive(device, int32(snaplen), true, 0)
	if h == nil {
		log.Fatal("tcpdump: %s\n", err)
		return
	}
	defer h.Close()

	// if expr != "" {
	// 	ferr := h.Setfilter(expr)
	// 	if ferr != nil {
	// 		log.Fatal("tcpdump: %s\n", ferr)
	// 		out.Flush()
	// 	}
	// }

	for pkt := h.Next(); pkt != nil; pkt = h.Next() {
		pkt.Decode()
		if pkt.IP != nil {
			if pkt.IP.Protocol == 1 && pkt.IP.DestAddr() == incomingIP.String() {

				// 	type Icmphdr struct {
				// 	Type     uint8
				// 	Code     uint8
				// 	Checksum uint16
				// 	Id       uint16
				// 	Seq      uint16
				// 	Data     []byte
				// }

				for level, headerr := range pkt.Headers {
					switch header := headerr.(type) {
					case *pcap.Icmphdr:
						if header.Type == 0 {
							log.Printf("What(%d) ICMP! %s %d %d %d %d %d", level, pkt.IP.SrcAddr(), header.Type, header.Code, header.Checksum, header.Id, header.Seq)
							LogPing(pkt.IP.SrcAddr())
						}
					case *pcap.Iphdr:
						//log.Printf("What(%d) ICMP! %d %d %d %d %d", level, header.Type, header.Code, header.Checksum, header.Id, header.Seq)
					default:
						log.Printf("Ahem %s ", header)
					}
				}
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func LogPing(ip string) {
	if ipptr+1 > len(lastIPs) {
		ipptr = 0
	}

	lastIPs[ipptr] = ip
	ipptr++
}
