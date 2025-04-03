# Get Static API Documentation

## Overview

This API is used to Static Data which we are using in Paydoh Application like DebitCard Purchase Amount.
## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
GET /api/static-parameters
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
curl --location 'https://staging-api.paydoh.in/bankapi/api/static-parameters' \
--header 'X-Device-Ip: 192.168.0.101' \
--header 'X-OS: android' \
--header 'X-OS-Version: 14.0' \
--header 'X-Lat-Long: 92.18,12.31' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzdlODAxMmQtZTczMC00fFpDNjVrezN0aDpva1A9cXEtLT12ezBSfXE4JE9mW3RIIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjM0NTg1Nzd9.FX4-byw-0WEAa6gHLFPWGGx8GJwJ3ifgwGYWN053UW4'
```

#### Success Response

**Status Code:** 200 OK

```json
{
    "data": {
        "debitcard_amt": "236"
    },
    "message": "Static Parameters get successfully",
    "status": 200
}


#### Security Considerations

1. The endpoint should be accessed over HTTPS in production
2. Implement rate limiting to prevent brute force attacks
3. Monitor failed attempts.
