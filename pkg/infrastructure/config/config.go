package config_chatNcallSvc

import "github.com/spf13/viper"

type PortManager struct {
	RunnerPort string `mapstructure:"PORTNO"`
	AuthSvcUrl string `mapstructure:"AUTH_SVC_URL"`
}

type MongoDataBase struct {
	MongoDbURL    string `mapstructure:"MONGODB_URL"`
	DataBase      string `mapstructure:"MONGODB_DATABASE"`
	MongoUsername string `mapstructure:"MONGODB_USERNAME"`
	MongoPassword string `mapstructure:"MONGODB_PASSWORD"`
}

type ApacheKafka struct {
	KafkaPort           string `mapstructure:"KAFKA_PORT"`
	KafkaTopicOneToOne  string `mapstructure:"KAFKA_TOPIC_1"`
	KafkaTopicOneToMany string `mapstructure:"KAFKA_TOPIC_2"`
}

type Config struct {
	PortMngr PortManager
	MongoDB  MongoDataBase
	Kafka    ApacheKafka
}

func LoadConfig() (*Config, error) {
	var portmngr PortManager
	var MongoDb MongoDataBase
	var kafka ApacheKafka

	viper.AddConfigPath("./")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&portmngr)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&MongoDb)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&kafka)
	if err != nil {
		return nil, err
	}

	config := Config{PortMngr: portmngr, MongoDB: MongoDb, Kafka: kafka}
	return &config, nil

}
