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
	ctx.BindJSON(&model)
	fmt.Println("model", model)
	collection, _, client := mongoHelper.GetCollection("stock", "twStock")
	fmt.Println(collection)
	res := mongoHelper.InsertOne(collection, model, "stock")
	client.Disconnect(ctx)

	// var res = gin.H{
	// 	"status": "ok",
	// 	"res":    model,
	// }
	// cancel()
	ctx.JSON(200, res)
}

func GetStock(ctx *gin.Context) {

}

func GetStocks(ctx *gin.Context) {
	collection, _, client := mongoHelper.GetCollection("stock", "twStock")
	// collection.Drop(ctx)

	res := mongoHelper.GetMany(collection, bson.M{}, "stocks")
	// var res = gin.H{
	// 	"status": "ok",
	// 	"res":    123,
	// }
	client.Disconnect(ctx)
	ctx.JSON(200, res)
}
