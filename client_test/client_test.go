package client

import (
	"io"
	"testing"

	pb "github.com/sikang99/grpc-example/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var address string = "127.0.0.1:11111"

func TestBaseForCustomerService(t *testing.T) {
	// connect to rpc server
	conn, err := grpc.Dial(address)
	if err != nil {
		t.Errorf("connect error %v\n", err)
	}
	defer conn.Close()

	client := pb.NewCustomerServiceClient(conn)

	// test functions of CRUD: Add, Get, Update, and Delete

	// for Add
	person := &pb.Person{
		Name: "Stoney",
		Age:  502,
	}

	res, err := client.AddPerson(context.Background(), person)
	t.Logf("Add: %v\n", res.Person)

	// for Get
	person.Id = res.Person.Id

	res, err = client.GetPerson(context.Background(), person)
	t.Logf("Get: %v\n", res.Person)

	// for Update
	person.Name = "Mandoo"
	person.Age = 109

	res, err = client.UpdatePerson(context.Background(), person)
	t.Logf("Update: %v\n", res.Person)

	// for Delete
	res, err = client.DeletePerson(context.Background(), person)
	t.Logf("Delete: %v\n", res.Person)

	// check if the res is nil
	res, err = client.GetPerson(context.Background(), person)
	t.Logf("Get: %v\n", res.Person)

	// there should be nothing.
	if res.Person != nil {
		t.Errorf("result error %v\n", res.Person)
	}
}

func TestUtilsForCustomerService(t *testing.T) {
	// connect to rpc server
	conn, err := grpc.Dial(address)
	if err != nil {
		t.Errorf("connect error %v\n", err)
	}
	defer conn.Close()

	client := pb.NewCustomerServiceClient(conn)

	// for List
	req := &pb.RequestType{}

	stream, err := client.ListPersons(context.Background(), req)
	if err != nil {
		t.Errorf("request error %v\n", err)
	}
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("recv error %v\n", err)
		}
		t.Logf("List: %v\n", person)
	}
}

func TestLoadForCustomerService(t *testing.T) {
	// connect to rpc server
	conn, err := grpc.Dial(address)
	if err != nil {
		t.Errorf("connect error %v\n", err)
	}
	defer conn.Close()

	client := pb.NewCustomerServiceClient(conn)

	person := &pb.Person{
		Name: "Stoney",
		Age:  502,
	}

	for i := 0; i < 100; i++ {
		res, err := client.AddPerson(context.Background(), person)
		if err != nil {
			t.Errorf("Add: %v\n", err)
		}

		person.Id = res.Person.Id

		res, err = client.DeletePerson(context.Background(), person)
		if err != nil {
			t.Errorf("Delete: %v\n", err)
		}

		res, err = client.GetPerson(context.Background(), person)
		if err != nil {
			t.Errorf("Get: %v\n", err)
		}
	}
}
