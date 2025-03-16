package service

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kolllaka/_img_uploader/internal/img"
	"github.com/kolllaka/_img_uploader/internal/model"
	"github.com/kolllaka/_img_uploader/internal/storage"
)

type Service interface {
	SaveImage(image io.Reader) (model.Image, error)
	GetAllImages() ([]model.Image, error)
	GetImage(id string) (model.Image, error)
	DeleteImage(id string) error

	DeleteExpiresFiles(chan error) int
}

type service struct {
	cfg       *model.Config
	storage   storage.Storage
	transform img.Transform
}

// GetAllImages implements Service.
func (s *service) GetAllImages() ([]model.Image, error) {
	return s.storage.GetAllImgPath()
}

func New(cfg *model.Config, storage storage.Storage, transform img.Transform) Service {
	return &service{
		cfg:       cfg,
		storage:   storage,
		transform: transform,
	}
}

func (s *service) SaveImage(image io.Reader) (model.Image, error) {
	m, err := s.transform.Resize(image, s.cfg.Images.Width, s.cfg.Images.Height)
	if err != nil {
		return model.Image{}, err
	}

	filePath := fmt.Sprintf("%s/%s", s.cfg.Images.SaveFolder, imageName())
	if err := s.transform.SavePng(filePath, m); err != nil {
		return model.Image{}, err
	}

	return s.storage.SaveImgPath(model.Image{
		Path: fmt.Sprintf("/%s.png", filePath),
	})
}
func (s *service) GetImage(id string) (model.Image, error) {
	return s.storage.GetImgPath(id)
}
func (s *service) DeleteImage(id string) error {
	image, err := s.storage.DeleteImgPath(id)
	if err != nil {
		return err
	}

	return os.Remove("." + image.Path)
}

func (s *service) DeleteExpiresFiles(chanErr chan error) int {
	const op = "service.DeleteExpiresFiles"
	count := 0

	list, err := s.storage.GetListExpiredFiles()
	if err != nil {
		chanErr <- fmt.Errorf("%s: get list expired files: %w", op, err)

		return 0
	}

	for _, img := range list {
		if err := os.Remove("." + img.Path); err != nil {
			chanErr <- fmt.Errorf("%s: remove file: %w", op, err)

			continue
		}

		if _, err := s.storage.DeleteImgPath(img.ID); err != nil {
			chanErr <- fmt.Errorf("%s: delete from db: %w", op, err)

			continue
		}

		count++
	}

	return count
}

func imageName() string {
	return time.Now().UTC().Format("20060102T150405Z")
}
