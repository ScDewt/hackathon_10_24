package urlwalker

import (
	"context"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v4"
)

type News struct {
	ID   string
	Url  string
	News string
}

type UrlsWalker struct {
	db *pgx.Conn
}

func NewUrlWalker(db *pgx.Conn) *UrlsWalker {
	return &UrlsWalker{db}
}

func (r *UrlsWalker) parseNews(ref string) (string, error) {
	newsList := make([]string, 0)
	// Выполняем HTTP запрос
	res, err := http.Get(ref)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", err
	}

	// Загружаем HTML документ в goquery
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	// Парсим и выводим текст сообщений
	doc.Find(".tgme_widget_message_text").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		newsList = append(newsList, text)
	})

	return strings.Join(newsList, ";"), nil
}

func (r *UrlsWalker) writeToDB(url, news string) error {
	rowNews := News{
		Url:  url,
		News: news,
	}
	query := `
		INSERT INTO news (url, news) 
		VALUES ($1, $2) 
		RETURNING id
	`
	err := r.db.QueryRow(context.Background(), query, rowNews.Url, rowNews.News).Scan(&rowNews.ID)
	return err
}

func (r *UrlsWalker) Walk(urlNews []string) error {
	for _, un := range urlNews {
		news, err := r.parseNews(un)
		if err != nil {
			return err
		}

		err = r.writeToDB(un, news)
		if err != nil {
			return err
		}
	}

	return nil
}
