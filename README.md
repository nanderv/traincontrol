# Readme
This is a prototype of a model train control system written in golang. It's an attempt to use a very modular structure.


# Services

## Bridge
The bridge service is responsible for sending and receiving messages from the physical actual railway.

## Hardware
The hardware module is responsible for keeping track (ha) of the state of the physical railway equipment on the railway.
It keeps track of the state of switches, power sections and sensors.

## Layout
The layout module is responsible for keeping track of the railway-operational structure of the layout.
This partly has the same concepts as hardware. The difference with these overlapping concepts is that the hardware package concerns itself only with the physical thing, while the Layout module reasons about it in operational terms.
This means that the hardware layer may be aware of sensor spots on the railway, but the Layout module is responsible for knowing where trains are / could be.

The layout module should ensure that no actively unsafe states are reached, such as powering deadend sections if the direction of power is towards the end.
