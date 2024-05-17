package models

type Recipe struct {
	Id          int
	CategoryId  int
	Name        string
	Url         string
	Description string
}

type RecipeDto struct {
	CategoryId   int    `json:categoryId`
	Name         string `json:name`
	Url          string `json:url`
	Description  string `json:description`
	IngedientIds []int  `json:ingedientIds`
}
