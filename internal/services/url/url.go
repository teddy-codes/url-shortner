package url

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

type Store interface {
	CreateURL(context.Context, string, *url.URL) error
	GetURL(ctx context.Context, id string) (string, error)
}

type Service struct {
	Store Store
}

func (s *Service) generateIDHash() string {
	return "thisisrandom"
}

func (s *Service) CreateURL(ctx context.Context, link, id string) (*url.URL, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse url")
	}

	if id == "" {
		id = s.generateIDHash()
	}

	if err = s.Store.CreateURL(ctx, id, u); err != nil {
		return nil, errors.Wrap(err, "cannot save url to database")
	}

	out := &url.URL{
		Scheme: "https",
		Host:   "somewhere-over-the.rainbow",
		Path:   id,
	}

	return out, nil
}

func (s *Service) CheckURLExists(ctx context.Context, id string) (bool, error) {
	out, err := s.Store.GetURL(ctx, id)
	if err != nil {
		return false, err
	}
	return out != "", nil
}

func (s *Service) GetURLRedirect(ctx context.Context, id string) (string, error) {
	out, err := s.Store.GetURL(ctx, id)
	if err != nil {
		return "", err
	}
	return out, nil
}
