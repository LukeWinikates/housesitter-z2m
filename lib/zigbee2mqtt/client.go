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
	SubscribeDeviceCatalog(func(devices []payloads.MessagePayload))
	SetDeviceState(getenv string, message devices.LightControl) error
}

type RealClient struct {
	mqttClient mqtt.Client
}

func (c *RealClient) SubscribeDeviceCatalog(f func(devices []payloads.MessagePayload)) {
	topic := "zigbee2mqtt/bridge/devices"
	token := c.mqttClient.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("Received message from topic: %s\n", message.Topic())
		deviceList, err := payloads.Parse(message.Payload())
		if err != nil {
			log.Default().Printf("err: %s", err.Error())
			return
		}
		f(deviceList)
	})
	token.Wait()
}

func (c *RealClient) SetDeviceState(deviceName string, message devices.LightControl) error {
	var payloadBytes, err = json.Marshal(message)
	if err != nil {
		return err
	}

	topic := fmt.Sprintf("zigbee2mqtt/%s/set", deviceName)
	publish := c.mqttClient.Publish(topic, 0, false, payloadBytes)
	if publish.Error() != nil {
		return publish.Error()
	}
	return nil
}

func NewClient(mqttHost string) Client {
	options := mqtt.NewClientOptions()
	options.AddBroker(mqttHost)
	options.SetClientID("hogepiyo")
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
