package middleware

import (
	"fmt"

	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/utils"
	"github.com/gin-gonic/gin"
)

// WebMiddleware creates a gin middleware for authorization
func WebMiddleware(ctx *gin.Context)  {
	access_token, err := ctx.Cookie(utils.KeyAccessToken)

	if err == nil && len(access_token) > 0 	{
		ctx.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	}
	ctx.Next()
}
