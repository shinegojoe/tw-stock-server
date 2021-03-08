package api

import (
	"fmt"

	"goServer/src/helper/mongoHelper"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//CreateAccount ...
func AddStock(ctx *gin.Context) {
	//CreateAccount ...
	var model map[string]interface{}
	err := ctx.BindJSON(&model)
	if err != nil {
		fmt.Println(err.Error())
		res := mongoHelper.MakeErrResp(err.Error())
		ctx.JSON(200, res)
		return
	}
	fmt.Println("model", model)
	mh := mongoHelper.MHBulider()
	collection := mh.GetCollection("stock", "twStock")
	res := mongoHelper.InsertOne(collection, model, "stock")
	mh.DisConnection()
	ctx.JSON(200, res)
}

func GetStocks(ctx *gin.Context) {
	// collection.Drop(ctx)
	var model map[string]interface{}
	err := ctx.BindJSON(&model)
	if err != nil {
		fmt.Println(err.Error())
		res := mongoHelper.MakeErrResp(err.Error())
		ctx.JSON(200, res)
		return
	}
	fmt.Println("model", len(model))
	mh := mongoHelper.MHBulider()
	collection := mh.GetCollection("stock", "twStock")
	// collection.Drop(ctx)

	filter := bson.M{}
	if len(model) != 0 {
		filter = model
	}
	res := mongoHelper.GetMany(collection, filter, "stocks")
	mh.DisConnection()
	ctx.JSON(200, res)
}
