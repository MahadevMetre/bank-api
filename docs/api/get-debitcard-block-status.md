# Get Debit Card Block Status API

## Overview
REST API endpoint to retrieve the block status of debit cards.

**Base URL**: `https://staging-api.paydoh.in/bankapi`
**Endpoint**: `GET /api/debitcard/get-card-status`

## Request Headers

| Header | Required | Description |
|--------|----------|-------------|
| `Content-Type` | Yes | `application/json` |
| `X-Device-IP` | Yes | Device IP address |
| `X-OS` | Yes | Operating system (Android, iOS) |
| `X-OS-Version` | Yes | OS version (e.g., 10.2.1) |
| `X-Lat-Long` | Yes | Location coordinates (e.g., 92.16,12.00) |
| `Authorization` | Yes | Bearer token |

## Request Example

```bash
curl --location '0.0.0.0:4100/api/debitcard/get-card-status' \
--header 'X-Device-IP: test' \
--header 'X-OS: test' \
--header 'X-OS-Version: test' \
--header 'X-Lat-Long: 92.16,12.00' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIs...' \
--data ''
```

## Responses

### Success Response (200 OK)

The API returns an encrypted response that must be decrypted.

**Encrypted:**
```json
{
    "data": "0b94e36b5ec92c73872908f03655f772...",
    "message": "Status retrieved successfully",
    "status": 200
}
```

**Decrypted:**
```json
{
    "domestic_block_status": "0",
    "international_block_status": "0",
    "is_permanently_blocked": "0"
}
```

#### Response Fields

| Field | Values | Description |
|-------|--------|-------------|
| `domestic_block_status` | 0/1 | Card status for domestic transactions |
| `international_block_status` | 0/1 | Card status for international transactions |
| `is_permanently_blocked` | 0/1 | Permanent block status (if 1, cannot be unblocked) |

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

## Security Requirements

1. Use HTTPS in production
2. Include valid Bearer token
3. Implement proper encryption/decryption
4. Use appropriate access controls
5. Validate all inputs
6. Implement rate limiting
7. Monitor API usage

## Troubleshooting

**Authentication Issues:**
- Verify token validity and format
- Check token expiration

**Request Errors:**
- Confirm all required headers
- Validate header format

**Response Handling:**
- Implement decryption
- Handle all status codes
- Log errors appropriately