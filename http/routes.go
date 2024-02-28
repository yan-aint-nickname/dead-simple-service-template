package main

import (
	"fmt"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
)

type RouterGroupV1 struct {
	*gin.RouterGroup
}

type Route interface {
	Pattern() string
	Method() string
	Service() gin.HandlerFunc
}

func NewApiV1Router(router *gin.Engine) *RouterGroupV1 {
	v1 := router.Group("/api/v1")
	return &RouterGroupV1{v1}
}

type TodosHandlerGet struct {}

func NewTodosHandlerGet() *TodosHandlerGet {
	return &TodosHandlerGet{}
}

func (*TodosHandlerGet) Pattern() string {
	return "/:id"
}

func (*TodosHandlerGet) Method() string {
	return http.MethodGet
}

func (*TodosHandlerGet) Service() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawId := c.Param("id")
		todoId, err := strconv.Atoi(rawId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"error": nil,
				"msg":    fmt.Sprintf("Some test todo id: %d", todoId),
			})
		}
	}
}

type TodosHandlerPost struct {}

func NewTodosHandlerPost() *TodosHandlerPost {
	return &TodosHandlerPost{}
}

func (*TodosHandlerPost) Pattern() string {
	return "/"
}

func (*TodosHandlerPost) Method() string {
	return http.MethodPost
}

func (*TodosHandlerPost) Service() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"msg":   "Some todo added",
		})
	}
}

func RegisterTodoApi(v1 *RouterGroupV1, routes []Route) {
	todos := v1.Group("/todos")
	for _, route := range routes {
		todos.Handle(route.Method(), route.Pattern(), route.Service())
	}
}
