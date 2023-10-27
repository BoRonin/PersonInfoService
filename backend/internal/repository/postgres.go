package repository

import (
	"context"
	"emtest/internal/models"
	"emtest/internal/storage/postgres"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Masterminds/squirrel"
)

type Postgres struct {
	pg *postgres.Postgres
}

func NewPostgresService(bd *postgres.Postgres) *Postgres {
	return &Postgres{
		pg: bd,
	}
}

func (p *Postgres) StorePersonToDB(ctx context.Context, req models.PersonFull) error {
	q := `INSERT INTO people (name, surname, patronymic, gender, age, nationality) values ($1, $2, $3, $4, $5, $6)`
	_, err := p.pg.Pool.Exec(ctx, q, req.Name, req.Surname, req.Patronymic, req.Gender, req.Age, req.Nationality)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) FilterPersons(ctx context.Context, filter models.HttpFilterRequest) (models.HttpSearchResponse, error) {
	lookingFor := `id, name, surname, patronymic, gender, age, nationality`
	sql := p.pg.Builder.Select(lookingFor + ` from people`).Limit(uint64(filter.PerPage))
	if len(filter.Name) > 0 {
		sql = sql.Where(squirrel.Eq{"name": filter.Name})
	}
	if len(filter.Gender) > 0 {
		sql = sql.Where(squirrel.Eq{"gender": filter.Gender})
	}
	if len(filter.Patronymic) > 0 {
		sql = sql.Where(squirrel.Eq{"patronymic": filter.Patronymic})
	}
	if len(filter.Surname) > 0 {
		sql = sql.Where(squirrel.Eq{"surname": filter.Surname})
	}
	if len(filter.Age) > 0 {
		sql = sql.Where(squirrel.Eq{"age": filter.Age})
	}
	if filter.AgeGT != 0 {
		sql = sql.Where(`age > ?`, filter.AgeGT)
	}
	if filter.AgeLT != 0 {
		sql = sql.Where(`age < ?`, filter.AgeLT)
	}
	sql2 := sql
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "age_asc":
			sql = sql.OrderBy("age ASC")
		case "age_desc":
			sql = sql.OrderBy("age DESC")
		case "name_desc":
			sql = sql.OrderBy("name DESC")
		case "name_asc":
			sql = sql.OrderBy("name ASC")
		case "nationality_asc":
			sql = sql.OrderBy("nationality ASC")
		case "nationality_desc":
			sql = sql.OrderBy("nationality DESC")
		case "gender_desc":
			sql = sql.OrderBy("gender DESC")
		case "gender_asc":
			sql = sql.OrderBy("gender ASC")
		case "surname_desc":
			sql = sql.OrderBy("surname DESC")
		case "surname_asc":
			sql = sql.OrderBy("surname ASC")
		default:
			break
		}
	}
	if filter.Page > 1 {
		sql = sql.Offset(uint64((filter.Page - 1) * filter.PerPage))
	}
	q, args, err := sql.ToSql()
	if err != nil {
		return models.HttpSearchResponse{}, err
	}
	q2, args2, _ := sql2.ToSql()
	q2 = strings.Replace(q2, fmt.Sprintf("SELECT %s", lookingFor), "SELECT COUNT(*)", 1)
	var result models.HttpSearchResponse
	result.PerPage = filter.PerPage
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_ = p.pg.Pool.QueryRow(ctx, q2, args2...).Scan(&result.Quantity)
		defer wg.Done()
		if result.Quantity != 0 {
			if result.Quantity%filter.PerPage > 0 {
				result.LastPage = result.Quantity/filter.PerPage + 1
			} else {
				result.LastPage = result.Quantity / filter.PerPage
			}
		}
		result.Page = filter.Page
	}()
	rows, err := p.pg.Pool.Query(ctx, q, args...)
	if err != nil {
		return models.HttpSearchResponse{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.PersonFull
		if err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Gender, &p.Age, &p.Nationality); err != nil {
			return models.HttpSearchResponse{}, err
		}
		result.Persons = append(result.Persons, p)
	}
	wg.Wait()

	return result, nil
}

func (p *Postgres) CheckPerson(ctx context.Context, person models.Person) bool {
	q := `SELECT EXISTS(SELECT 1 FROM people (name, surname, patronymic) = ($1, $2, $3))`
	var b bool
	if err := p.pg.Pool.QueryRow(ctx, q, person.Name, person.Surname, person.Patronymic).Scan(&b); err != nil {
		return false
	}
	return b
}

func (p *Postgres) DeletePersonByID(ctx context.Context, id int) error {
	q := `DELETE FROM people WHERE id = $1`
	_, err := p.pg.Pool.Exec(ctx, q, id)
	log.Println("Deleting", id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdatePersonById(ctx context.Context, person models.PersonFull) (models.PersonFull, error) {
	pf := models.PersonFull{}
	sql := p.pg.Builder.Update("people").Where("id = ?", person.ID).Suffix("returning id, name, surname, patronymic, age, gender, nationality")
	if person.Name != "" {
		sql = sql.Set("name", person.Name)
	}
	if person.Surname != "" {
		sql = sql.Set("surname", person.Surname)
	}
	if person.Gender != "" {
		sql = sql.Set("gender", person.Gender)
	}
	if person.Patronymic != "" {
		sql = sql.Set("patronymic", person.Patronymic)
	}
	if person.Nationality != "" {
		sql = sql.Set("nationality", person.Nationality)
	}
	if person.Age != 0 {
		sql = sql.Set("age", person.Age)
	}
	q, args, err := sql.ToSql()
	if err != nil {
		return models.PersonFull{}, err
	}
	if err := p.pg.Pool.QueryRow(ctx, q, args...).Scan(&pf.ID, &pf.Name, &pf.Surname, &pf.Patronymic, &pf.Age, &pf.Gender, &pf.Nationality); err != nil {
		return models.PersonFull{}, nil
	}
	return pf, nil
}
