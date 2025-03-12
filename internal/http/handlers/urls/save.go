package urls

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"seelochka/internal/pkg/utils"
	storage "seelochka/internal/storages"

	"github.com/go-playground/validator/v10"
)

// @Description	Alias data for creation
type AliasRequest struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// @Description	Response data for alias creation
type AliasResponse struct {
	Status string `json:"status" validate:"required" enums:"ok,error"`
	Alias  string `json:"alias,omitempty"`
	Error  string `json:"error,omitempty"`
}

type URLSaver interface {
	SaveURL(longURL, alias string) error
}

// @Summary	Create an alias
// @Accept		json
// @Produce	json
// @Param		request	body	urls.AliasRequest	true	"reqbody"
// @Success	200		{object}	urls.AliasResponse
// @Failure	400
// @Router		/ [post]
func NewURLSave(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("handler", "urlSave"))

		var aliasreq AliasRequest
		jsonEncoder := json.NewEncoder(w)

		err := json.NewDecoder(r.Body).Decode(&aliasreq)
		if err != nil {
			log.Error("url save error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := validator.New().Struct(aliasreq); err != nil {
			err = err.(validator.ValidationErrors)
			log.Error("invalid urlsave request", slog.String("error", err.Error()))

			w.WriteHeader(http.StatusBadRequest)
			jsonEncoder.Encode(AliasResponse{Status: statusError, Error: err.Error()})
			return
		}

		alias := aliasreq.Alias
		if alias == "" {
			alias = utils.GenerateRandomString(aliasDefSize)
		}

		err = urlSaver.SaveURL(aliasreq.URL, alias)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorText := "error"
			if errors.Is(err, storage.ErrAliasUsed) {
				errorText = err.Error()
			}

			jsonEncoder.Encode(AliasResponse{Status: statusError, Error: errorText})
			return
		}

		jsonEncoder.Encode(AliasResponse{Status: statusOK, Alias: alias})
	}
	return http.HandlerFunc(fn)
}
