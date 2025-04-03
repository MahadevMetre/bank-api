# DebitCard Transaction Limit API Documentation

## Overview

The DebitCard Transaction Limit API allows users to set the transaction limits associated with their debit card. This API can be used to fetch the maximum transaction limits, including ATM withdrawals and other transaction types. The API is particularly useful for managing and setting limits on debit card transactions.

## Base URL

`https://staging-api.paydoh.in/bankapi`

## Endpoint

### POST `/api/debitcard/set-txn-limit`

This endpoint sets the debit card transaction limits for a user.

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

## Request Body

The body of the request should include a JSON object with the following structure:

```json
{
    "data": "65adf92bf8ba0217178f8819a9cd97c69dcbd96729849bb596ff14e17f97e37d6035edefe5bfd1f6d6c6b2ab2b749374fb9d5dba5ff0c142f675f84066877a1ad4f78eb22858ec1fc34f796db317891c6c02583c323ba8a53cf383cebfe9399aed73e0fd536fd68d9b4d74b59f0a052e3ad8b7c16a766b4beb20e8fdc2587a2f38060461f5671aff17a183026f2111f9f79f506f4c36cef33fc25c0a3a4d9131a2d7ba7ce83ad18a9dd51b4a722039f0ed032378e562271b0eedae34f5e9632c6566a1b92d558723fb4fdc2440d90b7dcc417091368c0389d350a970c92ab82b8ef3c2618282d3b2bae7d2962eca299dc016878160c95211a62a53d796601eb21b81ee520823548adf61bd6139f99972862c2137cabcb15e20af578ee72bedb266b8bac558535ff972b0aad8c639ccc6582caa06d6995f960c54b4aa45a179172873b32dcccfd917c66d160223361fe476f932245b23fe9d2d6ccc77dde87f411fd07009e39df7d992424c2708222f3b10d5b465f4ec97e3266bdec21b5ba9e915aa9bf6a2820a1c590dd2cfe89a05c9421cde0f27aa0d9e45549df9a575bebda696a7f32a9c69dddb8a3d9c894849a35ae00e8a3aad067828d37b8ad017cb17444aa777acd9589d2740ad00b05520d7d4995d3f6bfc8edc46b5c604a36f8fa8c9d3794e67293549b3deedc742785d702a012db7606c8addaa3b406275913907de084ffe91c5c05b414a24ff5433a121304213b329759aaeb04b0a90aa3994822eacd0aa5c668b1ddbe1bfd176cc232c44194a704f97f3eb4b49149dd643269aacf98839b5084eb65fec0f042f4bbf27e4465bf6620e63c9b0e93a3134db3daa45fb0714519fdc35092090a1dd"
}
```

### Decrypted Request Example

```json
{
    "data": [
        {
            "dchBlkStatus": "0",
            "maxLimit": "5000",
            "name": "ATM Withdrawal",
            "setlimit": "2000",
            "statusDc": "0",
            "tranTypes": "01",
            "transtatus": "0",
            "type": "Domestic"
        }
    ]
}
```

- `data`: A string representing encrypted data related to the user's account.

## Example Request

```bash
curl --location 'https://staging-api.paydoh.in/bankapi/api/debitcard/set-txn-limit' \
--header 'X-Device-IP: test' \
--header 'X-OS: test' \
--header 'X-OS-Version: test' \
--header 'X-Lat-Long: 92.16,12.00' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTJiZjdhMWUtYzZkYS00fDNiaVB0P0N-VHIpWC03NjMwa2xiZ29QQCNGJTIjPztEIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjUxMjg3MTh9.2zb5jjNi4t-XnFaW7ZF73GJHLJ9MqIpsxzl0d7hT9Us' \
--data '{
    "data":"65adf92bf8ba0217178f8819a9cd97c69dcbd96729849bb596ff14e17f97e37d6035edefe5bfd1f6d6c6b2ab2b749374fb9d5dba5ff0c142f675f84066877a1ad4f78eb22858ec1fc34f796db317891c6c02583c323ba8a53cf383cebfe9399aed73e0fd536fd68d9b4d74b59f0a052e3ad8b7c16a766b4beb20e8fdc2587a2f38060461f5671aff17a183026f2111f9f79f506f4c36cef33fc25c0a3a4d9131a2d7ba7ce83ad18a9dd51b4a722039f0ed032378e562271b0eedae34f5e9632c6566a1b92d558723fb4fdc2440d90b7dcc417091368c0389d350a970c92ab82b8ef3c2618282d3b2bae7d2962eca299dc016878160c95211a62a53d796601eb21b81ee520823548adf61bd6139f99972862c2137cabcb15e20af578ee72bedb266b8bac558535ff972b0aad8c639ccc6582caa06d6995f960c54b4aa45a179172873b32dcccfd917c66d160223361fe476f932245b23fe9d2d6ccc77dde87f411fd07009e39df7d992424c2708222f3b10d5b465f4ec97e3266bdec21b5ba9e915aa9bf6a2820a1c590dd2cfe89a05c9421cde0f27aa0d9e45549df9a575bebda696a7f32a9c69dddb8a3d9c894849a35ae00e8a3aad067828d37b8ad017cb17444aa777acd9589d2740ad00b05520d7d4995d3f6bfc8edc46b5c604a36f8fa8c9d3794e67293549b3deedc742785d702a012db7606c8addaa3b406275913907de084ffe91c5c05b414a24ff5433a121304213b329759aaeb04b0a90aa3994822eacd0aa5c668b1ddbe1bfd176cc232c44194a704f97f3eb4b49149dd643269aacf98839b5084eb65fec0f042f4bbf27e4465bf6620e63c9b0e93a3134db3daa45fb0714519fdc35092090a1dd"
}'
```

## Success Response

### Status Code: 200 OK

```json
{
    "message": "OTP Sent Successfully",
    "status": 200
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

### Bad Request
**Status Code: 400 Bad Request**

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

## Additional Notes

1. The data field should always be encrypted before sending in the request body.
2. The Authorization header should contain a valid Bearer Token for the user.
3. All fields are required for successful request processing.

## Security Considerations

1. **HTTPS**: Ensure that the API is accessed over HTTPS in production to secure the communication.
2. **Bearer Token**: Always authenticate API requests using a valid Bearer token.
3. **Encryption**: The data parameter in the request is encrypted. Ensure proper handling of encryption and decryption.
4. **Access Control**: Implement proper access controls to restrict API usage to authorized users only.

## Important Reminders

- This API is designed to help users set and modify their debit card transaction limits.
- Always validate and sanitize input data before processing.
- Implement proper error handling and logging mechanisms.