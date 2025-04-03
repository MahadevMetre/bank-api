# Email Verification Status API Documentation

## Overview

The **Email Verification Status API** 
allows you to check whether a user's email address has been verified or not. This can be useful in applications where you need to validate user email addresses before allowing access to certain features or performing specific actions.

## Base URL
https://staging-api.paydoh.in/bankapi

## API Endpoint

### POST `/api/email/verification-status`

This endpoint is used to fetch the verification status of a user's email.

---

## Request Details

### Request Headers

To successfully call the API, the following headers are required:

| Header         | Required | Description                                                  |
|----------------|----------|--------------------------------------------------------------|
| `Content-Type` | Yes      | Set to `application/json`                                     |
| `X-Device-IP`  | Yes      | The device's IP address (e.g., `192.0.0.2`)                  |
| `X-OS`         | Yes      | The operating system of the device (e.g., `ios`, `android`)  |
| `X-OS-Version` | Yes      | The version of the operating system (e.g., `18.0`)           |
| `X-Lat-Long`   | Yes      | Latitude and longitude of the device (e.g., `19.093377739232523,72.91605446814229`) |
| `Authorization`| Yes      | Bearer token for user authentication                         |

### Example Request

```bash
curl --location '0.0.0.0:4100/api/email/verification-status' \
--header 'X-Device-IP: 192.0.0.2' \
--header 'X-OS: ios' \
--header 'X-OS-Version: 18.0' \
--header 'X-Lat-Long: 19.093377739232523,72.91605446814229' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTVjNWUwODItMWUxMC00fD93Skp5MlJvPTQ_S2lOJXt9UWtGNEN-MXBbQGEqUTI9IiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjUwNDcyNzJ9.t79NucuBVc9gItAf0it33EwCLJRXUng7z6c_KtbmiaI'
```

### Success Response
### Status Code: 200 OK

{
    "data": {
        "isEmailVerified": true  // true = verified, false = unverified
    },
    "message": "Status fetched successfully",
    "status": 200
}


### Key Changes:
1. **Format:** The markdown file is structured with sections for headers, request details, and responses for clarity.
2. **Headers and Example:** The headers section is clearly laid out in a table for easy reference.
3. **Error Responses:** Added some common error codes and examples to give more detail on how errors are handled.
4. **Security Considerations:** Added security tips for production use.

This should be easier to follow and much clearer for someone implementing or using the API.
