#!/bin/bash

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Bulb Talk API 테스트 시작${NC}"
echo "=============================="

# 테스트 환경 변수 설정
export TEST_MODE=true

# 테스트 실행
echo -e "${YELLOW}사용자 API 테스트 실행 중...${NC}"
go test -v ./test/user_handler_test.go
USER_TEST_RESULT=$?

echo -e "${YELLOW}친구 API 테스트 실행 중...${NC}"
go test -v ./test/friend_service_test.go
FRIEND_TEST_RESULT=$?

echo -e "${YELLOW}채팅 서비스 테스트 실행 중...${NC}"
go test -v ./test/chat_service_test.go
CHAT_SERVICE_TEST_RESULT=$?

echo -e "${YELLOW}채팅 핸들러 테스트 실행 중...${NC}"
go test -v ./test/chat_handler_test.go
CHAT_HANDLER_TEST_RESULT=$?

echo -e "${YELLOW}방 핸들러 테스트 실행 중...${NC}"
go test -v ./test/room_handler_test.go
ROOM_TEST_RESULT=$?

echo -e "${YELLOW}WebSocket 테스트 실행 중...${NC}"
go test -v ./test/websocket_test.go
WEBSOCKET_TEST_RESULT=$?

# 결과 요약
echo "=============================="
echo -e "${YELLOW}테스트 결과 요약${NC}"

if [ $USER_TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}✓ 사용자 API 테스트: 성공${NC}"
else
    echo -e "${RED}✗ 사용자 API 테스트: 실패${NC}"
fi

if [ $FRIEND_TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}✓ 친구 API 테스트: 성공${NC}"
else
    echo -e "${RED}✗ 친구 API 테스트: 실패${NC}"
fi

if [ $CHAT_SERVICE_TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}✓ 채팅 서비스 테스트: 성공${NC}"
else
    echo -e "${RED}✗ 채팅 서비스 테스트: 실패${NC}"
fi

if [ $CHAT_HANDLER_TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}✓ 채팅 핸들러 테스트: 성공${NC}"
else
    echo -e "${RED}✗ 채팅 핸들러 테스트: 실패${NC}"
fi

if [ $ROOM_TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}✓ 방 핸들러 테스트: 성공${NC}"
else
    echo -e "${RED}✗ 방 핸들러 테스트: 실패${NC}"
fi

if [ $WEBSOCKET_TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}✓ WebSocket 테스트: 성공${NC}"
else
    echo -e "${RED}✗ WebSocket 테스트: 실패${NC}"
fi

# 전체 결과
echo "=============================="
if [ $USER_TEST_RESULT -eq 0 ] && [ $FRIEND_TEST_RESULT -eq 0 ] && [ $CHAT_SERVICE_TEST_RESULT -eq 0 ] && [ $CHAT_HANDLER_TEST_RESULT -eq 0 ] && [ $ROOM_TEST_RESULT -eq 0 ] && [ $WEBSOCKET_TEST_RESULT -eq 0 ]; then
    echo -e "${GREEN}모든 테스트가 성공적으로 완료되었습니다!${NC}"
    exit 0
else
    echo -e "${RED}일부 테스트가 실패했습니다. 로그를 확인하세요.${NC}"
    exit 1
fi 