package postgres

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/marki-eriker/event-listener/db"
	"github.com/marki-eriker/event-listener/entity/user"
	"time"
)

type UserStore struct {
	storage     *pg.DB
	queryMetric func(success bool, duration time.Duration, query string)
}

func (us *UserStore) GetByID(ctx context.Context, id uint) (*user.User, error) {
	bt := time.Now()

	var u user.User

	err := us.storage.ModelContext(ctx, &u).Where("id = ?", id).Select()
	if err != nil {
		us.queryMetric(false, time.Since(bt), "get_user_by_id")
		return nil, err
	}

	us.queryMetric(true, time.Since(bt), "get_user_by_id")
	return &u, nil
}

func (us *UserStore) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	bt := time.Now()

	var u user.User

	err := us.storage.ModelContext(ctx, &u).Where("login = ?", login).Select()
	if err != nil {
		us.queryMetric(false, time.Since(bt), "get_user_by_login")
		return nil, err
	}

	us.queryMetric(true, time.Since(bt), "get_user_by_login")
	return &u, nil
}

func (us *UserStore) Insert(ctx context.Context, user *user.User) error {
	bt := time.Now()

	_, err := us.storage.ModelContext(ctx, user).Insert()
	if err != nil {
		pgErr, ok := err.(pg.Error)
		if ok && pgErr.IntegrityViolation() {
			us.queryMetric(false, time.Since(bt), "insert_user")
			return db.ErrDuplicateKey
		}

		us.queryMetric(false, time.Since(bt), "insert_user")
		return err
	}

	us.queryMetric(true, time.Since(bt), "insert_user")
	return nil
}

func (us *UserStore) UpdateVerify(ctx context.Context, id uint, verify bool) error {
	bt := time.Now()

	u := &user.User{
		ID:       id,
		Verified: verify,
	}

	_, err := us.storage.ModelContext(ctx, u).Column("verified").WherePK().Update()
	if err != nil {
		us.queryMetric(false, time.Since(bt), "verify_user")
		return err
	}

	us.queryMetric(true, time.Since(bt), "verify_user")
	return nil
}
