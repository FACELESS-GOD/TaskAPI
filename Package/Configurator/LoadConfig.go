package Configurator

import (
	"TaskManager/Helper/Startup"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type ConfiguratorInterface interface {
	LoadConfig()
	LoadDBInstance()
}

type configParser struct {
	DbDriver     string `mapstructure="DBDRIVER"`
	DbConnString string `mapstructure="DBCONNSTRING"`
	Address      string `mapstructure="ADDRESS"`
}

type ConfiguratorStruct struct {
	DbDriver     string
	DbConnString string
	SqlDBConn    *sql.DB
	DbCtx        context.Context
	DbCancelFunc context.CancelFunc
	Address      string
}

func NewConfigurator() *ConfiguratorStruct {
	return &ConfiguratorStruct{}
}

func (Conf *ConfiguratorStruct) LoadConfig(Mode int) {

	switch Mode {
	case Startup.DebugMode:

		var configParser configParser

		viper.AddConfigPath(".")
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		//viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}

		err = viper.Unmarshal(&configParser)

		if err != nil {
			log.Fatal(err)
		}

		Conf.DbDriver = configParser.DbDriver
		Conf.DbConnString = configParser.DbConnString
		Conf.Address = configParser.Address

	case Startup.QAMode:

	case Startup.QAMode:

	}
}

func (Conf *ConfiguratorStruct) LoadDBInstance() error {
	if len(Conf.DbConnString) <= 0 {
		return errors.New("Configuration is not Loaded Properly. Please run LoadConfig() before executing this function")
	}

	db, err := sql.Open(Conf.DbDriver, Conf.DbConnString)

	if err != nil {
		return err
	}

	Conf.SqlDBConn = db

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*2)

	Conf.DbCtx = ctx

	Conf.DbCancelFunc = cancelFunc

	return nil
}
