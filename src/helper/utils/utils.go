package utils

import (
	"errors"
	"goServer/src/helper/mongoHelper"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MakeErrResp(msg string) gin.H {
	var res gin.H = gin.H{
		"status":  "error",
		"message": msg,
	}
	return res
}

func ResponseHelper(ctx *gin.Context, statusCode int, res gin.H, err error) {
	if err != nil {
		var res gin.H = MakeErrResp(err.Error())
		ctx.JSON(statusCode, res)
		return
	}
	ctx.JSON(200, res)
}

func DeleteById(ctx *gin.Context, collection *mongo.Collection) (gin.H, error) {
	var id string = ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// mongoHelper.MakeErrResp(err.Error())
		return gin.H{}, err
	}
	// var collection *mongo.Collection = getCollection()
	var filter bson.M = bson.M{
		"_id": objectId,
	}
	var dataName string = "deleteCount"
	res, err := mongoHelper.DeleteOne(collection, filter, dataName)
	if err != nil {
		return gin.H{}, err
	}
	// var count int64 = res[dataName].(int64)
	if res[dataName] == int64(0) {
		// return mongoHelper.MakeErrResp("there is no user id " + id)
		err = errors.New("there is no user id " + id)
		return gin.H{}, err
	} else {
		delete(res, dataName)
		res["message"] = "user " + id + " is deleted"
		return res, err
	}
}

func GetById(ctx *gin.Context, collection *mongo.Collection, dataName string) (gin.H, error) {
	var id string = ctx.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// fmt.Println("Invalid id")
		// mongoHelper.MakeErrResp(err.Error())
		return gin.H{}, err
	}
	// collection := mongoHelper.GetCollection(dbName, userCollectionName)
	var filter bson.M = bson.M{
		"_id": objectId,
	}
	res, err := mongoHelper.GetOne(collection, filter, dataName)
	return res, err
}
