# Get Sim Verification Status API Documentation

## Overview

This api is used to get status of sim verification status of user as per bearer token.

## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
GET /api/sim-verification-status
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
curl --location 'https://staging-api.paydoh.in/bankapi/api/authentication/sim-verification-status' \
--header 'X-Device-Ip: 192.168.1.3' \
--header 'X-OS: android' \
--header 'X-OS-Version: 10' \
--header 'X-Lat-Long: 19.0933218,72.9161727' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzdlNzkyNGUtY2VkOC00fGQ1c3poNSs9P0owczpjaGlaLS8kQTloI0BEVFBnUUtsIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjE3MTc5Njh9.QIvofLOxCbv3ERP5MSQpA_mQJsOUku6IshGkDaGmiU8'
```

#### Success Response

**Status Code:** 200 OK

```json
{
  "data": {
    "is_sim_verified": true
  },
  "message": "success",
  "status": 200
}

//  in case false

{
  "data": {
    "is_sim_verified": false
  },
  "message": "success",
  "status": 200
}
```

#### Error Responses

##### Internal Server Error

**Status Code:** 500 Service Unavailable

```json
{
  "error": {
    "errors": {
      "body": "test"
    }
  },
  "message": "Internal Server Error",
  "status": 500
}
```

#### Security Considerations

1. The endpoint should be accessed over HTTPS in production
2. Implement rate limiting to prevent brute force attacks
3. Monitor failed attempts.

#### Notes

- The mobile number should be in a valid format without country code
- An OTP will be sent to the provided mobile number
- The session_id in the response should be stored and used for subsequent OTP verification

## Support

For any issues or queries, please contact:

- Email: indal@paydoh.money
