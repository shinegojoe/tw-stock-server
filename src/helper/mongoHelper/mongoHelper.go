package mongoHelper

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func MHBulider() MongoHelper {
	var mh MongoHelper = MongoHelper{}
	mh.Init()
	mh.Connection()
	return mh
}

type MongoHelper struct {
	client *mongo.Client
	url    string
	ctx    context.Context
	// collection mongo.Collection
}

func (mh *MongoHelper) Init() {
	var url string = "mongodb://167.179.80.227:5569"
	mh.url = url
}

func (mh *MongoHelper) Connection() {
	client, err := mongo.NewClient(options.Client().ApplyURI(mh.url))

	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	mh.client = client
	mh.ctx = ctx
}

func (mh *MongoHelper) GetCollection(dbName string, collectionName string) *mongo.Collection {
	db := mh.client.Database(dbName)
	collection := db.Collection(collectionName)
	return collection
}

func (mh *MongoHelper) DisConnection() {
	mh.client.Disconnect(mh.ctx)
}

func (mh *MongoHelper) DropCollection(collection *mongo.Collection) {
	collection.Drop(mh.ctx)
}

func MakeErrResp(err string) gin.H {
	var res gin.H = gin.H{
		"status":  "error",
		"message": err,
	}
	return res
}

func InsertOne(collection *mongo.Collection, data map[string]interface{}, dataName string) gin.H {

	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		// log.Fatal(err)
		// logHelper.LogToFile(err.Error())
		return MakeErrResp(err.Error())
	}

	fmt.Println("Inserted a single document: ", insertResult)
	return gin.H{
		"status": "ok",
		dataName: insertResult.InsertedID,
	}
}

func DeleteOne(collection *mongo.Collection, filter bson.M, dataName string) gin.H {
	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		return MakeErrResp(err.Error())
	}

	return gin.H{
		"status": "ok",
		dataName: res.DeletedCount,
	}

}

func UpdateOne(collection *mongo.Collection, filter bson.D, data bson.D,
	resName string) gin.H {
	result, err := collection.UpdateOne(context.TODO(), filter, data)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		return MakeErrResp(err.Error())
	}

	if result.MatchedCount == 0 {
		// fmt.Println("matched and replaced an existing document")
		return gin.H{
			"status": "ok",
			resName:  "user not found",
		}
	}
	// if result.UpsertedCount != 0 {
	// 	fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	// }
	return gin.H{
		"status": "ok",
		resName:  "updated",
	}
}

func GetOne(collection *mongo.Collection, filter bson.M, rowName string) gin.H {
	var result bson.M
	// var filter bson.D = bson.D{}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		return MakeErrResp(err.Error())
	}
	// fmt.Printf("Found a single document: %+v\n", result)
	return gin.H{
		"status": "ok",
		rowName:  result,
	}
}

func GetMany(collection *mongo.Collection, filter bson.M, rowName string) gin.H {
	// var results []bson.M
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		MakeErrResp(err.Error())
	}
	var results []bson.M = make([]bson.M, count)

	// var filter bson.D = bson.D{}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		// log.Fatal(err)
		// logHelper.LogToFile(err.Error())
		MakeErrResp(err.Error())
	}
	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	var i int = 0
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem bson.M
		err := cur.Decode(&elem)
		if err != nil {
			// log.Fatal(err)
			// logHelper.LogToFile(err.Error())

			return MakeErrResp(err.Error())
		}
		// fmt.Println("res", elem)

		// results = append(results, elem)
		results[i] = elem
		i++
	}
	// fmt.Println("get", results)
	// fmt.Println("get len", len(results))

	return gin.H{
		"status": "ok",
		rowName:  results,
	}
}
