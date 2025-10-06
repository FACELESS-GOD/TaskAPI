package Configurator

import (
	"TaskManager/Helper/Startup"
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ConfiguratorInterface interface {
	LoadConfig()
	LoadDBInstance()
}

type ConfiguratorStruct struct {
	DbDriver     string
	DbConnString string
	SqlDBConn    *sql.DB
	DbCtx        context.Context
	DbCancelFunc context.CancelFunc
}

func NewConfigurator() *ConfiguratorStruct {
	return &ConfiguratorStruct{}
}

func (Conf *ConfiguratorStruct) LoadConfig(Mode int) {

	switch Mode {
	case Startup.DebugMode:

		Conf.DbDriver = "mysql"
		Conf.DbConnString = ""

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
