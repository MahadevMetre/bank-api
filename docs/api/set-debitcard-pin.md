# Set DebitCard PIN API Documentation

## Overview

This API is used to set DebitCard Pin.
## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/debitcard/set-debitcard-pin
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
curl --location '0.0.0.0:4200/api/debitcard/set-debitcard-pin' \
--header 'X-Device-IP: test' \
--header 'X-OS: test' \
--header 'X-OS-Version: test' \
--header 'X-Lat-Long: 92.16,12.00' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTJiZjdhMWUtYzZkYS00fDNiaVB0P0N-VHIpWC03NjMwa2xiZ29QQCNGJTIjPztEIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjQ4NDQyMzB9.xZAgMJdCzTX_4C_qDGbAizQibcDmvHeMryg0emN8S1k' \
--data '{
    "data":"bf2ff7434714c464a69e24b1c54d9e8e734633e7f951f245d8c0048ebda1d931d083214323c274f7e2b113df5a8c9bc73a87ca9f509186b60899c8a39a965c"
}'
```
Below data should be pass as encrypted formate
// {
//     "pin":"1234",
//     "pin_set_type":"New" //New - Reset
// }

#### Success Response

**Status Code:** 200 OK

```json
{
    "message": "OTP sent Successfully",
    "status": 200
}


#### Error Response

**Status Code:** 200 OK

```json
{
    "error": {
        "errors": {
            "body": "Pin Setup Type is not valid"
        }
    },
    "message": "Bad Request",
    "status": 400
}

#### Error Response

**Status Code:** 200 OK

```json
{
    "error": {
        "errors": {
            "body": "DebitCard Pin Must be 4 Digit"
        }
    },
    "message": "Bad Request",
    "status": 400
}


#### Security Considerations

1. The endpoint should be accessed over HTTPS in production
2. Implement rate limiting to prevent brute force attacks
3. Monitor failed attempts.
