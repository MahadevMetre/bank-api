# DebitCard Transaction Limit API Documentation

## Overview

The DebitCard Transaction Limit API allows users to retrieve the transaction limits associated with their debit card. This API can be used to fetch the maximum transaction limits, including ATM withdrawals and other transaction types. The API is particularly useful for managing and setting limits on debit card transactions.

## Base URL

https://staging-api.paydoh.in/bankapi

## Endpoint

### POST `/api/debitcard/get-limit-list`

This endpoint retrieves the debit card transaction limits for a user.

## Request Headers

The following headers must be included in the request:

| Header          | Required | Description                                                      |
| -------------   | -------- | -----------------------------------------------------------------|
| `Content-Type`  | Yes      | Must be set to `application/json`                                |
| `X-Device-IP`   | Yes      | Must be set to the IP address of the device making the request   |
| `X-OS`          | Yes      | Must be set to the operating system name (e.g., Android, iOS)    |
| `X-OS-Version`  | Yes      | Must be set to the OS version (e.g., 10.2.1)                     |
| `X-Lat-Long`    | Yes      | Must be set to the latitude and longitude (e.g., `92.16,12.00`)  |
| `Authorization` | Yes      | Bearer Token of the user (format: `Bearer <token>`)              |

## Request Body

The body of the request should include a JSON object with the following structure:

```json
{
    "data": "47b6a8e315adb12c1f139b883c1d2be6fd163b155128a1a800d181e9bfceed1dd3a72ebade81d1de8dc63ba53899ced66172872d01cd8053e615be"
}
```

### Decrypted Request Data

```json
{
  "transaction_type": "Domestic" // Can be Domestic or International
}
```

- `data`: A string representing encrypted data related to the user's account.

## Example Request

```bash
curl --location 'https://staging-api.paydoh.in/bankapi/api/debitcard/get-limit-list' \
--header 'X-Device-IP: test' \
--header 'X-OS: test' \
--header 'X-OS-Version: test' \
--header 'X-Lat-Long: 92.16,12.00' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTJiZjdhMWUtYzZkYS00fDNiaVB0P0N-VHIpWC03NjMwa2xiZ29QQCNGJTIjPztEIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjUxMjg3MTh9.2zb5jjNi4t-XnFaW7ZF73GJHLJ9MqIpsxzl0d7hT9Us' \
--data '{
    "data": "47b6a8e315adb12c1f139b883c1d2be6fd163b155128a1a800d181e9bfceed1dd3a72ebade81d1de8dc63ba53899ced66172872d01cd8053e615be"
}'
```

## Success Response

### Status Code: 200 OK

#### Encrypted Response
```json
{
    "data": "c7b1aa46fe22c5d0817f8700df6a08c733d6623364702eb92c00f4c857cfa2584cc4ae71f8642ec27f1ff482b490a159a75a894e474bfe2ff36e3d92abafcaaa0f61641ede19baed99788138e7b9f7dd6746293c9521705dd4b8f2e5be25b322b0bbc81febbc0b7e05ed29e7ace92e91398fd9b70142920e2135f98cbbf7f2c6884ae0a57ae091e6ccbeb4c268c3d7420fd8368ee45576bad0bfd7d344d57f51f72686fc4691e8d3e532f951fff7c06953719584e5d77d12cf4ae280d33815e8dbb0672e268991d5040eb27aaea6aa3de9e64c4d963a1c66f4fde6d72e48dec2eca958df7c5525bf3f4b9d6da6958c0ea711574c9865ef718391552587ec0639b09877fa0a9ff2bfc7459730066a10a7a265fab2511c0e446107f52844037eda205435a75d1a16518c73d189e1ed37986b418502414ee3c414c9fb46eb235bfcdafa3730bf032a700645310305bda91c3a40bc2d9d0010f08b496464b2d77a21a8fa15b9330f7476c19d73f0b60d3c079962945d3e9e938bd44e6ef07a9d7e063abaacaed1a975eb6dce7aaf13c030c813bb17b5343027b1f0ce846bc300f4040fdbf425ef0f0827184e79056ab849f98b7190545cc58efdec4e7ce35497ba19c9147c6253dbc1772d69f0363f1c9910d62f5e7875114a1944afdf1401f87bdd1e30ca69f41f01c8458e6a0b80bfc1482d562d551d79346e2ee38dff99f93139a34ea95d8b17400f17f1fcf39974086e62a4f01dcec571751c5fd0f216708230086c22f0ee4be0e0dfd90d5a4b5c50e0c40c07129a3362323b67c30dc981ff527aead7452409ca11fa8120137f949faf3566d5946954b354bba782ebf58d8de64cc4f390596309fd183677d39374c0b6eb32d523e4d5fd3fa4488cf61aa048d0e4ee29740aa16ee33482f3277415412daf30aa0712",
    "message": "Successfully Get Transaction Limit",
    "status": 200
}
```

#### Decrypted Response Example
```json
{
    "data": [
        {
            "dchBlkStatus": "0",
            "maxLimit": "5000",
            "name": "ATM Withdrawal",
            "setlimit": "5000",
            "statusDc": "0",
            "tranTypes": "01",
            "transtatus": "0",
            "type": "Domestic"
        }
    ]
}
```

### Response Fields Explanation

| Field         | Description                                                   |
|---------------|---------------------------------------------------------------|
| `dchBlkStatus`| Block status of the debit card (0 for active)                 |
| `maxLimit`    | Maximum allowed transaction limit                             |
| `name`        | The type of transaction (e.g., ATM Withdrawal)                |
| `setlimit`    | The current limit set for the transaction type                |
| `statusDc`    | Status of the debit card (0 for active)                       |
| `tranTypes`   | The transaction type (e.g., 01 for ATM withdrawals)           |
| `transtatus`  | Status of the transaction (0 for active)                      |
| `type`        | Type of the transaction, such as Domestic or International    |

## Error Responses

### 1. Bad Request
**Status Code: 400 Bad Request**

This error occurs when the request body is malformed or missing required parameters.

**Example Response:**
```json
{
    "error": {
        "errors": {
            "body": "Technical issue"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

## Security Considerations

1. **HTTPS**: Ensure that the API is accessed over HTTPS in production to secure the communication.
2. **Bearer Token**: Always authenticate API requests using a valid Bearer token.
3. **Encryption**: The data parameter in the request is encrypted. Make sure to handle encryption and decryption correctly.
4. **Access Control**: Implement proper access controls to restrict API usage to authorized users only.

## Additional Notes

- This API is designed to retrieve debit card transaction limits.
- It helps users set up or modify their transaction limits.
- The data is encrypted to ensure security during transmission.