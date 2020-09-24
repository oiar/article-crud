package controller

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yhyddr/article-crud/admin/mysql"
)

type Controller struct {
	db *sql.DB
}

func New(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

//register router
func (c *Controller) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	name := "admin"
	password := "admin"
	err := mysql.CreateTable(c.db, &name, &password)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/create", c.create)
}

func (c *Controller) create(ctx *gin.Context) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"omitempty,min=5,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	//Default password
	if admin.Password == "" {
		admin.Password = "123456"
	}

	err = mysql.Create(c.db, &admin.Name, &admin.Password)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

//Login JWT validation
func (c *Controller) Login(ctx *gin.Context) (uint32, error) {
	var (
		admin struct {
			Name     string `json:"name"      binding:"required,alphanum,min=5,max=30"`
			Password string `json:"password"  binding:"omitempty,min=5,max=30"`
		}
	)

	err := ctx.ShouldBind(&admin)
	if err != nil {
		return 0, err
	}

	ID, err := mysql.Login(c.db, &admin.Name, &admin.Password)
	if err != nil {
		return 0, err
	}

	return ID, nil
}
