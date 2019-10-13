package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func connectionToBroker(ipAddress string, port string, name string) mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + ipAddress + ":" + port).SetClientID(name)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
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
	temperature := strconv.Itoa(rand.Intn(35 - 15) + 15)
	c.Publish("devices/envia", 0, false, "temperatureSensor;" + temperature + "°C;")
}

func identifyYourself(c mqtt.Client) {
	ipAddress := getIp()
	c.Publish("devices/envia", 0, false, "temperatureSensor;"+ ipAddress + ";")
}

func main() {
	// CONEXÃO COM O BROKER MQTT
	c := connectionToBroker("3.15.205.236", "1883", "temperatureSensor")
	if token := c.Subscribe("devices/recebe", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	} else {
		identifyYourself(c)
	}

	// LAÇO INFINITO DE FUNCIONAMENTO
	for {
		// De tempos em tempos ele envia sua temperatura
		go sendTemperature(c)
		time.Sleep(30 * time.Second)



	}

	//text := fmt.Sprintf("LA")
	//token := c.Publish("devices/recebe", 0, false, text)
	//token.Wait()
	//
	//if token := c.Unsubscribe("devices/recebe"); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}
	//
	//c.Disconnect(250)
	//
	//time.Sleep(1 * time.Second)

}
