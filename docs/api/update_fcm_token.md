# Update FCM Token API Request

This document outlines the `curl` request to update the FCM token for a device.

### Endpoint:
`POST /api/open/update-fcm-token`

### Request Headers:
- **X-Device-Ip**: `192.168.0.101`
  (The IP address of the device)

- **X-OS**: `android`
  (Operating System of the device)

- **X-OS-Version**: `14.0`
  (Version of the operating system)

- **X-Lat-Long**: `92.18,12.31`
  (Latitude and Longitude of the device)

- **Authorization**: `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiWkM2ODAxMmQtZTczMC00fDc3ZTVrezN0aDpva1A9cXEtLT12ezBSfXE4JE9mW3RIIiwiZGV2aWNlX2luZm8iOiJQb3N0bWFuUnVudGltZS83LjQzLjAiLCJ1c2VyX2lwIjoiMTkyLjE2Ny4wLjEwMSIsImlzcyI6InBheWRvaC1iYW5rIiwiZXhwIjoxNzQwMTQ2NTc1fQ.dOIoUNEeByHmXR5sDM8B8VDVWYQGG8FcUeCh5G4AJoU`
  (JWT Bearer token for authentication)

- **Content-Type**: `application/json`
  (Specifies that the request body contains JSON)

### Request Body:
```json
{
  "fcm_token": "12345"
}
```

### Description:
This API request updates the Firebase Cloud Messaging (FCM) token for a specific device using the provided details, including the device's IP, operating system, version, location, and the new FCM token.

---

### Example `curl` Command:
```bash
curl --location '0.0.0.0:4100/api/open/update-fcm-token' \
--header 'X-Device-Ip: 192.168.0.101' \
--header 'X-OS: android' \
--header 'X-OS-Version: 14.0' \
--header 'X-Lat-Long: 92.18,12.31' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiWkM2ODAxMmQtZTczMC00fDc3ZTVrezN0aDpva1A9cXEtLT12ezBSfXE4JE9mW3RIIiwiZGV2aWNlX2luZm8iOiJQb3N0bWFuUnVudGltZS83LjQzLjAiLCJ1c2VyX2lwIjoiMTkyLjE2Ny4wLjEwMSIsImlzcyI6InBheWRvaC1iYW5rIiwiZXhwIjoxNzQwMTQ2NTc1fQ.dOIoUNEeByHmXR5sDM8B8VDVWYQGG8FcUeCh5G4AJoU' \
--header 'Content-Type: application/json' \
--data '{
    "fcm_token":"12345"
}'
```

### Example Success Response:

```bash
{
    "message": "Successfully updated FCM token",
    "status": 200
}
```


### Example Error Response:

```bash
{
    "error": {
        "errors": {
            "body": "device token already exists"
        }
    },
    "message": "Internal Server Error",
    "status": 500
}
```
### Notes:
- Replace the `fcm_token` in the request body with the actual Firebase Cloud Messaging token for the device.
- Ensure the authorization token (Bearer token) is valid and not expired.

```

Let me know if you need anything else!