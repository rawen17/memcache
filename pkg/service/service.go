package service

import (
	"context"
	"errors"
	"time"

	"memcache/pkg/api"

	uuid "github.com/satori/go.uuid"
)

type Service interface {
	SetPerson(ctx context.Context, req *api.SetPersonRequest) (api.SetPersonResponse, error)
	GetPerson(ctx context.Context, req *api.GetPersonRequest) (api.GetPersonResponse, error)
	DeletePerson(ctx context.Context, req *api.DeletePersonRequest) error
}

type storage interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type service struct {
	storage storage
}

func NewService(storage storage) (Service, error) {
	if storage == nil {
		return nil, errors.New("storage is nil")
	}
	return &service{
		storage: storage,
	}, nil
}

func (s *service) SetPerson(ctx context.Context, req *api.SetPersonRequest) (api.SetPersonResponse, error) {
	id := uuid.NewV4().String()
	s.storage.Set(id, req.Person, time.Duration(req.TTL))
	return api.SetPersonResponse{ID: id}, nil
}

func (s *service) GetPerson(ctx context.Context, req *api.GetPersonRequest) (api.GetPersonResponse, error) {
	data, isExist := s.storage.Get(req.ID)
	if !isExist {
		return api.GetPersonResponse{}, errors.New("person is not exist")
	}

	person, ok := data.(api.Person)
	if !ok {
		return api.GetPersonResponse{}, errors.New("data in storage has different type")
	}

	return api.GetPersonResponse{Person: person}, nil
}

func (s *service) DeletePerson(ctx context.Context, req *api.DeletePersonRequest) error {
	s.storage.Delete(req.ID)
	return nil
}
