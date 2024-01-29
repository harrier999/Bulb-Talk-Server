package test

import (
	"encoding/json"
	"net/http"

	"testing"
	"time"

	"server/internal/handler/friends"
	"server/pkg/tutils"

	"github.com/stretchr/testify/assert"
)

var (
	_WRONG_USER_ID_1 = "wrong_user_id"
	_WRONG_USER_ID_2 = "b5b06de3-5a19-4420-a6ae-c84b4439a61"

	_VALID_USER_ID_1 = "b5b06de3-5a19-4420-a6ae-cc84b4439a61"
)

func TestNormalCase(t *testing.T) {
	r := tutils.CreateRouterWithMiddleware(friends.GetFriendList)
	requestData := friends.FriendListRequest{
		LastRequestTime: time.Now(),
	}
	reqBody, _ := json.Marshal(requestData)

	req, rr := tutils.CreateRequestAndResponse(reqBody)
	req.Header.Set("Authorization", tutils.CreateTokenByString(_VALID_USER_ID_1))

	r.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("failed to get friend list")
	}
}

func TestInvalidUUID(t *testing.T) {
	r := tutils.CreateRouterWithMiddleware(friends.GetFriendList)

	req, res := tutils.CreateRequestAndResponse(nil)
	req.Header.Set("Authorization", tutils.CreateTokenByString(_WRONG_USER_ID_1))
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)

	req2, res2 := tutils.CreateRequestAndResponse(nil)
	req2.Header.Set("Authorization", tutils.CreateTokenByString(_WRONG_USER_ID_2))
	r.ServeHTTP(res2, req2)
	assert.Equal(t, http.StatusUnauthorized, res2.Result().StatusCode)
}

func TestEmptyFriend(t *testing.T) {
	r := tutils.CreateRouterWithMiddleware(friends.GetFriendList)

	req, res := tutils.CreateRequestAndResponse(nil)
	req.Header.Set("Authorization", tutils.CreateTokenByString(_VALID_USER_ID_1))

	r.ServeHTTP(res, req)
	t.Logf("res: %v", res.Body)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
}
