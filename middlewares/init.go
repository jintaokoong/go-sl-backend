package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"rmrf-slash.com/backend/configurations/firebase"
	"rmrf-slash.com/backend/configurations/logger"
	"rmrf-slash.com/backend/utils/jwt"
)

func Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := firebase.GetInstance().Auth(context.Background())
		if err != nil {
			logger.GetInstance().Println(err)
			c.AbortWithStatusJSON(500, gin.H{"message": "internal server error"})
			return
		}
		token, err := jwt.DecodeBearer(c.Request.Header.Get("Authorization"))
		if err != nil {
			logger.GetInstance().Println(err)
			c.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
			return
		}
		verified, err := auth.VerifyIDToken(context.Background(), token)
		if err != nil {
			logger.GetInstance().Println(err)
			c.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
			return
		}
		c.Set("user", verified)
		c.Next()
	}
}
