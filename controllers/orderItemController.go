package controller

import (
	"golang-res/database"
	"golang-res/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id *string
	Order_item []models.OrderItem

}

var orderItemCollection *mongo.Collection = database.openCollection(database, "orderItem")



func GetOrderItems() gin.HandlerFunc{
	return func (c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)


		result, err := orderItemCollection.Find(context.TODO().bson.M{})
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H("error":"error occured while listing order items"))
		return
		}

		var allOrdersItems []bson.M
		if err = result.All(ctx, &allOrdersItems); err !=nil{
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, allOrderItems)


	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func (c *gin.Context){

		orderId := c.Params("order_id")
		allOrderItems, err := ItemsByOrder(orderId)

		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H("error":"error occured while listing order items by order Id"))
		return
		}

		c.JSON(http.StatusOK, allOrderItems)



	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error){
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D({"$match", bson.D{{"order_id", id}}})

	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foriegnField", "food_id"}, {"as", "food"}}}}

	
	
	unwindStage := bson.D{{"$unwind", bson.D{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"food", "order"},{"localField", "order_id"},{"foriegnField", "order_id"}, {"as", "order"}}}}

	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArraays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{"from", "table"}, {"localField", "order.table_id"}, {"foreignField","table_id"},{"as","table"}}}

	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArray", true}}}}

	projectStage := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"amount", "$food.price"},
			{"total_count", 1},
			{"food_name", "$food.name"},
			{"food_image", "food.food_image"}
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity",1},

		}

		}
	}

	groupStage := bson.D{{"$group", bson.D{{"id", bson.D{"order_id", "$order_id"}, {table_id", "$table_id}, {"table_number", "$table_number"}}}}}, {"payment_due", bson.D{{"$sum","$amount"}}}, {"toal_count", bson.D{{"$sum", 1}}}, {"order_items", }


	projectStage2 := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"table_number", "$_id.table_number"},
			{"order_items", 1}

		}}
	}


	result, errr := ordertemCollection.Aggregrate{ctx, mongo.Pipeline{
		matchStaage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,

		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2

	}}
	if err!= nil{
		panic{err}

	}

	defer cancel()

	return OrderItems, err



	result.All(ctx, &orderItems); err !=nil{
		panic(err)
	}



}


func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		OrderItemId := c.Params("order_item_id")
		var orderItem models.OrderItem


		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": OrderItemId}).Decode(&orderItem)

		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H("error":"error occured while listing order items"))
			return

		
		}
		c.JSON(httpStatusOK, orderItem)

	}
}


func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var orderItem models.OrderItem

		orderItemId = c.Param("order_item_id")

		filter := bson.M{"order_item_id": orderItemId}

		var updatObj primitive.D

		if orderItem.Unit_price != nil{
			updateObj = append(updateObj, bson.E{"unit_price", *&orderItem.Unit_price})
		}

		
		if orderItem. != nil{
			updatObj = append(updatObj, bson.E{"quantity", *&orderItem.Quantity})
		}


		if orderItem.Food_id != nil{

			updatObj = append(updatObj, bson.E{"food_id", *&orderItem.Updated_at})


		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at",orderItem.Updated_at})


		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},

			},

			&opt, 
		)

		if err!= nil{
			msg := "Order item update failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return

		}

		defer cancel()

		c.JSON(http.StatusOK, result)



	}
}



func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var OrderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err!= nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}




		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	//	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemsToBeIntersted : []interface{}{}
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, range orderItemPack.Order_items{
			orderItem.Order_id = order_id

			validationErr := validte.Struct(orderItem)
			if validationErr != nil{
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at,_ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price,2 )
			orderItem.Unit_price = &num
			orderItemsToBeIntersted = append(orderItemsToBeInserted, orderItem)

		}

		insertOrderitems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil{
			log.Fatal(err)
		}

		defer cancel()

		c.JSON(http.StatusOK, insertOrderitems)


	
	
	}
}
