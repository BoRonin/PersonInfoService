package service

import (
	"context"
	"emtest/internal/models"
	"emtest/pkg/logger"
	"fmt"
	"strconv"

	"log/slog"
)

type CacheWorker interface {
	GetName(ctx context.Context, name string) (models.PersonInfo, error)
	StoreName(ctx context.Context, name string, pi models.PersonInfo) error
}

type DBWorker interface {
	DeletePersonByID(ctx context.Context, id int) error
	StorePersonToDB(ctx context.Context, req models.PersonFull) error
	FilterPersons(ctx context.Context, filter models.HttpFilterRequest) (models.HttpSearchResponse, error)
	CheckPerson(ctx context.Context, person models.Person) bool
	UpdatePersonById(ctx context.Context, person models.PersonFull) (models.PersonFull, error)
}
type Fetcher interface {
	FetchInfo(ctx context.Context, name string) models.PersonInfo
}

type PersonService struct {
	DB      DBWorker
	Cache   CacheWorker
	Fetcher Fetcher
	Log     *slog.Logger
}

// NewPersonService создает новый объект с возможностью работы с базой данный, кешем и отправкой http запросов к стороннему api
func NewPersonService(cw CacheWorker, db DBWorker, f Fetcher, log *slog.Logger) *PersonService {
	return &PersonService{
		DB:      db,
		Cache:   cw,
		Fetcher: f,
		Log:     log,
	}
}

// StorePerson достает информацию и сохраняет человека в базу или убеждается, что он уже есть, и выходит.
func (ps *PersonService) StorePerson(ctx context.Context, req models.HttpPostRequest) error {
	log := ps.Log.With(slog.String("op", "ps.StorePerson"))
	//Проверяем, есть ли человек в базе
	if ps.DB.CheckPerson(ctx, req.Person) {
		log.Debug("person found in db:", logger.Val("person", fmt.Sprintf("%s %s %s", req.Name, req.Surname, req.Patronymic)))
		return nil
	}

	//Проверяем, есть ли человек в кеше
	info, err := ps.Cache.GetName(ctx, req.Name)
	if err != nil {
		//Если нет, то делаем http запрос
		log.Debug("fetching name", logger.Val("name", req.Name))
		info = ps.Fetcher.FetchInfo(ctx, req.Name)
		//Кладем в редис, если ошибка, логируем
		log.Debug("storing in redis", logger.Val("info", fmt.Sprintf("%s, %d, %s", info.Gender, info.Age, info.Nationality)))
		err := ps.Cache.StoreName(ctx, req.Name, info)
		if err != nil {
			log.Debug(err.Error())
		}
	}
	p := models.PersonFull{
		Person:     req.Person,
		PersonInfo: info,
	}
	log.Debug("storing in db", logger.Val("person", fmt.Sprintf("%s %s %s %s %d %s", p.Name, p.Surname, p.Patronymic, p.Gender, p.Age, p.Nationality)))
	//Созраняем человека с обогащенной информацией в базу
	err = ps.DB.StorePersonToDB(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

// FilterPersons делает запрос в БД и ищет людей по фильтрам
func (ps *PersonService) FilterPersons(ctx context.Context, filter models.HttpFilterRequest) (models.HttpSearchResponse, error) {
	flog := ps.Log.With(slog.String("op", "ps.FilterPerson"))
	flog.Debug("filtering in db", logger.Val("filter", fmt.Sprintf("%v", filter)))
	return ps.DB.FilterPersons(ctx, filter)
}

func (ps *PersonService) DeletePersonByID(ctx context.Context, id int) error {
	dplog := ps.Log.With(slog.String("op", "ps.DeletePersonByID"))
	dplog.Debug("deleting from db", logger.Val("id", strconv.Itoa(id)))
	if err := ps.DB.DeletePersonByID(ctx, id); err != nil {
		dplog.Debug("couldn't delete from db", logger.Val("id", strconv.Itoa(id)))
		return err
	}
	return nil
}
