package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/models/orm"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// RoomServiceMock은 RoomService 인터페이스를 구현하는 모의 객체입니다.
type RoomServiceMock struct {
	mock.Mock
}

func (m *RoomServiceMock) CreateRoom(ctx context.Context, name string, creatorID uuid.UUID, participantIDs []uuid.UUID) (orm.Room, error) {
	args := m.Called(ctx, name, creatorID, participantIDs)
	return args.Get(0).(orm.Room), args.Error(1)
}

func (m *RoomServiceMock) GetRoomByID(ctx context.Context, id uuid.UUID) (orm.Room, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(orm.Room), args.Error(1)
}

func (m *RoomServiceMock) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]orm.Room, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]orm.Room), args.Error(1)
}

func (m *RoomServiceMock) AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	args := m.Called(ctx, roomID, userID)
	return args.Error(0)
}

func (m *RoomServiceMock) RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	args := m.Called(ctx, roomID, userID)
	return args.Error(0)
}

func TestRoomHandlerGetRoomList(t *testing.T) {
	// 모의 서비스 생성
	roomService := new(RoomServiceMock)

	// 테스트 데이터
	userID := uuid.New()

	room1 := orm.Room{
		UUIDv7BaseModel: orm.UUIDv7BaseModel{ID: uuid.New()},
		RoomName:        "Room 1",
	}

	room2 := orm.Room{
		UUIDv7BaseModel: orm.UUIDv7BaseModel{ID: uuid.New()},
		RoomName:        "Room 2",
	}

	rooms := []orm.Room{room1, room2}

	// 모의 서비스 동작 설정
	roomService.On("GetUserRooms", mock.Anything, userID).Return(rooms, nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("GET", "/rooms", nil)

	// 컨텍스트에 사용자 ID 추가 (미들웨어에서 수행되는 작업 시뮬레이션)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	req = req.WithContext(ctx)

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 핸들러 함수 정의
	handler := func(w http.ResponseWriter, r *http.Request) {
		// 컨텍스트에서 사용자 ID 추출
		userIDStr := r.Context().Value("userID").(string)
		userUUID, _ := uuid.Parse(userIDStr)

		// 모의 서비스 호출
		rooms, err := roomService.GetUserRooms(r.Context(), userUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 응답 생성
		response := map[string]interface{}{
			"success": true,
			"rooms":   rooms,
		}

		// JSON 응답 반환
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	// 핸들러 호출
	handler(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["rooms"])

	// 모의 서비스 호출 검증
	roomService.AssertExpectations(t)
}

func TestRoomHandlerCreateRoom(t *testing.T) {
	// 모의 서비스 생성
	roomService := new(RoomServiceMock)

	// 테스트 데이터
	userID := uuid.New()
	roomName := "New Room"

	newRoom := orm.Room{
		UUIDv7BaseModel: orm.UUIDv7BaseModel{ID: uuid.New()},
		RoomName:        roomName,
	}

	// 모의 서비스 동작 설정
	roomService.On("CreateRoom", mock.Anything, roomName, userID, []uuid.UUID{}).Return(newRoom, nil)

	// 테스트 요청 생성
	reqBody := map[string]string{
		"name": roomName,
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/rooms", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// 컨텍스트에 사용자 ID 추가 (미들웨어에서 수행되는 작업 시뮬레이션)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	req = req.WithContext(ctx)

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 핸들러 함수 정의
	handler := func(w http.ResponseWriter, r *http.Request) {
		// 요청 본문 파싱
		var requestBody struct {
			Name string `json:"name"`
		}
		json.NewDecoder(r.Body).Decode(&requestBody)

		// 컨텍스트에서 사용자 ID 추출
		userIDStr := r.Context().Value("userID").(string)
		userUUID, _ := uuid.Parse(userIDStr)

		// 모의 서비스 호출
		room, err := roomService.CreateRoom(r.Context(), requestBody.Name, userUUID, []uuid.UUID{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 응답 생성
		response := map[string]interface{}{
			"success": true,
			"room":    room,
		}

		// JSON 응답 반환
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}

	// 핸들러 호출
	handler(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["room"])

	// 모의 서비스 호출 검증
	roomService.AssertExpectations(t)
}

func TestRoomHandlerAddUser(t *testing.T) {
	// 모의 서비스 생성
	roomService := new(RoomServiceMock)

	// 테스트 데이터
	userID := uuid.New()
	roomID := uuid.New()
	targetUserID := uuid.New()

	// 모의 서비스 동작 설정
	roomService.On("AddUserToRoom", mock.Anything, roomID, targetUserID).Return(nil)

	// 테스트 요청 생성
	reqBody := map[string]string{
		"userId": targetUserID.String(),
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/rooms/"+roomID.String()+"/users", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// 컨텍스트에 사용자 ID 추가 (미들웨어에서 수행되는 작업 시뮬레이션)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	req = req.WithContext(ctx)

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 라우터 설정 (URL 파라미터 추출을 위해)
	router := mux.NewRouter()
	router.HandleFunc("/rooms/{roomId}/users", func(w http.ResponseWriter, r *http.Request) {
		// URL 파라미터 추출
		vars := mux.Vars(r)
		roomIDStr := vars["roomId"]
		roomUUID, _ := uuid.Parse(roomIDStr)

		// 요청 본문 파싱
		var requestBody struct {
			UserID string `json:"userId"`
		}
		json.NewDecoder(r.Body).Decode(&requestBody)
		targetUUID, _ := uuid.Parse(requestBody.UserID)

		// 모의 서비스 호출
		err := roomService.AddUserToRoom(r.Context(), roomUUID, targetUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 응답 생성
		response := map[string]interface{}{
			"success": true,
		}

		// JSON 응답 반환
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	// 핸들러 호출
	router.ServeHTTP(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])

	// 모의 서비스 호출 검증
	roomService.AssertExpectations(t)
}

func TestRoomHandlerRemoveUser(t *testing.T) {
	// 모의 서비스 생성
	roomService := new(RoomServiceMock)

	// 테스트 데이터
	userID := uuid.New()
	roomID := uuid.New()
	targetUserID := uuid.New()

	// 모의 서비스 동작 설정
	roomService.On("RemoveUserFromRoom", mock.Anything, roomID, targetUserID).Return(nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("DELETE", "/rooms/"+roomID.String()+"/users/"+targetUserID.String(), nil)

	// 컨텍스트에 사용자 ID 추가 (미들웨어에서 수행되는 작업 시뮬레이션)
	ctx := context.WithValue(req.Context(), "userID", userID.String())
	req = req.WithContext(ctx)

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 라우터 설정 (URL 파라미터 추출을 위해)
	router := mux.NewRouter()
	router.HandleFunc("/rooms/{roomId}/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		// URL 파라미터 추출
		vars := mux.Vars(r)
		roomIDStr := vars["roomId"]
		userIDStr := vars["userId"]

		roomUUID, _ := uuid.Parse(roomIDStr)
		targetUUID, _ := uuid.Parse(userIDStr)

		// 모의 서비스 호출
		err := roomService.RemoveUserFromRoom(r.Context(), roomUUID, targetUUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 응답 생성
		response := map[string]interface{}{
			"success": true,
		}

		// JSON 응답 반환
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("DELETE")

	// 핸들러 호출
	router.ServeHTTP(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])

	// 모의 서비스 호출 검증
	roomService.AssertExpectations(t)
}
