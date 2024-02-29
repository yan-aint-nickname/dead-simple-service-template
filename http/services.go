package main

import (
	"encoding/json"
)

type Service interface {
	Call()
}

type TodosServiceGet struct {
	Redis *RedisClient
}

func NewTodosService(redis_client *RedisClient) *TodosServiceGet {
	return &TodosServiceGet{Redis: redis_client}
}

type Todo struct {
	Id string
	Name string
}

func (svc TodosServiceGet) Call(todoId string) (todo Todo, err error) {
	todoRaw, err := svc.Redis.Get(todoId)
	if err != nil {
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
