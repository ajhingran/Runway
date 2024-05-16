# Runway

## Overview
Runway is a simple application that provides daily notification service for available flights for requested city pairings. 
Requests can be configured using all available Google Flights toggle options and leverages Google Flights internally to run queries to obtain all possible flights based on a given request. 
Those are then filtered, with the cheapest flight's (or flights under a certain target price) information and purchase link sent via SMS to the user. 

## Deployment
Clone this repository and then launch the server side driver by running ```make build``` and then launching ```./runway``` with no arugments. 
You can specify a set of 12 arguments to the driver program if you want to run a singular request locally, and not through the client-server application. 

If you want to use the client-server approach, after deploying ```./runway``` you can connect to it and issue requests by simply running ```client.go``` and configuring your request as necessary. The driver will run by default on ```localhost:8080```. 

## Missing features
Currently there is no proper user client as no user input is requested. Additionally the service must be deployed locally and the driver must be running for texts to send. 
