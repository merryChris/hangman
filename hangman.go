package hangman

import (
	"encoding/json"
	"math/rand"

	"github.com/astaxie/beego/orm"
	"github.com/merryChris/gameLord/core"
	"github.com/merryChris/gameLord/types"
	"github.com/merryChris/gameLord/utils"
	"github.com/spf13/viper"
)

type HangmanManager struct {
	*core.Manager
}

var registeredHangmanOrm = false

func init() {
	orm.RegisterModel(new(Hangman))
	orm.RegisterModelWithPrefix("hangman_", new(Category), new(Dict))
}

func getCurrentWordByLetterStatus(word string, status int64) (string, bool) {
	cw := []byte(word)
	flag := true
	for i, x := range cw {
		if (1<<(x-'a'))&status == 0 {
			flag = false
			cw[i] = '*'
		}
	}
	return string(cw), flag
}

func NewHangmanManager(redisConf *viper.Viper) (*HangmanManager, error) {
	m, err := core.NewManager(redisConf)
	if err != nil {
		return nil, err
	}

	hm := &HangmanManager{Manager: m}
	hm.Initialized = true
	return hm, nil
}

func (this *HangmanManager) ListCategory(req HangmanListCategoryJsonRequest) (string, bool) {
	user := types.User{Id: req.UserId}
	if err := this.MysqlOrm.Read(&user); err != nil {
		return types.Error112(), false
	}
	if req.GameId != 1 {
		return types.Error103(), false
	}
	user.CurrentGameId = req.GameId
	user.CurrentDevice = req.DeviceName

	if ok := user.CheckGameToken(this.RedisClient, req.GameToken); !ok {
		return types.Error115(), false
	}

	categories := []*Category{}
	if _, err := this.MysqlOrm.QueryTable((*Category)(nil)).All(&categories, "Id", "Name"); err != nil {
		return types.Error101(), false
	}
	respBytes, _ := json.Marshal(&HangmanListCategoryJsonResponse{
		BaseJsonResponse: types.BaseJsonResponse{
			Code:        130,
			MessageType: "user_success"},
		Categories: categories,
	})
	return string(respBytes), true
}

func (this *HangmanManager) PickWord(req HangmanPickWordJsonRequest) (string, bool) {
	user := types.User{Id: req.UserId}
	if err := this.MysqlOrm.Read(&user); err != nil {
		return types.Error112(), false
	}
	if req.GameId != 1 {
		return types.Error103(), false
	}
	user.CurrentGameId = req.GameId
	user.CurrentDevice = req.DeviceName

	if ok := user.CheckGameToken(this.RedisClient, req.GameToken); !ok {
		return types.Error115(), false
	}

	dics := []*Dict{}
	if cnt, err := this.MysqlOrm.QueryTable((*Dict)(nil)).Filter("Category__Id", req.CategoryId).
		All(&dics, "Word"); err != nil || cnt == 0 {
		return types.Error101(), false
	}

	// Update or Insert `hangman` Table
	hangman := Hangman{UserId: req.UserId}
	if _, _, err := this.MysqlOrm.ReadOrCreate(&hangman, "UserId"); err != nil {
		return types.Error101(), false
	}
	hangman.Word = dics[rand.Intn(len(dics))].Word
	hangman.LetterStatus = 0
	if _, err := this.MysqlOrm.Update(&hangman); err != nil {
		return types.Error101(), false
	}
	// Update `game_status.status` Field
	if _, err := this.MysqlOrm.QueryTable((*types.GameStatus)(nil)).Filter("UserId", req.UserId).
		Filter("GameId", req.GameId).Update(orm.Params{"Status": 1}); err != nil {
		return types.Error101(), false
	}

	word, _ := getCurrentWordByLetterStatus(hangman.Word, hangman.LetterStatus)
	respBytes, _ := json.Marshal(&HangmanPickWordJsonResponse{
		BaseJsonResponse: types.BaseJsonResponse{
			Code:        130,
			MessageType: "user_success"},
		CurrentWord:  word,
		LetterStatus: hangman.LetterStatus,
	})
	return string(respBytes), true
}

func (this *HangmanManager) LoadStatus(req HangmanLoadStatusJsonRequest) (string, bool) {
	user := types.User{Id: req.UserId}
	if err := this.MysqlOrm.Read(&user); err != nil {
		return types.Error112(), false
	}
	if req.GameId != 1 {
		return types.Error103(), false
	}
	user.CurrentGameId = req.GameId
	user.CurrentDevice = req.DeviceName

	if ok := user.CheckGameToken(this.RedisClient, req.GameToken); !ok {
		return types.Error115(), false
	}

	hangman := Hangman{UserId: req.UserId}
	if err := this.MysqlOrm.Read(&hangman, "UserId"); err != nil {
		return types.Error131(), false
	}
	word, _ := getCurrentWordByLetterStatus(hangman.Word, hangman.LetterStatus)
	respBytes, _ := json.Marshal(&HangmanLoadStatusJsonResponse{
		BaseJsonResponse: types.BaseJsonResponse{
			Code:        130,
			MessageType: "user_success"},
		CurrentWord:  word,
		LetterStatus: hangman.LetterStatus,
	})
	return string(respBytes), true
}

func (this *HangmanManager) Validate(req HangmanValidateJsonRequest) (string, bool) {
	user := types.User{Id: req.UserId}
	if err := this.MysqlOrm.Read(&user); err != nil {
		return types.Error112(), false
	}
	if req.GameId != 1 {
		return types.Error103(), false
	}
	user.CurrentGameId = req.GameId
	user.CurrentDevice = req.DeviceName

	if ok := user.CheckGameToken(this.RedisClient, req.GameToken); !ok {
		return types.Error115(), false
	}

	hangman := Hangman{UserId: req.UserId}
	if err := this.MysqlOrm.Read(&hangman, "UserId"); err != nil {
		return types.Error131(), false
	}

	// Unified with Lowercase Letters
	if 'A' <= req.CurrentLetter && req.CurrentLetter <= 'Z' {
		req.CurrentLetter ^= 2
	}
	preWord, _ := getCurrentWordByLetterStatus(hangman.Word, hangman.LetterStatus)
	hangman.LetterStatus |= 1 << (req.CurrentLetter - 'a')
	curWord, ending := getCurrentWordByLetterStatus(hangman.Word, hangman.LetterStatus)
	cntOnes := utils.CountOneBits(hangman.LetterStatus)

	// Update `hangman.letter_status` Field
	if _, err := this.MysqlOrm.Update(&hangman); err != nil {
		return types.Error101(), false
	}
	// Update 'game_status.status' Field
	if cntOnes >= 7 {
		if _, err := this.MysqlOrm.QueryTable((*types.GameStatus)(nil)).Filter("UserId", req.UserId).
			Filter("GameId", req.GameId).Update(orm.Params{"Status": 0}); err != nil {
			return types.Error101(), false
		}
	}
	respBytes, _ := json.Marshal(&HangmanValidateJsonResponse{
		BaseJsonResponse: types.BaseJsonResponse{
			Code:        130,
			MessageType: "user_success"},
		TrueOrFalse:  preWord != curWord,
		Ending:       cntOnes >= 7 || ending,
		CurrentWord:  curWord,
		LetterStatus: hangman.LetterStatus,
	})
	return string(respBytes), true
}

func (this *HangmanManager) Reset(req HangmanResetJsonRequest) (string, bool) {
	return "", true
}

func (this *HangmanManager) Close() {
	if this.Initialized {
		this.Initialized = false
	}
}
