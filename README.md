# Code for nodeMCU
- Project started using: https://www.filipeflop.com/blog/controle-monitoramento-iot-nodemcu-e-mqtt/

## About
- This project was built for my IoT course in college.

## Features
- Broker MQTT on AWS
- Turn on and Turn off some leds
- Topics:
```
devices/recebe -> nodeMCU listen this topic
devices/envia  -> nodeMCU send information in this topic
```     
- Messages:
1. nodeMCU
```
"LR"   -> nodeMCU turn on red led
"DR"   -> nodeMCU turn off red led
"LG"   -> nodeMCU turn on green led
"DG"   -> nodeMCU turn off green led
"INFO" -> nodeMCU send information about himself
```
2. Simulated Temparature Sensor
```
Receiver                 Topic                         Respond                 Topic
"INFO"                   devices/recebe                about information       devices/envia
"up:10"                  devices/temperature/recebe    increase temperature.   devices/temperature/envia
"down:10"                devices/temperature/recebe    decrease temperature.   devices/temperature/envia
```
3. Simulated Sound System
```
Receiver                 Topic                         Respond                 Topic
"INFO"                   devices/recebe                about information       devices/envia
"turnIn"                 devices/soundSystem/recebe    turn in                 devices/soundSystem/envia
"turnOff"                devices/soundSystem/recebe    turn off                devices/soundSystem/envia
"increaseVolume"         devices/soundSystem/recebe    increase volume         devices/soundSystem/envia
"decreaseVolume"         devices/soundSystem/recebe    decrease volume         devices/soundSystem/envia

```