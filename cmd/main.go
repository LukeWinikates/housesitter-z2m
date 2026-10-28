package main

import (
	"LukeWinikates/january-twenty-five/lib/database"
	"LukeWinikates/january-twenty-five/lib/server"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	mqttHost := os.Getenv("MQTT_HOST")
	clientID := os.Getenv("MQTT_CLIENT_ID")
	noOpFlag := os.Getenv("MQTT_NO_OP_CLIENT") != ""
	if mqttHost == "" || clientID == "" {
		log.Fatal("missing required MQTT_HOST or MQTT_CLIENT_ID environment variables")
	}
	s, err := createServer(noOpFlag, mqttHost, clientID)
	if err != nil {
		log.Fatal(err.Error())
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("received signal: %s\n", sig.String())
		log.Println(s.Stop())
	}()

	log.Println("starting server")
	log.Println(s.Start())

}

func createServer(noOpFlag bool, mqttHost, mqttClientID string) (server.Server, error) {
	client := createMQTTClient(noOpFlag, mqttHost, mqttClientID)
	options, err := createServerOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to set up with configuration: %s", err.Error())
	}
	db, err := gorm.Open(sqlite.Open(options.DataDir+"/test.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database (%s): %s", options.DataDir, err.Error())
	}
	if err = database.AutoMigrate(db); err != nil {
		return nil, err
	}
	s, err := server.New(db, client, *options)
	return s, err
}

func createMQTTClient(noOpFlag bool, mqttHost, mqttClientID string) zigbee2mqtt.Client {
	if noOpFlag {
		log.Println("Using No-op MQTT client")
		return zigbee2mqtt.NoOpClient()
	}
	log.Printf("Using real MQTT client on %s\n\n", mqttHost)
	return zigbee2mqtt.NewClient(mqttHost, mqttClientID)
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
