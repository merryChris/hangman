package hangman

import (
	"io"
	"net/http"

	"github.com/merryChris/gameLord/api"
	"github.com/merryChris/gameLord/types"
	"github.com/spf13/viper"
)

type HangmanHandler struct {
	api.Handler
	hangmanManager *HangmanManager
}

var (
	V1HangmanRoutes = types.Routes{
		types.Route{
			"POST",
			"ListCategory",
			"/v1/hangman/list_category",
			RoutingHangmanHandler.ListCategory,
		},
		types.Route{
			"POST",
			"PickWord",
			"/v1/hangman/pick_word",
			RoutingHangmanHandler.PickWord,
		},
		types.Route{
			"POST",
			"LoadStatus",
			"/v1/hangman/load_status",
			RoutingHangmanHandler.LoadStatus,
		},
		types.Route{
			"POST",
			"Validate",
			"/v1/hangman/validate",
			RoutingHangmanHandler.Validate,
		},
		types.Route{
			"POST",
			"Reset",
			"/v1/hangman/reset",
			RoutingHangmanHandler.PickWord, // Same with `PickWord`
		},
	}
	RoutingHangmanHandler = &HangmanHandler{}
)

func (this *HangmanHandler) Init(redisConf *viper.Viper) error {
	hm, err := NewHangmanManager(redisConf)
	if err != nil {
		return err
	}
	this.hangmanManager = hm
	this.Initialized = true
	return nil
}

func (this *HangmanHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	if resp, ok := this.SelfCheck(); !ok {
		io.WriteString(w, resp)
		return
	}

	var requestObject HangmanListCategoryJsonRequest
	if resp, ok := api.ParseRequestJsonData(r.Body, &requestObject); !ok {
		io.WriteString(w, resp)
		return
	}
	if resp, ok := api.CheckRequestJsonData(requestObject); !ok {
		io.WriteString(w, resp)
		return
	}
	resp, _ := this.hangmanManager.ListCategory(requestObject)
	io.WriteString(w, resp)
}

func (this *HangmanHandler) PickWord(w http.ResponseWriter, r *http.Request) {
	if resp, ok := this.SelfCheck(); !ok {
		io.WriteString(w, resp)
		return
	}

	var requestObject HangmanPickWordJsonRequest
	if resp, ok := api.ParseRequestJsonData(r.Body, &requestObject); !ok {
		io.WriteString(w, resp)
		return
	}
	if resp, ok := api.CheckRequestJsonData(requestObject); !ok {
		io.WriteString(w, resp)
		return
	}
	resp, _ := this.hangmanManager.PickWord(requestObject)
	io.WriteString(w, resp)
}

func (this *HangmanHandler) LoadStatus(w http.ResponseWriter, r *http.Request) {
	if resp, ok := this.SelfCheck(); !ok {
		io.WriteString(w, resp)
		return
	}

	var requestObject HangmanLoadStatusJsonRequest
	if resp, ok := api.ParseRequestJsonData(r.Body, &requestObject); !ok {
		io.WriteString(w, resp)
		return
	}
	if resp, ok := api.CheckRequestJsonData(requestObject); !ok {
		io.WriteString(w, resp)
		return
	}
	resp, _ := this.hangmanManager.LoadStatus(requestObject)
	io.WriteString(w, resp)
}

func (this *HangmanHandler) Validate(w http.ResponseWriter, r *http.Request) {
	if resp, ok := this.SelfCheck(); !ok {
		io.WriteString(w, resp)
		return
	}

	var requestObject HangmanValidateJsonRequest
	if resp, ok := api.ParseRequestJsonData(r.Body, &requestObject); !ok {
		io.WriteString(w, resp)
		return
	}
	if resp, ok := api.CheckRequestJsonData(requestObject); !ok {
		io.WriteString(w, resp)
		return
	}
	resp, _ := this.hangmanManager.Validate(requestObject)
	io.WriteString(w, resp)
}

func (this *HangmanHandler) Close() {
	if this.Initialized {
		this.hangmanManager.Close()
		this.Initialized = false
	}
}
