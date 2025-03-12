package urls

import (
	"errors"
	"log/slog"
	"net/http"
	storage "seelochka/internal/storages"

	"github.com/go-chi/chi/v5"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

// @Summary	Redirect by an alias
// @Accept		json
// @Produce	json
// @Param		alias	path	string	true	"Alias"
// @Success	301
// @Failure	404
// @Router		/{alias} [get]
func NewURLRedirect(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log = log.With(slog.String("handler", "urlRedirect"))

		alias := chi.URLParam(r, "alias")
		longURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			log.Error("url redirect error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, longURL, http.StatusFound)
	}

	return http.HandlerFunc(fn)
}
