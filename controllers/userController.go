package controller

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang-res/databse"
	"golang-res/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang-res/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


var userCollection *mongo/Collection = database.OpenCollection(database.Client, "user")

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){

		var ctx, cancel = context.WithTimeout(context.Background(). 100*time.Second)
		
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))

		if err != nil || recordPerPage < 1 {
			recordPerPage	= 10

		}



		page, err := string.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1

		} 

		startIndex := (page-1)*recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))



		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count" , 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}}

			}
			}
		}

		reult, err := Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage
		})

		defer cancel()
		if err!- nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"} )


		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil{
			log.Fatal(err)
		}
		c.JSON(gttp.StatusOK, allUsers[0])

	

		// 

	}
}


func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userID := c.Params("user_id")

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}.Decode(&user))

		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error": "error occured while listing user items"})
		}

		c.JSON(http.StatusOK, user)


	}
}


func Signup() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		if err := c.BindJSON(&user) : err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return


		}


		validationErr := validate.Struct(user)
		if validationErr != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error", validationErr.Error()})
			return
		}

		count , err := userCollection.CountDocuments(ctx, bson.{"email": user.Email})

		defer cancel()

		if err!= nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the mail"})

			return
		}


		password: HashPassword((*user.Password))
		user.Password = &password


		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phoen})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		if count > 0{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return

		} 

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC339))

		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC339))

		user.ID = primitive.NewObjectID()

		user.User_id = user.ID.Hex()


		token, refereshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)

		user.Token = &token
		user.Refresh_Token = &refereshToken

		resultInsertionNumber, insertErr :=  userCollection.InsertOne(ctx, user)

		if insertErr != nil {
			msg := fmt.Sprintf("user item was not created")
			c.JSON(http.StatusInternalServerErrorm gin.H{"error": msg})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

		// convert the JSON data coming from postman to something that golang understand

		// validator the data based on user struct

		// you'll chack if email has already been used by another user

		// has password

		// check phone number is already used

		// get some extra details - created_at, updated_at, ID

		// generate token and referesh token

		// if all ok, then you insert this new user collection

		// return status OK and send the result back
	}
}


func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		var foundUser models.User

		if err := c.BindJSON(&user) : err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		err := userCollection.FindOne(ctx,bson.M{"email": user.Email}).Decode(&foundUser)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})

			return
	
	
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})

			return
		}

		token, refereshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		helper.UpdateAllTokens(token, refereshToken, foundUser.user_id)

		c.JSON(http.StatusOK, foundUser)

		
	}
}


func HashPassword(password string) string{
	bytes,, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	iff err != nil {
		log.Panic(err)

	}

	 return string(bytes)



}

func VerifyPassword(userPassword string, providePassword string)(bool, string){

	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}

	return check , msg
}