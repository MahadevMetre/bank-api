# Re-Mapping Upi_Id API Documentation

## Overview

This API is used to remapping upi_id if user is change the mobile device or change sim from one carrier to another carrier 
or else same sim vendor but different sim_id 
## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/upi/remapping-upi-id
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
curl --location 'https://staging-api.paydoh.in/bankapi/api/upi/remapping-upi-id' \
--header 'X-Device-Ip: 192.168.1.3' \
--header 'X-OS: android' \
--header 'X-OS-Version: 10' \
--header 'X-Lat-Long: 19.0933218,72.9161727' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzdlNzkyNGUtY2VkOC00fGQ1c3poNSs9P0owczpjaGlaLS8kQTloI0BEVFBnUUtsIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjE3MTc5Njh9.QIvofLOxCbv3ERP5MSQpA_mQJsOUku6IshGkDaGmiU8'\
--data '{

    "data": "54f2e0d6e53fea23af966ec25de318a6d8d92e37a35ebc6ec3b460a9e2d541893bd10f4004d8ee6c85efc0"

}'
```

#### Success Response

**Status Code:** 200 OK

```json
{
    "message": "Successfully remapped upi id for the user",
    "status": 200
}


#### Error Responses

##### BadRequest Server Error 

**Status Code:** 400 Service Badrequest

#### Case 1 if user is uninstall application

```json
{
    "error": {
        "errors": {
            "body": "MobileNo Already Exist with Same Device"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Case 2 if user is trying to install in another device

```json
{
    "error": {
        "errors": {
            "body": "MobileNo Already Exist with Different Device"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```
#### Case 3 Technical issue

```json
{
    "error": {
        "errors": {
            "body": "Due to technical issue, please try after sometime"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Case 4 Please Try Again Later issue

```json
{
    "error": {
        "errors": {
            "body": "Please try again later"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Case 5 Longcode smsfail to hit ourservice

```json
{
    "error": {
        "errors": {
            "body": "MobileNo not registered"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Case 6 Updation Failure

```json
{
    "error": {
        "errors": {
            "body": "Failure"
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
