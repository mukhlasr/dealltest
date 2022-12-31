package psqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"simpleblog/posts"
)

type PostStorage sql.DB

func (s *PostStorage) GetAll(ctx context.Context) ([]posts.StoredPost, error) {
	db := (*sql.DB)(s)

	var res []posts.StoredPost
	rows, err := db.QueryContext(ctx, "select id, title, content, timestamp from posts")
	if err != nil {
		return res, fmt.Errorf("failed to query data: %w", err)
	}

	for rows.Next() {
		var post posts.StoredPost
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		res = append(res, post)
	}

	return res, rows.Err()
}

func (s *PostStorage) GetByID(ctx context.Context, id int64) (posts.StoredPost, error) {
	db := (*sql.DB)(s)

	var res posts.StoredPost

	err := db.QueryRowContext(ctx, "select id, title, content, timestamp from posts where id = $1", id).
		Scan(&res.ID, &res.Title, &res.Content, &res.Timestamp)

	if errors.Is(err, sql.ErrNoRows) {
		return res, posts.ErrNotFound
	}

	return res, err
}

func (s *PostStorage) Update(ctx context.Context, p posts.StoredPost) error {
	db := (*sql.DB)(s)

	res, err := db.ExecContext(ctx, "update posts set title=$1, content=$2, timestamp=$3 where id = $4",
		p.Title, p.Content, p.Timestamp, p.ID)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get RowsAffected: %w", err)
	}

	if affected == 0 {
		return posts.ErrNotFound
	}

	return nil
}

func (s *PostStorage) DeletePostByID(ctx context.Context, id int64) error {
	db := (*sql.DB)(s)

	res, err := db.ExecContext(ctx, "delete from posts where id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get RowsAffected: %w", err)
	}

	if affected == 0 {
		return posts.ErrNotFound
	}

	return nil
}

func (s *PostStorage) Create(ctx context.Context, p posts.StoredPost) (int64, error) {
	db := (*sql.DB)(s)

	var id int64
	err := db.QueryRowContext(ctx,
		"insert into posts(title, content, timestamp) values($1, $2, $3) returning id",
		p.Title, p.Content, p.Timestamp).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("failed to insert: %w", err)
	}

	return id, nil
}
