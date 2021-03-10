package api

import (
	"fmt"

	"goServer/src/helper/logHelper"
	"goServer/src/helper/mongoHelper"
	"goServer/src/helper/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//CreateAccount ...
func AddStock(ctx *gin.Context) {
	//CreateAccount ...
	var model map[string]interface{}
	err := ctx.BindJSON(&model)
	if err != nil {
		logHelper.LogToFile(err.Error())
		ctx.JSON(200, utils.MakeErrResp(err.Error()))
		return
	}
	fmt.Println("model", model)
	connectionHelper, err := mongoHelper.MHBulider()
	if err != nil {
		logHelper.LogToFile(err.Error())
		ctx.JSON(200, utils.MakeErrResp(err.Error()))
		return
	}
	connectionHelper.Init()
	defer connectionHelper.DisConnection()
	collection := connectionHelper.GetCollection("stock", "twStock")
	var stock bson.M
	err = collection.FindOne(ctx, model).Decode(&stock)
	if err != nil {
		res, err := mongoHelper.InsertOne(collection, model, "stock")

		utils.ResponseHelper(ctx, 200, res, err)
	} else {
		logHelper.LogToFile("data is exists")
		ctx.JSON(200, gin.H{
			"status":  "ok",
			"message": "data is exists",
		})
	}

}

func GetStocks(ctx *gin.Context) {
	// collection.Drop(ctx)
	var model map[string]interface{}
	err := ctx.BindJSON(&model)
	if err != nil {
		res := utils.MakeErrResp(err.Error())
		ctx.JSON(200, res)
		return
	}
	fmt.Println("model", len(model), model)
	connectionHelper, err := mongoHelper.MHBulider()
	if err != nil {
		ctx.JSON(200, utils.MakeErrResp(err.Error()))
		return
	}
	connectionHelper.Init()
	defer connectionHelper.DisConnection()
	collection := connectionHelper.GetCollection("stock", "twStock")
	// collection.Drop(ctx)

	filter := bson.M{}
	if len(model) != 0 {
		filter = model
	}
	res, err := mongoHelper.GetMany(collection, filter, "stocks")
	utils.ResponseHelper(ctx, 200, res, err)

}
