package save

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/zzvanq/seelochka/internal/http/handlers/url/consts"
	"github.com/zzvanq/seelochka/internal/pkg/random"
	"github.com/zzvanq/seelochka/internal/storage"

	"github.com/go-playground/validator/v10"
)

//	@Description	Alias data for creation
type AliasRequest struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty" validate:"min=5,max=16"`
}

//	@Description	Response data for alias creation
type AliasResponse struct {
	Alias string `json:"alias,omitempty"`
	Error string `json:"error,omitempty"`
}

type URLSaver interface {
	SaveURL(longURL, alias string) error
}

//	@Summary	Create an alias
//	@Accept		json
//	@Produce	json
//	@Param		request	body		save.AliasRequest	true	"Alias"
//	@Success	200		{object}	save.AliasResponse
//	@Failure	400
//	@Router		/ [post]
func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("handler", "urlSave"))

		var data AliasRequest
		respEnc := json.NewEncoder(w)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			log.Error("url save error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := validator.New().Struct(data); err != nil {
			validErr := err.(validator.ValidationErrors)[0]
			log.Error("invalid urlsave request", slog.String("error", validErr.Error()))

			validationText := fmt.Sprintf(
				"field '%s' wrong '%s%s'", validErr.Field(), validErr.ActualTag(), validErr.Param())

			w.WriteHeader(http.StatusBadRequest)
			respEnc.Encode(AliasResponse{Error: validationText})
			return
		}

		alias := data.Alias
		if alias == "" {
			alias = random.GenerateRandomString(consts.AliasDefSize)
		}

		err = urlSaver.SaveURL(data.URL, alias)
		if err != nil {
			errorText := "error"
			if errors.Is(err, storage.ErrAliasUsed) {
				errorText = err.Error()
			}

			w.WriteHeader(http.StatusBadRequest)
			respEnc.Encode(AliasResponse{Error: errorText})
			return
		}

		w.WriteHeader(http.StatusCreated)
		respEnc.Encode(AliasResponse{Alias: alias})
	}
	return http.HandlerFunc(fn)
}
