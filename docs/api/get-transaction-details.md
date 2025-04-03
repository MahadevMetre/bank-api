# Get Transaction Details API Documentation

## Overview

This api is used to get transaction details as per given input parameters.

## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/transaction/history/get-details
```

#### Request Headers

| Header       | Required | Description                       |
| ------------ | -------- | --------------------------------- |
| Content-Type | Yes      | Must be set to `application/json` |
| X-Device-Ip  | Yes      | Must be set to added              |
| X-OS         | Yes      | Must be set to added              |
| X-OS-Version | Yes      | Must be set to added              |
| X-Lat-Long   | Yes      | Must be set to added              |
| Bearer       | Yes      | Must be a Bearer Token of user.   |

#### Example Request

```curl
curl --location 'http://localhost:4100/api/transaction/history/get-details' \
--header 'Content-Type: application/json' \
--header 'X-Device-Ip: 192.168.1.4' \
--header 'X-Lat-Long: 19.0935922,72.9163083' \
--header 'X-Os: android' \
--header 'X-Os-Version: 12' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNmYzZTYxYzktNjI5Ni00fENdLSNJW05CRzlcdTAwMjZMMipTR3ouQ14jTWNRV31IYlBCYUgiLCJpc3MiOiJwYXlkb2gtYmFuayIsImV4cCI6MTc2MjIzMTI2M30.uPCYjCpzlkZRd1ON4FGAXiM2VXPv8Mt1ftV52bYXjW0' \
--data '{
    "CodeDRCR": "D"
    "TransactionAmount": "1500.00"
    "TransactionDate": "13-11-2024 15:11:19"
    "TransactionDate": "07-11-2024 14:17:06",
    "TransactionDescription": "UPI-DR-431245961608-PAVAN KHAIRE-KVBL-1770155000072605-HK"
}'
```

#### Success Response

**Status Code:** 200 OK

```json
{
  "data": {
    "id": "ee4fa7f6-4b9a-40db-9229-74f475b29998",
    "user_id": "77e7924e-ced8-4",
    "transaction_id": "PAYDOHb3a3c0760b02073d73bf62cccefb",
    "transaction_desc": "testing",
    "beneficiary_id": "",
    "payment_mode": "UPI",
    "amount": "1.00",
    "utr_ref_number": "",
    "upi_payee_addr": "j8108.paydoh@kvb",
    "upi_payee_name": "",
    "beneficiary_name": "",
    "beneficiary_ifsc": ""
  },
  "message": "successfully fetched transaction details.",
  "status": 200
}
```

#### Error Responses

##### Internal Server Error

**Status Code:** 500 Service Unavailable

```json
{
  "error": {
    "errors": {
      "body": ""
    }
  },
  "message": "Internal Server Error",
  "status": 500
}
```

## Support

For any issues or queries, please contact:

- Email: indal@paydoh.money
