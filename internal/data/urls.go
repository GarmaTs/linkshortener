package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/GarmaTs/linkshortener/internal/validator"
	"github.com/itchyny/base58-go"
)

var (
	ErrDuplicateUrl = errors.New("duplicate url")
)

type Url struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ShortUrl  string    `json:"short_url"`
	FullUrl   string    `json:"full_url"`
}

func ValidateUrl(v *validator.Validator, url *Url) {
	v.Check(url.FullUrl != "", "full url", "must be provided")
	v.Check(len(url.FullUrl) <= 1000, "full url", "must be less than 1000 characters long")
	v.Check(strings.HasPrefix(url.FullUrl, "http"), "full url", "has wrong format")
}

type UrlModel struct {
	DB *sql.DB
}

func (m UrlModel) Insert(url *Url, userName, fullUrl string) error {
	shortUrl := generateShortLink(fullUrl, userName)
	url.ShortUrl = shortUrl

	query :=
		`insert into public.urls (user_id, short_url, full_url)
		select id, $1, $2 from public.users where name = $3
		returning id, created_at`

	args := []interface{}{shortUrl, fullUrl, userName}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&url.ID, &url.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "urls_user_id_full_url_key"`:
		case err.Error() == `pq: повторяющееся значение ключа нарушает ограничение уникальности "urls_user_id_full_url_key"`:
			err2 := m.getShortUrlByUserNameAndFullUrl(url, userName, fullUrl)
			if err2 != nil {
				return err
			} else {
				return nil
			}
		default:
			return err
		}
	}

	return nil
}

func generateShortLink(srcLink string, userName string) string {
	urlHashBytes := sha256Of(srcLink + userName)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return strings.ToLower(finalString[:8])
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		panic(err.Error())
	}
	return string(encoded)
}

func (m UrlModel) getShortUrlByUserNameAndFullUrl(url *Url, userName, fullUrl string) error {
	query :=
		`select id, short_url, created_at from public.urls
		where
			user_id in (select id from public.users where name = $1)
			and full_url = $2
		order by created_at desc
		limit 1`
	args := []interface{}{userName, fullUrl}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&url.ID, &url.ShortUrl, &url.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (m UrlModel) GetOne(url *Url, shortUrl string) error {
	query :=
		`select id, created_at, full_url from public.urls
		where short_url = $1`
	args := []interface{}{shortUrl}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&url.ID, &url.CreatedAt, &url.FullUrl)
	if err != nil {
		return err
	}

	return nil
}

func (m UrlModel) GetList(username string) ([]*Url, error) {
	query :=
		`select id, full_url, short_url, created_at from public.urls
		where
			user_id in (select id from public.users where name = $1)
		order by id desc`
	args := []interface{}{username}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := []*Url{}
	for rows.Next() {
		u := &Url{}
		err = rows.Scan(&u.ID, &u.FullUrl, &u.ShortUrl, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (m UrlModel) Delete(id int) error {
	query := `delete from public.urls where id = $1`

	args := []interface{}{id}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}
