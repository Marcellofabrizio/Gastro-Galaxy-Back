package models

type Recipe struct {
	Id              int
	CategoryId      int
	Name            string
	Url             string
	Description     string
	LongDescription string
}

type RecipeInputDto struct {
	CategoryId      int    `json:categoryId`
	Name            string `json:name`
	Url             string `json:url`
	Description     string `json:description`
	LongDescription string `json:longDescription`
	IngedientIds    []int  `json:ingedientIds`
}

type RecipeWithIngredientsDto struct {
	Recipe      Recipe
	Ingredients []Ingedient
}
