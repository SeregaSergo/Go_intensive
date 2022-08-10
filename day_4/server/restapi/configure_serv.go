// This file is safe to edit. Once it exists it will not be overwritten

package restapi

//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//
//unsigned int i;
//unsigned int argscharcount = 0;
//
//char *ask_cow(char phrase[]) {
//	int phrase_len = strlen(phrase);
//	char *buf = (char *)malloc(sizeof(char) * (160 + (phrase_len + 2) * 3));
//	strcpy(buf, " ");
//
//	for (i = 0; i < phrase_len + 2; ++i) {
//	strcat(buf, "_");
//	}
//
//	strcat(buf, "\n< ");
//	strcat(buf, phrase);
//	strcat(buf, " ");
//	strcat(buf, ">\n ");
//
//	for (i = 0; i < phrase_len + 2; ++i) {
//		strcat(buf, "-");
//	}
//	strcat(buf, "\n");
//	strcat(buf, "        \\   ^__^\n");
//	strcat(buf, "         \\  (oo)\\_______\n");
//	strcat(buf, "            (__)\\       )\\/\\\n");
//	strcat(buf, "                ||----w |\n");
//	strcat(buf, "                ||     ||\n");
//	return buf;
//}
import "C"
import (
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"net/http"
	calculator "swagger/code-gen/server"
	operations2 "swagger/code-gen/server/restapi/operations"
	"unsafe"
)

//go:generate swagger generate server --target ../../day_4 --name Serv --spec ../schema.yaml --principal interface{}

func configureFlags(api *operations2.ServAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations2.ServAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.BuyCandyHandler = operations2.BuyCandyHandlerFunc(func(params operations2.BuyCandyParams) middleware.Responder {
		change, err := calculator.GetChange(*params.Order.CandyType, *params.Order.Money, *params.Order.CandyCount)
		if err != nil {
			return operations2.NewBuyCandyBadRequest().WithPayload(&operations2.BuyCandyBadRequestBody{
				Error: "some error in input data",
			})
		}
		if change < 0 {
			return operations2.NewBuyCandyPaymentRequired().WithPayload(&operations2.BuyCandyPaymentRequiredBody{
				Error: "not enough money",
			})
		} else {
			phrase := C.CString("Thank you!")
			defer C.free(unsafe.Pointer(phrase))
			ptr := C.ask_cow(phrase)
			defer C.free(unsafe.Pointer(ptr))
			thanks := C.GoString(ptr)
			return operations2.NewBuyCandyCreated().WithPayload(&operations2.BuyCandyCreatedBody{
				Change: change,
				Thanks: thanks,
			})
		}
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
