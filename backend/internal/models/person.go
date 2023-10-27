package models

type Person struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}
type PersonInfo struct {
	Gender      string `json:"gender"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
}

type PersonFull struct {
	Person
	PersonInfo
}
type PersonFilter struct {
	Name        []string `json:"name"`
	Surname     []string `json:"surname"`
	Patronymic  []string `json:"patronymic"`
	Gender      string   `json:"gender"`
	Age         []int    `json:"age"`
	AgeGT       int      `json:"age_gt"`
	AgeLT       int      `json:"age_lt"`
	Nationality []string `json:"nationality"`
}
