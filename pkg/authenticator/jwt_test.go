package authenticator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/golang-jwt/jwt/v5"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestJWTMiddleware(t *testing.T) {
	m := mux.NewRouter()
	m.Use(JWTMiddleware)
	m.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		user_id := r.Context().Value(ContextKeyUserID)
		assert.NotNil(t, user_id)
		assert.Equal(t, "test_user_id", user_id)

	})
	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)
	token, _ := CreateToken("test_user_id", 100)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestExpirationTime(t *testing.T) {
	token, _ := CreateToken("test_user_id", -1)
	m := mux.NewRouter()
	m.Use(JWTMiddleware)
	m.HandleFunc("/test", dummyHandler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	token2, _ := CreateToken("test_user_id", 1)
	req2, err2 := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err2)
	req.Header.Set("Authorization", token2)
	rr2 := httptest.NewRecorder()
	m.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusUnauthorized, rr2.Code)
}

func TestCreateAndValidateToken(t *testing.T) {
	// 테스트를 위한 환경 설정
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret) // 테스트 후 원래 값으로 복원

	// 테스트용 시크릿 설정
	testSecret := "epqmdhqtm"
	os.Setenv("JWT_SECRET", testSecret)

	// hmacSecret 변수를 직접 설정 (패키지 변수 재설정)
	hmacSecret = []byte(testSecret)

	// 테스트용 사용자 ID
	userID := "test-user-id"

	// 토큰 생성
	token, err := CreateToken(userID, 1)
	assert.NoError(t, err, "토큰 생성 중 오류 발생")
	assert.NotEmpty(t, token, "생성된 토큰이 비어 있음")

	// 토큰 파싱 및 검증
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// 서명 방식 확인
		assert.IsType(t, &jwt.SigningMethodHMAC{}, token.Method, "예상치 못한 서명 방식")
		assert.Equal(t, "HS256", token.Header["alg"], "예상치 못한 알고리즘")

		// 시크릿 키 반환
		return []byte(testSecret), nil
	})

	assert.NoError(t, err, "토큰 파싱 중 오류 발생")
	assert.True(t, parsedToken.Valid, "토큰이 유효하지 않음")

	// 클레임 확인
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok, "클레임을 맵으로 변환할 수 없음")
	assert.Equal(t, userID, claims["user_id"], "사용자 ID가 일치하지 않음")

	// 만료 시간 확인
	expTime := time.Unix(int64(claims["exp"].(float64)), 0)
	assert.True(t, expTime.After(time.Now()), "토큰이 이미 만료됨")
}

func TestInvalidSecret(t *testing.T) {
	// 테스트를 위한 환경 설정
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret) // 테스트 후 원래 값으로 복원

	// 테스트용 시크릿 설정
	testSecret := "epqmdhqtm"
	os.Setenv("JWT_SECRET", testSecret)

	// hmacSecret 변수를 직접 설정 (패키지 변수 재설정)
	hmacSecret = []byte(testSecret)

	// 토큰 생성
	userID := "test-user-id"
	token, err := CreateToken(userID, 1)
	assert.NoError(t, err, "토큰 생성 중 오류 발생")

	// 잘못된 시크릿으로 토큰 검증
	wrongSecret := "wrong-secret"
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(wrongSecret), nil
	})

	assert.Error(t, err, "잘못된 시크릿으로 토큰이 검증됨")
	assert.False(t, parsedToken.Valid, "잘못된 시크릿으로 토큰이 유효함")
}

func TestEmptySecret(t *testing.T) {
	// 테스트를 위한 환경 설정
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret) // 테스트 후 원래 값으로 복원

	// 빈 시크릿 설정
	os.Setenv("JWT_SECRET", "")

	// hmacSecret 변수를 직접 설정 (패키지 변수 재설정)
	hmacSecret = []byte("")

	// 토큰 생성 시도
	userID := "test-user-id"
	token, err := CreateToken(userID, 1)

	// 빈 시크릿으로도 토큰이 생성될 수 있지만, 보안상 좋지 않음
	// 이 테스트는 빈 시크릿으로 토큰이 생성되는지 확인
	if err == nil && token != "" {
		t.Log("경고: 빈 시크릿으로 토큰이 생성됨")

		// 빈 시크릿으로 토큰 검증
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(""), nil
		})

		if err == nil && parsedToken.Valid {
			t.Log("경고: 빈 시크릿으로 토큰이 검증됨")
		}
	}
}
