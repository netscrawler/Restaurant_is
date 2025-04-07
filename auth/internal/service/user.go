package service

import "github.com/netscrawler/Restaurant_is/auth/internal/repository"

type UserService struct {
	clientRepo repository.ClientRepository
	staffRepo  repository.StaffRepository
}
