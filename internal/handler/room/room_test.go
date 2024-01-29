package room

import (
	"encoding/json"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"server/pkg/tutils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	_USER_1 = orm.User{
		UserName:    "test",
		PhoneNumber: "01012345678",
		CountryCode: "82",
		UserID:      uuid.New(),
	}
	_USER_2 = orm.User{
		UserName:    "test2",
		PhoneNumber: "01022222222",
		CountryCode: "82",
		UserID:      uuid.New(),
	}
	_ROOM_1 = orm.Room{
		RoomName: "test",
	}
	_ROOM_2 = orm.Room{
		RoomName: "test2",
	}
)

func TestMain(m *testing.M) {

	postgresClient := postgres_db.GetTestPostgresCleint()
	postgresClient.Migrator().DropTable(&orm.User{}, &orm.Room{}, &orm.RoomUser{})
	postgresClient.AutoMigrate(&orm.User{}, &orm.Room{}, &orm.RoomUser{})
	postgresClient.Create(&_USER_1)
	postgresClient.Create(&_USER_2)

	m.Run()
}

func TestCreateRoomNormalCase(t *testing.T) {
	t.Log("Test Create Room Normal Case")
	r := tutils.CreateRouterWithMiddleware(CreateRoomHandler)
	reqBody, _ := json.Marshal(createRoomRequest{[]uuid.UUID{_USER_2.UserID}})
	req, res := tutils.CreateRequestAndResponse(reqBody)
	req.Header.Set("Authorization", tutils.CreateToken(_USER_1.UserID))

	r.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
	t.Logf("Response Body: %s", res.Body.String())
}
