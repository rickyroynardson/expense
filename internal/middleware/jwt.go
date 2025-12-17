package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickyroynardson/expense/utils"
)

func AuthenticatedJWT(cfg *utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		cookie, err := c.Cookie("access_token")
		if err != nil {
			authorizationToken, err := utils.GetAuthorizationToken(c.Request.Header)
			if err != nil {
				utils.RespondJSON(c, http.StatusUnauthorized, err.Error(), nil)
				c.Abort()
				return
			}
			token = authorizationToken
		} else {
			token = cookie
		}

		userID, err := utils.ValidateJWT(token, cfg.JwtSecret)
		if err != nil {
			utils.RespondJSON(c, http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
