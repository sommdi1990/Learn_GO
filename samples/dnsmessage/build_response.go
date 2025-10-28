package main

import (
	"fmt"

	dns "golang.org/x/net/dns/dnsmessage"
)

// ساخت پاسخ DNS با یک رکورد A و رکورد TXT اضافی.
func main() {
	var b dns.Builder
	h := dns.Header{Response: true}
	b = dns.NewBuilder(nil, h)

	// Answers
	_ = b.StartAnswers()
	name, _ := dns.NewName("example.com.")
	_ = b.AResource(dns.ResourceHeader{Name: name, TTL: 120}, dns.AResource{A: [4]byte{1, 2, 3, 4}})

	// Additionals (مثلاً TXT)
	_ = b.StartAdditionals()
	_ = b.TXTResource(dns.ResourceHeader{Name: name, TTL: 60}, dns.TXTResource{TXT: []string{"hello", "world"}})

	pkt, err := b.Finish()
	if err != nil {
		panic(err)
	}
	fmt.Printf("response bytes: % x\n", pkt)
}
