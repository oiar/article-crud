package controller

import (
	"database/sql"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/yhyddr/article-crud/article/mysql"
)

type ArticleController struct {
	db        *sql.DB
	DBName    string
	tableName string
}

func New(db *sql.DB, DBName string, tableName string) *ArticleController {
	return &ArticleController{
		db:        db,
		DBName:    DBName,
		tableName: tableName,
	}
}

// RegisterRouter
func (a *ArticleController) RegisterRouter(r gin.IRouter, JWT *jwt.GinJWTMiddleware) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := mysql.CreateDB(a.db, a.DBName)

	err = mysql.CreateTable(a.db, a.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.Use(JWT.MiddlewareFunc())
	{
		r.POST("/create", a.create)
		r.POST("/delete", a.deleteByID)
		r.POST("/update", a.updateByID)
		r.POST("/query", a.queryByID)
	}

}

// create
func (a *ArticleController) create(c *gin.Context) {
	var (
		req struct {
			UserId      int    `json:"userid"         binding:"required"`
			ArticleName string `json:"articlename"    binding:"required"`
			Author      string `json:"author"         binding:"required"`
			Content     string `json:"content"        binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	id, err := mysql.CreateArticle(a.db, a.tableName, req.UserId, req.ArticleName, req.Author, req.Content)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ID": id})
}

// delete
func (a *ArticleController) deleteByID(c *gin.Context) {
	var (
		req struct {
			ArticleID int `json:"articleid"    binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.DeleteArticleByID(a.db, a.tableName, req.ArticleID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

// update
func (a *ArticleController) updateByID(c *gin.Context) {
	var (
		req struct {
			ArticleID int    `json:"articleid"     binding:"required"`
			Content   string `json:"content"       binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = mysql.UpdateArticleByID(a.db, a.tableName, req.Content, req.ArticleID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

// query
func (a *ArticleController) queryByID(c *gin.Context) {
	var (
		req struct {
			ArticleID int `json:"articleid"     binding:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	art, err := mysql.QueryArticleByID(a.db, a.tableName, req.ArticleID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "art": art})
}
