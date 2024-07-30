package repo

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rlapenok/messagio/internal/config"
	"github.com/rlapenok/messagio/internal/models"
	"github.com/rlapenok/messagio/internal/utils"
)

type Repository interface {
	Save(ctx context.Context, msg *models.MessageForRepo) error
	Update()
	Close() error
	GetSats(ctx context.Context) (*models.Response, error)
}

type repository struct {
	wg      *sync.WaitGroup
	db      *sqlx.DB
	channel chan uuid.UUID
}

func New(cfg config.DataBaseConfig, channel chan uuid.UUID) Repository {
	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Db)
	db, err := sqlx.Connect("postgres", url)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	utils.Logger.Info("Connet to postgres - OK")
	createTableMessages(db)
	repo := &repository{db: db, channel: channel, wg: &sync.WaitGroup{}}
	repo.Update()
	return repo

}

func (r repository) Save(ctx context.Context, msg *models.MessageForRepo) error {

	queryString := "INSERT INTO messages VALUES (:id, :message)"
	query := func(ctx context.Context, tx *sqlx.Tx) (sql.Result, error) {
		return tx.NamedExecContext(ctx, queryString, msg)
	}
	_, err := withTransaction(ctx, r, query)
	return err

}

func (r repository) Update() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		for msg := range r.channel {
			func() {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				queryString := "UPDATE messages SET processed=true WHERE id=$1"

				query := func(ctx context.Context, tx *sqlx.Tx) (sql.Result, error) {

					return tx.ExecContext(ctx, queryString, msg)
				}
				_, err := withTransaction(ctx, r, query)
				if err != nil {
					logMsg := fmt.Sprintf("Cannot update record with id:%v", msg)
					utils.Logger.Error(logMsg)

				}
			}()
		}
	}()
}

func (r repository) GetSats(ctx context.Context) (*models.Response, error) {

	queryString := "SELECT SUM(CASE WHEN processed = true THEN 1 ELSE 0 END) AS processed_count,SUM(CASE WHEN processed = false THEN 1 ELSE 0 END) AS notprocessed_count,COUNT(*) AS total_count FROM messages"
	query := func(ctx context.Context, tx *sqlx.Tx) (*models.Response, error) {
		var data models.Response
		err := tx.GetContext(ctx, &data, queryString)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}
	return withTransaction(ctx, r, query)
}
func (r repository) Close() error {
	r.wg.Wait()
	return r.db.Close()
}
