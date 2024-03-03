package main

import (
	"fmt"
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

type TodosHandlerGet struct {
	Svc *TodosServiceGet
}

func NewTodosHandlerGet(svc *TodosServiceGet) *TodosHandlerGet {
	return &TodosHandlerGet{Svc: svc}
}

func (*TodosHandlerGet) Pattern() string {
	return "/:id"
}

func (*TodosHandlerGet) Method() string {
	return http.MethodGet
}

func (h *TodosHandlerGet) Service() gin.HandlerFunc {
	return func(c *gin.Context) {
		todo, err := h.Svc.Call(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint: errcheck
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"msg":   todo.Name,
		})
	}
}

type TodosHandlerPost struct {
	Svc *TodosServicePost
}

func NewTodosHandlerPost(svc *TodosServicePost) *TodosHandlerPost {
	return &TodosHandlerPost{Svc: svc}
}

func (*TodosHandlerPost) Pattern() string {
	return "/"
}

func (*TodosHandlerPost) Method() string {
	return http.MethodPost
}

func (h *TodosHandlerPost) Service() gin.HandlerFunc {
	return func(c *gin.Context) {
		todoNew := Todo{}
		if err := c.BindJSON(&todoNew); err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint: errcheck
			return
		}
		fmt.Println("RIGHT BEFORE", todoNew)
		todo, err := h.Svc.Call(todoNew)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint: errcheck
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"msg":   todo.Name,
		})
	}
}

// NOTE: Interface Compliance Verification
var _ Route = (*TodosHandlerGet)(nil)
var _ Route = (*TodosHandlerPost)(nil)

func RegisterTodosApi(v1 *RouterGroupV1, routes []Route) {
	todos := v1.Group("/todos")
	for _, route := range routes {
		todos.Handle(route.Method(), route.Pattern(), route.Service())
	}
}
