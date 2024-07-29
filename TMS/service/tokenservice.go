package service

import (
	"context"
	"fmt"
	"time"
	"tmsservice/proto"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("jwt-secret-key")

type TokenServer struct {
	proto.UnimplementedTokenServiceServer
}

func CreateToken(username string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"userid":   id,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func IsTokenVerified(tokenString string) (jwt.MapClaims,error){
	token,err:=jwt.Parse(tokenString,func(t *jwt.Token) (interface{}, error) {
		return jwtSecret,nil
	})
	if err!=nil{
		return nil,err
	}
	if !token.Valid{
		return nil, fmt.Errorf("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok{
		return claims, nil
	}
	return nil,fmt.Errorf("invalid token")
}

func (s *TokenServer) GetToken(ctx context.Context, req *proto.TokenRequest) (*proto.TokenResponse, error) {
	fmt.Println("get request from client : ", req)
	
	token,err:=CreateToken(req.Username,int(req.UserId))
	
	if err!=nil{
		fmt.Println("failed to create token")
		return &proto.TokenResponse{},err
	}
	
	return &proto.TokenResponse{
		JwtToken:token,
	}, nil
}

func (s *TokenServer) VerifyToken(ctx context.Context, req *proto.VerifyRequest) (*proto.VerifedTokenResponse, error) {
	fmt.Println("got token to verify: ",req.JwtToken)
	
	claims,err:=IsTokenVerified(req.JwtToken);
	
	if err!=nil{
		return nil,err
	}

	username,_:=claims["username"].(string)
	id,_:=claims["userid"].(float64)

	fmt.Println(" user id is :  ",int64(id))
	return &proto.VerifedTokenResponse{
		Username: username,
		UserId: int64(id),
	},nil

}