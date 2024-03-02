package main

import (
	"encoding/json"
	"fmt"
)

type Service interface {
	Call()
}

type TodosServiceGet struct {
	Redis    *RedisClient
	Postgres *Postgres
}

func NewTodosService(redis_client *RedisClient, postgres *Postgres) *TodosServiceGet {
	return &TodosServiceGet{Redis: redis_client, Postgres: postgres}
}

type Todo struct {
	Id   string
	Name string
}

func (svc TodosServiceGet) Call(todoId string) (todo Todo, err error) {
	todoRaw, err := svc.Redis.Get(todoId)
	if err != nil {
		err = svc.Postgres.Pool.QueryRow(
			svc.Postgres.Ctx,
			"select id, name from todos where id=$1",
			todoId,
		).Scan(&todo.Id, &todo.Name)
		fmt.Println("GOT FROM POSTGRES")
		return
	}

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
