package persongrpc

import (
	"context"
	"errors"
	"fmt"

	"memcache/pkg/api"
	"memcache/pkg/proto/person"
	"memcache/pkg/service"
)

type server struct {
	person.UnsafePersonServiceServer
	service service.Service
}

func NewServer(service service.Service) (*server, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}
	return &server{
		service: service,
	}, nil
}

func (s *server) SetPerson(ctx context.Context, req *person.SetPersonRequest) (*person.SetPersonResponse, error) {
	fmt.Printf("SetPerson function was invoked with %v\n", req)

	resp, err := s.service.SetPerson(ctx, &api.SetPersonRequest{
		Person: api.Person{
			Name: req.GetPerson().Name,
			Age:  req.GetPerson().GetAge(),
		},
		TTL: req.Ttl,
	})
	if err != nil {
		return nil, err
	}

	return &person.SetPersonResponse{
		Id: resp.ID,
	}, nil
}

func (s *server) GetPerson(ctx context.Context, req *person.GetPersonRequest) (*person.GetPersonResponse, error) {
	fmt.Printf("GetPerson function was invoked with %v\n", req)

	resp, err := s.service.GetPerson(ctx, &api.GetPersonRequest{ID: req.Id})
	if err != nil {
		return nil, err
	}

	return &person.GetPersonResponse{
		Person: &person.Person{
			Name: resp.Person.Name,
			Age:  resp.Person.Age,
		},
	}, nil
}

func (s *server) DeletePerson(ctx context.Context, req *person.DeletePersonRequest) (*person.DeletePersonResponse, error) {
	fmt.Printf("DeletePerson function was invoked with %v\n", req)

	s.service.DeletePerson(ctx, &api.DeletePersonRequest{ID: req.Id})

	return &person.DeletePersonResponse{}, nil
}
