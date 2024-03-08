package main

import (
	"encoding/json"
)

type Service[T, K any] interface {
	Call(T) (K, error)
}

type Todo struct {
	Id   string
	Name string
}

type TodosServiceGet struct {
	Redis    *RedisClient
	Postgres *Postgres
}

func NewTodosServiceGet(redis_client *RedisClient, postgres *Postgres) *TodosServiceGet {
	return &TodosServiceGet{Redis: redis_client, Postgres: postgres}
}

func (svc TodosServiceGet) Call(todoId string) (todo Todo, err error) {
	todo, err = svc.getFromCache(todoId)

	if err != nil {
		if err = svc.Postgres.Pool.QueryRow(
			svc.Postgres.Ctx,
			"select id, name from todos where id=$1",
			todoId,
		).Scan(&todo.Id, &todo.Name); err != nil {
			return
		}
		errCh := make(chan error)
		go func() {
			errCh <- svc.setToCache(todo)
		}()
		// NOTE: syntax is such a disaster here, but it is what it is
		err = <-errCh
		return
	}

	return
}

func (svc TodosServiceGet) setToCache(todo Todo) (err error) {
	todoRaw, err := json.Marshal(todo)
	if err != nil {
		return
	}
	return svc.Redis.Set(todo.Id, todoRaw)
}

func (svc TodosServiceGet) getFromCache(todoId string) (todo Todo, err error) {
	todoRaw, err := svc.Redis.Get(todoId)
	if err != nil {
		return
	}
	if err = json.Unmarshal(todoRaw, &todo); err != nil {
		return
	}
	return
}

type TodosServicePost struct {
	Redis    *RedisClient
	Postgres *Postgres
}

func NewTodosServicePost(redis_client *RedisClient, postgres *Postgres) *TodosServicePost {
	return &TodosServicePost{Redis: redis_client, Postgres: postgres}
}

func (svc TodosServicePost) Call(todoNew Todo) (todo Todo, err error) {
	err = svc.Postgres.Pool.QueryRow(
		svc.Postgres.Ctx,
		"insert into todos(name) values ($1) returning id, name",
		todoNew.Name,
	).Scan(&todo.Id, &todo.Name)
	return
}

type TodosServiceDelete struct {
	Redis    *RedisClient
	Postgres *Postgres
}

func NewTodosServiceDelete(redis *RedisClient, postgres *Postgres) *TodosServiceDelete {
	return &TodosServiceDelete{Redis: redis, Postgres: postgres}
}

func (svc TodosServiceDelete) Call(todoId string) (status string, err error) {

	// TODO: add redis deleting of a given id
	// errCh := make(chan error)
	// go func() {
	// 	errCh <- svc.Redis.Delete(todo)
	// }()
	// err = <-errCh
	// var status string

	commandTag, err := svc.Postgres.Pool.Exec(
		svc.Postgres.Ctx,
		"delete from todos where id=$1",
		todoId,
	)
	status = commandTag.String()
	return
}

var _ Service[string, Todo] = (*TodosServiceGet)(nil)
var _ Service[Todo, Todo] = (*TodosServicePost)(nil)
var _ Service[string, string] = (*TodosServiceDelete)(nil)


type ProjectsServiceGet struct {
	Api *ProjectsAPI
}

func NewProjectsServiceGet(api *ProjectsAPI) *ProjectsServiceGet {
	return &ProjectsServiceGet{Api: api}
}

// I didn't figure out how to pass a "nil" or "empty" type to a function
// Optional parameters isn't very helpfull :(
func (svc *ProjectsServiceGet) Call(params map[string]string) ([]Project, error) {
	projects, err := svc.Api.GetProjects()
	if err != nil {
		return []Project{}, err
	}
	return projects.Projects, nil
}

var _ Service[map[string]string, []Project] = (*ProjectsServiceGet)(nil)
