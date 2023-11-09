package config

import "github.com/spf13/viper"

type Config struct {
	GrpcAddr  string `mapstructure:"GRPC_ADDR"`
	HttpAddr  string `mapstructure:"HTTP_ADDR"`
	WorkDir   string `mapstructure:"WORK_DIR"`
	RootFsDir string `mapstructure:"ROOTFS_DIR"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
