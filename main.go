// Example SVG parser using a combination of xml.Unmarshal and the
// xml.Unmarshaler interface to handle an unknown combination of group
// elements where order is important.

package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/beevik/etree"
	"github.com/libesz/pingboard/pkg/config"
	"github.com/libesz/pingboard/pkg/svgmanip"
	"github.com/tatsushid/go-fastping"
	"gopkg.in/yaml.v2"
)

func main() {
	filename := os.Args[1]
	configSource, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	config := config.Config{}

	err = yaml.Unmarshal([]byte(configSource), &config)
	if err != nil {
		panic(err)
	}

	svg := etree.NewDocument()
	if err := svg.ReadFromFile(config.SvgPath); err != nil {
		panic(err)
	}
	svgmanip.UpdateDoc(svg, config)
	svg.WriteToFile("new.svg")
	//ping()
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
