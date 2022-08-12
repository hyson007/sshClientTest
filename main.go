package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/melbahja/goph"
)

type result struct {
	ip  string
	out []byte
	err error
}

func main() {
	ph := os.Getenv("passphrase")

	// Start new ssh connection with private key.
	auth, err := goph.Key("/Users/jackyao/.ssh/id_rsa", ph)
	if err != nil {
		log.Fatal(err)
	}

	ipaddress := []string{"13.229.155.44", "13.214.204.163"}

	ch := make(chan result, len(ipaddress))
	now := time.Now()

	wg := &sync.WaitGroup{}
	wg.Add(len(ipaddress))
	for _, ip := range ipaddress {

		go func(ip string) {
			client, err := goph.New("ec2-user", ip, auth)
			if err != nil {
				log.Fatal(err)
			}

			// Defer closing the network connection.
			defer client.Close()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
			defer cancel()

			// Execute your command.
			out, err := client.RunContext(ctx, "timeout 5 ping 8.8.8.8 -i 0.2")
			ch <- result{ip, out, err}
			wg.Done()
		}(ip)
	}

	wg.Wait()
	close(ch)

	for c := range ch {
		log.Println("ip:", c.ip, "out:", string(c.out), "err:", c.err)
		log.Println("----------------------------------------------------")
	}
	log.Println("time:", time.Since(now))

}
