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

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	command := string(msg.Payload())
	index := strings.Index(command, ":")
	if index != -1 {
		if command[:index] == "up" {
			increase, err := strconv.Atoi(command[index+1:])
			if err != nil {
				fmt.Println(err)
			}
			currentTemperature += increase
		} else if command[:index] == "down" {
			decrease, err := strconv.Atoi(command[index+1:])
			if err != nil {
				fmt.Println(err)
			}
			currentTemperature -= decrease
		}
	} else if string(msg.Payload()) == "INFO" {
		identifyYourself(client)
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

func sendTemperature(c mqtt.Client) {
	temperature := strconv.Itoa(currentTemperature)
	c.Publish("devices/envia", 0, false, "temperatureSensor;" + temperature + "°C;")
	time.Sleep(15 * time.Second)
}

func identifyYourself(c mqtt.Client) {
	ipAddress := getIp()
	c.Publish("devices/envia", 0, false, "temperatureSensor;"+ ipAddress + ";")
}


func main() {
	// CONEXÃO COM O BROKER MQTT
	c := connectionToBroker("3.15.205.236", "1883", "temperatureSensor")
	if token := c.Subscribe("devices/recebe", 2, f); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	} else {
		identifyYourself(c)
	}

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
