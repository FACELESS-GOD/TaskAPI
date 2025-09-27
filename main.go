package main

import (
	"TaskManager/Helper/Startup"
	"TaskManager/Package/Configurator"
)

func main() {
	config := Configurator.NewConfigurator()
	config.LoadConfig(Startup.DebugMode)

	err := config.LoadDBInstance()

	if err != nil {
		panic(err)
	}

}
