<div align="center">
    <h1>Assignment 2</h1>
    <i>A project by: Sander Hauge, Kjetil Indrehus & Martin Johannessen</i>
</div>

<div align="center">
    <sup>Status</sup>
    <br />
    <a href="https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2023/-/wikis/Assignments/Assignment-2">
        <img alt="Assignment" src="https://img.shields.io/badge/Assignment-Click%20me-orange" />
    </a>
    
</div>

---

API which allows for searching of reports on percentage of renewable energy in different countries' energy mix over time.
```
/energy/v1/renewables/current 
/energy/v1/renewables/history/
/energy/v1/notifications/ 
/energy/v1/status/
```

---

# Index #
* [Current endpoint](#Current endpoint)
  * [Current test](#current-test)
* [History endpoint](#history-endpoint)
  * [History test](#history-test)
* [Notification endpoint](#notification-endpoint)
  * [Setting up a new notification subscription](#setting-up-a-new-notification-subscription--br)
  * [Deleting subscription](#deleting-a-notification-subscription--br)
  * [Information about a subscription](#retrieving-information-about-a-notification-subscription--br)
  * [Information about all subscriptions](#retrieving-information-about-all-notification-subscriptions-br)
  * [Notifications being sent](#when-are-you-notified)
  * [Purging mechanism](#purging-mechanism)
  * [Invocation incrementation](#when-is-invocations-incremented)
  * [Storing notifications](#how-are-notifications-stored)
  * [Notification tests](#notification-endpoint-and-firebase-tests)
* [Deployment](#deployment)
  * [OpenStack Configurations: Instance resources](#openstack-configurations-instance-resources)
  * [OpenStack Configurations: Security and Access](#openstack-configurations-security-and-access)
  * [About Docker](#docker-and-its-purpose)
  * [Deploying with Docker](#how-to-deploy-the-docker-service)

---


## Current endpoint ##

This endpoint retrieves the elements of the latest year currently available. The newest data in renewable-share-energy
is from 2021, and is therefore the current year of this project. Features of this endpoint:
* Search for country by name and country code.
* Add-on to get neighbouring countries.
* Cache for reducing amount of calls to countries API.
* Sorting of results.

It uses the file renewable-share-energy.csv and REST Countries API, which is retrieved from: http://129.241.150.113:8080/v3.1.
The endpoint only supports GET requests.

<br>

**Using the endpoint**
```
REQUEST: GET
PATH: /energy/v1/renewables/current/{country?}{?neighbours=bool?}{sortbyvalue/sortalphabetically=bool?, descending=bool?}
```
Using no extra parameters will print all countries of the year=2021, to the client. The year found is based on the
highest year found in the csv file.
If an optional parameter: /{country?}, is passed the corresponding country will be printed. This variable could be both
country codes, and country name.
The query: {?neighbours=bool?}, may also be used, and will print information about the neighbouring countries of the
country passed. This query is dependent on the optional parameter: country.

Three other queries may also be used:
* {sortbyvalue=bool?} is a query which sorts the results by percentage.
* {sortalphabetically=bool?} is a query which sorts the results alphabetically.
* {descending=bool?} is a query to sort descending, requires either sortbyvalue or alphabetically.

<br>

**Example Request:**
```
REQUEST: GET
PATH: /energy/v1/renewables/current/norway/neighbours=true
```
This is returned to the client.
```json
[
  {
    "name": "Norway",
    "isoCode": "NOR",
    "year": 2021,
    "percentage": 71.558365
  },
  {
    "name": "Finland",
    "isoCode": "FIN",
    "year": 2021,
    "percentage": 34.61129
  },
  {
    "name": "Sweden",
    "isoCode": "SWE",
    "year": 2021,
    "percentage": 50.924007
  },
  {
    "name": "Russia",
    "isoCode": "RUS",
    "year": 2021,
    "percentage": 6.6202893
  }
]
```

<br>

### Current Test ###

There is created a test class for the current endpoint.

To use the test, print into command line when in root project folder:
> go test .\internal\webserver\handlers\current_test.go

<br>

## History endpoint ##

This endpoint retrieves all elements from renewable-share-energy. When no query is passed it will return the mean of all
data based on each country. Functionality of history endpoint:
* Search for specific countries based on country code and name. 
* Allows for searching for specific years.
* Allows for searching to, from and between specific years.
* Sort by percentage and alphabetically, both descending and ascending. 
* Calculating the mean of a country.

It uses the file renewable-share-energy.csv and API for countries.
The endpoint only supports GET requests.

<br>

**Using the endpoint**
```
REQUEST: GET
PATH: /energy/v1/renewables/history/{country?}{?begin=year?}{?end=year?}{?mean=bool?}{?sortbyvalue=bool?}
```
These can also be combined, using "&" after "?". Begin and end query combined will find countries between the ones written.

<br>

**Example request:**
```
REQUEST: GET
PATH: /energy/v1/renewables/history/nor?begin=2011&end=2014&sortbyvalue=true
```
This request returns.
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

<br>

### History test ###

There is a test class for the history endpoint, which covers most of the history endpoints' functions.

To use the test, print into command line when in root project folder:
```terminal
go test .\internal\webserver\handlers\history_test.go
```

<br>

# Notification Endpoint #

To get notified by a given amount of calls a country has, register a webhook with this service.
To the body make sure to add: <br>
    - the url that should be invoked<br>
    - the alpha code of the country that you want to be notified by<br>
    - the number of calls to be notified for. See _How are notifications sent?_ for more details <br>
<br>

## Setting up a new notification subscription: <br>
Provide the following details to get notifications to the given url. The standard way is that the user will receive a GET request for the given url in the body. Here is how you register a notification: 

```
    REQUEST: Post
    PATH: "/energy/v1/notification" 
    BODY: 
    {
        "url": <The given url for the webhook to call>,
        "country": <Alpha code of the country>,
        "calls": <Number of calls for notification>
    }
```

The response should be **201 Created** if all went well. See the error message for more details. <br> 
You should also the the webhook ID in the body of the response. This ID is important, so save it for either deletion or retrieving details about it. Here is an example response: <br>

```json
    {
        "webhook_id": "OIdksUDwveiwe"
    }
```

## Deleting a notification subscription: <br>
To delete a webhook, send a DELETE request to the following endpoint, including the ID of the webhook in the URL:

```
REQUEST: DELETE
PATH: /energy/v1/notifications/{webhook_id}
```
<br>
Look at the status code for how the request for deletion went. If the status was: <br>
-  400: Please make sure that you added a id the the url. <br>
-  200: Webhook was either found and deleted, or not found (so nothing happend) <br>
-  500: Internal error while trying to delete the webhook. See the status endpoint to check if all services are running
<br><br>

## Retrieving information about a notification subscription: <br>

To get only information for a single given notification, use the id in the request: 

```
REQUEST: GET
PATH: /energy/v1/notifications/{webhook_id}
```

Look at the status code and message if no webhook was received. <br>
If there is a webhook with the given ID, the response could look like this:

```json
{
    "webhook_id": <ID_of_the_webhook>,
    "url": <Url_of_the_registration>,
    "country": <Alpha_code_of_the_country>,
    "calls": <The_amount_of_calls_that_needs_to_be_for_invoking>,
    "created_timestamp": <Server_timestamp_when_the_webhook_was_created>,
    "invocations": <The_amount_of_times_the_country_with_the_given_alpha_code_has_been_invoked>
}
```

## Retrieving information about all notification subscriptions <br>

To get all the notifications that are stored in the register: 

```
REQUEST: GET
PATH: /energy/v1/notifications/
```

Should return a list of all webhooks. Could also be empty if non are registered yet. Expected response would look like this; <br>

```json
[
    {
        "webhook_id": <ID_of_the_webhook>,
        "url": <Url_of_the_registration>,
        "country": <Alpha_code_of_the_country>,
        "calls": <The_amount_of_calls_that_needs_to_be_for_invoking>,
        "created_timestamp": <Server_timestamp_when_the_webhook_was_created>,
        "invocations": <The_amount_of_times_the_country_with_the_given_alpha_code_has_been_invoked>
    },
    ...
]
```

## When are you notified? 

You will receive a notification if the number of API calls made to a country's endpoint is divisible by the number of calls set for the notification. For example, if you set the number of calls to be notified for a country to 100 and 1000 API calls are made to that country's endpoint, you will be notified since 1000 % 100 = 0. In simpler terms, if you set **call** to be **5**, every 5th time the country have been called, you get notified.  

Here is an example of the JSON response you will receive when a notification is triggered:

```json
    {
        "webhook_id": "OIdksUDwveiwe",
        "country": "USA",
        "calls": 100,
        "invocations": 1000,
        "message": "Notification triggered: 1000 invocations made to USA endpoint."
    }
```    

## Purging mechanism

When the user adds a notification, a method called PurgeWebhooks is called. It checks if the amount of webhooks is now over the limit. If it is, then it starts removing the oldest notifications. It only removes enough webhooks so that the total amount of notifications are stored is under the limit. By default, the total amount of webhooks allowed is **40**. This could also be changed in the `constants` file. 

- Did not choose to delete webhooks that has the least amount of invocations, because new notifications would be deleted. 
- Improvements could be to update the creation time, whenever information has changed. Disregarded this due to conflict of naming: creation does not imply last updated, and also another field to keep track on would lead to unnecessary storage of data.

## When is invocations incremented?

Whenever there is a request to the third party api, `restcounties`, we increment all the webhooks for that country with one. In the same function we notify the user, if the condition of being notified is met. See paragraph above for more details. 

## How are notifications stored?

Using the firebase cloud storage called: [Firestore](https://firebase.google.com/docs/firestore). The application uses firestore to store webhooks in form of documents. See document databases for more information on how this works. The technical aspects for getting this to work is; <br>

> 1) Having a firestore credentials file in the root folder of the project. <br>
> 2) The credentials file MUST be called **cloud-assignment-2.json**
> 3) Manually created two collections called: **test_collection** and **webhooks**

<br>
<b>Note:</b> changes might lead to errors, so don't. <br>
This is also located in the constant code: <br>

```go
const (
	....
	FIRESTORE_COLLECTION = "webhooks" 
	FIRESTORE_COLLECTION_TEST = "test_collection" 
	FIREBASE_CREDENTIALS_FILE = "cloud-assignment-2.json" 
    ...
)

```

## Notification endpoint and Firebase tests.

The notification test are highly coupled with the Firebase test. Therefore are the notification test only to check that the endpoint works as it is supposed to. This means that it may be lacking. However, the Firebase test should have no issue if Firestore is correctly setup. From this if: <br>

1) **FIRESTORE && NOTIFICATION ENDPOINT TEST FAIL** -> Most likely just incorrectly setup the firestore 
2) **ONLY NOTIFICATION ENDPOINT FAIL** -> Logical error in the code in **notification.go** 

To test the firebase methods only: <br>

```terminal
go test ./db
```

To test the endpoints only: <br>

```terminal
go test ./internal/webserver/handlers/
```
<br>

# Deployment 

This service is deployed with OpenStack. OpenStack is a IaaS where the user define what resources is needed. It vitalizes resources to serve all end users. More information here: [Openstack Link](https://www.ntnu.no/wiki/display/skyhigh) <br>

## OpenStack Configurations: Instance resources
This service has the following resources predefined: 

1) Ubuntu Server 22.04 LTS (Jammy Jellyfish) amd64 
2) gx1.1c1r flavor 


## OpenStack Configurations: Security and Access

Security group prevents all communication with the server, but this is allowed: 

1) Allows to any ICMP package (Ping is allowed)
2) SSH (Port 22)
3) Http (Port 8000)

To access our service you need be connected to the NTNU network. This could also use the Cisco VPN to connect to the campus network. [More information about the VPN here!](https://i.ntnu.no/wiki/-/wiki/Norsk/Cisco+AnyConnect+VPN)

The service can be access with: 

<br>
```http
http://10.212.169.162:8000/energy/v1/status/
```

**Note:** In the case of self hosting, use the floating IP. 

## Docker and its purpose 
Docker is a set of platform as a service (PaaS) products that use OS-level virtualization to deliver software in packages called containers. [This project used this docs for setting up docker on the OpenStack server.](https://docs.docker.com/engine/install/ubuntu/#set-up-the-repository)

<br>
Note that the project uses both a Docker file and a docker compose file. The docker file contains introductions for building the Docker Image. By defining a base image, additional packages and other versions, the Docker image can be created more deterministic. The base image that we use is golang:1.18. The Docker File also gets all packages from the go project.

[Read more about how dockerfile works here](https://docs.docker.com/engine/reference/builder/)

<br>

Docker Compose is for running multi-container Docker application. However, this project uses it to define volumes for the credentials file for firebase. Also the renewable energy **.csv** file should also stay in a volume. 

[More on docker compose here](https://docs.docker.com/compose/)

## How to deploy the docker service? 

The following steps are
1. Connect to the NTNU Campus network (Via VPN or direct connection)
2. Create an OpenStack Instance has the same flavor and OS that has been specified in this README
3. Create and add Security Policy to the Instance with
    - SSH to port 22 (For logging on to the instance)
    - Ingress at port 8000 (For accessing the service)
4. Create and add SSH key to the Instance
    - Store the `.pem` file for logging in to the server
5. Login the the server using the floating IP and the `.pem` file with:
    ```terminal
    ssh -i ./name-of-ssh-key.pem ubuntu@<YOUR_FLOATING_IP>
    ```
    - **Common errors:** Not correct permissions to the `.pem`file or that the ssh key has not been set to the instance
    - Other errors is due to not correctly deploying an instance, see [OpenStack introduction docs here](https://www.ntnu.no/wiki/display/skyhigh/Using+the+webinterface)
6. Installed docker: [Docker installation manual for ubuntu here.](https://docs.docker.com/engine/install/ubuntu/#set-up-the-repository)
7. Clone the repo to the machine using `git clone`
8. Use the `scp` linux command for adding the firestore credential file:
    - Secure copy command article 
    - See the section on [Storing notifications](#how-are-notifications-stored) for how Firestore should be setup
    - File **must** be named:
      ```terminal
      ./cloud-assignment-2.json
      ```
    - File **must** be moved inside the repo at the root of the project 
7. **OPTIONAL** Set docker to the group of sodu privileges. The rest of the steps assumes docker can be used without sudo. By default, docker needs sudo privileges to run. Can also use `sudo docker ....` when using docker commands. 
8. Build and deploy with this command. Uses the compose file. Also detaches form the :
    ```terminal
    docker compose up -d 
    ```
9. Verify that compose the service has been deployed by using the docker command for checking on the service:
    ```terminal
    docker ps -a 
    ```

