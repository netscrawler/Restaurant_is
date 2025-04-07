package models

import (
	"github.com/golang-jwt/jwt/v5"
	pb "github.com/netscrawler/RispProtos/proto/gen/go/v1/auth"
)

type Claims struct {
	UserID   string
	UserType string // "client" или "staff"
	Roles    []pb.Role
	jwt.RegisteredClaims
}
