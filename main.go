package main

import (
	"flag"
	"log"
	"os"
)

var AnycastIP string = ""

func main() {

	var device *string = flag.String("i", "", "interface")
	var snaplen *int = flag.Int("s", 65535, "snaplen")
	var anycastIP *string = flag.String("a", "1.2.3.4", "anycastip")
	var password *string = flag.String("p", "wow", "password for http interface")

	flag.Usage = func() {
		log.Printf("usage: %s [ -i interface ] [ -a anycastip ] [ -s snaplen ] [ -X ] [ expression ]\n", os.Args[0])
		os.Exit(1)
	}

	flag.Parse()
	AnycastIP = *anycastIP

	go StartListeningForPings(*device, *anycastIP, *snaplen)
	StartServer(*password)

}
