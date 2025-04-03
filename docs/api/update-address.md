# Update Address API Documentation

## Overview

This API allows users to update their address in the bank's system.

---

## Base URL

https://staging-api.paydoh.in/bankapi


---

## Endpoint

### POST `/api/address/update`

This endpoint is used to update a user's address.

---

## Request Headers

The following headers must be included in the request:

| Header          | Required | Description                                                      |
|-----------------|----------|------------------------------------------------------------------|
| `Content-Type`  | Yes      | Must be set to `application/json`                                |
| `X-Device-IP`   | Yes      | Must be set to the IP address of the device making the request   |
| `X-OS`          | Yes      | Must be set to the operating system name (e.g., Android, iOS)    |
| `X-OS-Version`  | Yes      | Must be set to the OS version (e.g., 10.2.1)                     |
| `X-Lat-Long`    | Yes      | Must be set to the latitude and longitude (e.g., `92.16,12.00`)  |
| `Authorization` | Yes      | Bearer Token of the user (format: `Bearer <token>`)              |

---

## Request Body

The body of the request should include a **Multipart Form Data** object with the following structure:

| Field            | Description                                                         |
|------------------|---------------------------------------------------------------------|
| `address1`       | First line of the address (e.g., `903, 9th floor`)                  |
| `address2`       | Second line of the address (e.g., `Ajmera Sikova, Near Ashok Mills`)|
| `address3`       | Third line of the address (e.g., `Ghatkopar west`)                  |
| `city`           | City (e.g., `Mumbai`)                                               |
| `pincode`        | Postal code (e.g., `400086`)                                        |
| `state`          | State (e.g., `Maharashtra`)                                         |
| `country`        | Country (e.g., `India`)                                             |
| `addressType`    | Address type (e.g., `communication`)                                |
| `address_proof`  | Type of address proof (e.g., `AadharCard`)                          |
| `document_front` | Front image of the address proof (file upload)                      |
| `document_back`  | Back image of the address proof (file upload)                       |

---

## Example Request

```bash
curl --location 'https://staging-api.paydoh.in/bankapi/api/address/update' \
--header 'X-Device-IP: 192.0.0.4' \
--header 'X-OS: android' \
--header 'X-OS-Version: 10' \
--header 'X-Lat-Long: 97.15,78.26' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTJiZjdhMWUtYzZkYS00fDNiaVB0P0N-VHIpWC03NjMwa2xiZ29QQCNGJTIjPztEIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjQ2NzIxMzV9.fVADFNE5r-unCOV9oMDKPf2VRSTxfag0hOZQ-3sZkjs' \
--form 'address1="903, 9th floor"' \
--form 'address2="Ajmera Sikova, Near Ashok Mills"' \
--form 'address3="Ghatkopar west"' \
--form 'city="Mumbai"' \
--form 'pincode="400086"' \
--form 'state="Maharashtra"' \
--form 'country="India"' \
--form 'addressType="communication"' \
--form 'document_front=@"/home/abcd/Desktop/WhatsApp Image 2024-11-25 at 11.41.35 AM.jpeg"' \
--form 'document_back=@"/home/abcd/Desktop/WhatsApp Image 2024-11-25 at 11.41.35 AM.jpeg"' \
--form 'address_proof="AadharCard"'


Success Response
Status Code: 200 OK

{
    "data": {
        "Req_Ref_No": "4908P2501002137",
        "ErrorCode": "00",
        "ErrorMessage": "SUCCESS"
    },
    "message": "Address update request has been submitted successfully",
    "status": 200
}

Error Responses
1. Bad Request
Status Code: 400 Bad Request
{
    "error": {
        "errors": {
            "body": "Address Update Request already in process"
        }
    },
    "message": "Bad Request",
    "status": 400
}

Notes
The Authorization header should contain a valid Bearer token for the user.
All fields in the request are required for successful request processing.
Make sure that the request is sent as multipart/form-data.


Security Considerations
HTTPS: Ensure that the API is accessed over HTTPS to secure the communication.
Bearer Token: Always authenticate API requests using a valid Bearer token.

Additional Information
This API is used to update the user's address in the bank system.
The address proof can be an image file uploaded as part of the request (front and back images of the document).