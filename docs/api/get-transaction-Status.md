# Get Transaction Status API Documentation

## Overview

This API allows you to retrieve transaction status via Transaction ID. The response is encrypted, and it needs to be decrypted to access the actual data.

### Base URL
```
http://0.0.0.0:4100/api/beneficiary/payment-status
```

## Request

### HTTP Method
```
POST
```

### Headers

| Header Name         | Description                            |
|---------------------|----------------------------------------|
| `X-Device-IP`       | The IP address of the device (example: `test`) |
| `X-OS`              | The operating system of the device (example: `test`) |
| `X-OS-Version`      | Version of the OS (example: `test`)    |
| `X-Lat-Long`        | Latitude and Longitude (example: `92.16,12.00`) |
| `Content-Type`      | `application/json`                     |
| `Authorization`     | `Bearer <token>` - Token required to authenticate the request |

### Request Example

```bash
curl --location '0.0.0.0:4100/api/beneficiary/payment-status' \
--header 'X-Device-Ip: 192.167.0.101' \
--header 'X-OS: android' \
--header 'X-OS-Version: 14.0' \
--header 'X-Lat-Long: 92.168,13.01' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiWkM2ODAxMmQtZTczMC00fDc3ZTVrezN0aDpva1A9cXEtLT12ezBSfXE4JE9mW3RIIiwiZGV2aWNlX2luZm8iOiJQb3N0bWFuUnVudGltZS83LjQzLjAiLCJ1c2VyX2lwIjoiMTkyLjE2Ny4wLjEwMSIsImlzcyI6InBheWRvaC1iYW5rIiwiZXhwIjoxNzM4NzYwMzQxfQ.s2wpJndq0FdVA2_K1T9Wb0anjbv6IcGfqwqTgG-dZ9I' \
--data '{
    "data":"b33039d79dad02a83825b3356ab47191fd31b6951a3a77f24c67db0bde60a5aea90a192d2d0445439509958355e8338011cc8fb34f93f5"
}'
```

### Response

The response will contain the encrypted data. Here is an example of the response:

```json
{
    "data": "6b3ba9435a157ffee6dba84bd520e4025639ac80f0f178eef134e3a1f7f7edf032291361ad10bd56ab0e0d47cace5bb4eb242a54e1f0f9260c812279623e30e1d97d5fcf3baef5931158196a80c4e595ee835da970b90b5c80caf066004dcfdb458bb6bc5113f3cda15d5ce25d2d04bf1c74ef6a4cc1001c075ee4410197f9eed26b9c7dd5e228eb090a7c3ff9a6c4fec63da0437ba78e41a8baf34803812647322ed28e35b34839457f0157f5d52b75616c5e73920f00",
    "message": "successfully get transaction status",
    "status": 200
}
```

Below data should be pass as encrypted formate
//{
//    "txnId":"b53d9e48-12d7-4"
//}

### Decrypted Response

Once decrypted, the response data will contain the following fields:

```json
{
            "PaydohTransactionId": "a3ce3792-ccc7-426c-b817-29c70654d373",
            "TxnAmount": "1.00",
            "TxnIdentifier": "c6935903-ec7e-4",
            "TxnRefNo": "61d2936b-207a-416f-9469-4ef5808d1d8f",
            "TxnStatus": "Success"
        }
```

```json
{
            "PaydohTransactionId": "db1f2f64-b2c0-4ad6-8de3-4c82e15a64c0",
            "TxnAmount": "12360.00",
            "TxnIdentifier": "c694336c-a3c4-4",
            "TxnRefNo": "",
            "TxnStatus": "Failure"
        }
```

### Response Fields

| Field Name                | Description                                                   |
|---------------------------|---------------------------------------------------------------|
| `TxnIdentifier`               | This is transaction ID            |
| `TxnRefNo`             | This Transaction Refrence ID.                              |
| `TxnStatus`      | The status of the Transaction         |                         |
| `PaydohTransactionId`      | The Paydoh Transaction ID         |                         |
| `TxnAmount`      | The Transaction Amount         |                         |

## Authentication

The API requires an authentication token to access the data. You must include a valid Bearer token in the `Authorization` header.

### Example Token

```
Authorization: Bearer <your-token>
```

## Error Body
```
{
    "error": {
        "errors": {
            "body": "no data found"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

## Error Handling

In case of any errors, the API will return an appropriate HTTP status code along with a message indicating the nature of the error.


---


### Notes:
- Ensure that you decrypt the `data` field using the correct decryption mechanism before displaying or using the sensitive data.
- The API requires proper handling of the encrypted data and token for security purposes.
