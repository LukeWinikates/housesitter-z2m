package main

import (
	"LukeWinikates/january-twenty-five/lib/server"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s, err := createServer()
	if err != nil {
		log.Fatal(err.Error())
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Printf("received signal: %s\n", sig.String())
		fmt.Println(s.Stop())
	}()

	fmt.Println("starting server")
	fmt.Println(s.Start())

}

func createServer() (server.Server, error) {
	client := zigbee2mqtt.NewClient(os.Getenv("MQTT_HOST"), os.Getenv("MQTT_CLIENT_ID"))
	options, err := createServerOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to set up with configuration: %s", err.Error())
	}
	db, err := gorm.Open(sqlite.Open(options.DataDir+"/test.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %s", err.Error())
	}
	s, err := server.New(db, client, *options)
	return s, err
}

func createServerOptions() (*server.Options, error) {
	dataPath := os.Getenv("DATA_PATH")
	hostname := os.Getenv("HOUSESITTER_HOST")
	if hostname == "" {
		hostname = "localhost:8998"
	}
	location, err := time.LoadLocation(os.Getenv("HOUSESITTER_TIME_ZONE"))
	if err != nil {
		return nil, err
	}
	options := &server.Options{
		DataDir:  dataPath,
		Hostname: hostname,
		Location: location,
	}
	return options, nil
}
