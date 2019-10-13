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
"up:10"   -> increase temperature. Value specified after colon
"down:10" -> decrease temperature. Value specified after colon
"INFO"    -> Sensor sendo information about himself
```
