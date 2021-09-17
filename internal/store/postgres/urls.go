package postgres

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

func (s *Store) CreateURLTable(ctx context.Context) error {
	table := `CREATE TABLE IF NOT EXISTS urls (
		id TEXT NOT NULL UNIQUE,
		link TEXT NOT NULL
	)`
	_, err := s.DB.ExecContext(ctx, table)
	return err
}

func (s *Store) CreateURL(ctx context.Context, id string, u *url.URL) error {
	stmt, err := s.PrepareContext(ctx, "INSERT INTO urls (id, link) VALUES ($1, $2)")
	if err != nil {
		return errors.Wrap(err, "cannot prepare statement")
	}

	_, err = stmt.ExecContext(ctx, id, u.String())
	if err != nil {
		return errors.Wrapf(err, "cannot execute sql statement with values: id:%v, url:%v", id, u.String())
	}

	return nil
}

func (s *Store) GetURL(ctx context.Context, id string) (string, error) {
	stmt, err := s.PrepareContext(ctx, "SELECT link FROM urls WHERE id = $1")
	if err != nil {
		return "", errors.Wrap(err, "cannot prepare select query")
	}

	r, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return "", errors.Wrap(err, "cannot query for id")
	}

	if !r.Next() {
		return "", nil
	}

	var link string
	if err := r.Scan(&link); err != nil {
		return "", errors.Wrap(err, "cannot read link from scan")
	}

	return link, nil
}
