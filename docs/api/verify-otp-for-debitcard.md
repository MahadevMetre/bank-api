# Verify OTP for DebitCard API Documentation

## Overview

This API is used to Verify OTP related to DebitCard.
## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/debitcard/verify-otp
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
curl --location '0.0.0.0:4200/api/debitcard/verify-otp' \
--header 'X-Device-IP: test' \
--header 'X-OS: test' \
--header 'X-OS-Version: test' \
--header 'X-Lat-Long: 92.16,12.00' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTJiZjdhMWUtYzZkYS00fDNiaVB0P0N-VHIpWC03NjMwa2xiZ29QQCNGJTIjPztEIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjQ4NDQyMzB9.xZAgMJdCzTX_4C_qDGbAizQibcDmvHeMryg0emN8S1k' \
--data '{
    "data":"89b2f745fc038665dc3077ac85943a3fc077206f3cf847f1d5022071c43da550b2652dc0033ea0b91c7090ae6dae3a9fb971a958feb562529b1c6f76274bfd4e31922ca389bb"
}'
```
Below data should be pass as encrypted formate
//{
//  "otp": "09812812",
//  "otp_type":"New" //New - Reset
//}

#### Success Response

**Status Code:** 200 OK

```json
{
    "message": "OTP sent Successfully for DebitCard Pin SetUp",
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
            "body": "OTP must be 8 Digit"
        }
    },
    "message": "Bad Request",
    "status": 400
}


#### Security Considerations

1. The endpoint should be accessed over HTTPS in production
2. Implement rate limiting to prevent brute force attacks
3. Monitor failed attempts.
