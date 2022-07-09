package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rmrf-slash.com/backend/configurations/environments"
	"rmrf-slash.com/backend/configurations/firebase"
	"rmrf-slash.com/backend/middlewares"
	"rmrf-slash.com/backend/routers"
)

type Song struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name   string             `bson:"name,omitempty" json:"name"`
	Artist string             `bson:"artist,omitempty" json:"artist"`
	Genres []string           `bson:"genres,omitempty" json:"genres"`
}

type SongListParams struct {
	Page     int64  `form:"page" binding:"gt=0"`
	PageSize int64  `form:"pageSize" binding:"gt=0"`
	Search   string `form:"search"`
}

type Configuration struct {
	MongoUsername   string `env:"MONGO_USERNAME"`
	MongoPassword   string `env:"MONGO_PASSWORD"`
	MongoConnection string `env:"MONGO_CONN"`
	FirebaseCreds   string `env:"FIREBASE_CREDS"`
}

func decodeB64(value string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return decoded
}

func main() {
	cfg := environments.GetVariables()
	firebase.Instantiate(decodeB64(environments.GetVariables().FirebaseCreds))

	ctx := context.TODO()
	opts := options.Client().SetAuth(options.Credential{Username: cfg.MongoUsername, Password: cfg.MongoPassword}).ApplyURI(cfg.MongoConnection)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	defer client.Disconnect(ctx)
	fmt.Printf("%T\n", client)

	origins := strings.Split(cfg.AllowOrigins, ",")

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowOrigins = origins
	r.Use(cors.New(corsConfig))
	r.SetTrustedProxies(nil)
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "service is up!",
		})
	})

	apis := r.Group("/api")
	{
		apis.GET("/songs", routers.ListSongs(client))
		apis.PUT("/songs/:id", routers.UpdateSong(client))
		apis.DELETE("/songs/:id", routers.DeleteSong(client))
		apis.POST("/songs", middlewares.Secure(), routers.CreateSong(client))
	}

	r.Run()
}
