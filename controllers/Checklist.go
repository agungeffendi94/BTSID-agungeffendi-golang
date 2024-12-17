package controllers

import (
	"BTSID-agungeffendi-golang/helpers"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
)

type Checklist struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func Checklist_Create(c *gin.Context, dbs *sqlx.DB) {
	var postdata Checklist
	if errBind := c.ShouldBindJSON(&postdata); errBind != nil {
		c.JSON(403, gin.H{"status": "error", "message": "Invalid request.", "data": nil})
		return
	}

	user_id, errClaims := helpers.GetJWTClaims(c.Request.Header["Authorization"], "user_id")
	if errClaims != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Invalid JWT Token." + errClaims.Error(),
			"data":    nil,
		})
		return
	}

	today := time.Now()

	_, err := dbs.Exec("INSERT INTO checklist (id, title, description, create_by, create_date) VALUES ($1, $2, $3, $4, $5)", uuid.New().String(), postdata.Title, postdata.Description, user_id, today.Format(FormatDateTime))

	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Failed to create checklist.",
			"data":    nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Success to create checklist.",
		"data":    nil,
	})
}

func Checklist_GetDetail(c *gin.Context, dbs *sqlx.DB) {
	dataID := c.Params.ByName("id")
	get := helpers.DatabaseQuerySingleRow(dbs, `SELECT id, title, description, create_by, create_date FROM checklist WHERE id = $1`, dataID)

	get_items := helpers.DatabaseQueryRows(dbs, `SELECT id, title, status FROM items WHERE checklist_id = $1`, dataID)

	get["items"] = get_items
	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Success to get all checklists.",
		"data":    get,
	})
}

func Checklist_Delete(c *gin.Context, dbs *sqlx.DB) {
	dataID := c.Params.ByName("id")

	// delete all item
	_, errDelItems := dbs.Exec("DELETE FROM items WHERE checklist_id = $1", dataID)
	if errDelItems != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Failed to delete items.",
			"data":    nil,
		})
		return
	}

	// delete checklist
	_, errDelChecklist := dbs.Exec("DELETE FROM checklist WHERE id = $1", dataID)
	if errDelChecklist != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Failed to delete checklist.",
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

func Checklist_GetAll(c *gin.Context, dbs *sqlx.DB) {
	user_id, errClaims := helpers.GetJWTClaims(c.Request.Header["Authorization"], "user_id")
	if errClaims != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"message": "Invalid JWT Token." + errClaims.Error(),
			"data":    nil,
		})
		return
	}

	get_checklists := helpers.DatabaseQueryRows(dbs, `SELECT * FROM checklist WHERE create_by=$1`, user_id)

	arrIDS := []string{}

	for _, v := range get_checklists {
		arrIDS = append(arrIDS, "'"+cast.ToString(v["id"])+"'")
	}

	get_items := helpers.DatabaseQueryRows(dbs, `SELECT * FROM items WHERE checklist_id IN (`+strings.Join(arrIDS, ", ")+`)`)

	resp := []map[string]interface{}{}
	for _, v := range get_checklists {
		key := cast.ToString(v["id"])
		data := map[string]interface{}{}
		for _, z := range get_items {
			if cast.ToString(z["checklist_id"]) == key {
				data[cast.ToString(z["id"])] = map[string]interface{}{
					"id":     cast.ToString(z["id"]),
					"title":  cast.ToString(z["title"]),
					"status": cast.ToInt(z["status"]),
				}
			}
		}
		newdata := map[string]interface{}{
			"id":          cast.ToString(v["id"]),
			"title":       cast.ToString(v["title"]),
			"items":       data,
			"create_by":   cast.ToString(v["create_by"]),
			"create_date": cast.ToString(v["create_date"]),
		}
		resp = append(resp, newdata)
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Success to get all checklists.",
		"data":    resp,
	})
}
