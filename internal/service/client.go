package service

import (
	"context"
	"net/http"
)

type ClientRepository interface {
	Parse(ctx context.Context) error
}

type ClientService struct {
	client *http.Client
}

func NewClientService(client *http.Client) *ClientService {
	return &ClientService{
		client: client,
	}
}

func (p ClientService) Get(ctx context.Context) error {
	panic("implement me")
}
