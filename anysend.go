package main

import (
	"flag"
	"log"
	"net"
	"os"
)

func main() {
	var anycastIP *string = flag.String("a", "1.2.3.4", "anycastip")
	var targetIP *string = flag.String("t", "1.2.3.4", "targetIP")

	flag.Usage = func() {
		log.Printf("usage: %s [ -i interface ] [ -a anycastip ] [ -s snaplen ] [ -X ] [ expression ]\n", os.Args[0])
		os.Exit(1)
	}

	flag.Parse()

	if *anycastIP == "1.2.3.4" || *targetIP == "1.2.3.4" {
		flag.Usage()
	}

	raddr := &net.IPAddr{IP: net.ParseIP(*targetIP).To4()}
	laddr := &net.IPAddr{IP: net.ParseIP(*anycastIP)}

	con, err := net.DialIP("ip4:1", laddr, raddr)
	if err != nil {
		log.Fatalf("unable to make raw socket to dial out from, err was %s", err.Error())
	}

	// Now to hand craft a ICMP echo request packet!

	// pkt := p.Payload
	// icmp := new(Icmphdr)
	// icmp.Type = pkt[0]
	// icmp.Code = pkt[1]
	// icmp.Checksum = binary.BigEndian.Uint16(pkt[2:4])
	// icmp.Id = binary.BigEndian.Uint16(pkt[4:6])
	// icmp.Seq = binary.BigEndian.Uint16(pkt[6:8])
	// p.Payload = pkt[8:]
	// p.Headers = append(p.Headers, icmp)
	// return icmp

	payload := []byte("ANYCATCH")
	packet := make([]byte, 7+len(payload)) // 7 for the packet itself, 8 for the "ANYCATCH" string

	packet[0] = 8 // Type, in this case a echo request
	packet[1] = 0 // Code, in this case there is no sub code for this packet type
	/* packet[2-3] = checksum of packet, we will do this later when we are done */
	packet[4] = 69 // ICMP ID of the request, in this case I am just filling this in
	packet[5] = 69 // ICMP ID of the request
	packet[6] = 69 // ICMP Seq of the request, in this case I am just filling this in
	packet[7] = 69 // ICMP Seq of the request
	// ANYCATCH
	for i := 0; i < len(payload); i++ {
		packet[8+i] = payload[i]
	}

	con.Write(packet)
	log.Printf("Done.")

}
