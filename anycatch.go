package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/akrennmair/gopcap"
	"log"
	"net"
	"os"
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

var out *bufio.Writer
var errout *bufio.Writer

func main() {

	var device *string = flag.String("i", "", "interface")
	var snaplen *int = flag.Int("s", 65535, "snaplen")
	var anycastIP *string = flag.String("a", "1.2.3.4", "anycastip")
	expr := ""

	out = bufio.NewWriter(os.Stdout)
	errout = bufio.NewWriter(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(errout, "usage: %s [ -i interface ] [ -a anycastip ] [ -s snaplen ] [ -X ] [ expression ]\n", os.Args[0])
		os.Exit(1)
	}

	flag.Parse()

	if len(flag.Args()) > 0 {
		expr = flag.Arg(0)
	}

	var incomingIP net.IP = net.ParseIP(*anycastIP)

	if incomingIP == nil || *anycastIP == "1.2.3.4" {
		fmt.Fprintf(errout, "Incorrect Anycast IP given")
		flag.Usage()
	}

	if *device == "" {
		devs, err := pcap.Findalldevs()
		if err != nil {
			fmt.Fprintf(errout, "tcpdump: couldn't find any devices: %s\n", err)
		}
		if 0 == len(devs) {
			flag.Usage()
		}
		*device = devs[0].Name
	}

	h, err := pcap.Openlive(*device, int32(*snaplen), true, 0)
	if h == nil {
		fmt.Fprintf(errout, "tcpdump: %s\n", err)
		errout.Flush()
		return
	}
	defer h.Close()

	if expr != "" {
		ferr := h.Setfilter(expr)
		if ferr != nil {
			fmt.Fprintf(out, "tcpdump: %s\n", ferr)
			out.Flush()
		}
	}

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
						}
					case *pcap.Iphdr:
						//log.Printf("What(%d) ICMP! %d %d %d %d %d", level, header.Type, header.Code, header.Checksum, header.Id, header.Seq)
					default:
						log.Printf("Ahem %s ", header)
					}
				}
			}
		}
		out.Flush()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
