// Example SVG parser using a combination of xml.Unmarshal and the
// xml.Unmarshaler interface to handle an unknown combination of group
// elements where order is important.

package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/svgmanip"
	"github.com/tatsushid/go-fastping"
)

func main() {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile("drawing2.svg"); err != nil {
		panic(err)
	}
	config := svgmanip.Configs{[]svgmanip.Config{{ID: "path10", Fill: "#00ff00"}}}
	svgmanip.UpdateDoc(doc, config)
	doc.WriteToFile("new.svg")
	ping()
}

func ping() {
	//pinger.SetPrivileged(true)
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	p.MaxRTT = time.Millisecond * 1000
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
	//time.Sleep(10 * time.Second)
}
