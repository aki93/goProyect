package handlers

import (
	"context"
	"fmt"
	"go/goProyect/models"

	"github.com/aws/aws-lambda-go/events"
)

func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.ResponseApi {

	//loggeo el path y tipo de request que se ejecuto -> ej : /login > post
	fmt.Println("Voy a procesar " + ctx.Value(models.Key("path")).(string) + " > " + ctx.Value(models.Key("method")).(string))

	var response models.ResponseApi
	response.Status = 400

	switch ctx.Value(models.Key("method")).(string) {
	case "POST":
		switch ctx.Value(models.Key("path")).(string) {
		//POST endpoints
		}

	case "GET":
		switch ctx.Value(models.Key("path")).(string) {
		//GET endpoints
		}

	case "PUT":
		switch ctx.Value(models.Key("path")).(string) {
		//PUT endpoints
		}

	case "DELETE":
		switch ctx.Value(models.Key("path")).(string) {
		//DELETE endpoints
		}
	}

	response.Message = "Method Invalid"
	return response

}
