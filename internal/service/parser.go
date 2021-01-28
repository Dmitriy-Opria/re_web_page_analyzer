package service

import (
	"context"
	"net/http"
)

type ParserRepository interface {
	Parse(ctx context.Context) error
}

type ParserService struct{
	client             *http.Client
}

func NewParserService(client *http.Client,) *ParserService {
	return &ParserService{
		client:             client,
	}
}

func (p ParserService) Parse(ctx context.Context) error {
	panic("implement me")
}
