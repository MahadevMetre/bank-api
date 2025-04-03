# change upi pin API Documentation

## Overview

This API is used to change upi pin for the user with old-upi-pin and new-upi-pin
## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/upi/change-upi-pin
```

#### Request Headers

| Header       | Required | Description                       |
| ------------ | -------- | --------------------------------- |
| Content-Type | Yes      | Must be set to `application/json` |
| X-Device-Ip  | Yes      | Must be set to added              |
| X-OS         | Yes      | Must be set to added              |
| X-OS-Version | Yes      | Must be set to added              |
| X-Lat-Long   | Yes      | Must be set to added              |
| Bearer       | Yes      | Must be a Bearer Token of user.   |

#### Example Request

```curl
curl --location 'https://staging-api.paydoh.in/bankapi/api/upi/change-upi-pin' \
--header 'X-Device-Ip: 192.168.1.3' \
--header 'X-OS: android' \
--header 'X-OS-Version: 10' \
--header 'X-Lat-Long: 19.0933218,72.9161727' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzdlNzkyNGUtY2VkOC00fGQ1c3poNSs9P0owczpjaGlaLS8kQTloI0BEVFBnUUtsIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjE3MTc5Njh9.QIvofLOxCbv3ERP5MSQpA_mQJsOUku6IshGkDaGmiU8'\
--data '{

    "data": "54f2e0d6e53fea23af966ec25de318a6d8d92e37a35ebc6ec3b460a9e2d541893bd10f4004d8ee6c85efc0"

}'
```
Request Body:

{
  "old_upi_pin": "",
  "new_upi_pin": "",
  "trans_id": ""
}


| Field        |  Type    |              Description                          |  validation |
| ------------ | -------- | --------------------------------------------------|-------------|
| old_upi_pin  | string   | The current UPI PIN of the user.                  |Required     |
| new_upi_pin  | string   | The new UPI PIN that the user wants to set.       |Required     |
| trans_id     | string   | The transaction ID for the UPI PIN change request.|Required     |


#### Success Response

**Status Code:** 200 OK

```json

{
    "message": "Successfully user changed upi pin",
    "status": 200
}


#### Error Responses

##### BadRequest Server Error 

**Status Code:** 400 Service Badrequest

#### Case 1 

```json
{
    "error": {
        "errors": {
            "body": "wrong pin entered"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Security Considerations

1. The endpoint should be accessed over HTTPS in production
2. Implement rate limiting to prevent brute force attacks
3. Monitor failed attempts.

## Support

For any issues or queries, please contact:

- Email: yeshwanth@paydoh.money
