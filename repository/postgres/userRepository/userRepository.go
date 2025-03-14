package userRepository

import (
	"context"
	"database/sql"

	"github.com/mohammaderm/rootext/entity"
	"github.com/mohammaderm/rootext/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	// user
	IsUserUnique(ctx context.Context, username string) (bool, error)
	IsUserExistsById(ctx context.Context, id uint) (entity.User, bool, error)
	Register(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, bool, error)
}

type Repository struct {
	conn *postgres.PostgresDB
}

func New(conn *postgres.PostgresDB) UserRepo {
	return &Repository{
		conn: conn,
	}
}

func (r Repository) IsUserExistsById(ctx context.Context, id uint) (entity.User, bool, error) {
	row := r.conn.Conn().QueryRowxContext(ctx, "select * from users where id = $1", id)
	user, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, err
	}
	return user, true, nil
}

func (r Repository) GetUserByUsername(ctx context.Context, username string) (entity.User, bool, error) {
	row := r.conn.Conn().QueryRowxContext(ctx, "select * from users where username = $1;", username)
	user, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}
		return entity.User{}, false, err
	}
	return user, true, nil
}

func (r Repository) IsUserUnique(ctx context.Context, username string) (bool, error) {
	row := r.conn.Conn().QueryRowxContext(ctx, "select * from users where username = $1;", username)
	_, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func (r Repository) Register(ctx context.Context, user entity.User) (entity.User, error) {

	result, err := r.conn.Conn().ExecContext(ctx, "insert into users (username, password) values ($1,$2);", user.Username, user.Password)
	if err != nil {
		return entity.User{}, err
	}
	id, _ := result.LastInsertId()
	user.ID = uint(id)
	return user, nil
}

func scanUser(scanner sqlx.ColScanner) (entity.User, error) {
	var user entity.User
	err := scanner.Scan(&user.ID, &user.Username, &user.Password, &user.Created_at)
	return user, err
}
