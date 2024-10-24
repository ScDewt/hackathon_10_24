package refwalker

import "github.com/jackc/pgx/v4"

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

func GetNews() error {
	return nil
}

func WriteToDB() error {
	return nil
}
