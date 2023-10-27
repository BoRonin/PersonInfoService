package fetcher

import (
	"context"
	"emtest/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Fetcher struct {
	AgeUrl         string
	GenderUrl      string
	NationalityURL string
}
type Country struct {
	CountryId   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func NewFetcher(age string, gender string, nationality string) *Fetcher {
	return &Fetcher{
		AgeUrl:         age,
		GenderUrl:      gender,
		NationalityURL: nationality,
	}
}

// FetchInfo делает запрос на API для получения информации по имени
func (f *Fetcher) FetchInfo(ctx context.Context, name string) models.PersonInfo {
	start := time.Now()
	info := models.PersonInfo{}
	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		start1 := time.Now()
		resp, err := http.Get(fmt.Sprintf("%s?name=%s", f.AgeUrl, name))
		if err != nil {
			log.Println(err)
		}
		response := struct {
			Age int `json:"age"`
		}{}
		json.NewDecoder(resp.Body).Decode(&response)
		defer resp.Body.Close()
		info.Age = response.Age
		log.Println("goroutine 1 time", time.Since(start1))
		wg.Done()
	}()
	go func() {
		start2 := time.Now()
		resp, err := http.Get(fmt.Sprintf("%s?name=%s", f.GenderUrl, name))
		if err != nil {
			log.Println(err)
		}
		response := struct {
			Gender string `json:"gender"`
		}{}
		json.NewDecoder(resp.Body).Decode(&response)
		defer resp.Body.Close()
		info.Gender = response.Gender
		log.Println("goroutine 2 time", time.Since(start2))
		wg.Done()
	}()
	go func() {
		start3 := time.Now()
		resp, err := http.Get(fmt.Sprintf("%s?name=%s", f.NationalityURL, name))
		if err != nil {
			log.Println(err)
		}
		response := struct {
			Country []Country `json:"country"`
		}{}
		json.NewDecoder(resp.Body).Decode(&response)
		defer resp.Body.Close()
		if len(response.Country) > 0 {
			info.Nationality = response.Country[0].CountryId
		}
		log.Println("goroutine 3 time", time.Since(start3))
		wg.Done()
	}()
	wg.Wait()
	log.Println("time to fetch:", time.Since(start))
	return info
}
