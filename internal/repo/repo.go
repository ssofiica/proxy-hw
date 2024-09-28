package repo

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ssofiica/proxy-hw/internal/proxy/utils"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) Repo {
	return Repo{db: db}
}

// reponses
func (r *Repo) SaveResponse(ctx context.Context, response []byte, id int) error {
	query := `insert into response(request_id, data) values ($1, $2)`
	_, err := r.db.Exec(ctx, query, id, response)
	if err != nil {
		return err
	}
	return nil
}

// requestes
func (r *Repo) SaveRequest(ctx context.Context, request []byte) (int, error) {
	query := `insert into request(data) values ($1) returning id`
	var id int
	err := r.db.QueryRow(ctx, query, request).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repo) GetRequestList(ctx context.Context) ([]utils.RequestInfo, error) {
	query := `select data from request`
	var res []utils.RequestInfo
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []utils.RequestInfo{}, err
	}
	for rows.Next() {
		var req []byte
		err := rows.Scan(&req)
		if err != nil {
			return []utils.RequestInfo{}, err
		}
		var r utils.RequestInfo
		err = json.Unmarshal(req, &r)
		if err != nil {
			return []utils.RequestInfo{}, err
		}
		res = append(res, r)
	}
	return res, nil
}

func (r *Repo) GetRequestByID(ctx context.Context, id int) ([]byte, error) {
	query := `select data from request where id=$1`
	var res []byte
	err := r.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return []byte{}, err
	}
	return res, nil
}
