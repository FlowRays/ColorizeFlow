package dao

import (
	"io/ioutil"
)

type LocImageDAO struct {
	StoragePath string
}

func NewLocImageDAO(storagePath string) *LocImageDAO {
	return &LocImageDAO{
		StoragePath: storagePath,
	}
}

func (dao *LocImageDAO) Save(path string, image []byte) error {
	filePath := dao.StoragePath + "/" + path
	err := ioutil.WriteFile(filePath, image, 0644)
	return err
}

func (dao *LocImageDAO) Get(path string) ([]byte, error) {
	filePath := dao.StoragePath + "/" + path
	image, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return image, nil
}
