package main

import (
	"google.golang.org/grpc"
	"log"
	pb "wiloon.com/addressX/proto"
	"context"
	"os"
	"net/http"
	"io/ioutil"
	"time"
	"sync"
)

const address = "localhost:7000"

func main() {
	ticker := time.NewTicker(time.Second * 3)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for _ = range ticker.C {

			log.Printf("ticked at %v", time.Now())
			ip := getIp()
			send(ip)
		}

		wg.Done()
	}()

	wg.Wait()

}

func send(ip string) {
	log.Println("sending ip:" + ip)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewAddressClient(conn)

	log.Printf("ticked at %v", time.Now())
	reply, err := c.SetIp(context.Background(), &pb.AddressRequest{Ip: ip})
	log.Printf("success: %v", reply.Reply)

}
func getIp() string {
	url := "http://members.3322.org/dyndns/getip"

	resp, err := http.Get(url)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	ip := string(body)

	log.Println("ip:", ip)

	return ip
}
