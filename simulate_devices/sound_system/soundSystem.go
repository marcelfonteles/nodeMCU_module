package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net"
	"os"
	_"runtime"
	"strconv"
	_"strings"
	"time"
)

var currentStatus bool = false
var volume int = 0
var identified bool = false
var info mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if string(msg.Payload()) == "INFO" {
		identifyYourself(client)
	}
}

var actions mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if string(msg.Payload()) == "turnOn" {
		if !currentStatus {
			currentStatus = true
			client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;turning on the sound system;")
		} else {
			client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;already on;")
		}

	} else if string(msg.Payload()) == "turnOff" {
		if !currentStatus {
			client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;already off;")
		} else {
			currentStatus = false
			client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;turning off the sound system;")
		}
	} else if string(msg.Payload()) == "increaseVolume" {
		if currentStatus {
			if volume <= 10 {
				volume += 1
				client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;increased volume;volume:" + strconv.Itoa(volume) + ";")
			} else {
				client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;max volume;volume:" + strconv.Itoa(volume) + ";")
			}
		} else {
			client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;the sound system is off;")
		}
	} else if string(msg.Payload()) == "decreaseVolume" {
		if currentStatus {
			if volume > 0 {
				volume -= 1
				client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;decreased volume;volume:" + strconv.Itoa(volume) + ";")
			} else {
				client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;min volume;volume:" + strconv.Itoa(volume) + ";")
			}
		} else {
			client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;the sound system is off;")
		}
	} else {
		client.Publish("devices/soundSystem/envia", 0, false, "soundSystem;wrong command;")
	}
}

func connectionToBroker(ipAddress string, port string, name string) (mqtt.Client) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + ipAddress + ":" + port).SetClientID(name)
	opts.SetKeepAlive(3600 * time.Minute)
	//opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

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

func sendStatus(c mqtt.Client) {
	if currentStatus {
		c.Publish("devices/envia", 0, false, "soundSystem;On;Playing Peppa Pig;volume:" + strconv.Itoa(volume) + ";")
	} else {
		c.Publish("devices/envia", 0, false, "soundSystem;Off;")
	}
	time.Sleep(15 * time.Second)
}

func identifyYourself(c mqtt.Client) {
	ipAddress := getIp()
	c.Publish("devices/envia", 0, false, "soundSystem;"+ ipAddress + ";")
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
	c := connectionToBroker("3.15.205.236", "1883", "soundSystem")
	subscribe(c, "devices/recebe", info)
	subscribe(c, "devices/soundSystem/recebe", actions)

	// LAÇO INFINITO DE FUNCIONAMENTO
	for {
		time.Sleep(2 * time.Second)
	}
}
