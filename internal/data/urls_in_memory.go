package data

type InMemUrlModel struct {
	urls map[string]Url
}

func (m InMemUrlModel) Insert(url *Url, userName, fullUrl string) error {
	m.urls[url.ShortUrl] = *url

	return nil
}

func (m InMemUrlModel) GetOne(url *Url, shortUrl string) error {
	u, exists := m.urls[shortUrl]
	if !exists {
		return ErrDoesNotExist
	}

	url.ID = u.ID
	url.CreatedAt = u.CreatedAt
	url.ShortUrl = shortUrl
	url.FullUrl = u.FullUrl
	return nil
}

// Stab method
func (m InMemUrlModel) GetList(username string) ([]*Url, error) {
	return nil, nil
}

// Stab method
func (m InMemUrlModel) Delete(id int) error {
	return nil
}
