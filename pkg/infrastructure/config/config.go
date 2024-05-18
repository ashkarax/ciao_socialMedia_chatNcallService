package config_postNrelSvc

import "github.com/spf13/viper"

type PortManager struct {
	RunnerPort string `mapstructure:"PORTNO"`
	AuthSvcUrl string `mapstructure:"AUTH_SVC_URL"`
}

type DataBase struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBName     string `mapstructure:"DBNAME"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost     string `mapstructure:"DBHOST"`
	DBPort     string `mapstructure:"DBPORT"`
}

type AWS struct {
	Region     string `mapstructure:"AWS_REGION"`
	AccessKey  string `mapstructure:"AWS_ACCESS_KEY_ID"`
	SecrectKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	Endpoint   string `mapstructure:"AWS_ENDPOINT"`
}

type Config struct {
	PortMngr PortManager
	DB       DataBase
	AwsS3    AWS
}

func LoadConfig() (*Config, error) {
	var portmngr PortManager
	var db DataBase
	var awsS3 AWS

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
	err = viper.Unmarshal(&db)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&awsS3)
	if err != nil {
		return nil, err
	}

	config := Config{PortMngr: portmngr, DB: db, AwsS3: awsS3}
	return &config, nil

}
