# Logout API Documentation

## Overview

This API is used to log out a user from the system.

## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/authentication/logout
```

#### Request Headers

| Header       | Required | Description                       |
| ------------ | -------- | --------------------------------- |
| Content-Type | Yes      | Must be set to `application/json` |
| X-Device-Ip  | Yes      | IP address of the device          |
| X-OS         | Yes      | Operating system of the device    |
| X-OS-Version | Yes      | Version of the operating system   |
| X-Lat-Long   | Yes      | Latitude and longitude of device  |
| Bearer       | Yes      | Bearer Token of the user          |

#### Example Request

```curl
curl --location 'https://staging-api.paydoh.in/bankapi/api/authentication/logout' \
--header 'Content-Type: application/json' \
--header 'X-Device-Ip: 192.168.1.4' \
--header 'X-Lat-Long: 19.0935922,72.9163083' \
--header 'X-Os: android' \
--header 'X-Os-Version: 12' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
```

#### Success Response

**Status Code:** 200 OK

```json
{
  "message": "User logged out successfully",
  "status": 200
}
```

#### Error Responses

##### Unauthorized

**Status Code:** 401 Unauthorized

```json
{
  "error": {
    "message": "Invalid or expired token"
  },
  "status": 401
}
```

##### Internal Server Error

**Status Code:** 500 Internal Server Error

```json
{
  "error": {
    "message": "Internal Server Error"
  },
  "status": 500
}
```

## Support

For any issues or queries, please contact:

- Email: indal@paydoh.money
