package server

import (
	"context"
	"fmt"
	"github.com/kic/users/pkg/database"
	"github.com/kic/users/pkg/logging"
	pbcommon "github.com/kic/users/pkg/proto/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
	"time"
)

var log *zap.SugaredLogger

var repo *database.MockRepository

func prepDBForTests() {
	usersToAdd := []*pbcommon.User{
		{
			UserID:   0,
			UserName: "testUserName",
			Email:    "test@gmail.com",
			Birthday: &pbcommon.Date{
				Year:  1998,
				Month: 2,
				Day:   2,
			},
			City:     "Philadelphia",
			Bio:      "Hello",
		},

		{
			UserID:   1,
			UserName: "ryan",
			Email:    "ryan@gmail.com",
			Birthday: &pbcommon.Date{
				Year:  1998,
				Month: 5,
				Day:   19,
			},
			City:     "Scranton",
			Bio:      "Yo what up",
		},

		{
			UserID:   2,
			UserName: "priya",
			Email:    "priya@gmail.com",
			Birthday: &pbcommon.Date{
				Year:  1999,
				Month: 5,
				Day:   18,
			},
			City:     "Somewhere in New Jersey idk",
			Bio:      "hey ya'll",
		},

		{
			UserID:   3,
			UserName: "quentin",
			Email:    "quentin@gmail.com",
			Birthday: &pbcommon.Date{
				Year:  1998,
				Month: 8,
				Day:   21,
			},
			City:     "Wherever Wyoming Seminary is",
			Bio:      "I hate Java",
		},

		{
			UserID:   4,
			UserName: "jaime",
			Email:    "jaime@gmail.com",
			Birthday: &pbcommon.Date{
				Year:  1999,
				Month: 3,
				Day:   24,
			},
			City:     "Um prolly like Bucks County?",
			Bio:      "Gonna enjoy the sun today!",
		},

		{
			UserID:   5,
			UserName: "aszliah",
			Email:    "aszliah@gmail.com",
			Birthday: &pbcommon.Date{
				Year:  1998,
				Month: 3,
				Day:   16,
			},
			City:     "Somewhere in PA idk",
			Bio:      "Aszliah never told me her birthday, grrr",
		},

	}

	for _, user := range usersToAdd {
		id, err := repo.AddUser(user)
		log.Debugf("inserted id: %v", id)
		if err != nil {
			log.Debugf("insertion error: %v", err)
		}
	}
}

func TestMain(m *testing.M) {
	time.Sleep(1 * time.Second)
	log = logging.CreateLogger(zapcore.DebugLevel)

	// r, mongoClient := setup.DBRepositorySetup(log, "test-mongo-storage")

	usersCollection := &map[int]*pbcommon.User{}
	r := database.NewMockRepository(usersCollection, log)

	repo = r

	prepDBForTests()
	user, _ := r.GetUser(&pbcommon.User{UserName: "ryan"})
	fmt.Println(user.Email)

	// defer mongoClient.Disconnect(context.Background())

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestMongoRepository_GetUserWithName(t *testing.T) {
	usersToCheck := []*pbcommon.User{
		{
			UserName:    "ryan",
		},
		{
			UserName:    "quentin",
		},
		{
			UserName:    "priya",
		},
	}

	notThereUsers := []*pbcommon.User{
		{
			UserName:    "notThere123",
		},
	}

	for i, user := range usersToCheck {
		userResult, err := repo.GetUser(user)

		if err != nil || userResult == nil {
			t.Errorf("Test %v failed with err: %v", i, err)
		}
	}

	for i, user := range notThereUsers {
		userResult, err := repo.GetUser(user)

		if err == nil && userResult != nil {
			t.Errorf("Test %v succeeded but should not have", i)
		}
	}
}

func TestMongoRepository_GetUserWithID(t *testing.T) {
	usersToCheck := []*pbcommon.User{
		{
			UserID:    0,
		},
		{
			UserID:    1,
		},
		{
			UserID:    2,
		},
	}

	notThereUsers := []*pbcommon.User{
		{
			UserID:    100,
		},
	}

	for i, user := range usersToCheck {
		userResult, err := repo.GetUserByID(user.UserID)

		if err != nil || userResult == nil {
			t.Errorf("Test %v failed with err: %v", i, err)
		}
	}

	for i, user := range notThereUsers {
		userResult, err := repo.GetUserByID(user.UserID)

		if err == nil && userResult != nil {
			t.Errorf("Test %v succeeded but should not have", i)
		}
	}
}

func TestMongoRepository_DeleteUser(t *testing.T) {
	err := repo.DeleteUserByID(1)

	if err != nil {
		t.Errorf("Test 1 of delete failed with err: %v", err)
	}

	user, err := repo.GetUserByID(1)

	if err != nil || user != nil{
		t.Errorf("Test 2 of delete failed with err: %v", err)
	}
}

func TestMongoRepository_UpdateUser(t *testing.T) {
	userQuery := &pbcommon.User{
		UserID:   0,
		UserName: "jimmy",
		Email:    "jimmy@gmail.com",
		Birthday: &pbcommon.Date{
			Year:  2000,
			Month: 2,
			Day:   2,
		},
		City:     "Hollywood",
		Bio:      "I'm Hollywood Jimmy",
	}

	err := repo.UpdateUserInfo(context.TODO(), userQuery)

	if err != nil{
		t.Errorf("Test 1 of update failed with err: %v", err)
	}

	user, err := repo.GetUserByID(0)

	if err != nil || user.UserName != "jimmy"{
		t.Errorf("Test 2 of update failed with err: %v", err)
	}
}

