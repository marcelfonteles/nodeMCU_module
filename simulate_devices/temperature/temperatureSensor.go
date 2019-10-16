package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var currentTemperature int = 23
var identified bool = false

var maxTemp int = 25
var minTemp int = 16

var info mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if string(msg.Payload()) == "INFO" {
		identifyYourself(client)
	}
}

var actions mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	command := string(msg.Payload())
	index := strings.Index(command, ":")
	if index != -1 {
		if command[:index] == "up" {
			increase, err := strconv.Atoi(command[index+1:])
			if err != nil {
				fmt.Println(err)
			}
			if currentTemperature + increase > maxTemp {
				currentTemperature = 25
				client.Publish("devices/temperature/envia", 0, false, "temperatureSensor;max temp;25°C;")

			} else {
				currentTemperature += increase
				client.Publish("devices/temperature/envia", 0, false, "temperatureSensor;"+ strconv.Itoa(currentTemperature) +"°C;")
			}

		} else if command[:index] == "down" {
			decrease, err := strconv.Atoi(command[index+1:])
			if err != nil {
				fmt.Println(err)
			}
			if currentTemperature - decrease < minTemp {
				currentTemperature = 16
				client.Publish("devices/temperature/envia", 0, false, "temperatureSensor;min temp;16°C;")

			} else {
				currentTemperature -= decrease
				client.Publish("devices/temperature/envia", 0, false, "temperatureSensor;"+ strconv.Itoa(currentTemperature) +"°C;")
			}
		}
	}
}

func connectionToBroker(ipAddress string, port string, name string) (mqtt.Client) {
	//opts := mqtt.NewClientOptions().AddBroker("tcp://soldier.cloudmqtt.com:10290").SetClientID(name)
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + ipAddress + ":" + port).SetClientID(name)
	opts.SetKeepAlive(3600 * time.Minute)
	//opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	//opts.SetUsername("slhhlhba")
	//opts.SetPassword("cX2o65sDKOK_")

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return c
}

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "err_ip"
}

func sendTemperature(c mqtt.Client) {
	temperature := strconv.Itoa(currentTemperature)
	c.Publish("devices/temperatureSensor/envia", 0, false, "temperatureSensor;" + temperature + "°C;")
	time.Sleep(15 * time.Second)
}

func identifyYourself(c mqtt.Client) {
	ipAddress := getIp()
	c.Publish("devices/envia", 0, false, "temperatureSensor;"+ ipAddress + ";")
}

func subscribe(c mqtt.Client, topic string, mHan mqtt.MessageHandler) mqtt.Client {
	if token := c.Subscribe(topic, 2, mHan); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	} else {
		if !identified {
			identifyYourself(c)
			identified = true
		}
	}

	return c
}

func main() {
	// CONEXÃO COM O BROKER MQTT
	c := connectionToBroker("3.15.205.236", "1883", "temperatureSensor")
	subscribe(c, "devices/recebe", info)
	subscribe(c, "devices/temperature/recebe", actions)

	// LAÇO INFINITO DE FUNCIONAMENTO
	for {
		// 7 go routines até aqui
		if runtime.NumGoroutine() == 7 {
			// De tempos em tempos ele envia sua temperatura
			go sendTemperature(c)
		}
		time.Sleep(2 * time.Second)
	}
}
