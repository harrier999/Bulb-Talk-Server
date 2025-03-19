# Bulb Talk API 및 WebSocket 문서

## 목차
1. [개요](#개요)
2. [인증](#인증)
3. [REST API](#rest-api)
   - [사용자 관리](#사용자-관리)
   - [친구 관리](#친구-관리)
   - [채팅방 관리](#채팅방-관리)
4. [WebSocket](#websocket)
   - [연결 방법](#연결-방법)
   - [메시지 형식](#메시지-형식)
   - [이벤트 유형](#이벤트-유형)
5. [데이터 모델](#데이터-모델)
6. [오류 처리](#오류-처리)

## 개요

Bulb Talk 서버는 REST API와 WebSocket을 통해 클라이언트와 통신합니다. REST API는 사용자 관리, 친구 관리, 채팅방 관리 등의 기능을 제공하며, WebSocket은 실시간 채팅을 위해 사용됩니다.

기본 URL: `https://api.wasabi-labs.com` (예시)

## 인증

대부분의 API 엔드포인트는 인증이 필요합니다. 인증은 JWT 토큰을 사용하며, 토큰은 로그인 시 발급됩니다.

인증이 필요한 요청의 경우, HTTP 헤더에 다음과 같이 토큰을 포함해야 합니다:

```
Authorization: Bearer {token}
```

## REST API

### 사용자 관리

#### 회원가입

```
POST /signup
```

**요청 본문**:
```json
{
  "username": "사용자이름",
  "password": "비밀번호",
  "phoneNumber": "전화번호",
  "countryCode": "국가코드"
}
```

**응답**:
```json
{
  "success": true,
  "user": {
    "id": "사용자ID",
    "username": "사용자이름",
    "phoneNumber": "전화번호",
    "countryCode": "국가코드"
  }
}
```

#### 로그인

```
POST /login
```

**요청 본문**:
```json
{
  "phoneNumber": "전화번호",
  "password": "비밀번호"
}
```

**응답**:
```json
{
  "success": true,
  "token": "JWT 토큰"
}
```

#### 인증번호 요청

```
POST /authenticate
```

**요청 본문**:
```json
{
  "phoneNumber": "전화번호",
  "countryCode": "국가코드",
  "deviceId": "기기ID"
}
```

**응답**:
```json
{
  "success": true
}
```

#### 인증번호 확인

```
POST /checkauth
```

**요청 본문**:
```json
{
  "phoneNumber": "전화번호",
  "countryCode": "국가코드",
  "deviceId": "기기ID",
  "authNumber": "인증번호"
}
```

**응답**:
```json
{
  "success": true,
  "verified": true
}
```

### 친구 관리

#### 친구 목록 조회

```
POST /auth/getfriends
```

**요청 본문**:
```json
{}
```

**응답**:
```json
{
  "success": true,
  "friendList": [
    {
      "id": "친구ID",
      "username": "친구이름",
      "phoneNumber": "전화번호",
      "isBlocked": false
    }
  ]
}
```

#### 친구 추가

```
POST /auth/addfriends
```

**요청 본문**:
```json
{
  "phoneNumber": "친구전화번호"
}
```

**응답**:
```json
{
  "success": true
}
```

#### 친구 차단

```
POST /auth/blockfriend
```

**요청 본문**:
```json
{
  "friendId": "친구ID"
}
```

**응답**:
```json
{
  "success": true
}
```

#### 친구 차단 해제

```
POST /auth/unblockfriend
```

**요청 본문**:
```json
{
  "friendId": "친구ID"
}
```

**응답**:
```json
{
  "success": true
}
```

### 채팅방 관리

#### 채팅방 목록 조회

```
POST /auth/rooms
```

**요청 본문**:
```json
{}
```

**응답**:
```json
{
  "success": true,
  "rooms": [
    {
      "id": "채팅방ID",
      "name": "채팅방이름",
      "participants": [
        {
          "id": "사용자ID",
          "username": "사용자이름"
        }
      ]
    }
  ]
}
```

#### 채팅방 생성

```
POST /auth/createrooms
```

**요청 본문**:
```json
{
  "roomName": "채팅방이름",
  "roomUserList": ["사용자ID1", "사용자ID2"]
}
```

**응답**:
```json
{
  "success": true,
  "roomId": "채팅방ID"
}
```

#### 채팅방에 사용자 추가

```
POST /auth/adduser
```

**요청 본문**:
```json
{
  "roomId": "채팅방ID",
  "userId": "사용자ID"
}
```

**응답**:
```json
{
  "success": true
}
```

#### 채팅방에서 사용자 제거

```
POST /auth/removeuser
```

**요청 본문**:
```json
{
  "roomId": "채팅방ID",
  "userId": "사용자ID"
}
```

**응답**:
```json
{
  "success": true
}
```

#### 채팅 메시지 조회

```
GET /messages?roomId={roomId}&lastMessageId={lastMessageId}
```

**매개변수**:
- `roomId`: 채팅방 ID
- `lastMessageId` (선택사항): 마지막으로 받은 메시지 ID. 이 ID 이후의 메시지만 반환됩니다.

**응답**:
```json
{
  "success": true,
  "messages": [
    {
      "id": "메시지ID",
      "roomId": "채팅방ID",
      "type": "메시지타입",
      "author": {
        "id": "사용자ID"
      },
      "content": "메시지내용",
      "timestamp": "타임스탬프"
    }
  ]
}
```

## WebSocket

### 연결 방법

WebSocket 연결은 다음 URL을 통해 이루어집니다:

```
GET /chat?roomId={roomId}&token={token}
```

**매개변수**:
- `roomId`: 연결할 채팅방 ID
- `token`: 인증 토큰

### 메시지 형식

WebSocket을 통해 주고받는 메시지는 JSON 형식이며, 다음과 같은 구조를 가집니다:

```json
{
  "id": "메시지ID",
  "roomId": "채팅방ID",
  "type": "메시지타입",
  "author": {
    "id": "사용자ID"
  },
  "content": "메시지내용",
  "timestamp": "타임스탬프"
}
```

### 이벤트 유형

WebSocket을 통해 다음과 같은 이벤트를 주고받을 수 있습니다:

#### 클라이언트에서 서버로 보내는 이벤트

- **메시지 전송**: 채팅방에 메시지를 전송합니다.
  ```json
  {
    "type": "message",
    "roomId": "채팅방ID",
    "content": "메시지내용"
  }
  ```

- **타이핑 상태**: 사용자가 타이핑 중임을 알립니다.
  ```json
  {
    "type": "typing",
    "roomId": "채팅방ID",
    "isTyping": true
  }
  ```

#### 서버에서 클라이언트로 보내는 이벤트

- **메시지 수신**: 다른 사용자가 보낸 메시지를 수신합니다.
  ```json
  {
    "id": "메시지ID",
    "roomId": "채팅방ID",
    "type": "message",
    "author": {
      "id": "사용자ID"
    },
    "content": "메시지내용",
    "timestamp": "타임스탬프"
  }
  ```

- **타이핑 상태 수신**: 다른 사용자의 타이핑 상태를 수신합니다.
  ```json
  {
    "type": "typing",
    "roomId": "채팅방ID",
    "userId": "사용자ID",
    "isTyping": true
  }
  ```

- **사용자 입장**: 사용자가 채팅방에 입장했음을 알립니다.
  ```json
  {
    "type": "userJoined",
    "roomId": "채팅방ID",
    "userId": "사용자ID",
    "timestamp": "타임스탬프"
  }
  ```

- **사용자 퇴장**: 사용자가 채팅방에서 퇴장했음을 알립니다.
  ```json
  {
    "type": "userLeft",
    "roomId": "채팅방ID",
    "userId": "사용자ID",
    "timestamp": "타임스탬프"
  }
  ```

## 데이터 모델

### 사용자 (User)

```json
{
  "id": "UUID",
  "username": "문자열",
  "phoneNumber": "문자열",
  "countryCode": "문자열",
  "profileImage": "문자열 (URL)",
  "email": "문자열",
  "createdAt": "타임스탬프",
  "updatedAt": "타임스탬프"
}
```

### 친구 (Friend)

```json
{
  "userId": "UUID",
  "friendId": "UUID",
  "isBlocked": "불리언",
  "createdAt": "타임스탬프",
  "updatedAt": "타임스탬프"
}
```

### 채팅방 (Room)

```json
{
  "id": "UUID",
  "name": "문자열",
  "createdAt": "타임스탬프",
  "updatedAt": "타임스탬프"
}
```

### 메시지 (Message)

```json
{
  "id": "UUID",
  "roomId": "문자열",
  "type": "문자열",
  "author": {
    "id": "문자열"
  },
  "content": "문자열",
  "timestamp": "타임스탬프"
}
```

## 오류 처리

API 요청이 실패하면 다음과 같은 형식의 응답이 반환됩니다:

```json
{
  "success": false,
  "error": {
    "code": "오류코드",
    "message": "오류메시지"
  }
}
```

### 공통 오류 코드

- `400`: 잘못된 요청
- `401`: 인증 실패
- `403`: 권한 없음
- `404`: 리소스를 찾을 수 없음
- `500`: 서버 내부 오류

### 특정 오류 코드

- `1001`: 사용자를 찾을 수 없음
- `1002`: 잘못된 비밀번호
- `1003`: 이미 존재하는 사용자
- `2001`: 채팅방을 찾을 수 없음
- `2002`: 채팅방에 접근할 권한 없음
- `3001`: 친구를 찾을 수 없음

---
