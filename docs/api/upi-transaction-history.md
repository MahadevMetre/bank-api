# Upi Transaction History API Documentation

## Overview

This API is used to get all the upi transaction history for the user given specific time period. 
## Base URL

```
https://staging-api.paydoh.in/bankapi
```

## Endpoints

```
POST /api/upi/transaction-history
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
curl --location 'https://staging-api.paydoh.in/bankapi/api/upi/transaction-history' \
--header 'X-Device-Ip: 192.168.1.3' \
--header 'X-OS: android' \
--header 'X-OS-Version: 10' \
--header 'X-Lat-Long: 19.0933218,72.9161727' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNzdlNzkyNGUtY2VkOC00fGQ1c3poNSs9P0owczpjaGlaLS8kQTloI0BEVFBnUUtsIiwiaXNzIjoicGF5ZG9oLWJhbmsiLCJleHAiOjE3NjE3MTc5Njh9.QIvofLOxCbv3ERP5MSQpA_mQJsOUku6IshGkDaGmiU8'\
--data '{

    "data":
    "238dfdfd081c9118b4f71fca6c7e4126e5468df52a19a33bc060c7953931b70687baa846106858063214830117372ac50a2d0b9855869111d56b9d4762bbadf4baea8af3669972ecb67aee4bd1"

}'
```

#### Success Response

**Status Code:** 200 OK

```json
{
    "data": "7572718f9fc17ba310aee85b03d88940e0b6fcd019b84af67b383299687f4062e465f72940fe9631b1ca7da0913f77564fa93652f10a17c8d773d45ab591c06d5d5f748cc87b64f1d46360286fb83feb5d255c5b752df98c00d42e9c31fe3b4916d3beb88c2bfffae4c13d60555ad44a772ade9556dd24d1b2917dc26a793cfe87e15a78fbf4147c0cd1980c4df8d55b01586d753880e1fc99161e4f4e3cc5906e1897a7ce1a6a87366d669faf8a713612ddf4bffa00152fe06fbdb32ef0bcbbe317088c441e3654e46f4b97b3379aa67239f6e96035cfc35375bb3496088ae59302480158441569cd61cced59df021136ea5d630ec77705dc4da09b54fc711ca097563b4fbabdf251351289266841a07d380c97fbb5f0d19ec85e65253ea064ac7bb1e9395e457346f524b2e422d31eadbfe296a553a4d9256fe539547a126ab36ddd0d49d9510306cfb10d6a7175beb10ba80dadf94448e5363e2b5cf23f7f7c6349845229474e2ac53aa4fb2961b08ae00e9182188280799b53219cf6aa7cbe5a89ee4728f7f6f2b87fda74c0e980baadc7f2c8bb6ee6100fe74aa0b65b26003c96bb6592efd44adb7a6c8bde30ffacf42a486f450aec09cebce428daad9d3a708c5bd51b093bbc9475da744df6f5a78403c65fc34239ad34f0ee01bb3f6e1871759e03ae2f38f84957f705b4c2410c6d48d844743dda017ba867668ecd064ea13e8404527abb15258a8388ab550e96841c72647cd2c1092b5263f46626b26d464e2d9df7c71571effcfea7bede97aa38891135544f4a30818ec94f757f62a594e5fb4c5e8b59d4a1867fcf0cc9ed5c8d2bce514ce2e55267ce7cf59b25cb336409acf8b05a08f16fff07a15dfea36c1d0e7dbab17667b7d796dc6911b7d6a6ac5cc95a161af07ccba640bd502f2bcc4fc222aa9f55b08da122195ad7ffd8f958b06e1d2e5b3e016664d67f662bd48599ce0a9488b9045555ca181644cb5f8dc86afa7afdde01dc32b1c35105a59031a32965b76a1f4b77d3100cf196b463d56a8b0a8f7618242b24b312bb8bd1cb46ff6026fd08d25a53090059f7ce88d723d8f5387ad6f574de2769014ec76f333c8795b30eceab1d0000ad3776787e3909fb18d86a817e1e0c46bc673a4e03a4fe2a55d1124114e443e8715925a1de",
    "message": "Successfully got the upi transaction history details",
    "status": 200
}


#### Error Responses

##### BadRequest Server Error 

**Status Code:** 400 Service Badrequest

#### Case 1 if fromdate is missing

```json
{
    "error": {
        "errors": {
            "body": "FromDate is required"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Case 2 if todate is missing

```json
{
    "error": {
        "errors": {
            "body": "ToDate is required"
        }
    },
    "message": "Bad Request",
    "status": 400
}
```

#### Security Considerations

1. The endpoint should be accessed over HTTPS in production
2. Implement rate limiting to prevent brute force attacks
3. Monitor failed attempts.

## Support

For any issues or queries, please contact:

- Email: yeshwanth@paydoh.money
