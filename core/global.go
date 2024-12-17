package core

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var FormatDateTime = "2006-01-02 15:04:05"
var FormatDate = "2006-01-02"
var FormatDateHuman = "01 Jan 2006, 15:04"

// HelloWorld response message hello world
func HelloWorld(c *gin.Context) {
	middlewareData, exist := c.Get("middleware_data")
	data := map[string]interface{}{}
	if exist {
		mapMiddlewareData := cast.ToStringMap(middlewareData)
		data["token"] = cast.ToString(mapMiddlewareData["new_jwt"])
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "hello world",
	})
}

// InvalidJWTRes response this if jwt is not valid
func InvalidJWTRes(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"data":    nil,
		"message": "invalid token",
	})
}

// InternalServerErrorRes response this if there is internal server error
func InternalServerErrorRes(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"data":    nil,
		"message": err.Error(),
	})
}

// EncString response enc string
// func EncString(c *gin.Context) {
// 	type ReqStruct struct {
// 		Text string `form:"text" json:"text" binding:"required"`
// 	}

// 	reqJSON := new(ReqStruct)
// 	errBind := c.Bind(reqJSON)
// 	if errBind != nil {
// 		fmt.Println(errBind)
// 		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{
// 			"message": errBind.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": cdc.EncString(reqJSON.Text),
// 	})
// }

// DecString response enc string
// func DecString(c *gin.Context) {
// 	type ReqStruct struct {
// 		Text string `form:"text" json:"text" binding:"required"`
// 	}

// 	reqJSON := new(ReqStruct)
// 	errBind := c.Bind(reqJSON)
// 	if errBind != nil {
// 		fmt.Println(errBind)
// 		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{
// 			"message": errBind.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": cdc.DecString(reqJSON.Text),
// 	})
// }

// CreateJWTToken create token
func CreateJWTToken(claims jwt.MapClaims) (string, error) {
	JWTTokenSecret := viper.GetString("jwt.token_secret")
	JWTExp := viper.GetInt("jwt.expaired_duration")

	var err error
	// extend claims data
	// claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Second * time.Duration(JWTExp)).Unix()
	JWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := JWT.SignedString([]byte(JWTTokenSecret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyTokenSecretKey verify token secret key
func VerifyTokenSecretKey(JWTToken string) (*jwt.Token, error) {
	JWTTokenSecret := viper.GetString("jwt.token_secret")

	token, err := jwt.Parse(JWTToken, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTTokenSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
