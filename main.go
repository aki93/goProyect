package main

import (
	"context"
	"go/goProyect/awsgo"
	"go/goProyect/database"
	"go/goProyect/handlers"
	"go/goProyect/models"
	"go/goProyect/secretmanager"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
)

func main() {
	//Llamada y comienzo de la LAMBDA
	lambda.Start(EjecutoLambda)
}

func EjecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var response *events.APIGatewayProxyResponse

	awsgo.InicializadorAWS()

	if !ValidarParametros() {
		response = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en las variables de entorno, deben incluir 'SecretName','BucketName','UrlPrefix'",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		return response, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		response = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error en la lectura del Secret " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return response, nil
	}

	//limpiar el path parameter de la url del request
	path := strings.Replace(request.PathParameters["goProyect"], os.Getenv("UrlPrefix"), "", -1)
	//models.Key("method") -> el lenguaje sugiere que no utilicemos un string literal para claves
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtsign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketname"), os.Getenv("BucketName"))

	//Chequeo Conexion a la DB
	err = database.DBConnection(awsgo.Ctx)
	if err != nil {
		response = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error conectando la DB " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return response, nil
	}

	responseAPI := handlers.Handlers(awsgo.Ctx, request)
	if responseAPI.CustomResp == nil {
		response = &events.APIGatewayProxyResponse{
			StatusCode: responseAPI.Status,
			Body:       responseAPI.Message,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return response, nil
	} else {
		return responseAPI.CustomResp, nil
	}

}

func ValidarParametros() bool {
	_, traeParametro := os.LookupEnv("SecretName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("BucketName")
	if !traeParametro {
		return traeParametro
	}

	_, traeParametro = os.LookupEnv("UrlPrefix")
	if !traeParametro {
		return traeParametro
	}

	return true
}
