package di_chatNcallSvc

import (
	"fmt"

	client_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/client"
	config_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/config"
	db_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/db"
	server_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/server"
	repository_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/repository"
	usecase_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/usecase"
)

func InitializeChatNCallSvc(config *config_chatNcallSvc.Config) (*server_chatNcallSvc.ChatNCallSvc, error) {

	DB, err := db_chatNcallSvc.ConnectDatabaseMongo(&config.MongoDB)
	if err != nil {
		fmt.Println("ERROR CONNECTING DB FROM DI.GO")
		return nil, err
	}

	client, err := client_chatNcallSvc.InitAuthServiceClient(&config.PortMngr)
	if err != nil {
		fmt.Println("ERROR SETTING-UP AUTH CLIENT")
	}

	chatRepo := repository_chatNcallSvc.NewCharRepo(DB)
	chatUseCase := usecase_chatNcallSvc.NewChatUseCase(chatRepo, client, &config.Kafka)

	go chatUseCase.KafkaMessageConsumer()

	return server_chatNcallSvc.NewChatNCallServiceServer(chatUseCase), nil
}
