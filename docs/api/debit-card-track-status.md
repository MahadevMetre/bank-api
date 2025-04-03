# Track Debit Card Status

Retrieves the current status of a user's debit card.

**URL** : `/api/debitcard/track-status`

**Method** : `GET`

**Auth required** : YES

#### Request Headers

| Header       | Required | Description                       |
| ------------ | -------- | --------------------------------- |
| Content-Type | Yes      | Must be set to `application/json` |
| X-Device-Ip  | Yes      | Must be set to added              |
| X-OS         | Yes      | Must be set to added              |
| X-OS-Version | Yes      | Must be set to added              |
| X-Lat-Long   | Yes      | Must be set to added              |
| Bearer       | Yes      | Must be a Bearer Token of user.   |

## Success Response

**Code** : `200 OK`

**Content example**

```json
{
  "data": "173ca289e95399e65c0f6b1d9b71a576f970a6070d77fc7fe5ad80b49db795abd11afbe6a7d10b32353329e31ef641d5bbe8d4a6d86ebd290db7da37118d3fd3f5fcac4dbcd541f0fee2ff9213b7e1a71730a2b6e8b8eb67a8d844feb7fafc6ae7b834ed731a59ce06b7a30bd2bd41138110151a69cf3fea8841e1603a589ba231a729fd14e775e2c1b13a26a2c9e8c317e2a752b662b8ebc3fdf2a1d6faf54932e1e069fcd7040e84e706ee237e61337822c63c5c273020324bf26f379641b766f5af7a58935d214477eaa5b7fb484684368516eafcdb7e4ac95649f8ed21d0b3ff02a927d6968d155cf8fa3a52c21828c88d01a9f2a1d8f90e24c2113c299e65c417147e1aa43e20a3189f78abc44048cace2146b38597f11c80",
  "message": "Debit Card Status Retrieved Successfully",
  "status": 200
}
```

```json, example decrypted data:-
{
  "awb": "32627543012",
  "card_holder_name": "PAVAN KHAIRE",
  "card_number": "4214 XXXX XXXX 1234",
  "card_type": "CTB_VISA_PLATINUM",
  "date": "5/4/2024",
  "dispatch_date": "5/6/2024",
  "dispatch_mode": "BLUEDART",
  "dispatch_status": "DISPATCH",
  "reference_number": "150858937"
}
```

## Error Responses

**Status Code:** 400 Service Badrequest

```json
{
  "error": {
    "errors": {
      "body": "bad request"
    }
  },
  "message": "Bad Request",
  "status": 400
}
```

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

## Notes

- The dispatch status can vary depending on the current state of the card.
- This endpoint requires authentication and will return an error if the user is not authorized.
