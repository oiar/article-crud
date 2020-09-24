package main

import (
	"database/sql"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	admin "github.com/yhyddr/article-crud/admin/controller"
	"github.com/yhyddr/article-crud/article/controller"
)

var (
	// JWTMiddleware should be exported for user authentication.
	JWTMiddleware *jwt.GinJWTMiddleware
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3307)/?parseTime=true")
	if err != nil {
		panic(err)
	}

	adminCon := admin.New(dbConn)
	articleCon := controller.New(dbConn, "article", "article.article_crud")

	adminCon.RegisterRouter(router.Group("/api/v1/admin"))

	JWTMiddleware = &jwt.GinJWTMiddleware{
		Realm:   "Template",
		Key:     []byte("hydra"),
		Timeout: 24 * time.Hour,
	}

	adminCon.ExtendJWTMiddleWare(JWTMiddleware)

	router.POST("/api/v1/admin/login", JWTMiddleware.LoginHandler)
	articleCon.RegisterRouter(router.Group("/api/v1/article"), JWTMiddleware)

	router.Run(":8000")
}
