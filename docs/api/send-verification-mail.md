# Send Verification Email API Documentation

## Overview

This API is used to send a verification email to the user to verify their account.

## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoint

### `POST /api/email/sendverification`

This endpoint is responsible for sending a verification email to the user.

#### Request Headers

| Header        | Required | Description                                                                          |
|---------------|----------|--------------------------------------------------------------------------------------|
| Content-Type  | Yes      | Must be set to `application/json`.                                                   |
| X-Device-IP   | Yes      | The IP address of the device from which the request is made.                         |
| X-OS          | Yes      | The operating system of the device (e.g., `ios`, `android`).                         |
| X-OS-Version  | Yes      | The version of the operating system of the device.                                   |
| X-Lat-Long    | Yes      | The latitude and longitude of the device, used for geolocation.                      |
| Authorization | Yes      | A Bearer token for authentication. The token should be provided as `Bearer <token>`. |

#### Request Body

The request body should contain the encrypted user data.

```json
{
    "data": "27340bc29a2f4cae6673ae197d507a9561b81f666a00cf1be70ed8617178dd9321dc4cfb77efb7fde132a34d03685f304c0321723189e3bc784ae1626f7017540edbf51da639f056f1"
}


#### Example Request

```curl
curl --location '0.0.0.0:4100/api/email/sendverification' \
--header 'X-Device-IP: 192.0.0.2' \
--header 'X-OS: ios' \
--header 'X-OS-Version: 18.0' \
--header 'X-Lat-Long: 19.093377739232523,72.91605446814229' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzdlODAxMmQtZTczMC00fFpDNjVrezN0aDpva1A9cXEtLT12ezBSfXE4JE9mW3RIIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjUzNDY0Nzd9.V6SMnCHwuPuYkQqMpNohopJY_vtbuMBXtjUiTSd8F8Y' \
--data '{
    "data":"27340bc29a2f4cae6673ae197d507a9561b81f666a00cf1be70ed8617178dd9321dc4cfb77efb7fde132a34d03685f304c0321723189e3bc784ae1626f7017540edbf51da639f056f1"
}'

```

### Decryptd Data
The data field should be an encrypted string containing the following information and the both field are mandatory:
{
    "email":"testmail@gmail.com",
    "name":"Ruchika Solanki"
}

#### Success Response
**If the verification email is successfully sent, the API will respond with the following:**
**Status Code:** 200 OK

```json
{
    "message": "verification mail sent successfully",
    "status": 200
}

```

**Status Code:** 400 OK

```json
{
    "error": {
        "errors": {
            "body": "Name is required"
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
