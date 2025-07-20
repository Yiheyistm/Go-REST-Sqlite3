package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (app *application) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Autherization Header is required"})
			c.Abort()
			return
		}
		fmt.Println("Auth Header:", authHeader)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Bearer token is required"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, gin.Error{
					Err:  http.ErrNotSupported,
					Type: gin.ErrorTypePublic,
				}
			}
			return []byte(app.JwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid token claims"})
			c.Abort()
			return
		}
		userID := claims["userId"].(float64)

		user, err := app.Model.Users.Get(int(userID))
		fmt.Println("User:", user, userID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": "User not found", "user": user})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
