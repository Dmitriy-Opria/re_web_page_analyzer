package service

import "context"

type StaffRepository interface {
	Health(ctx context.Context) error
}

type StaffService struct{}

func NewStaffService() *StaffService {
	return &StaffService{}
}

func (s StaffService) Health(ctx context.Context) error {
	panic("implement me")
}
