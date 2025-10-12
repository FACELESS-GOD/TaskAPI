package main

import (
	"TaskManager/Helper/Startup"
	"TaskManager/Package/Configurator"
	"TaskManager/Package/Controller"
	"TaskManager/Package/Model"
	"log"
)

func main() {
	config := Configurator.NewConfigurator()
	config.LoadConfig(Startup.DebugMode)

	err := config.LoadDBInstance()

	if err != nil {
		log.Fatal(err)
	}

	mdl := Model.NewModel(*config)

	controller := Controller.NewController(&mdl)

	err = controller.StartServer(config.Address)

	if err != nil {
		log.Fatal(err)
	}

}
