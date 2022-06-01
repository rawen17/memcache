package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"memcache/pkg/proto/person"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I'm a client")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := person.NewPersonServiceClient(cc)
	fmt.Printf("Created client: %f", c)

	id := setPerson(c)
	time.Sleep(100 * time.Second)
	getPerson(c, id)
	deletePerson(c, id)
}

func setPerson(c person.PersonServiceClient) string {
	fmt.Println("Starting to do a SetPerson RPC...")
	req := &person.SetPersonRequest{
		Person: &person.Person{
			Name: "Dan",
			Age:  40,
		},
		Ttl: int64(2 * time.Second),
	}
	res, err := c.SetPerson(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling SetPerson RPC: %v", err)
	}
	log.Printf("Response from SetPerson: %v", res.Id)
	return res.Id
}

func getPerson(c person.PersonServiceClient, id string) {
	fmt.Println("Starting to do a GetPerson RPC...")
	req := &person.GetPersonRequest{
		Id: id,
	}
	res, err := c.GetPerson(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetPerson RPC: %v", err)
	}
	log.Printf("Response from GetPerson: %v", res.Person)
}

func deletePerson(c person.PersonServiceClient, id string) {
	fmt.Println("Starting to do a DeletePerson RPC...")
	req := &person.DeletePersonRequest{
		Id: id,
	}
	res, err := c.DeletePerson(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling DeletePerson RPC: %v", err)
	}
	log.Printf("Response from DeletePerson: %v", res)
}
