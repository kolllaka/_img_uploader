package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kolllaka/_img_uploader/internal/model"
)

type imgStore struct {
	cfg *model.Config
	db  *sql.DB
}

func NewImgStore(cfg *model.Config, db *sql.DB) *imgStore {
	return &imgStore{
		cfg: cfg,
		db:  db,
	}
}

func (s *imgStore) GetAllImgPath() ([]model.Image, error) {
	const op = "storage.sqlite.GetAllImgPath"

	stmt, err := s.db.Prepare("SELECT id, path, created_at, expires_at FROM image WHERE expires_at > ?")
	if err != nil {
		return []model.Image{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var images []model.Image

	rows, err := stmt.Query(time.Now())
	if err != nil {
		return []model.Image{}, fmt.Errorf("%s: query statement: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var image model.Image
		if err := rows.Scan(&image.ID, &image.Path, &image.CreatedAt, &image.ExpiresAt); err != nil {
			return []model.Image{}, fmt.Errorf("%s: scan row: %w", op, err)
		}

		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return []model.Image{}, fmt.Errorf("%s: iterate rows: %w", op, err)
	}

	return images, nil
}

func (s *imgStore) GetImgPath(id string) (model.Image, error) {
	const op = "storage.sqlite.GetImgPath"

	stmt, err := s.db.Prepare("SELECT id, path, created_at, expires_at FROM image WHERE id = ?")
	if err != nil {
		return model.Image{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var image model.Image

	if err := stmt.QueryRow(id).Scan(&image.ID, &image.Path, &image.CreatedAt, &image.ExpiresAt); err != nil {
		return model.Image{}, fmt.Errorf("%s: query row statement: %w", op, err)
	}

	return image, nil
}

func (s *imgStore) SaveImgPath(image model.Image) (model.Image, error) {
	const op = "storage.sqlite.SaveImgPath"

	stmt, err := s.db.Prepare("INSERT INTO image (path, created_at, expires_at) VALUES (?, ?, ?)")
	if err != nil {
		return model.Image{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	res, err := stmt.Exec(image.Path, time.Now(), time.Now().Add(s.cfg.Images.DelayExpires*time.Hour))
	if err != nil {

		return model.Image{}, fmt.Errorf("%s: exec statement: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Image{}, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	image.ID = fmt.Sprint(id)

	return image, nil
}

func (s *imgStore) DeleteImgPath(id string) (model.Image, error) {
	const op = "storage.sqlite.DeleteImgPath"

	image, err := s.GetImgPath(id)
	if err != nil {
		return model.Image{}, fmt.Errorf("%s: get img path: %w", op, err)
	}

	stmt, err := s.db.Prepare("DELETE FROM image WHERE id = ?")
	if err != nil {
		return model.Image{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	if _, err := stmt.Exec(id); err != nil {
		return model.Image{}, fmt.Errorf("%s: query row statement: %w", op, err)
	}

	return image, nil
}

func (s *imgStore) GetListExpiredFiles() ([]model.Image, error) {
	const op = "storage.sqlite.GetListExpiredFiles"

	stmt, err := s.db.Prepare("SELECT id, path FROM image WHERE expires_at < ?")
	if err != nil {
		return []model.Image{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var list []model.Image

	rows, err := stmt.Query(time.Now())
	if err != nil {
		return []model.Image{}, fmt.Errorf("%s: query statement: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var img model.Image
		if err := rows.Scan(&img.ID, &img.Path); err != nil {
			return []model.Image{}, fmt.Errorf("%s: scan row: %w", op, err)
		}

		list = append(list, img)
	}

	if err := rows.Err(); err != nil {
		return []model.Image{}, fmt.Errorf("%s: iterate rows: %w", op, err)
	}

	return list, nil
}
