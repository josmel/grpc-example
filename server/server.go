package main

import (
	"flag"
	"log"
	"net"
	"sort"
	"sync"

	"github.com/boltdb/bolt"

	//pb "../proto"
	pb "github.com/sikang99/grpc-example/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type customerService struct {
	customers  map[int]*pb.Person
	id         int // unique id number to assign
	db         *bolt.DB
	sync.Mutex // later change to sync.RWMutex
}

func NewCustomerService() *customerService {
	cs := customerService{
		id:        1, // start id number
		customers: make(map[int]*pb.Person),
	}
	cs.db, _ = pb.BoltOpen("person.db", nil)
	return &cs
}

func (cs *customerService) ListPersons(p *pb.RequestType, stream pb.CustomerService_ListPersonsServer) error {
	log.Printf("list (%v)\n", p.Person)

	cs.Lock()
	defer cs.Unlock()

	// check the condition
	var age int
	if p.Person != nil {
		age = int(p.Person.Age)
	}

	var keys []int
	for k := range cs.customers {
		keys = append(keys, k)
	}

	//fmt.Println(keys)
	sort.Sort(sort.IntSlice(keys))
	//fmt.Println(keys)

	for _, key := range keys {
		ps := cs.customers[key]

		// conditional listing
		if age > 0 && int(ps.Age) != age {
			continue
		}
		if err := stream.Send(ps); err != nil {
			return err
		}
	}

	return nil
}

func (cs *customerService) PurgePersons(p *pb.RequestType, stream pb.CustomerService_PurgePersonsServer) error {
	log.Printf("purge (%v)\n", p.Person)

	cs.Lock()
	defer cs.Unlock()

	// check the condition
	var age int
	if p.Person != nil {
		age = int(p.Person.Age)
	}

	for _, ps := range cs.customers {
		// conditional listing
		if age > 0 && int(ps.Age) == age {
			if err := stream.Send(ps); err != nil {
				return err
			}
			delete(cs.customers, int(ps.Id))
		}
	}

	return nil
}

func (cs *customerService) GetPerson(c context.Context, p *pb.Person) (*pb.ResponseType, error) {
	log.Printf("get (%d)\n", p.Id)
	resp := new(pb.ResponseType)

	cs.Lock()
	defer cs.Unlock()

	ps, ok := cs.customers[int(p.Id)]
	if ok {
		resp.Person = ps
	}

	log.Printf("%v\n", resp.Person)
	return resp, nil
}

func (cs *customerService) AddPerson(c context.Context, p *pb.Person) (*pb.ResponseType, error) {
	log.Printf("add (%d)\n", cs.id)
	resp := new(pb.ResponseType)

	cs.Lock()
	defer cs.Unlock()

	p.Id = int32(cs.id)
	cs.customers[int(p.Id)] = p
	cs.id++

	ps, ok := cs.customers[int(p.Id)]
	if ok {
		resp.Person = ps
	}

	log.Printf("%v\n", resp.Person)
	return resp, nil
}

func (cs *customerService) DeletePerson(c context.Context, p *pb.Person) (*pb.ResponseType, error) {
	log.Printf("delete (%d)\n", p.Id)
	resp := new(pb.ResponseType)

	cs.Lock()
	defer cs.Unlock()

	ps, ok := cs.customers[int(p.Id)]
	if ok {
		resp.Person = ps
		delete(cs.customers, int(p.Id))
	}

	log.Printf("%v\n", resp.Person)
	return resp, nil
}

func (cs *customerService) UpdatePerson(c context.Context, p *pb.Person) (*pb.ResponseType, error) {
	log.Printf("update (%d)\n", p.Id)
	resp := new(pb.ResponseType)

	cs.Lock()
	defer cs.Unlock()

	_, ok := cs.customers[int(p.Id)]
	if ok {
		cs.customers[int(p.Id)] = p
		ps, ok := cs.customers[int(p.Id)]
		if ok {
			resp.Person = ps
		}
	}

	log.Printf("%v\n", resp.Person)
	return resp, nil
}

func init() {
	_ = flag.String("-db", "person.db", "BoltDB name")
}

// server function
func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", ":11111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()

	pb.RegisterCustomerServiceServer(server, NewCustomerService())
	server.Serve(lis)
}
