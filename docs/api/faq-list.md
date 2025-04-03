### Get FAQ List

## Overview

Retrieves the list of frequently asked questions based on the specified platform.

## Base URL

https://staging-api.paydoh.in/bankapi

## API Endpoint

**URL**: `/api/faq/list`

**Method**: `GET`

**Query Parameters**:

- `platform` (optional): Specifies the platform for which to retrieve FAQs. Possible values:
  - `web`: Retrieves FAQs for web platform
  - `ios`: Retrieves FAQs for iOS platform
  - `android`: Retrieves FAQs for Android platform
    If not specified, defaults to app-specific FAQs (iOS or Android).

**Authentication**: Required

**Response**:

- Status Code: 200 OK
- Body: JSON object containing an array of FAQ categories and their corresponding questions and answers.

**Example Request**:

```
GET http://localhost:4100/api/faq/list?platform=web
```

### Request Headers

To successfully call the API, the following headers are required:

| Header          | Required | Description                                                                         |
| --------------- | -------- | ----------------------------------------------------------------------------------- |
| `Content-Type`  | Yes      | Set to `application/json`                                                           |
| `X-Device-IP`   | Yes      | The device's IP address (e.g., `192.0.0.2`)                                         |
| `X-OS`          | Yes      | The operating system of the device (e.g., `ios`, `android`)                         |
| `X-OS-Version`  | Yes      | The version of the operating system (e.g., `18.0`)                                  |
| `X-Lat-Long`    | Yes      | Latitude and longitude of the device (e.g., `19.093377739232523,72.91605446814229`) |
| `Authorization` | Yes      | Bearer token for user authentication                                                |

### Example Request

```bash
curl --location 'http://localhost:4100/api/faq/list?platform=app' \
--header 'X-Device-Ip: test' \
--header 'X-Os: android' \
--header 'X-Lat-Long: test,test' \
--header 'X-Os-Version: aerawe' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOWMyMzM3OTYtNDc4NC00fHdiKVtTP0Y6cn4pdzttUzdHK009SVx1MDAzY1tFdTJ9aTZiOS8iLCJpc3MiOiJwYXlkb2gtYmFuayIsImV4cCI6MTczNDk1ODQ2NX0.xQJBLSJFMznbH-4eZ-NIDAeqsxXFiml7HsDh8cI5Lt0'
```

**Example Response**:

### Status Code: 200 OK

```json
{
  "data": [
    {
      "id": "8596a9a2-c196-4d33-9183-f1db2ffcf639",
      "name": "Split bill and Budget book",
      "faqs": [
        {
          "id": "3524e5cb-4ee2-4eb8-8a96-a73159d2b218",
          "question": "How does the split bill feature work?",
          "answer": "The split bill feature allows you to divide expenses among friends easily.",
          "video_url": ""
        }
      ]
    },
    {
      "id": "c30d9569-c5ee-4f3f-95ab-2605aacce1de",
      "name": "Referral",
      "faqs": []
    }
  ],
  "message": "successfully fetched faq list",
  "status": 200
}
```

**Error Responses**:

- 401 Unauthorized: If the user is not authenticated
- 500 Internal Server Error: If there's an error retrieving the FAQ list
