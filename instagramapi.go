package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	time "time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}
func insertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	// select database and collection ith Client.Database method
	// and Database.Collection method
	collection := client.Database(dataBase).Collection(col)

	// InsertOne accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}
func insertMany(client *mongo.Client, ctx context.Context, dataBase, col string, docs []interface{}) (*mongo.InsertManyResult, error) {

	// select database and collection ith Client.Database
	// method and Database.Collection method
	collection := client.Database(dataBase).Collection(col)

	// InsertMany accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertMany(ctx, docs)
	return result, err
}

const UrlUserInfo = `https://reqres.in/api/users`

type User struct {
	Name     string `json:"name", db:"name"`
	Id       string `json:"id", db:"id"`
	Email    string `json:"email", db:"email"`
	Password string `json:"password", db:"password"`
}

type Post struct {
	Id             string `json:"id", db:"id"`
	Caption        string `json:"caption", db:"caption"`
	ImageExtension string `form:"image_extension" json:"image_extension, db:"image_extension"`
	Seconds        int64  `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
}

//get user id in url
func getSource(id string) (b []byte, err error) {
	resp, err := http.Get("https://reqres.in/api/users" + id + "/")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

//get post in url
func getPost(id string) (b []byte, err error) {
	resp, err := http.Get("https://reqres.in/api/posts" + id + "/")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func main() {

	// get Client, Context, CancelFunc and err from connect method.
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	// Release resource when main function is returned.
	defer close(client, ctx, cancel)

	// Create  a object of type interface to  store
	// the bson values, that  we are inserting into database.
	var user interface{}

	//Creating user struct which need to post.
	user := bson.User{
		Name:     "Test User",
		Id:       "1234test",
		Email:    "asdf@gmail.com",
		Password: "123456789",
	}
	// insertOne accepts client , context, database
	// name collection name and an interface that
	// will be inserted into the  collection.
	// insertOne returns an error and aresult of
	// insertina single document into the collection.
	insertOneResult, err := insertOne(client, ctx, "gfg",
		"userdata", user)

	// handle the error
	if err != nil {
		panic(err)
	}

	var post interface{}
	//Creating post
	post := bson.Post{
		Id:             "1234test",
		Caption:        "asdfghjkmnbvcxvbnjuytre",
		ImageExtension: "png",
	}

	insertOneResult, err := insertOne(client, ctx, "gfg",
		"postinfo", post)

	// handle the error
	if err != nil {
		panic(err)
	}

	//Converting User to byte using Json.Marshal
	//Ignoring error.
	body, _ := json.Marshal(user)
	body, _ := json.Marshal(post)

	//Passing new buffer for request with URL to post.
	//This will make a post request and will share the JSON data
	resp, err := http.Post("https://reqres.in/api/users", "application/json", bytes.NewBuffer(body))
	resp, err := http.Post("https://reqres.in/api/posts", "application/json", bytes.NewBuffer(body))

	// An error is returned if something goes wrong
	if err != nil {
		panic(err)
	}
	//Need to close the response stream, once response is read.
	//Hence defer close. It will automatically take care of it.
	defer resp.Body.Close()

	//Check response code, if New user is created then read response.
	if resp.StatusCode == http.StatusCreated {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//Failed to read response.
			panic(err)
		}

		//Convert bytes to String and print
		jsonStr := string(body)
		fmt.Println("Response: ", jsonStr)

	} else {
		//The status is not Created. print the error.
		fmt.Println("Get failed with error: ", resp.Status)
	}

	const b = 0
	//get user id in url
	b, err := getSource("id")
	if err != nil {
		panic(err)
	}
	const d = 0
	//get post in url
	d, err := getPost("id")
	if err != nil {
		panic(err)
	}
	//time stamp
	t := time.Now()
	fmt.Println("Location : ", t.Location(), " Time : ", t) // local time

}
