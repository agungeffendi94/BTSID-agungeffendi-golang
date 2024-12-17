package controllers

import (
	"BTSID-agungeffendi-golang/helpers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var FormatDateTime = "2006-01-02 15:04:05"

type UserRegister struct {
	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
	PasswordRetype string `json:"password_retype" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Mail           string `json:"email" binding:"required"`
}

func User_Register(c *gin.Context, dbs *sqlx.DB) {
	var userPost UserRegister
	if errBind := c.ShouldBindJSON(&userPost); errBind != nil {
		c.JSON(403, gin.H{"status": "error", "message": "Invalid request.", "data": nil})
		return
	}

	username := userPost.Username
	password := userPost.Password
	password_retype := userPost.PasswordRetype
	name := userPost.Name
	mail := userPost.Mail

	user_id := "USR-" + uuid.New().String()
	localtime := time.Now().Format(FormatDateTime)

	if password != password_retype {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Passwords does not match",
			"data":    nil,
		})
		return
	}

	check := helpers.DatabaseQuerySingleRow(dbs, `SELECT * FROM users WHERE username=$1 AND is_active=1`, username)
	if len(check) > 0 {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "User already exist.",
			"data":    nil,
		})
		return
	}

	password = helpers.Pass2Hash(password)
	query := `INSERT INTO users(user_id, username, password, nama, create_date, email, is_active) VALUES ($1,$2,$3,$4,$5,$6,1)`
	_, err := dbs.Exec(query, user_id, username, password, name, localtime, mail)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Data saved.",
			"data":    nil,
		})
	}
}
