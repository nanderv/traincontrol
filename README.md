# Readme
This is a prototype of a model train control system written in golang. It's an attempt to use a very modular structure.


# Services

## bridge
The bridge service is responsible for sending and receiving messages from the physical actual railway. 

## hwconfig
The hardwareconfig service is responsible for configuring the hardware (=microcontrollers), making sure that they are linked to the required peripherals

## traintracks
The traintracks module is responsible for keeping track (ha) of the state of the physical railway equipment on the railway. 
It keeps track of the state of switches, power sections and sensors. 