# DebitCard Set Block-Unblock Status API Documentation

## Overview

This API enables users to manage their debit card status by setting block, unblock, or permanent block configurations. This API will send the OTP and After OTP verification Your DebitCard Status will Update.

**Base URL**:`https://staging-api.paydoh.in/bankapi`

**Endpoint**:`POST /api/debitcard/set-card-status`

Sets the block or unblock status of debit cards.

## Request Headers

| Header          | Required | Description                                                      |
|-----------------|----------|------------------------------------------------------------------|
| `Content-Type`  | Yes      | Must be set to `application/json`                                |
| `X-Device-IP`   | Yes      | Must be set to the device IP address                            |
| `X-OS`          | Yes      | Must be set to the OS name (e.g., Android, iOS)                 |
| `X-OS-Version`  | Yes      | Must be set to the OS version (e.g., 10.2.1)                    |
| `X-Lat-Long`    | Yes      | Must be set to the latitude and longitude (e.g., `92.16,12.00`) |
| `Authorization` | Yes      | Bearer Token format: `Bearer <token>`                            |

## Request Body

### Encrypted Request Format
```json
{
    "data": "2f4e956b346063ab5770c2cd8f977e5f83a84a0455f2120502751eb9c2e24019a87b71dac17c6db16c93956ba352cf1fe9a201c1debb70528c617fd066d166e5b443e3a7c1f6eba6a014d87fe6720292ef4308ee0f8c0024f8f3d59cea3588e3a076d559e2374baae0a0df57962b225cd71ece4f992ab1"
}
```

### Decrypted Request Format
```json
{
    "domestic_block_status": "0",
    "international_block_status": "0",
    "is_permanently_blocked":"0"
}
```

### Request Parameters

| Field | Values | Description |
|-------|--------|-------------|
| `domestic_block_status` | 0/1 | Card status for domestic transactions |
| `international_block_status` | 0/1 | Card status for international transactions |
| `is_permanently_blocked` | 0/1 | Permanent block status (if 1, cannot be unblocked) |

## Responses

### Success Response (200 OK)
```json
{
    "message": "OTP Sent Successfully",
    "status": 200
}
```

### Error Response For Permanently Blocked Cards (400 Bad Request)
```json
{
    "error": {
        "errors": {
            "body": "your debit card is permanently blocked"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

### Error Response (400 Bad Request)
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

1. Always use HTTPS in production
2. Include valid Bearer token
3. Handle request data encryption/decryption properly
4. Implement proper access controls

## Best Practices

1. Validate block status values (0/1 only)
2. Implement rate limiting
3. Monitor for suspicious activities
4. Log errors for debugging
5. Rotate security keys regularly

## Implementation Notes

- API requires encrypted request data
- Both domestic and international status can be set independently
- Block status values: 0 (unblock), 1 (block)
- Permanent block cannot be reversed

## Troubleshooting

1. **Authentication Issues**
   - Verify Bearer token validity and format
   - Check token expiration

2. **Encryption Issues**
   - Verify encryption algorithm
   - Validate data structure

3. **Request Failures**
   - Check all required headers
   - Validate request body format