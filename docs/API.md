# API Documentation

<!-- TOC -->

- [API Documentation](#api-documentation)
    - [1. Join Waiting Room](#1-join-waiting-room)
    - [2. List Devices in Waiting Room](#2-list-devices-in-waiting-room)
    - [3. Check for profile](#3-check-for-profile)

<!-- /TOC -->

## 1. Join Waiting Room

This endpoint enables devices, which have obtained credentials (Polaris device key, device ID, and organization ID) via FDO T02, to join a virtual waiting room.

**POST** `/v1/device/{org_id}/waiting_room`

### Query Parameters

Parameter | Location | Required | Default | Description
--------- | ------- |------- |------- | -----------
Authorization | header | yes | n/a | Bearer token signed using the Polaris device key obtained through FDO T02
org_id | URL | yes | n/a | Organization ID obtained through FDO T02
device_id | Authorization header | yes | n/a | Device ID obtained through FDO T02

### Response Codes
Code | Description
---- | -----------
200 | Success
401 | Authorization failed
404 | Resource not found


## 2. List Devices in Waiting Room

List all devices in waiting room.

**GET** `/v1/device/{org_id}/waiting_room`

### Query Parameters

Parameter | Location | Required | Default | Description
-------- | ------- |------- |------- | -----------
Authorization | header | yes | n/a | Bearer token signed using the Polaris device key obtained through FDO T02
org_id | URL | yes | n/a | Organization ID
device_id | body | yes | n/a | Device ID received during FDO T02

### Response body

```json
[{
  "device_id": "string",
  "org_id": "string",
  "hostname": " string",
  "first_seen": 11111111,
  "last_seen": 11111111,
  "properties": "JSON object"
}]
```

### Response Codes
Code | Description
--------- | -----------
200 | Success
401 | Authorization failed
404 | Resource not found

## 3. Check for profile

This endpoint is used by devices, that have joined the waiting room, to check if a Profile has been assigned. The device shall periodically check this endpoint until a profile is assigned.

**POST** `/v1/device/{org_id}/profile`

### Query Parameters

Parameter | Location | Required | Default | Description
--------- | ------- |------- |------- | -----------
Authorization | header | yes | n/a | Bearer token signed using the Polaris device key obtained through FDO T02
org_id | URL | yes | n/a | Organization ID obtained through FDO T02
device_id | Authorization header | yes | n/a | Device ID obtained through FDO T02

### Response body

If the endpoint HTTP return code is 203, then a profile has not been assigned and the device should continue to periodically check this endpoint until a Profile URL is provided.

```json
{
    "profile_url": "https://github.com/amazingprofiles/unixprofile"
}
```

### Response Codes
Code | Description
---- | -----------
200 | Success
203 | Success, no profile assign
401 | Authorization failed
404 | Resource not found

## 3. Google login check

Used by CLI to poll RMS to check if login was successful. Responses should be success, failed, pending. Success results in access & refresh tokens.

**POST** `user/login/{login_id}`

### Query Parameters

Parameter | Location | Required | Default | Description
--------- | ------- |------- |------- | -----------
user_id | query | yes | n/a | User ID is the user's Gmail address.

### Response body

```json
{
  "user_code": 212123123,
  "URL": "https:/asdfasdf.com"
}
```

### Response Codes
Code | Description
--------- | -----------
200 | Success
403 | Failed, blah blah blah

## 4. Keep this as an example

This is an API example.

**POST** `user/example/{data}`

### Query Parameters

Parameter | Location | Required | Default | Description
--------- | ------- |------- |------- | -----------
user_id | query | yes | n/a | User ID is the user's Gmail address.
group_id | body | yes | n/a | If set to true, the result will also include cats.
swarm_id | header | yes | n/a | If set to true, the result will also include cats.
swarm_id | query | yes | n/a | If set to true, the result will also include cats.
swarm_id | URL | yes | n/a | If set to true, the result will also include cats.

### Response body

```json
[
  {
    "id": 1,
    "name": "Fluffums",
    "breed": "calico",
    "fluffiness": 6,
    "cuteness": 7
  },
  {
    "id": 2,
    "name": "Max",
    "breed": "unknown",
    "fluffiness": 5,
    "cuteness": 10
  }
]
```

### Response Codes
Code | Description
--------- | -----------
200 | Success
403 | If set to false, the result will include kittens that have already been adopted.