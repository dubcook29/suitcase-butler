package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// basic runtime parameter configuration

var (
	db_connect_url      string = "mongodb://localhost:27017"
	db_connect_type     string = "mongodb"
	db_connect_host     string = "localhost"
	db_connect_port     string = "27017"
	db_connect_username string = ""
	db_connect_password string = ""
	db_connect_client   *mongo.Client
)

var (
	CURRENT_WORK_PATH    string // current runtime directory
	CURRENT_HOME_PATH    string // user home directory
	CURRENT_CACHE_PATH   string // rootpath + .butler directory
	CURRENT_DEBUG_STATUS bool
)

func SetRuntimeDevMode(debug bool, path string) {
	CURRENT_WORK_PATH = path
	CURRENT_DEBUG_STATUS = debug
	fmt.Println("the (dev) current runtime root path:", CURRENT_WORK_PATH)
	if err := initialCurrentConfigPath(CURRENT_WORK_PATH); err != nil {
		panic(err)
	}
}

func initialCurrentWorkPath() error {
	if path, err := os.Executable(); err != nil {
		return err
	} else {
		CURRENT_WORK_PATH = path
		return nil
	}
}

func initialCurrentHomePath() error {
	if u, err := user.Current(); err != nil {
		return err
	} else if u.HomeDir != "" {
		CURRENT_HOME_PATH = u.HomeDir
		return nil
	} else if home := os.Getenv("HOME"); home != "" {
		CURRENT_HOME_PATH = home
		return nil
	} else if home := os.Getenv("USERPROFILE"); home != "" {
		CURRENT_HOME_PATH = home
		return nil
	} else {
		hd := os.Getenv("HOMEDRIVE")
		hp := os.Getenv("HOMEPATH")
		if hd != "" && hp != "" {
			CURRENT_HOME_PATH = hd + hp
			return nil
		}
	}
	return errors.New("the current user's home directory cannot be found")
}

func initialCurrentConfigPath(path string) error {
	CURRENT_CACHE_PATH = filepath.Join(path, ".butler")
	if _, err := os.Stat(CURRENT_CACHE_PATH); os.IsNotExist(err) {
		if err = os.MkdirAll(CURRENT_CACHE_PATH, os.ModePerm); err != nil {
			return fmt.Errorf("the runtime data cache path is created failed: %v", err)
		} else {
			fmt.Println("the runtime data cache path is created succeeded:", CURRENT_CACHE_PATH)
		}
	} else {
		fmt.Println("the runtime data cache path exists:", CURRENT_CACHE_PATH)
	}
	return nil
}

func init() {
	if err := initialCurrentHomePath(); err != nil {
		panic(err)
	}

	if err := initialCurrentWorkPath(); err != nil {
		panic(err)
	}

	if err := initialCurrentConfigPath(CURRENT_HOME_PATH); err != nil {
		panic(err)
	}
}

func GetCurrentMongoClient() *mongo.Client {
	return db_connect_client
}

func InitialRuntimeDBConnect(ctx context.Context, host, port, user, pass string) error {
	db_connect_host = host
	db_connect_port = port
	db_connect_username = user
	db_connect_password = pass

	if db_connect_username != "" && db_connect_password != "" {
		db_connect_url = fmt.Sprintf("%s://%s:%s@%s:%s", db_connect_type, db_connect_username, db_connect_password, db_connect_host, db_connect_port)
	} else {
		db_connect_url = fmt.Sprintf("%s://%s:%s", db_connect_type, db_connect_host, db_connect_port)
	}

	if db_connect, err := InitialMongoDBConnect(ctx, db_connect_url); err != nil {
		return err
	} else {
		if err := PingMongoDBConnect(db_connect); err != nil {
			return err
		} else {
			db_connect_client = db_connect
			fmt.Println("mongo db connect initial done ...")
		}
	}

	return nil
}

func InitialMongoDBConnect(ctx context.Context, mongodb_url string) (*mongo.Client, error) {
	if mongodb_url == "" {
		return nil, fmt.Errorf("mongodb_url is empty or unrecognized: %v", mongodb_url)
	}

	if mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb_url)); err != nil {
		return nil, err
	} else {
		return mongoClient, nil
	}
}

// PingMongoDBConnect: Mongo Unicom test, the delay is set to 5 seconds
func PingMongoDBConnect(db_connect *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db_connect.Ping(ctx, nil); err != nil {
		fmt.Println("connect to mongodb ping exception:", err)
		return err
	} else {
		fmt.Println("connect to mongodb ping succeeded")
	}
	return nil
}

func ConnectToDatabaseAndCollection(client *mongo.Client, database string, collection string) *mongo.Collection {
	if client != nil {
		return client.Database(database).Collection(collection)
	} else if db_connect_client != nil {
		return db_connect_client.Database(database).Collection(collection)
	} else {
		return nil
	}
}

func MongoId(d interface{}, id string) interface{} {
	originalVal := reflect.ValueOf(d)
	originalType := originalVal.Type()
	newFields := []reflect.StructField{}
	for i := 0; i < originalType.NumField(); i++ {
		field := originalType.Field(i)
		newFields = append(newFields, field)
	}
	newField := reflect.StructField{
		Name: "MongoId",
		Type: reflect.TypeOf(id),
		Tag:  reflect.StructTag(`bson:"_id"`),
	}
	newFields = append(newFields, newField)
	newStructType := reflect.StructOf(newFields)
	newStructValue := reflect.New(newStructType).Elem()
	for i := 0; i < originalVal.NumField(); i++ {
		newStructValue.Field(i).Set(originalVal.Field(i))
	}
	newStructValue.Field(len(newFields) - 1).Set(reflect.ValueOf(id))
	return newStructValue.Interface()
}
