# Assignment 2 #

> Group members: Kjetil Indrehus, Martin Johannessen, Sander Hauge.

This is an API which allows for: searching of reports on percentage of renewable energy in different countries' energy mix over time.

> /energy/v1/renewables/current 
> /energy/v1/renewables/history/
> /energy/v1/notifications/ 
> /energy/v1/status/


## Current endpoint ##

This endpoint retrieves the elements of the latest year currently available. The newest data in renewable-share-energy
is from 2021, and is therefore the current year of this project. The endpoint allows for searching of countries based on 
country code, as well as country name.

It uses the file renewable-share-energy.csv and REST Countries API, which is retrieved from: http://129.241.150.113:8080/v3.1. 

** Using the endpoint **
> Path: /energy/v1/renewables/current/{country?}{?neighbours=bool?}

Using no extra parameters will print all countries to the client.
If an optional parameter: /{country?}, is passed the corresponding country will be printed. This variable could be both
country codes, and also country name.
The query: {?neighbours=bool?}, may also be used, and will print information about the neighbouring countries of the
country passed. This query is dependent on the optional parameter country.

# Current Test #

There is created a test class for the current endpoint.

To use the test, print into command line when in root project folder:
> go test .\internal\webserver\handlers\current_test.go


## History endpoint ##

This endpoint retrieves all elements from renewable-share-energy. The endpoint allows for searching for specific countries,
years and also has the function of sorting by percentage and also calculating the mean of a country throughout.

It uses the file renewable-share-energy.csv 

** Using the endpoint **
> Path: /energy/v1/renewables/history/{country?}{?begin=year?}{?end=year?}{?mean=bool?}{?sortbyvalue=bool?}

These can also be combined, using "&" after "?". Begin and end query combined will find countries between the ones written.

Example request:
/energy/v1/renewables/history/nor?begin=2011&end=2014&sortbyvalue=true
```json
[
{
"name": "Norway",
"isoCode": "NOR",
"year": 2012,
"percentage": 70.095116
},
{
"name": "Norway",
"isoCode": "NOR",
"year": 2014,
"percentage": 68.88728
},
{
"name": "Norway",
"isoCode": "NOR",
"year": 2013,
"percentage": 67.50864
},
{
"name": "Norway",
"isoCode": "NOR",
"year": 2011,
"percentage": 66.30012
}
]
```

# History test #

There is created a test class for the history endpoint.

To use the test, print into command line when in root project folder:
> go test .\internal\webserver\handlers\history_test.go


<br>

# Deployment 

This service is deployed with OpenStack. OpenStack is a IaaS where the user define what resources is needed. It vitalizes resources to serve all end users. More information here: [OpenstackLink](https://www.ntnu.no/wiki/display/skyhigh) <br>

## Predefined resources
This service has the following resources predefined: 

1) Ubuntu Server 22.04 LTS (Jammy Jellyfish) amd64 
2) gx1.1c1r flavor 


## Security and Access

Security group prevents all communication with the server, but this is allowed: 

1) Allows to any ICMP package (Ping is allowed)
2) SSH (Port 22)
3) Http (Port 81)

To access our service you need be connected to the NTNU network. This could also use the Cisco VPN to connect to the campus network. [More information about the VPN here!](https://i.ntnu.no/wiki/-/wiki/Norsk/Cisco+AnyConnect+VPN)

The service can be access with: 
<br>
```http
http://10.212.169.162:81/energy/v1/status/
```

## Docker and its purpose 
Docker is a set of platform as a service (PaaS) products that use OS-level virtualization to deliver software in packages called containers. [This project used this docs for setting up docker on the OpenStack server.](https://docs.docker.com/engine/install/ubuntu/#set-up-the-repository)

<br>
Note that the project uses both a Docker file and a docker compose file. The docker file contains introductions for building the Docker Image. By defining a base image, additional packages and other versions, the Docker image can be created more deterministic. The base image that we use is `golang:1.18`. The Docker File also gets all packages from the go project.

[Read more about how dockerfile works here](https://docs.docker.com/engine/reference/builder/)

<br>

Docker Compose is for running multi-container Docker application. However, this project uses it to define volumes for the credentials file for firebase. Also the renewable energy .csv file should also stay in a volume. 

[More on docker compose here](https://docs.docker.com/compose/)