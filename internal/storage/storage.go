package storage

import "github.com/kolllaka/_img_uploader/internal/model"

type Storage interface {
	GetAllImgPath() ([]model.Image, error)
	GetImgPath(id string) (model.Image, error)
	SaveImgPath(image model.Image) (model.Image, error)
	DeleteImgPath(id string) (model.Image, error)

	GetListExpiredFiles() ([]model.Image, error)
}
