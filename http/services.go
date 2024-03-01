package main

import (
	"fmt"
	"encoding/json"
)

type Service interface {
	Call()
}

type TodosServiceGet struct {
	Redis *RedisClient
	Postgres *Postgres
}

func NewTodosService(redis_client *RedisClient, postgres *Postgres) *TodosServiceGet {
	return &TodosServiceGet{Redis: redis_client, Postgres: postgres}
}

type Todo struct {
	Id string
	Name string
}

func (svc TodosServiceGet) Call(todoId string) (todo Todo, err error) {
	todoRaw, err := svc.Redis.Get(todoId)

	var greeting string
	err = svc.Postgres.Pool.QueryRow(svc.Postgres.Ctx, "select 'Hello world!'").Scan(&greeting)

	if err != nil {
		return
	}
	fmt.Println("HELLO THERE:", greeting)
	if err = json.Unmarshal(todoRaw, &todo); err != nil {
		return
	}
	return
}


// func (svc *TodosService) Set(todo Todo) error {
// 	todoRaw, err := json.Marshal(todo)
// 	if err != nil {
// 		return err
// 	}
// 	if err := svc.Redis.Set(todo.Id, todoRaw); err != nil {
// 		return err
// 	}
// 	return nil
// }
