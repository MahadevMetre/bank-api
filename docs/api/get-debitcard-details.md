# Get DebitCard Details API Documentation

## Overview

This API allows you to retrieve details about a user's debit card. The response is encrypted, and it needs to be decrypted to access the actual data.

### Base URL
```
http://0.0.0.0:4100/api/debitcard/detail
```

## Request

### HTTP Method
```
GET
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
curl --location '0.0.0.0:4100/api/debitcard/detail' \
--header 'X-Device-IP: test' \
--header 'X-OS: test' \
--header 'X-OS-Version: test' \
--header 'X-Lat-Long: 92.16,12.00' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <your-token>'
```

### Response

The response will contain the encrypted data. Here is an example of the response:

```json
{
    "data": "7833f2fa1654aa4b914b46890c59d2ec302cdbc82fdb820047ad166a42e3e123dca8b61a7ef6b75048c023ccb48acb404d2b431df1d24512f14cc58112f2ebfa867153ceecbcd2b3d884efd3fbdc877c46890007d27a438cf7546e1671eee3f4c93cb7682eb836b5d671bdc3118b4837db1e3555426b84540cdc14a4d86f184f8207e0dddd010f39c90995d106587d40cad0d30e835231d6d4cba5bab9bce63333ebffc965d4a392f6202e6d2c620f99cf7a354dfa799fdb65b5c954167958f40222e5c6a506268ebce9cf8e163105a3c916497878a47f5d38318ff51a5b47558c545237f032bcc9ed9b6e4e034104c5c60cc54e6e599d0b2aa9f1bc7a5fc492d1121ece6ce299940e270404078889d06d747fd0f33c970b20f4f402f7903d389bc5edd9e739c5ff1bfd6363d3ea34bfc3a5e7a55917f5f294c600f9a98f876219d3e091c9b5e491b0fa3bb07e971036",
    "message": "Successfull",
    "status": 200
}
```

### Decrypted Response

Once decrypted, the response data will contain the following fields:

```json
 {
        "AccountNo": "1770155000072572",
        "ApplicantId": "PAYDOH77e801",
        "cardIssuanceStatus": 1,
        "cardholderName": "JAKIR SAYYED",
        "cvvValue": "445",
        "encryptedPAN": "6079410000001135",
        "expiryDate": "10/31",
        "has_physical_debit_card": false,
        "has_virtual_debit_card": false,
        "is_permanently_blocked": false,
        "proxyNumber": "123436778"
    }
```

### Response Fields

| Field Name                | Description                                                   |
|---------------------------|---------------------------------------------------------------|
| `AccountNo`               | The account number associated with the debit card.            |
| `ApplicantId`             | The unique ID for the applicant.                              |
| `cardIssuanceStatus`      | The issuance status of the card (e.g., 1 for issued).         |
| `cardholderName`          | Name of the cardholder.                                       |
| `cvvValue`                | The CVV value of the debit card.                              |
| `encryptedPAN`            | The encrypted primary account number (PAN).                   |
| `expiryDate`              | The expiration date of the debit card (MM/YY).                |
| `has_physical_debit_card` | Boolean flag indicating if the user has a physical debit card.|
| `has_virtual_debit_card`  | Boolean flag indicating if the user has a virtual debit card. |
| `is_permanently_blocked`  | Boolean flag indicating if the debit card is permanently blocked. |
| `proxyNumber`             | Proxy number for the debit card.                              |

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
