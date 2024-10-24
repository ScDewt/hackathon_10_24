package refwalker

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v4"
)

type News struct {
	ID   string
	Ref  string
	News string
}

type RefWalk struct {
	db       *pgx.Conn
	refsNews map[string]string
}

func NewRefWalk(db *pgx.Conn) *RefWalk {
	refs := make(map[string]string)

	return &RefWalk{db, refs}
}

func (r *RefWalk) AddRef(refs []string) error {
	for _, ref := range refs {
		if _, ok := r.refsNews[ref]; !ok {
			r.refsNews[ref] = ""
		}
	}

	return nil
}

func (r *RefWalk) getNews(ref string) (string, error) {
	// Запрос на получение HTML страницы
	resp, err := http.Get(ref)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Парсинг HTML с использованием goquery
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing HTML: %v", err)
	}

	var newsList []string

	doc.Find(".news-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".news-title").Text()
		newsList = append(newsList, title)
	})

	return strings.Join(newsList, ";"), nil
}

func (r *RefWalk) inserNews(ref, news string) error {
	rowNews := News{
		Ref:  ref,
		News: news,
	}
	query := `
		INSERT INTO news (id, ref, news) 
		VALUES (uuid_generate_v4(), $1, $2) 
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, rowNews.Ref, rowNews.News).Scan(&rowNews.ID)
	return err
}

func (r *RefWalk) WriteToDB() error {
	for rn := range r.refsNews {
		news, err := r.getNews(rn)
		if err != nil {
			return err
		}

		err = r.inserNews(rn, news)
		if err != nil {
			return err
		}
	}

	return nil
}
