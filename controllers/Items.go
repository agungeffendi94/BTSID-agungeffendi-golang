package controllers

import (
	"BTSID-agungeffendi-golang/helpers"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Items struct {
	Title       string `json:"title" binding:"required"`
	ChecklistID string `json:"checklist_id" binding:"required"`
}

type ItemsUpdate struct {
	Title string `json:"title" binding:"required"`
}

func Items_Create(c *gin.Context, dbs *sqlx.DB) {
	var postdata Items
	if errBind := c.ShouldBindJSON(&postdata); errBind != nil {
		fmt.Println(errBind)
		c.JSON(403, gin.H{"status": "error", "message": "Invalid request.", "data": nil})
		return
	}

	today := time.Now()

	_, err := dbs.Exec("INSERT INTO items (id, title, checklist_id, create_date) VALUES ($1, $2, $3, $4)", uuid.New().String(), postdata.Title, postdata.ChecklistID, today.Format(FormatDateTime))

	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Failed to create items.",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Success to create items.",
		"data":    nil,
	})
}

func Items_Update(c *gin.Context, dbs *sqlx.DB) {
	var postdata ItemsUpdate
	if errBind := c.ShouldBindJSON(&postdata); errBind != nil {
		fmt.Println(errBind)
		c.JSON(403, gin.H{"status": "error", "message": "Invalid request.", "data": nil})
		return
	}
	dataID := c.Params.ByName("id")
	_, errUp := dbs.Exec("UPDATE items SET title = $1, last_update=$2 WHERE id = $3", postdata.Title, time.Now().Format(FormatDateTime), dataID)

	if errUp != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Failed to update items.",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Success to update items.",
		"data":    nil,
	})
}

func Items_Delete(c *gin.Context, dbs *sqlx.DB) {
	dataID := c.Params.ByName("id")

	_, errDelItems := dbs.Exec("DELETE FROM items WHERE id = $1", dataID)
	if errDelItems != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Failed to delete items.",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Success to delete checklist.",
		"data":    nil,
	})
}

func Items_GetDetail(c *gin.Context, dbs *sqlx.DB) {
	dataID := c.Params.ByName("id")
	get := helpers.DatabaseQuerySingleRow(dbs, `SELECT id, title, create_date FROM items WHERE id = $1`, dataID)

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Success to get all checklists.",
		"data":    get,
	})
}

func Items_setDone(c *gin.Context, dbs *sqlx.DB) {
	dataID := c.Params.ByName("id")
	_, errUp := dbs.Exec("UPDATE items SET status = 1 WHERE id = $1", dataID)
	if errUp != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Failed to update items.",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Success to update items.",
		"data":    nil,
	})
}
