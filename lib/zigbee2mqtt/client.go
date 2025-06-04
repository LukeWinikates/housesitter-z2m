package zigbee2mqtt

import (
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt/devices"
	"LukeWinikates/january-twenty-five/lib/zigbee2mqtt/payloads"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

func init() {
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

type Client interface {
	SetDeviceState(getenv string, message devices.LightControl) error
	DeviceUpdates() (chan payloads.MessagePayload, chan error)
}

type RealClient struct {
	mqttClient mqtt.Client
}

func (c *RealClient) DeviceUpdates() (chan payloads.MessagePayload, chan error) {
	topic := "zigbee2mqtt/bridge/devices"
	payloadsChannel := make(chan payloads.MessagePayload)
	errorsChannel := make(chan error)

	fmt.Println("getting messages")
	token := c.mqttClient.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("Received message from topic: %s\n", message.Topic())
		deviceList, err := payloads.Parse(message.Payload())
		if err != nil {
			errorsChannel <- err
		}
		for _, device := range deviceList {
			payloadsChannel <- device
		}
	})
	go func() {
		token.Wait()
	}()

	return payloadsChannel, errorsChannel
}

func (c *RealClient) SetDeviceState(deviceName string, message devices.LightControl) error {
	var payloadBytes, err = json.Marshal(message)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("zigbee2mqtt/%s/set", deviceName)
	fmt.Printf("sending message for topic: %s\n", topic)
	publish := c.mqttClient.Publish(topic, 0, false, payloadBytes)
	if publish.Error() != nil {
		return publish.Error()
	}
	return nil
}

func NewClient(mqttHost, clientID string) Client {
	options := mqtt.NewClientOptions()
	options.AddBroker(mqttHost)
	fmt.Printf("client id: %s\n", clientID)
	options.SetClientID(clientID)
	client := mqtt.NewClient(options)

	options.OnConnect = connectHandler
	options.OnConnectionLost = connectLostHandler
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return &RealClient{
		mqttClient: client,
	}
}
