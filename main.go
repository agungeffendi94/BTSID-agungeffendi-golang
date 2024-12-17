package main

import (
	"BTSID-agungeffendi-golang/controllers"
	"BTSID-agungeffendi-golang/core"
	"BTSID-agungeffendi-golang/helpers"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var MySecret = []byte("")

type MyClaims struct {
	UserName string `json:"username"`
	Name     string `json:"name"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}

type LoginPost struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigName("app.conf")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	MySecret = []byte(viper.GetString("jwt.token_secret"))

	router := gin.New()
	router.Use(cors.Default())

	sessionstore := memstore.NewStore([]byte("SuperDuperRahasia1234567890!@#$%"))
	sessionstore.Options(sessions.Options{
		Path:   "/",
		MaxAge: 60 * 60,
	})
	router.Use(sessions.Sessions("BTSID", sessionstore))
	router.Use(gin.Recovery())

	db := core.ConnectDB()
	defer db.Close()

	routing(router, db)

	tmphttpreadheadertimeout, _ := time.ParseDuration(viper.GetString("server.readheadertimeout") + "s")
	tmphttpreadtimeout, _ := time.ParseDuration(viper.GetString("server.readtimeout") + "s")
	tmphttpwritetimeout, _ := time.ParseDuration(viper.GetString("server.writetimeout") + "s")
	tmphttpidletimeout, _ := time.ParseDuration(viper.GetString("server.idletimeout") + "s")

	s := &http.Server{
		Addr:              ":" + viper.GetString("server.port"),
		Handler:           router,
		ReadHeaderTimeout: tmphttpreadheadertimeout,
		ReadTimeout:       tmphttpreadtimeout,
		WriteTimeout:      tmphttpwritetimeout,
		IdleTimeout:       tmphttpidletimeout,
	}
	fmt.Println("Server running on port:", viper.GetString("server.port"))
	s.ListenAndServe()
}

func routing(router *gin.Engine, db *sqlx.DB) {
	router.POST("/login", func(ctx *gin.Context) { Login(ctx, db) })
	router.POST("/register", func(ctx *gin.Context) { controllers.User_Register(ctx, db) })

	checklist := router.Group("checklist")
	checklist.Use(JWTAuthMiddleware())
	{
		checklist.POST("/create", func(ctx *gin.Context) { controllers.Checklist_Create(ctx, db) })
		checklist.GET("/all", func(ctx *gin.Context) { controllers.Checklist_GetAll(ctx, db) })
		checklist.GET("/detail/:id", func(ctx *gin.Context) { controllers.Checklist_GetDetail(ctx, db) })
		checklist.DELETE("/delete/:id", func(ctx *gin.Context) { controllers.Checklist_Delete(ctx, db) })

		item := checklist.Group("item")
		{
			item.POST("/create", func(ctx *gin.Context) { controllers.Items_Create(ctx, db) })
			item.GET("/detail/:id", func(ctx *gin.Context) { controllers.Items_GetDetail(ctx, db) })
			item.PUT("/update/:id", func(ctx *gin.Context) { controllers.Items_Update(ctx, db) })
			item.PUT("/set-status/:id", func(ctx *gin.Context) { controllers.Items_setDone(ctx, db) })
			item.DELETE("/delete/:id", func(ctx *gin.Context) { controllers.Items_Delete(ctx, db) })
		}
	}
}

func Login(c *gin.Context, db *sqlx.DB) {
	var postdata LoginPost
	if errBind := c.ShouldBindJSON(&postdata); errBind != nil {
		c.JSON(403, gin.H{"status": "error", "message": "Invalid request.", "data": nil})
		return
	}

	username := postdata.Username
	password := postdata.Password

	if username == "" || password == "" {
		c.JSON(200, gin.H{
			"data":    nil,
			"message": "Login failed, please complete all inputs",
			"status":  "error",
		})
		return
	}

	pass := helpers.Pass2Hash(password)
	stmt := `SELECT user_id, username, nama FROM users WHERE username = $1 AND password = $2 AND is_active = 1`
	login := helpers.DatabaseQuerySingleRow(db, stmt, username, pass)

	// validate
	if len(login) == 0 {
		c.JSON(200, gin.H{
			"data":    nil,
			"message": "Login failed, Check username and password",
			"status":  "error",
		})
		return
	}

	userName := cast.ToString(login["username"])
	userID := cast.ToString(login["user_id"])
	name := cast.ToString(login["nama"])

	secreet := viper.GetString("jwt.token_secret")
	expired := viper.GetInt("jwt.expired_duration")
	JWTExp := expired
	TokenExpireDuration := time.Now().Add(time.Second * time.Duration(JWTExp)).Unix()
	token, errCJT := GenToken(userName, name, userID, secreet, TokenExpireDuration)
	if errCJT != nil {
		c.JSON(200, gin.H{
			"message": "Error occured (JWT) : " + errCJT.Error(),
			"status":  "error",
			"data":    nil,
		})
		return
	}

	data := map[string]interface{}{
		"user_id":      userID,
		"jwt":          token,
		"nama_lengkap": name,
		"token_expiry": expired,
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"message": "success",
		"status":  "ok",
	})
}

func GenToken(username, name, user_id, secreet string, TokenExpireDuration int64) (string, error) {
	// Create our own statement
	JWTTokenSecret := MySecret
	c := MyClaims{
		username, // Custom field
		name,
		user_id,
		jwt.StandardClaims{
			ExpiresAt: TokenExpireDuration, // Expiration time
			Issuer:    "SIAVA",             // Issuer
		},
	}
	// Creates a signed object using the specified signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// Use the specified secret signature and obtain the complete encoded string token
	return token.SignedString(JWTTokenSecret)
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(200, gin.H{
				"status":  "error",
				"message": "Request header auth Empty",
				"data":    nil,
			})
			c.Abort()
			return
		}
		// Split by space
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(200, gin.H{
				"status":  "error",
				"message": "Request header auth Incorrect format",
				"data":    nil,
			})
			c.Abort()
			return
		}

		token, err := core.VerifyTokenSecretKey(parts[1])
		if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
			c.Abort()
			core.InvalidJWTRes(c)
			return
		}

		// if token exp already expired
		if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
			fmt.Println("C")
			c.Abort()
			core.InvalidJWTRes(c)
			return
		}
		// parts[1] is the obtained tokenString. We use the previously defined function to parse JWT to parse it
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(200, gin.H{
				"status":  "error",
				"message": "invalid Token",
				"data":    nil,
			})
			c.Abort()
			return
		}
		// Save the currently requested username information to the requested context c
		c.Set("user_name", mc.UserName)
		c.Set("user_id", mc.UserID)
		c.Set("name", mc.Name)
		c.Next() // Subsequent processing functions can use c.Get("username") to obtain the currently requested user information
	}
}

func ParseToken(tokenString string) (*MyClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // Verification token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func checkTokenValidity(tokenString string) bool {
	JWT := tokenString
	// check is jwt have valid token secret key
	token, err := core.VerifyTokenSecretKey(JWT)
	if err != nil {
		return false
	}

	// if token exp already expired
	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return false
	}

	return true
}
