package mongoHelper

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"context"
	"log"

	"github.com/gin-gonic/gin"
)

func MHBulider() (ConnectionHelper, error) {
	var mh ConnectionHelper = ConnectionHelper{}
	mh.Init()
	err := mh.Connection()
	return mh, err
}

type ConnectionHelper struct {
	client *mongo.Client
	url    string
	ctx    context.Context
	// collection mongo.Collection
}

func (mh *ConnectionHelper) Init() {
	var url string = "mongodb://167.179.80.227:5569"
	mh.url = url
}

func (mh *ConnectionHelper) Connection() error {
	var err error
	client, err := mongo.NewClient(options.Client().ApplyURI(mh.url))

	if err != nil {
		log.Fatal(err)
		return err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	mh.client = client
	mh.ctx = ctx
	return err
}

func (mh *ConnectionHelper) GetCollection(dbName string, collectionName string) *mongo.Collection {
	db := mh.client.Database(dbName)
	collection := db.Collection(collectionName)
	return collection
}

func (mh *ConnectionHelper) DisConnection() {
	mh.client.Disconnect(mh.ctx)
}

func (mh *ConnectionHelper) DropCollection(collection *mongo.Collection) {
	collection.Drop(mh.ctx)
}

func MakeErrResp(err string) gin.H {
	var res gin.H = gin.H{
		"status":  "error",
		"message": err,
	}
	return res
}

func InsertOne(collection *mongo.Collection,
	data map[string]interface{}, dataName string) (gin.H, error) {
	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		return gin.H{}, err
	}

	// fmt.Println("Inserted a single document: ", insertResult)
	return gin.H{
		"status": "ok",
		dataName: insertResult.InsertedID,
	}, err
}

func DeleteOne(collection *mongo.Collection, filter bson.M, dataName string) (gin.H, error) {
	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		// return MakeErrResp(err.Error())
		return gin.H{}, err
	}
	return gin.H{
		"status": "ok",
		dataName: res.DeletedCount,
	}, err
}

func UpdateOne(collection *mongo.Collection, filter bson.D, data bson.D,
	resName string) (gin.H, error) {
	result, err := collection.UpdateOne(context.TODO(), filter, data)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		return gin.H{}, err
	}

	if result.MatchedCount == 0 {
		// fmt.Println("matched and replaced an existing document")
		err = errors.New("there is no user in db")
		return gin.H{}, err
	}
	return gin.H{
		"status": "ok",
		resName:  "updated",
	}, err
}

func GetOne(collection *mongo.Collection, filter bson.M, rowName string) (gin.H, error) {
	var result bson.M
	// var filter bson.D = bson.D{}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		// return MakeErrResp(err.Error())
		return gin.H{}, err
	}
	// fmt.Printf("Found a single document: %+v\n", result)
	return gin.H{
		"status": "ok",
		rowName:  result,
	}, err
}

func GetMany(collection *mongo.Collection, filter bson.M, rowName string) (gin.H, error) {
	// var results []bson.M
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		// logHelper.LogToFile(err.Error())
		// MakeErrResp(err.Error())
		return gin.H{}, err
	}
	var results []bson.M = make([]bson.M, count)

	// var filter bson.D = bson.D{}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		// log.Fatal(err)
		// logHelper.LogToFile(err.Error())
		// MakeErrResp(err.Error())
		return gin.H{}, err
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

			return gin.H{}, err
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
	}, err
}
