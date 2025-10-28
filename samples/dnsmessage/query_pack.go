package main

import (
	"fmt"

	dns "golang.org/x/net/dns/dnsmessage"
)

// ساخت یک DNS Query برای A record example.com.
func main() {
	var b dns.Builder
	h := dns.Header{RecursionDesired: true}
	b = dns.NewBuilder(nil, h)
	_ = b.StartQuestions()
	name, _ := dns.NewName("example.com.")
	_ = b.Question(dns.Question{Name: name, Type: dns.TypeA, Class: dns.ClassINET})
	msg, err := b.Finish()
	if err != nil {
		panic(err)
	}
	fmt.Printf("packed query bytes: % x\n", msg)
}
