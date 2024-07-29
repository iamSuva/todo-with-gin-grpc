package grpclient ///grpc client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"todowithgin/proto"
	"google.golang.org/grpc"
)

var (
	ErrorFailedTogenerateToken = errors.New("failed to generate token")
	ErrorInvalidToken          = errors.New("Invalid token")
)

type GrpcInterface interface {
	GetTokenHandler(string, int) (string, error)
	VeriFyTokenHandler(string) (*VerifiedDecodedToken, error)
}
type GrpcService struct{}

func NewGrpcService() *GrpcService {
	return &GrpcService{}
}

type VerifiedDecodedToken struct {
	Username string
	UserId   int
}
var GrpcConn *grpc.ClientConn
func GrpcConnection(){
	opt := grpc.WithInsecure()
	conn, err := grpc.Dial("localhost:4000", opt)
    
	if err != nil {
		log.Fatalf("failed to connect %v ", err)
	}

	GrpcConn=conn
}
func (g *GrpcService) GetTokenHandler(username string, id int) (string, error) {
	
	client := proto.NewTokenServiceClient(GrpcConn)

	res, err := client.GetToken(context.Background(), &proto.TokenRequest{
		UserId:   int64(id),
		Username: username,
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println("received res: ", res)
	return res.JwtToken, nil
}

func (g *GrpcService) VeriFyTokenHandler(token string) (*VerifiedDecodedToken, error) {
	
	client := proto.NewTokenServiceClient(GrpcConn)
	res, err := client.VerifyToken(context.Background(), &proto.VerifyRequest{
		JwtToken: token,
	})

	if err != nil {
		fmt.Println("token verify error: ", err)
		return nil, ErrorInvalidToken
	}

	return &VerifiedDecodedToken{
		Username: res.Username,
		UserId:   int(res.UserId),
	}, nil
}
