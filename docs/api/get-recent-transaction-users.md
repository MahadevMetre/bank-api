# Get Recent Transaction Users API Documentation

## Overview

This API endpoint retrieves a list of recent transaction users.

## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoint

```
GET /api/transaction/recent-users
```

## Request Headers

| Header        | Required | Description                     |
| ------------- | -------- | ------------------------------- |
| X-Device-Ip   | Yes      | IP address of the device        |
| X-Os          | Yes      | Operating system of the device  |
| X-Os-Version  | Yes      | Version of the operating system |
| X-Lat-Long    | Yes      | Latitude and longitude          |
| Authorization | Yes      | Bearer token for authentication |

## Example Request

```bash
curl --location 'localhost:4200/api/transaction/recent-users' \
--header 'X-Device-Ip: test' \
--header 'X-Os: test' \
--header 'X-Os-Version: test' \
--header 'X-Lat-Long: test,test' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTVmMWY1NmQtYTJkMS00fHAzSENDbyt3e0o3czddRTE7bExXRX1afm8hdlpCPVx1MDAyNnMiLCJpc3MiOiJwYXlkb2gtYmFuayIsImV4cCI6MTc0MDc3NDEwM30.i_eoe5XlxwCHKJGt2B-OGvqMUjjpfpxaVJAkBZif0l4'
```

## Success Response

### Status Code: 200 OK

```json
{
  "data": [
    {
      "first_name": "ASMITA",
      "last_name": "DALVI",
      "middle_name": "SUMEET",
      "mobile_number": "9920740561",
      "upi_addr": "j8108.paydoh@kvb",
      "user_id": "77e7924e-ced8-4"
    }
  ],
  "message": "Recent Transaction Users Retrieved",
  "status": 200
}
```

## Error Responses

### Status Code: 400 Bad Request

```json
{
  "error": {
    "errors": {
      "body": "Bad request"
    }
  },
  "message": "Bad Request",
  "status": 400
}
```

### Status Code: 401 Unauthorized

```json
{
  "error": "Unauthorized",
  "message": "Invalid or expired token",
  "status": 401
}
```

### Status Code: 500 Internal Server Error

```json
{
  "error": "Internal Server Error",
  "message": "An unexpected error occurred",
  "status": 500
}
```
