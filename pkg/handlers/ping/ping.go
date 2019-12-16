package ping

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

type PingConfig struct {
	Target string
}

func (pingConfig *PingConfig) Run(ctx context.Context) bool {
	//pinger.SetPrivileged(true)
	recvChan := make(chan bool, 2)
	p := fastping.NewPinger()
	p.Network("udp")
	ra, err := net.ResolveIPAddr("ip4:icmp", pingConfig.Target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		recvChan <- true
	}
	p.OnIdle = func() {
		recvChan <- false
	}
	p.MaxRTT = time.Millisecond * 1000
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
	select {
	case result := <-recvChan:
		return result
	case <-ctx.Done():
		return false
	}
}
