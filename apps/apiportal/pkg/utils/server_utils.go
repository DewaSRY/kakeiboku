package utils

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func CommonResponse(message string) gin.H {
	return gin.H{"message": message, "success": true}
}

// BindJSON wraps ShouldBindJSON and returns a descriptive error when the
// request body is empty, listing all fields that are marked as required.
func BindJSON(ctx *gin.Context, obj interface{}) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		if err == io.EOF {
			fields := requiredJSONFields(obj)
			if len(fields) > 0 {
				return fmt.Errorf("request body is required with fields: %s", strings.Join(fields, ", "))
			}
			return fmt.Errorf("request body is required")
		}
		return err
	}
	return nil
}

// requiredJSONFields returns the JSON field names that carry binding:"required".
func requiredJSONFields(obj interface{}) []string {
	var fields []string
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fields
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if strings.Contains(field.Tag.Get("binding"), "required") {
			name := strings.Split(field.Tag.Get("json"), ",")[0]
			if name == "" || name == "-" {
				name = strings.ToLower(field.Name)
			}
			fields = append(fields, name)
		}
	}
	return fields
}
