<div align="center">
    <h1>Assignment 2</h1>
    <i>A project by: Sander Hauge, Kjetil Indrehus & Martin Johannessen</i>
</div>

<div align="center">
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


# Table of content

* [Current endpoint](#current-endpoint)
  + [Request](#request)
  + [Response](#response)
* [History endpoint](#history-endpoint)
  + [Request](#request-1)
  + [Response](#response-1)
  + [Example request:](#example-request-)
- [Notification Endpoint](#notification-endpoint)
  * [Setting up a new notification subscription: <br>](#setting-up-a-new-notification-subscription---br-)
  * [Deleting a notification subscription: <br>](#deleting-a-notification-subscription---br-)
  * [Retrieving information about a notification subscription: <br>](#retrieving-information-about-a-notification-subscription---br-)
  * [Retrieving information about all notification subscriptions <br>](#retrieving-information-about-all-notification-subscriptions--br-)
  * [Notification Event Types](#notification-event-types)
  * [Purging mechanism](#purging-mechanism)
  * [When is invocations incremented?](#when-is-invocations-incremented-)
  * [How are notifications stored?](#how-are-notifications-stored-)
- [Status endpoint](#status-endpoint)
      - [Response](#response-2)
      - [Example request and response](#example-request-and-response)
- [Default endpoint](#default-endpoint)
  * [Endpoint Tests](#endpoint-tests)
    + [Current Tests](#current-tests)
    + [History test](#history-test)
  * [Notification endpoint and Firebase tests.](#notification-endpoint-and-firebase-tests)
  * [Status tests](#status-tests)
  * [Default tests](#default-tests)
- [Deployment](#deployment)
  * [OpenStack Configurations: Instance resources](#openstack-configurations--instance-resources)
  * [OpenStack Configurations: Security and Access](#openstack-configurations--security-and-access)
  * [Docker and its purpose](#docker-and-its-purpose)
  * [How to deploy the docker service?](#how-to-deploy-the-docker-service-)
- [Advanced functionality](#advanced-functionality)
  * [Cache](#cache)
  * [Event types](#event-types)
  * [Current endpoint](#current-endpoint-1)
  * [History endpoint](#history-endpoint-1)
  * [Country searching](#country-searching)
    + [Drawbacks](#drawbacks)
    + [Examples](#examples)
  * [Wiki](#wiki)
- [Design](#design)
    + [Project structure](#project-structure)
- [Further development](#further-development)
  * [Administrator user](#administrator-user)
  * [Other](#other)

---


## Current endpoint ##
This endpoint retrieves the elements of the latest year currently available. The newest data in renewable-share-energy
is from 2021, and is therefore the current year of this project.

<u>Features of this endpoint:</u> <br>
* Search for country by name and country code.
* Add-on to get neighbouring countries.
* Cache for reducing amount of calls to countries API.
* Sorting of results.

The endpoint uses a file: "renewable-share-energy.csv" and REST Countries API, which is retrieved from: http://129.241.150.113:8080/v3.1.
The file contains historical data from each countries' share of renewable sources.

### Request ###
```
Method: GET
Path: /energy/v1/renewables/current/{country?}{?neighbours=bool?}{sortbyvalue/sortalphabetically=bool?, descending=bool?}
```
Using no extra parameters will print all countries of the year=2021, to the client. The year found is based on the
highest year found in the csv file.

`{country?}` is an optional parameter which could be passed to the API, which will print information about the country 
as long as it is found. It could be a 3-letter country code, or country name.

`{?neighbours=bool?}` is an optional query parameter which will print information about the neighbouring countries of the
country passed. It therefore, dependent on the optional parameter: `country`.

`{?sortbyvalue=bool?}` is an optional query parameter which will sort the results by percentage.

`{?sortalphabetically=bool?}` is an optional query parameter which sorts the results alphabetically.

`{?descending=bool?}` is an optional query parameter which is dependent on the other sorting queries. When used it will sort descending instead.

The endpoint only supports GET requests.

<br>

### Response ###
Content type: `application/json`

Status codes:
* 200: Success, everything works as intended.
* 400: Bad request, error in the request. For example sent a number when it should have been a string.
* 404: Not found, for example did not find country code based on search.
* 500: Internal server error, error due to the server. For example, faulty file reading.
* 501: Not implemented, function is not implemented.

Body:
```
{
	"Name":       <country_name>,               (string)
	"IsoCode":    <country_code>,               (string)
	"Year"        <year_recorded>,              (int)
	"Percentage": <percentage_of_renewables>    (float64)
}
```

<br>

**Example Requests & Responses:**

Request sent: `/energy/v1/renewables/current/swe`
Response:
```json
{
        "name": "Sweden",
        "isoCode": "SWE",
        "year": 2021,
        "percentage": 50.924007
}
```

Request sent: `/energy/v1/renewables/current/norway?neighbours=true`
Response:
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

## History endpoint ##

This endpoint retrieves all elements from renewable-share-energy. When no query is passed it will return the mean of all
data based on each country. 
<u> Functionality of history endpoint: </u> <br>
* Search for specific countries based on country code and name. 
* Allows for searching for specific years.
* Allows for searching to, from and between specific years.
* Sort by percentage and alphabetically, both descending and ascending. 
* Calculating the mean of a country.

The endpoint uses a file: "renewable-share-energy.csv" and REST Countries API, which is retrieved from: http://129.241.150.113:8080/v3.1.
The file contains historical data from each countries' share of renewable sources.

<br>

### Request ###
```
REQUEST: GET
PATH: /energy/v1/renewables/history/{country?}{?begin=year?}{?end=year?}{?mean=bool?}{?sortbyvalue=bool?}
```
When you use this endpoint with no parameters or queries, it will print the mean of all historical entries for each country.
The data is retrieved from renewable share energy.

`{country?}` is an optional parameter which could be passed to the API, which will print information about the country
as long as it is found. It could be a 3-letter country code, or country name.

`{?begin=year?}` is an optional query parameter used to filter results from a specific year. 

`{?end=year?}` is an optional query parameter used to filter results to a specific year.

`{?begin=year&end=year?}` using both begin and end it will return results between the years written.

`{?mean=bool?}` is an optional query parameter, which will only work in tandem with `country`, `begin` or `end`. It 
calculates the mean of the elements returned. This is done if no queries is presented.

`{?sortbyvalue=bool?}` is an optional query parameter which will sort the results by percentage.

`{?sortalphabetically=bool?}` is an optional query parameter which sorts the results alphabetically.

`{?descending=bool?}` is an optional query parameter which is dependent on the other sorting queries. When used it will sort descending instead.

The endpoint only supports GET requests.

<br>

### Response ###
Content type: `application/json`

Status codes:
* 200: Success, everything works as intended.
* 400: Bad request, error in the request. For example no country matching search parameter.
* 404: Not found, for example did not find country code based on search.
* 405: Method not allowed, writing invalid value in query. Example: ?query=notSupposedToBeLikeThis.
* 411: Length required, need more information to work.
* 500: Internal server error, error due to the server. For example, faulty file reading.
* 501: Not implemented, function is not implemented.

Body:
```
{
	"Name":       <country_name>,               (string)
	"IsoCode":    <country_code>,               (string)
	"Year"        <year_recorded>,              (int)
	"Percentage": <percentage_of_renewables>    (float64)
}
```

### Example request: ###

Request sent: `/energy/v1/renewables/history/sverige?mean=true`
Response:
```json
[
    {
        "name": "Sweden",
        "isoCode": "SWE",
        "percentage": 33.970860684210535
    }
]
```

Request sent: `/energy/v1/renewables/history/nor?begin=2011&end=2014&sortbyvalue=true`
Response:
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

# Notification Endpoint #

To get notified by a given amount of calls a country has, register a webhook with this service.
To the body make sure to add: <br>
    - the url that should be invoked<br>
    - the alpha code of the country that you want to be notified by<br>
    - the number of calls to be notified if the event is for calls. <br>
    - the type of event to be notified on. See [the notification types here.](#notification-event-types)
<br>

## Setting up a new notification subscription: <br>
Provide the following details to get notifications to the given url. The standard way is that the user will receive a GET request for the given url in the body. Here is how you register a notification: 

```
    REQUEST: Post
    PATH: "/energy/v1/notification" 
    BODY: 
    {
        "url": "The given url for the webhook to call",
        "country": "Alpha code of the country",
        "calls": "Number of calls for notification"
        "event": "Type of event"
    }
```

The response should be **201 Created** if all went well. See the error message for more details. <br> 
You should also the webhook ID in the body of the response. This ID is important, so save it for either deletion or retrieving details about it. Here is an example response: <br>

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
-  400: Please make sure that you added an ID the url. <br>
-  200: Webhook was either found and deleted, or not found (so nothing happened) <br>
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
    "webhook_id": "ID_of_the_webhook",
    "url": "Url_of_the_registration",
    "country": "Alpha_code_of_the_country",
    "calls": "The_amount_of_calls_that_needs_to_be_for_invoking",
    "event" : "The event type of the notification",
    "created_timestamp": "Server_timestamp_when_the_webhook_was_created",
    "invocations": "The_amount_of_times_the_country_with_the_given_alpha_code_has_been_invoked"
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
        "webhook_id": "ID_of_the_webhook",
        "url": "Url_of_the_registration",
        "country": "Alpha_code_of_the_country",
        "calls": "The_amount_of_calls_that_needs_to_be_for_invoking",
        "event": "The_type_of_event",
        "created_timestamp": "Server_timestamp_when_the_webhook_was_created",
        "invocations": "The_amount_of_times_the_country_with_the_given_alpha_code_has_been_invoked"
    },
    ...
]
```

## Notification Event Types

This service offers three types of of events:

- **PURGE:** 
  - **Description:** Notification if the service had to purge webhooks (due to stepping over the limit)
  - **Does it delete itself after invocation?** No, the notification is saved.
  - **Example of registration:**
      ```
      REQUEST: Post
      PATH: "/energy/v1/notification" 
      BODY: 
      {
        "url": "https://webhook.site/url-stuff",
        "event": "purge"
      }
    ```
- **CALLS:**  
  - **Description:** When the number of invocations is dividable by the calls number
  - **Does it delete itself after invocation?** No, the notification is saved.
  - **Example of registration:**
      ```
      REQUEST: Post
      PATH: "/energy/v1/notification" 
      BODY: 
      {
        "url": "https://webhook.site/url-stuff",
        "country": "NOR",
        "calls": 4,
        "event": "calls"
      }
    ```
- **COUNTRY_DOWN:** 
  - **Description:** If the country API goes down, a notification is sent. Only happens if the status endpoint gives other status code for the country API then **200**
  - **Does it delete itself after invocation?** Yes, once invoked, a new notification is needed.
  - **Example of registration:**
      ```
      REQUEST: Post
      PATH: "/energy/v1/notification" 
      BODY: 
      {
          "url": "https://webhook.site/url-stuff",
          "event": "country_down"
      }
      ```

Here is an example of the JSON response you will receive when a notification is triggered based on the calls event:

```json
    {
      "webhook_id": "32b184e5bc9e9bee7fdff1362dc2e05bc7174290a4cc3622cd39f5b1803c97e6",
      "url": "https://webhook.site/sample_url",
      "country": "NOR",
      "calls": 2,
      "event": "CALLS",
      "invocations": 38,
      "message": "Notification triggered: 38 invocations made to NOR endpoint."
    }
```    

## Purging mechanism

When the user adds a notification, a method called PurgeWebhooks is called. It checks if the amount of webhooks is now over the limit. If it is, then it starts removing the oldest notifications. It only removes enough webhooks so that the total amount of notifications are stored is under the limit. By default, the total amount of webhooks allowed is **40**. This could also be changed in the `constants` file. 

- Did not choose to delete webhooks that has the least amount of invocations, because new notifications would be deleted. 
- Improvements could be to update the creation time, whenever information has changed. Disregarded this due to conflict of naming: creation does not imply last updated, and also another field to keep track on would lead to unnecessary storage of data.

## When is invocations incremented?

Whenever there is a request to the third party api, `restcounties`, we increment all the webhooks for that country with one. In the same function we notify the user, if the condition of being notified is met.  

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

<br>

# Status endpoint #
The status endpoint provides the availability of all individual services this service depends on.
The reporting occurs based on status codes returned by the dependent services. The status interface
further provides information about the number of registered webhooks, and the uptime of the service.
It also provides the total memory usage of the computer in use.

```
Method: GET
Path: /energy/v1/status/
```

#### Response
Content type: `application/json`

Status codes
* 200: Everything is OK.
* 404: Not found.
* 500: Internal server error.

Status content
* countries_api: the http status code for the "REST Countries API".
* notification_db: the http status code for "Notification DB" in Firebase.
* webhooks: the number of registered webhooks.
* version: set to "v1".
* uptime: the time since the last service restart.

#### Example request and response

Request: `/status`

Response:

```
{
   "countries_api": "http status code for restcountries API",
   "notification_db": "http status code for notification DB in Firebase",
   "webhooks": "amount of registered webhooks",
   "version": "v1",
   "uptime": "time elapsed from the last service restart"
   "total_memory_usage": "percent of total memory usage on the user's computer"
}
```

Note: `"some value"` indicates placeholders for values to be populated by the service.
An example response is provided underneath.

Example response:
```
{
    "countries_api": 200,
    "notification_db": 200,
    "webhooks": 2,
    "version": "v1",
    "uptime": "10 seconds",
    "total_memory_usage": "78%"
}
```

<br>

# Default endpoint
This endpoint is the server's root path level. It does not provide any functionality, but assists the user to navigate
in the server. The HTML file in linked up with a css file in order to provide a more clean look to the page, with the 
endpoints being displayed in an organized and easy-to-use format. It is possible to press the different endpoints to 
navigate to their respective endpoints.

```
REQUEST: GET
PATH: /energy/
```

<br>


## Endpoint Tests ##
To run all endpoint tests write the following command in root folder: 
```
go test .\internal\webserver\handlers
```


### Current Tests ###
There is created a test class for the current endpoint.

To use the test, print into command line when in root project folder: 
```
go test .\internal\webserver\handlers\current_test.go
```

<br>

### History test ###

There is a test class for the history endpoint, which covers most of the history endpoints' functions.

To use the test, print into command line when in root project folder:
``` 
go test .\internal\webserver\handlers\history_test.go
```

<br>

## Notification endpoint and Firebase tests.

The notification test are highly coupled with the Firebase test. Therefore are the notification test only to check that the endpoint works as it is supposed to. This means that it may be lacking. However, the Firebase test should have no issue if Firestore is correctly setup. From this if: <br>

1) **FIRESTORE && NOTIFICATION ENDPOINT TEST FAIL** -> Most likely just incorrectly setup the firestore
2) **ONLY NOTIFICATION ENDPOINT FAIL** -> Logical error in the code in **notification.go**

To test the firebase methods only: <br>
``` 
go test ./db
```

<br>

## Status tests ##
There is created a test class for the status endpoint.

To use the test, print into command line when in root project folder:
```
go test .\internal\webserver\handlers\status_test.go
```

<br>

## Default tests ##
There is created a test class for the default endpoint.

To use the test, print into command line when in root project folder:
``` 
go test .\internal\webserver\handlers\default_test.go
```

<br>

# Deployment 

This service is deployed with OpenStack. OpenStack is a IaaS where the user define what resources is needed. It vitalizes resources to serve all end users. More information here: [Openstack Link](https://www.ntnu.no/wiki/display/skyhigh) <br>

## Deployed Service 

The service is deployed with openstack. <br>
Access it with this floating IP:

http://10.212.169.162:8000/energy/v1/status/


**Note:** In the case of self-hosting, use the floating IP of the instance.

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
5. Login the server using the floating IP and the `.pem` file with:
    ```terminal
    ssh -i ./name-of-ssh-key.pem ubuntu@"YOUR_FLOATING_IP"
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

## How to locally run the service?

1. Have go installed on the local machine. See [download versions here (use go.1.18)](https://go.dev/dl/)
2. Clone the repo.
3. Run the project by cd into the project folder, then run:
```terminal
  go run ./cmd/main.go
```
4. See logs for the port of running service. (Usually port 8080)
5. Access the service with local host here:
```
  http://localhost:PORT_NUMBER/energy/v1/status/
```

# Advanced functionality
This assignment introduced the following advanced tasks:

## Cache
Description: Implement purging of cached information for requests older than a given number of hours/days.

When searching for a country/ies, the country/ies will be cached for a set period of time. This is done to ensure an
updated cache. This result in less frequent requests for the API and shorter response time.

## Different Event types
The notification supports different types of events. There are currently 3 types of events. Read more about them [here](#notification-event-types)

## Purging of notifications
When the number of notifications are over the max limit defined, webhook will be deleted. Read more [here](#purging-of-notifications)

## Current endpoint
Description: Extend {?country} to support country name (e.g., norway) as input.

## History endpoint
Description: Selective use of only begin or end as single parameter (e.g., ?begin=1980 only consider data from 1980
onwards; ?end=1980, values from the first time entry until 1980 only).

Description: Extend the history for all countries with a time constraint. Where {?begin=year&end=year?} is specified
(e.g., ?begin=1960&end=1970), only calculate mean values for these years (not for all years).

Description: Additional optional parameter {?sortByValue=bool?} to support sorting of output by percentage
value (e.g., ?sortByValue=true).

## Country searching ##
Both current and history endpoint has the functionality of searching by country code and also name, which is an advanced
functionality. We have implemented another bonus functionality, which searches the API if it does not find a country
in the csv file. It will use the API: http://129.241.150.113:8080/v3.1/name/. 
It will then search for any type of name in country body, which could be common, official and nativeName.
This allows for searches of `/history/Kongeriket Norge`, which will return information about norway.

### Drawbacks ###
This feature will increase the amount of calls to API, which is sometimes unnecessary. If a user writes gibberish into
our API, it will search the country API. However, using the API a user may even search the countries native name and receive
the correct country. 

### Examples ###
Request: `/current/España`
Response:
```json
[
  {
        "name": "Spain",
        "isoCode": "ESP",
        "year": 2021,
        "percentage": 22.341663
  }
]
```

## Wiki
We have also created a [wiki](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2023-workspace/marhjoh/assignment-2/-/wikis/home)
for this project which includes more information about the assignment, like for example an overview of the
applications' use-case examples and the group dynamic during the assignment.

# Design
Throughout the implementation of this application, the focus points on the design has been loose coupling, high 
cohesion and modularity as close to Golang convention as possible. This has been done through constants, different 
files for handlers and generic functions.

### Project structure
The project structure was created with the goal of responsibility driven design, and to minimize code duplication overall.

The endpoint-handlers got one file each, and are all located in the "handlers" package.

In order to limit API-requests, when countries are requested all borders from the country is retrieved and for each 
border request the country. By doing this the API workload can get large. The API-server side's workload is reduced. 
In this way the REST-principles are met.

# Further development
These are further improvements we did not have time to resolve.

## Administrator user
Let the user get extra functionality based on their user role, for example an administrator. An administrator would
have certain privileges to data related to the server health, examples are provided below:
* Response time: Measure the time it takes for the service to respond to requests.
* Error rate: An idea of the number of errors that occur in the service.
* Request count: The number of requests made to the service.

This could be solved through making requests with a HEAD field including a passphrase or password. A user without
the credentials in the HEAD field would not be able to access the endpoint that offers the privileged functionality since
the user would not be authenticated.

## Other
* Use a middleware to set the content-type header for all response.
* Implement Gorilla Mux to define URL routes and extract variables from them instead of doing it manually.