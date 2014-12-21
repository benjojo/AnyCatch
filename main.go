package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/akrennmair/gopcap"
	"github.com/oschwald/geoip2-golang"
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
	db, err := geoip2.Open("GeoIP2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var device *string = flag.String("i", "", "interface")
	var snaplen *int = flag.Int("s", 65535, "snaplen")
	expr := ""

	out = bufio.NewWriter(os.Stdout)
	errout = bufio.NewWriter(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(errout, "usage: %s [ -i interface ] [ -s snaplen ] [ -X ] [ expression ]\n", os.Args[0])
		os.Exit(1)
	}

	flag.Parse()

	if len(flag.Args()) > 0 {
		expr = flag.Arg(0)
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
			if pkt.IP.Protocol == 17 {
				if pkt.UDP.DestPort == 48563 {
					record, err := db.City(net.IP(pkt.IP.SrcIp))
					if err == nil {
						fmt.Fprintf(out, "%s | %f | %f\n", record.Country.Names["en"], record.Location.Latitude, record.Location.Longitude)
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
