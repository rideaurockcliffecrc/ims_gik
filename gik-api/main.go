package main

import (
	"GIK_Web/database"
	"GIK_Web/env"
	"GIK_Web/src/routers"
	"fmt"
	"net/http"
)

func main() {
	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	// router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	// 	// your custom format
	// 	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
	// 		param.ClientIP,
	// 		param.TimeStamp.Format(time.RFC1123),
	// 		param.Method,
	// 		param.Path,
	// 		param.Request.Proto,
	// 		param.StatusCode,
	// 		param.Latency,
	// 		param.Request.UserAgent(),
	// 		param.ErrorMessage,
	// 	)
	// }))
	// router.Run(":8080")

	env.SetEnv()

	routersInit := routers.InitRouter()

	database.ConnectDatabase()

	server := &http.Server{
		Handler: routersInit,
		Addr:    env.WebserverHost + ":" + env.WebserverPort,
	}

	fmt.Printf("\nServer running at %s:%s\n", env.WebserverHost, env.WebserverPort)

	if env.HTTPS {
		server.ListenAndServeTLS(".cert/server.crt", ".cert/server.key")

	} else {
		server.ListenAndServe()
	}

}
