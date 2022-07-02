package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rmrf-slash.com/backend/configurations/logger"
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
}

func setup() {
	godotenv.Load()
}

func main() {
	setup()

	cfg := Configuration{}
	if err := env.Parse(&cfg); err != nil {
		logger.GetInstance().Println("Failed to parse env vars")
	}

	ctx := context.TODO()
	opts := options.Client().SetAuth(options.Credential{Username: cfg.MongoUsername, Password: cfg.MongoPassword}).ApplyURI(cfg.MongoConnection)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	defer client.Disconnect(ctx)
	fmt.Printf("%T\n", client)

	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "service is up!",
		})
	})

	apis := r.Group("/api")
	{
		apis.GET("/songs", routers.ListSongs(client))
	}
	r.Run()
}
