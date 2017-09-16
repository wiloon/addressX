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
	"github.com/wiloon/app-config"
)

const (
	duration = time.Second * 5
	url      = "http://members.3322.org/dyndns/getip"
)

var server_address = config.GetStringWithDefaultValue("server_address", "localhost:7000")

func main() {
	server_address = "localhost:7000"
	ticker := time.NewTicker(duration)
	var wg sync.WaitGroup
	wg.Add(1)

	for range ticker.C {

		log.Printf("ticked at %v", time.Now())
		ip := getIp()
		send(ip)
	}

	wg.Wait()
}

func send(ip string) {
	log.Println("sending ip:" + ip)

	conn, err := grpc.Dial(server_address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	defer func() {
		conn.Close()
		if r := recover(); r != nil {
			log.Println("error:", r)
		}
	}()

	c := pb.NewAddressClient(conn)

	log.Printf("ticked at %v", time.Now())
	reply, err := c.SetIp(context.Background(), &pb.AddressRequest{Ip: ip})
	log.Printf("success: %v", reply.Reply)

}

func getIp() string {

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
