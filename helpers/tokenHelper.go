package helper

import (
	"golang-res/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")



func GenerateAllTokens(email string, firstName string, lastName string, uid string) (signedTokon string, refreshToken string, err error){
	claims := &SignedDetails{
		Email: email,
		First_name: firstName, 
		Last_name: lastName,
		Uid: uid,
		StandardClaims : jwt.STandardClaims{
			ExpiredAt : time.Now().Local().Add(time.Hour*time.Duration(24)).Unix(),


		}, 




	}

	refershClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiredAt: time.Now().Local().Add(time.Hour *time.Duration(168)).Unix()
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodH5256,claims).SignedString([]byte(SECRET_KEY))

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodH5256, refreshClaims).SignedString([]byte(SECRET_KEY))
	
	if err!= nil {
		log.Panic(err)
		return
	}

	return token , refreshToken, err 


}

func UpdateAllTokens() {

}

func ValidateToken() {

}
