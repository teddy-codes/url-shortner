package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type handlersService interface {
	CreateURL(ctx context.Context, link, id string) (*url.URL, error)
	CheckURLExists(ctx context.Context, id string) (bool, error)
	GetURLRedirect(ctx context.Context, id string) (string, error)
}

type handlers struct {
	log     zerolog.Logger
	service handlersService
}

type createURLRequest struct {
	ID   string `json:"id" xml:"id"`
	Link string `json:"link" xml:"link"`
}

func (h *handlers) createURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createURLRequest
	switch r.Header.Get("content-type") {
	case "application/x-www-form-urlencoded":
		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusTeapot)
			_, _ = w.Write([]byte("we don't provide no bodies around here"))
			return
		}
		values, err := url.ParseQuery(string(b))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid url params in body"))
			return
		}
		req.Link = values.Get("link")
		req.ID = values.Get("id")
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid json was provided"))
			return
		}
	case "text/xml":
		if err := xml.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("invalid json was provided"))
			return
		}
	}

	existingURL, err := h.service.CheckURLExists(ctx, req.ID)
	if err != nil || existingURL {
		h.log.Info().Err(err).Str("id", req.ID).Str("link", req.Link).Msg("link id already exists")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte("id already exists"))
		return
	}

	out, err := h.service.CreateURL(ctx, req.Link, req.ID)
	if err != nil {
		h.log.Err(err).Str("id", req.ID).Str("link", req.Link).Msg("cannot create url")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("cannot save url"))
		return
	}

	_, _ = w.Write([]byte(out.String()))
}

func (h *handlers) getURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "path")
	log.Info().Msg(id)

	link, err := h.service.GetURLRedirect(r.Context(), id)
	if err != nil {
		h.log.Err(err).Str("path", id).Msg("cannot retrieve link")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	}
	if link == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link, http.StatusSeeOther)
}
