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
`
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
`

# History test #

There is created a test class for the history endpoint.

To use the test, print into command line when in root project folder:
> go test .\internal\webserver\handlers\history_test.go

# Notification Endpoint #

To get notified by a given amount of calls a country has, register a webhook with this service.
To the body make sure to add: <br>
    - the url that should be invoked<br>
    - the alpha code of the country that you want to be notified by<br>
    - the number of calls to be notified for. See _How are notifications sent?_ for more details <br>
<br>
### **To setup a new notification subscription:** <br>
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

### **To Delete a notification:** <br>

There is not hard to delete a webhook. Simply use the Webhook ID in the url as shown below. ID has to be given to be able to delete. Do this by sending a DELETE request; 

```
REQUEST: DELETE
PATH: /energy/v1/notifications/{id_of_the_webhook}
```
<br>
Look at the status code for how the request for deletion went. If the status was: <br>
-  400: Please make sure that you added a id the the url. <br>
-  200: Webhook was either found and deleted, or not found (so nothing happend) <br>
-  500: Internal error while trying to delete the webhook. See the status endpoint to check if all services are running

### **To get a notification (with an ID):** <br>

To get only information for a single given notification, use the id in the request: 

```
REQUEST: GET
PATH: /energy/v1/notifications/{id_of_the_webhook}
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




### **To get all notifications:** <br>

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

