
### /api/authentication/initiate-sim-verification/

#### POST
##### Summary:

Api get message token by sending encrypted data and crc value

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| Authorization | header | With the bearer started | Yes | string |
| authenticationRequest | body | Authentication Request | Yes | [requests.AuthenticationRequest](https://github.com/Grapesberry-Technologies-Pvt-Ltd/bank-api/blob/bf67ab86f88cb0f4eaa5e5ed1a06496dee4e4f01/requests/authentication.go#L14) |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ok |

### /api/authorization/

#### POST
##### Summary:

Api get authorization Token for all the api

##### Description:

Using the mobile_number get the authorization token for the other apis.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| authorizationRequest | body | Authorization Request | Yes | [requests.AuthorizationRequest](https://github.com/Grapesberry-Technologies-Pvt-Ltd/bank-api/blob/97fb1518cf28eab83169eec6c9f3ec884642715f/requests/authorization.go#L9) |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ok |

### /api/webhook/route-mobile

#### GET
##### Summary:

Web hook for getting api senders mobile data from Route mobile.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| sender | query | Sender | Yes | string |
| message | query | Message | Yes | string |
| operator | query | Operator | No | string |
| circle | query | Circle | No | string |

##### Responses

| Code | Description |
| ---- | ----------- |
| 200 | ok |

### Models


#### requests.AuthenticationRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| crc_value | string |  | Yes |
| data | string |  | Yes |

#### requests.AuthorizationRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| mobile_number | string |  | Yes |