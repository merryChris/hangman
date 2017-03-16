package hangman

import (
	"github.com/merryChris/gameLord/types"
)

type Category struct {
	Id   int64   `orm:"column(id)" json:"id"`
	Name string  `orm:"column(name)" json:"name"`
	Dict []*Dict `orm:"reverse(many)" json:"-"`
}

type Dict struct {
	Id       int64     `orm:"column(id)" json:"id"`
	Category *Category `orm:"rel(fk)" json:"-"`
	Word     string    `orm:"column(word)" json:"word"`
}

type Hangman struct {
	Id           int64  `orm:"column(id)"`
	UserId       int64  `orm:"column(user_id);unique"`
	Word         string `orm:"column(word)"`
	LetterStatus int64  `orm:"column(letter_status)"`
}

type HangmanListCategoryJsonRequest struct {
	types.BaseJsonRequest
	UserId    int64  `json:"user_id"`
	GameId    int64  `json:"game_id"`
	GameToken string `json:"game_token"`
}

type HangmanListCategoryJsonResponse struct {
	types.BaseJsonResponse
	Categories []*Category `json:"categories"`
}

type HangmanPickWordJsonRequest struct {
	types.BaseJsonRequest
	UserId     int64  `json:"user_id"`
	GameId     int64  `json:"game_id"`
	CategoryId int64  `json:"category_id"`
	GameToken  string `json:"game_token"`
}

type HangmanPickWordJsonResponse struct {
	types.BaseJsonResponse
	CurrentWord  string `json:"current_word"`
	LetterStatus int64  `json:"status"`
}

type HangmanLoadStatusJsonRequest HangmanListCategoryJsonRequest

type HangmanLoadStatusJsonResponse HangmanPickWordJsonResponse

type HangmanValidateJsonRequest struct {
	types.BaseJsonRequest
	UserId        int64  `json:"user_id"`
	GameId        int64  `json:"game_id"`
	CurrentLetter byte   `json:"current_letter"`
	GameToken     string `json:"game_token"`
}

type HangmanValidateJsonResponse struct {
	types.BaseJsonResponse
	TrueOrFalse  bool   `json:"true_or_false"`
	Ending       bool   `json:"ending"`
	CurrentWord  string `json:"current_word"`
	LetterStatus int64  `json:"letter_status"`
}

type HangmanResetJsonRequest HangmanPickWordJsonRequest

type HangmanResetJsonResponse HangmanPickWordJsonResponse
