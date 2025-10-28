package main

import (
	"fmt"

	dns "golang.org/x/net/dns/dnsmessage"
)

// پارس پاسخ شبیه‌سازی‌شده (برای مثال).
func main() {
	// برای سادگی، ابتدا یک پاسخ ساختگی می‌سازیم و سپس آن را پارس می‌کنیم.
	var b dns.Builder
	h := dns.Header{Response: true}
	b = dns.NewBuilder(nil, h)
	_ = b.StartAnswers()
	name, _ := dns.NewName("example.com.")
	rh := dns.ResourceHeader{Name: name, TTL: 60}
	_ = b.AResource(rh, dns.AResource{A: [4]byte{93, 184, 216, 34}}) // 93.184.216.34
	pkt, _ := b.Finish()

	var p dns.Parser
	if _, err := p.Start(pkt); err != nil {
		panic(err)
	}
	answers, err := p.AllAnswers()
	if err != nil {
		panic(err)
	}
	for _, a := range answers {
		switch body := a.Body.(type) {
		case *dns.AResource:
			ip := body.A
			fmt.Printf("A: %d.%d.%d.%d TTL=%d\n", ip[0], ip[1], ip[2], ip[3], a.Header.TTL)
		default:
			fmt.Printf("other rr: %T\n", body)
		}
	}
}
