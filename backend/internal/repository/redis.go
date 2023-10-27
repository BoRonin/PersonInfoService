package repository

import (
	"context"
	"emtest/internal/models"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rc *redis.Client
}

const (
	sessionTime = time.Hour * 24 * 2
)

func NewRedis(r *redis.Client) *Redis {
	return &Redis{
		rc: r,
	}
}

// GetName проверяет, если в редисе есть данные
func (r *Redis) GetName(ctx context.Context, name string) (models.PersonInfo, error) {
	pi := models.PersonInfo{}
	//Смотрим результат в редисе. Если его нет, возвращаем ошибку
	result, err := r.rc.Get(ctx, name).Result()
	if err == redis.Nil {
		return pi, err
	}
	//Продлеваем время жизни ключа, так как он пользуется спросом
	r.rc.ExpireGT(ctx, name, sessionTime)
	//Приводим объект в PersonInfo и возвращаем
	if err = json.Unmarshal([]byte(result), &pi); err != nil {
		return pi, err
	}
	return pi, nil
}

// StoreName вводит информацию с ключом по имени в редис
func (r *Redis) StoreName(ctx context.Context, name string, pi models.PersonInfo) error {
	//Маршалим информацию. Если не выходит, возвращаем ошибку
	info, err := json.Marshal(pi)
	if err != nil {
		return err
	}
	//Кладем в редис. Если не выходит, возвращаем ошибку
	err = r.rc.Set(ctx, name, info, sessionTime).Err()
	if err != nil {
		return err
	}
	return nil
}
