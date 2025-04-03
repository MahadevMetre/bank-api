# Verify Mpin API Documentation

## Overview

This API is used to verify the M-PIN provided by the user, ensuring it matches the M-PIN stored in the backend for the user authenticated via the bearer token.
## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/open/verify-mpin
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
curl --location 'https://staging-api.paydoh.in/bankapi/api/open/verify-mpin' \
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
    "message": "successfully verified mpin",
    "status": 200
}


#### Error Responses

##### BadRequest Server Error 

**Status Code:** 400 Service Badrequest

#### Case 1 incorrect mpin

```json
{
    "error": {
        "errors": {
            "body": "incorrect MPIN. You have 4 attempts remaining"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Case 2 account locked reset mpin

```json
{
    "error": {
        "errors": {
            "body": "account locked. please reset your MPIN"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Rate Limiting

- Users are allowed a limited number of attempts (between 3 to 5) to verify their MPIN within an hour. Once this limit is reached, the account will be temporarily locked, and the user will need to reset their MPIN to regain access.

#### Security Considerations

1. The endpoint should be accessed over HTTPS in production
2. Implement rate limiting to prevent brute force attacks
3. Monitor failed attempts.

#### Notes

- After multiple failed attempts, users will receive a prompt to reset their M-PIN for account security.
- Users locked out due to rate limiting must reset their M-PIN or wait for the temporary lockout to expire.
- To reattempt verification after lockout, users will need to initiate a new verification session.

## Support

For any issues or queries, please contact:

- Email: yeshwanth@paydoh.money
