package firebase

import (
	"context"
	"fmt"

	fb "firebase.google.com/go"
	"google.golang.org/api/option"
)

var firebase *fb.App

func Instantiate(json []byte) {
	var err error
	firebase, err = fb.NewApp(context.Background(), nil, option.WithCredentialsJSON(json))
	if err != nil {
		fmt.Println("errored loading firebase admin sdk")
	}
}

func GetInstance() *fb.App {
	return (firebase)
}
