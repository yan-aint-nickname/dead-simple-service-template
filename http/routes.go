package main

import (
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

// NOTE: I don't know which is better
// Svc *TodosServiceGet tight bonding
// Svc Service[string, Todo] tight bonding also? Because I specify in and out types
// But at least developer will know that Svc has method `Call`
type TodosHandlerGet struct {
	Svc Service[string, Todo]
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

// NOTE: see TodoHanlderGet
type TodosHandlerPost struct {
	Svc Service[Todo, Todo]
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

// NOTE: see TodoHanlderGet
type TodosHandlerDelete struct {
	Svc Service[string, string]
}

func NewTodosHandlerDelete(svc *TodosServiceDelete) *TodosHandlerDelete {
	return &TodosHandlerDelete{Svc: svc}
}

func (*TodosHandlerDelete) Pattern() string {
	return "/:id"
}

func (*TodosHandlerDelete) Method() string {
	return http.MethodDelete
}

func (h *TodosHandlerDelete) Service() gin.HandlerFunc {
	return func(c *gin.Context) {
		idDel, err := h.Svc.Call(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err) // nolint: errcheck
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"msg":   idDel,
		})
	}
}

// NOTE: Interface Compliance Verification
var _ Route = (*TodosHandlerGet)(nil)
var _ Route = (*TodosHandlerPost)(nil)
var _ Route = (*TodosHandlerDelete)(nil)

func RegisterTodosApi(v1 *RouterGroupV1, routes []Route) {
	todos := v1.Group("/todos")
	for _, route := range routes {
		todos.Handle(route.Method(), route.Pattern(), route.Service())
	}
}
