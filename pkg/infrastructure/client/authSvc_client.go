package client_chatNcallSvc

import (
	"fmt"

	config_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/config"
	"github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAuthServiceClient(config *config_chatNcallSvc.PortManager) (*pb.AuthServiceClient, error) {
	cc, err := grpc.NewClient(config.AuthSvcUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("-------", err)
		return nil, err
	}

	Client := pb.NewAuthServiceClient(cc)

	return &Client, nil
}
