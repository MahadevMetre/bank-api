# Get Pincode Details API Documentation

## Overview

This API is used to fetch details associated with a given pincode.

## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoint

```
POST /api/onboarding/pincode/details
```

## Request Headers

| Header        | Required | Description                       |
| ------------- | -------- | --------------------------------- |
| Content-Type  | Yes      | Must be set to `application/json` |
| X-Device-Ip   | Yes      | Device IP address                 |
| X-OS          | Yes      | Operating system                  |
| X-OS-Version  | Yes      | OS version                        |
| X-Lat-Long    | Yes      | Latitude and longitude            |
| Authorization | Yes      | Bearer token for authentication   |

## Request Body

```json
{
  "data": "encrypted_data_here"
}
```

The `data` field should contain the encrypted form of:

```json
{
  "pincode": "400001"
}
```

## Example Request

```curl
curl --location 'https://staging-api.paydoh.in/bankapi/api/onboarding/pincode/details' \
--header 'Content-Type: application/json' \
--header 'X-Device-Ip: 192.168.1.1' \
--header 'X-OS: ios' \
--header 'X-OS-Version: 14.0' \
--header 'X-Lat-Long: 19.0760,72.8777' \
--header 'Authorization: Bearer your_token_here' \
--data '{
    "data": "encrypted_data_here"
}'
```

## Success Response

**Status Code:** 200 OK

```json
{
  "data": [
    {
      "state_name": "Uttar Pradesh",
      "city_name": "BADALPUR"
    }
  ],
  "message": "Pincode details fetched successfully",
  "status": 200
}
```

## Error Responses

### Bad Request

**Status Code:** 200 Bad Request

```json
{
  "error": {
    "errors": {
      "body": "Pincode must be less than or equal to 6"
    }
  },
  "message": "Bad Request",
  "status": 400
}
```

### Unauthorized

**Status Code:** 200 Unauthorized

```json
{
  "error": {
    "errors": {
      "body": "Unauthorized access"
    }
  },
  "message": "Unauthorized",
  "status": 401
}
```

### Internal Server Error

**Status Code:** 500 Internal Server Error

```json
{
  "error": {
    "errors": {
      "body": "An unexpected error occurred"
    }
  },
  "message": "Internal Server Error",
  "status": 500
}
```
