# DebitCard Generation API Documentation

## Overview

This API is used to generate virtual and/or physical debit cards for users based on their preferences during or after onboarding.

## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/debitcard/generate
```

#### Request Headers

| Header        | Required | Description                                |
| ------------- | -------- | ------------------------------------------ |
| Content-Type  | Yes      | Must be set to `application/json`          |
| X-Device-IP   | Yes      | Device IP address must be provided         |
| X-OS          | Yes      | Operating system must be specified         |
| X-OS-Version  | Yes      | OS version must be specified               |
| X-Lat-Long    | Yes      | Latitude and longitude must be provided    |
| Authorization | Yes      | Must be a Bearer Token of user             |

#### Request Body

The request body must be encrypted:

```json
{
    "data": "b3345cd3c654cbf80eae942b24d9ffc190a45668e25c42edb1f8ca4609b7dd6e0358c379e4f5eeaed3a194c107ccc7dee48fad8e9b249c7fc916b334233e911f"
}
```

Decrypted format of the data:
```json
{
    "debitcard_generation_type": "virtual"  // Options: "virtual", "physical", "both"
}
```

#### Example Request

```curl
curl --location '0.0.0.0:4100/api/debitcard/generate' \
--header 'X-Device-IP: test' \
--header 'X-OS: test' \
--header 'X-OS-Version: test' \
--header 'X-Lat-Long: 92.16,12.00' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOWUwZTYxMTItYmE5MC00fCNTLFx1MDAyNi4qa2ZGMzpjSntZWDZJaGdsWlIxXHUwMDNjTUtkMjBHMSIsImlzcyI6InBheWRvaC1iYW5rIiwiZXhwIjoxNzM1MjA4ODMxfQ.ZFyKLdye3nq8Mi5I-JtahIIC8Cy_oGgoSxcLgDFSlfM' \
--data '{
    "data":"b3345cd3c654cbf80eae942b24d9ffc190a45668e25c42edb1f8ca4609b7dd6e0358c379e4f5eeaed3a194c107ccc7dee48fad8e9b249c7fc916b334233e911f"
}'
```

#### Error Responses

##### Internal Server Error

**Status Code:** 500 Internal Server Error

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

#### Use Case Scenarios

1. During onboarding:
   - For virtual card only: Use `"debitcard_generation_type": "virtual"`
   - For both types: Use `"debitcard_generation_type": "both"`
2. Post-onboarding:
   - To request physical card: Use `"debitcard_generation_type": "physical"`

#### Security Considerations

1. The endpoint must be accessed over HTTPS in production
2. Request body must be encrypted
3. Valid Bearer token authentication is mandatory
4. Monitor for suspicious activities
5. Implement appropriate rate limiting

#### Notes

- This API previously supported GET method but has been updated to POST
- The request body must always be encrypted
- Authorization token must be valid and not expired
- All headers are mandatory
