package dao

import (
	"bytes"
	"io/ioutil"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSSImageDAO struct {
	BucketName string
	Client     *oss.Client
}

func NewOSSImageDAO(bucketName string, client *oss.Client) *OSSImageDAO {
	return &OSSImageDAO{
		BucketName: bucketName,
		Client:     client,
	}
}

func (dao *OSSImageDAO) Save(path string, image []byte) error {
	bucket, err := dao.Client.Bucket(dao.BucketName)
	if err != nil {
		return err
	}
	err = bucket.PutObject(path, bytes.NewReader([]byte(image)))
	return err
}

func (dao *OSSImageDAO) Get(path string) ([]byte, error) {
	bucket, err := dao.Client.Bucket(dao.BucketName)
	if err != nil {
		return nil, err
	}
	body, err := bucket.GetObject(path)
	if err != nil {
		return nil, err
	}
	image, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return image, nil
}
