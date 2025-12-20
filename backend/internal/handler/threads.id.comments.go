package handler

import (
	"fmt"
	"microblog/backend/internal/filter"
	"microblog/backend/internal/helper"
	"microblog/backend/internal/model"
	"microblog/backend/pkg/util"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Ensure this not column name draw, start, length, sort, schema
func GET_THREADS_ID_COMMENTS_HANDLER(db *gorm.DB, preload []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelStruct := &model.Comment{}
		// =============================
		// 🔹 Build schema dari struct
		// =============================
		schema := filter.BuildSchemaFromStruct(modelStruct)

		// =============================
		// 🔹 Expose schema ke UI
		// =============================
		if c.Query("schema") == "true" {
			c.JSON(http.StatusOK, gin.H{
				"schema": schema,
			})
			return
		}

		// =============================
		// 🔹 DataTables params
		// =============================
		var req struct {
			Draw   int    `form:"draw"`
			Start  int    `form:"start"`
			Length int    `form:"length"`
			Sort   string `form:"sort"`
			Fields string `form:"fields"`
		}
		_ = c.BindQuery(&req)

		if req.Length <= 0 {
			req.Length = 20
		}
		if req.Length > 2000 {
			req.Length = 2000
		}

		// =============================
		// 🔹 Base query + preload
		// =============================
		query := db.Model(modelStruct)
		for _, p := range preload {
			query = query.Preload(p)
		}

		// Track which fields to select
		// var selectedFields []string
		var selectedJSONKeys []string

		if req.Fields != "" {
			fields := strings.Split(req.Fields, ",")
			dbColumns := []string{}

			for _, f := range fields {
				f = strings.TrimSpace(f)
				fieldSchema, exists := schema[f]
				if !exists {
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"error":   "Invalid field: " + f,
					})
					return
				}
				dbColumns = append(dbColumns, fieldSchema.DBColumn)
				selectedJSONKeys = append(selectedJSONKeys, f)
			}

			// Always include ID for reference
			if !util.Contains(dbColumns, "id") {
				dbColumns = append([]string{"id"}, dbColumns...)
				if !util.Contains(selectedJSONKeys, "id") {
					selectedJSONKeys = append([]string{"id"}, selectedJSONKeys...)
				}
			}

			// selectedFields = dbColumns
			query = query.Select(dbColumns)
		}

		// =============================
		// 🔹 Apply filtering DSL
		// =============================
		var err error
		query, err = filter.ApplyQueryFilters(
			query,
			c.Request.URL.Query(),
			schema,
		)
		if err != nil {
			if filterErr, ok := err.(*filter.FilterError); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": filterErr.Message,
					"error":   filterErr,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": err.Error(),
					"error":   err.Error(),
				})
			}
			return
		}

		// sort=-created_at  [desc created_at] | sort=created_at,name [asc created_at, asc name]
		if req.Sort == "" {
			req.Sort = "-id"
		}
		query, err = filter.ApplySorting(query, req.Sort, schema)
		if err != nil {
			if filterErr, ok := err.(*filter.FilterError); ok {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": filterErr.Message,
					"error":   filterErr,
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": err.Error(),
					"error":   err.Error(),
				})
			}
			return
		}

		// =============================
		// 🔹 Count filtered
		// =============================
		var recordsFiltered int64
		query.Count(&recordsFiltered)

		// =============================
		// 🔹 Pagination (DataTables)
		// =============================
		query = query.
			Offset(req.Start).
			Limit(req.Length)

		// =============================
		// 🔹 Execute
		// =============================
		var results []model.Comment

		if err := query.Find(&results).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// =============================
		// 🔹 Total records (tanpa filter)
		// =============================
		var recordsTotal int64
		db.Model(modelStruct).Count(&recordsTotal)

		// =============================
		// 🔹 Format response with field selection
		// =============================
		var responseData []gin.H

		if len(selectedJSONKeys) > 0 {
			// Only return selected fields
			t := reflect.TypeOf(modelStruct).Elem()
			sliceValue := reflect.ValueOf(results).Elem()
			data := make([]gin.H, 0, sliceValue.Len())

			// Build map of JSON keys to field indices
			fieldMap := make(map[string]int)
			for i := 0; i < t.NumField(); i++ {
				jsonTag := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
				if jsonTag != "" && jsonTag != "-" {
					fieldMap[jsonTag] = i
				}
			}

			for i := 0; i < sliceValue.Len(); i++ {
				row := sliceValue.Index(i)
				rowData := gin.H{}

				for _, jsonKey := range selectedJSONKeys {
					if fieldIdx, exists := fieldMap[jsonKey]; exists {
						fieldValue := row.Field(fieldIdx)
						rowData[jsonKey] = fieldValue.Interface()
					}
				}

				data = append(data, rowData)
			}
			responseData = data
		} else {
			// Return all fields, but clean up empty relations
			responseData = cleanupEmptyRelations(&results, preload)
		}

		if user, err := helper.GetFirebaseUser(c); err == nil {
			var ids []string
			for _, comment := range results {
				ids = append(ids, comment.ID)
			}
			var votes []model.CommentVote
			db.Where("user_id = ? AND comment_id IN ?", user.ID, ids).Find(&votes)
			voteMap := make(map[string]string) // commentID -> vote
			for _, vote := range votes {
				voteMap[vote.CommentID] = vote.VoteType
			}
			fmt.Println("votes")
			fmt.Println(votes)
			for i, comment := range responseData {
				if voteMap[comment["id"].(string)] == "up" {
					responseData[i]["up_voted_by_me"] = true
				}
				if voteMap[comment["id"].(string)] == "down" {
					responseData[i]["down_voted_by_me"] = true
				}
			}
		}

		// =============================
		// 🔹 Response (DataTables)
		// =============================
		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"draw":            req.Draw,
			"recordsTotal":    recordsTotal,
			"recordsFiltered": recordsFiltered,
			"data":            responseData,
		})
	}
}
