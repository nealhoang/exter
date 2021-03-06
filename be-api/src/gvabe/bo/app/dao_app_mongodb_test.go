package app

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	"main/src/gvabe/bo/user"
)

func _createMongoConnect(t *testing.T, testName string) *prom.MongoConnect {
	mongoDb := strings.ReplaceAll(os.Getenv("MONGO_DB"), `"`, "")
	mongoUrl := strings.ReplaceAll(os.Getenv("MONGO_URL"), `"`, "")
	if mongoDb == "" || mongoUrl == "" {
		t.Skipf("%s skipped", testName)
		return nil
	}
	mc, err := prom.NewMongoConnect(mongoUrl, mongoDb, 10000)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewMongoConnect", err)
	}
	return mc
}

const collectionNameMongo = "exter_test_app"

func TestNewAppDaoMongo(t *testing.T) {
	name := "TestNewAppDaoMongo"
	mc := _createMongoConnect(t, name)
	appDao := NewAppDaoMongo(mc, collectionNameMongo)
	if appDao == nil {
		t.Fatalf("%s failed: nil", name)
	}
}

func _initAppDaoMongo(t *testing.T, testName string, mc *prom.MongoConnect) AppDao {
	mc.GetCollection(collectionNameMongo).Drop(nil)
	henge.InitMongoCollection(mc, collectionNameMongo)
	return NewAppDaoMongo(mc, collectionNameMongo)
}

func TestAppDaoMongo_Create(t *testing.T) {
	name := "TestAppDaoMongo_Create"
	mc := _createMongoConnect(t, name)
	appDao := _initAppDaoMongo(t, name, mc)
	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	ok, err := appDao.Create(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}
}

func TestAppDaoMongo_Get(t *testing.T) {
	name := "TestAppDaoMongo_Get"
	mc := _createMongoConnect(t, name)
	appDao := _initAppDaoMongo(t, name, mc)
	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	if app, err := appDao.Get("not_found"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app != nil {
		t.Fatalf("%s failed: app %s should not exist", name, "not_found")
	}

	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := app.GetId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
		}
		if v := app.GetTagVersion(); v != 1357 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 1357, v)
		}
		if v := app.GetOwnerId(); v != "btnguyen2k" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "btnguyen2k", v)
		}
		if v := app.GetAttrsPublic().Description; v != "System application (do not delete)" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "System application (do not delete)", v)
		}
	}
}

func TestAppDaoMongo_Delete(t *testing.T) {
	name := "TestAppDaoMongo_Delete"
	mc := _createMongoConnect(t, name)
	appDao := _initAppDaoMongo(t, name, mc)

	appDao.Create(NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)"))
	app, err := appDao.Get("exter")
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	}

	ok, err := appDao.Delete(app)
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if !ok {
		t.Fatalf("%s failed: cannot delete app [%s]", name, app.GetId())
	}

	app, err = appDao.Get("exter")
	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app != nil {
		t.Fatalf("%s failed: app %s should not exist", name, "exter")
	}
}

func TestAppDaoMongo_Update(t *testing.T) {
	name := "TestAppDaoMongo_Update"
	mc := _createMongoConnect(t, name)
	appDao := _initAppDaoMongo(t, name, mc)

	app := NewApp(1357, "exter", "btnguyen2k", "System application (do not delete)")
	appDao.Create(app)

	app.SetOwnerId("nbthanh")
	app.SetTagVersion(2468)
	app.attrsPublic.Description = "App description"
	ok, err := appDao.Update(app)
	if err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name, ok, err)
	}

	if app, err := appDao.Get("exter"); err != nil {
		t.Fatalf("%s failed: %s", name, err)
	} else if app == nil {
		t.Fatalf("%s failed: nil", name)
	} else {
		if v := app.GetId(); v != "exter" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "exter", v)
		}
		if v := app.GetTagVersion(); v != 2468 {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, 2468, v)
		}
		if v := app.GetOwnerId(); v != "nbthanh" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "nbthanh", v)
		}
		if v := app.GetAttrsPublic().Description; v != "App description" {
			t.Fatalf("%s failed: expected [%#v] but received [%#v]", name, "App description", v)
		}
	}
}

func TestAppDaoMongo_GetUserApps(t *testing.T) {
	name := "TestAppDaoMongo_GetUserApps"
	mc := _createMongoConnect(t, name)
	appDao := _initAppDaoMongo(t, name, mc)

	for i := 0; i < 10; i++ {
		app := NewApp(uint64(i), strconv.Itoa(i), strconv.Itoa(i%3), "App #"+strconv.Itoa(i))
		appDao.Create(app)
	}

	u := user.NewUser(123, "2")
	appList, err := appDao.GetUserApps(u)
	if err != nil {
		t.Fatalf("%s failed: %s", name, err)
	}
	if len(appList) != 3 {
		t.Fatalf("%s failed: expected %#v apps but received %#v", name, 3, len(appList))
	}
	for _, app := range appList {
		if app.GetOwnerId() != "2" {
			t.Fatalf("%s failed: app %#v does not belong to user %#v", name, app.GetId(), "2")
		}
	}
}
