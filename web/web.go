package web

import (
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rekey/go-club/dao"
)

//go:embed front
var frontFS embed.FS

func Run() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.RequestLogger())

	// 嵌入式静态文件服务
	e.GET("/", func(c echo.Context) error {
		data, err := frontFS.ReadFile("front/index.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return c.HTMLBlob(http.StatusOK, data)
	})

	g := e.Group("/api")
	g.GET("/task/all", func(c echo.Context) error {
		tasks, err := dao.GetAllTask()
		if err != nil {
			return c.JSON(200, echo.Map{"code": 3, "msg": "get tasks error"})
		}
		return c.JSON(200, tasks)
	})
	g.GET("/task/item", func(c echo.Context) error {
		u := c.QueryParam("url")
		task, err := dao.FindTaskByURL(u)
		if err != nil {
			return c.JSON(200, echo.Map{"code": 2, "msg": "task not found"})
		}
		return c.JSON(200, task)
	})
	g.GET("/task/add", func(c echo.Context) error {
		u := c.QueryParam("url")
		task := dao.CreateTask(u)
		if task == nil {
			return c.JSON(200, echo.Map{"code": 1, "msg": "url not support"})
		}
		return c.JSON(200, task)
	})
	g.GET("/task/start", func(c echo.Context) error {
		u := c.QueryParam("url")
		task, err := dao.FindTaskByURL(u)
		if err != nil {
			return c.JSON(200, echo.Map{"code": 2, "msg": "task not found"})
		}
		err = task.Start()
		if err != nil {
			return c.JSON(200, echo.Map{"code": 4, "msg": "start task failed"})
		}
		return c.JSON(200, task)
	})
	err := e.Start(":8888")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
