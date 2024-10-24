package source_manager

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"net/http"
	"time"
)

type SourceManagerHandler struct {
	db *pgx.Conn
}

type Source struct {
	Token     string    `json:"token"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}
type List struct {
	Token        string   `json:"token"`
	Description  string   `json:"description"`
	ScheduleTime string   `json:"schedule_time"`
	Sources      []Source `json:"sources,omitempty"`
}

func NewSourceManagerHandler(conn *pgx.Conn) *SourceManagerHandler {
	return &SourceManagerHandler{conn}
}

func (h *SourceManagerHandler) RegisterHandler(r gin.IRoutes) {
	r.GET("/api/list/:token", h.apiGetList)
	r.POST("/api/list/:token", h.apiUpdateList)
	r.POST("/api/list/:token/source", h.apiAddSource)
	r.DELETE("/api/list/:token/source", h.apiDeleteSource)
}

func (h *SourceManagerHandler) apiGetList(r *gin.Context) {
	token := r.Param("token")
	if token == "" {
		r.AbortWithStatus(http.StatusNotFound)
		return
	}

	list, err := h.getList(r.Request.Context(), token)
	if err != nil {
		handleError(r, err)
		return
	}

	if list == nil {
		r.AbortWithStatus(http.StatusNotFound)
		return
	}

	list.Sources, err = h.getSources(r.Request.Context(), token)
	if err != nil {
		handleError(r, err)
		return
	}

	r.JSON(http.StatusOK, list)
}

func (h *SourceManagerHandler) apiUpdateList(r *gin.Context) {
	token := r.Param("token")
	if token == "" {
		r.AbortWithStatus(http.StatusNotFound)
		return
	}

	list, err := h.getOrCreateList(r.Request.Context(), token)
	if err != nil {
		handleError(r, err)
		return
	}

	if description, exist := r.GetQuery("description"); exist {
		err = h.updateListDescription(r.Request.Context(), token, description)
		if err != nil {
			handleError(r, err)
			return
		}
		list.Description = description
	}
	if scheduleTime, exist := r.GetQuery("schedule_time"); exist {
		err = h.updateListScheduleTime(r.Request.Context(), token, scheduleTime)
		if err != nil {
			handleError(r, err)
			return
		}
		list.ScheduleTime = scheduleTime
	}

	r.JSON(http.StatusOK, list)
}

func (h *SourceManagerHandler) apiAddSource(r *gin.Context) {
	token := r.Param("token")
	if token == "" {
		r.AbortWithStatus(http.StatusNotFound)
		return
	}

	source := r.Query("source")
	if source == "" {
		handleBadRequest(r, "source is empty")
		return
	}

	_, err := h.getOrCreateList(r.Request.Context(), token)
	if err != nil {
		handleError(r, err)
		return
	}

	s, err := h.addSourceToList(r.Request.Context(), token, source)
	if err != nil {
		handleError(r, err)
		return
	}

	r.JSON(http.StatusOK, s)
}

func (h *SourceManagerHandler) apiDeleteSource(r *gin.Context) {
	token := r.Param("token")
	if token == "" {
		r.AbortWithStatus(http.StatusNotFound)
		return
	}

	source := r.Query("source")
	if source == "" {
		handleBadRequest(r, "source is empty")
		return
	}

	err := h.removeSourceFromList(r.Request.Context(), token, source)
	if err != nil {
		handleError(r, err)
		return
	}

	r.JSON(http.StatusOK, nil)
}

func (h *SourceManagerHandler) getList(ctx context.Context, token string) (list *List, err error) {
	list = &List{Token: token}
	err = h.db.QueryRow(ctx, "SELECT token, description, schedule_list FROM lists WHERE token = $1", list.Token).
		Scan(&list.Token, &list.Description, &list.ScheduleTime)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return
}

func (h *SourceManagerHandler) getSources(ctx context.Context, token string) ([]Source, error) {
	sources := make([]Source, 0)
	rows, err := h.db.Query(ctx, "SELECT token, source, created_at FROM sources WHERE token = $1", token)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	for rows.Next() {
		s := Source{}
		err = rows.Scan(&s.Token, &s.Source, &s.CreatedAt)
		if err != nil {
			return nil, err
		}

		sources = append(sources, s)
	}
	return sources, nil
}
func (h *SourceManagerHandler) getOrCreateList(ctx context.Context, token string) (list *List, err error) {
	list, err = h.getList(ctx, token)
	if err != nil {
		return nil, err
	}

	if list == nil {
		list = &List{Token: token}
		_, err = h.db.Exec(ctx, "INSERT INTO lists(token) VALUES ($1)", token)
	}
	return
}

func (h *SourceManagerHandler) addSourceToList(ctx context.Context, token string, source string) (s *Source, err error) {
	s = &Source{
		Token:  token,
		Source: source,
	}
	err = h.db.QueryRow(ctx, "INSERT INTO sources(token, source) VALUES ($1, $2) RETURNING created_at", s.Token, s.Source).
		Scan(&s.CreatedAt)

	return
}

func (h *SourceManagerHandler) removeSourceFromList(ctx context.Context, token string, source string) (err error) {
	_, err = h.db.Exec(ctx, "DELETE FROM sources WHERE token = $1 AND source = $2", token, source)

	return
}

func (h *SourceManagerHandler) updateListDescription(ctx context.Context, token string, description string) (err error) {
	_, err = h.db.Exec(ctx, "UPDATE lists SET description = $1 WHERE token = $2", description, token)

	return
}

func (h *SourceManagerHandler) updateListScheduleTime(ctx context.Context, token string, time string) (err error) {
	_, err = h.db.Exec(ctx, "UPDATE lists SET schedule_time = $1 WHERE token = $2", time, token)

	return
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details any    `json:"details,omitempty"`
}

func handleError(r *gin.Context, err error) {
	r.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   err.Error(),
		Details: err,
	})
}
func handleBadRequest(r *gin.Context, msg string) {
	r.JSON(http.StatusBadRequest, ErrorResponse{Error: msg})
}
