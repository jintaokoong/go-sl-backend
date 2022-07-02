package routers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rmrf-slash.com/backend/configurations/logger"
)

type Song struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name   string             `bson:"name,omitempty" json:"name"`
	Artist string             `bson:"artist,omitempty" json:"artist"`
	Genres []string           `bson:"genres" json:"genres"`
}

type SongListParams struct {
	Page     int64  `form:"page" binding:"gt=0"`
	PageSize int64  `form:"pageSize" binding:"gt=0"`
	Search   string `form:"search"`
}

type CreateSongDto struct {
	Name   string   `json:"name" binding:"required"`
	Artist string   `json:"artist" binding:"required"`
	Genres []string `json:"genres" binding:"dive,required"`
}

type UpdateSongDto struct {
	Name   string   `json:"name" binding:"required"`
	Artist string   `json:"artist" binding:"required"`
	Genres []string `json:"genres" binding:"dive,required"`
}

type UpdateSongURI struct {
	ID string `uri:"id"`
}

var searchFields = []string{"artist", "name"}

func ListSongs(client *mongo.Client) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		log := logger.GetInstance()
		params := SongListParams{
			Page:     1,
			PageSize: 10,
		}
		if ctx.BindQuery(&params) != nil {
			ctx.AbortWithStatusJSON(400, gin.H{"message": "bad query params"})
			return
		}

		var songs []Song = make([]Song, 0)
		songCol := client.Database("sldb").Collection("songs")
		skip := (params.Page - 1) * params.PageSize
		opts := options.FindOptions{
			Skip:  &skip,
			Limit: &params.PageSize,
		}

		filter := bson.M{}
		if len(params.Search) > 0 {
			log.Println("set search")
			conds := []bson.M{}
			for _, v := range searchFields {
				f := bson.M{v: bson.M{"$regex": params.Search, "$options": "i"}}
				conds = append(conds, f)
			}
			filter["$or"] = conds
		}

		songCursor, err := songCol.Find(ctx, filter, &opts)
		if err != nil {
			log.Println("error:", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error querying data",
			})
			return
		}
		if err = songCursor.All(ctx, &songs); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error querying data",
			})
			return
		}
		count, err := songCol.CountDocuments(ctx, filter)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error querying data",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"data":    songs,
			"records": count,
		})
	}
}

func CreateSong(client *mongo.Client) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		body := CreateSongDto{
			Genres: []string{},
		}
		if err := ctx.BindJSON(&body); err != nil {
			logger.GetInstance().Println(err)
			ctx.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
			return
		}

		songCol := client.Database("sldb").Collection("songs")
		res, err := songCol.InsertOne(context.TODO(), Song{
			Name:   body.Name,
			Artist: body.Artist,
			Genres: body.Genres,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{"message": "Create song failed"})
			return
		}

		findRes := songCol.FindOne(context.TODO(), bson.M{"_id": res.InsertedID})
		created := Song{}
		err = findRes.Decode(&created)
		if err != nil {
			ctx.AbortWithStatusJSON(203, gin.H{"_id": res.InsertedID})
			return
		}

		ctx.JSON(200, created)
	}
}

func UpdateSong(client *mongo.Client) func(ctx *gin.Context) {
	col := client.Database("sldb").Collection("songs")
	return func(ctx *gin.Context) {
		body := UpdateSongDto{
			Genres: []string{},
		}
		uri := UpdateSongURI{}
		log := logger.GetInstance()
		if err := ctx.BindUri(&uri); err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
			return
		}

		if err := ctx.BindJSON(&body); err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
			return
		}

		log.Println(uri.ID)
		id, err := primitive.ObjectIDFromHex(uri.ID)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{"message": "internal server error"})
			return
		}

		after := options.After
		opt := options.FindOneAndReplaceOptions{
			ReturnDocument: &after,
		}
		r := col.FindOneAndReplace(context.TODO(), bson.M{"_id": id}, body, &opt)

		updated := Song{}
		err = r.Decode(&updated)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(200, updated)
	}
}
