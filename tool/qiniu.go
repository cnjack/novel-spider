package tool

import (
	"errors"
	"git.oschina.net/cnjack/novel-spider/config"
	"github.com/satori/go.uuid"
	"qiniupkg.com/api.v7/kodo"
)

var bucket *kodo.Bucket

func newBucket() (*kodo.Bucket, error) {
	ak := config.GetHttpConfig().AccessKey
	sk := config.GetHttpConfig().SecretKey
	b := config.GetHttpConfig().BucketName
	if ak == "" || sk == "" || b == "" {
		return nil, errors.New("qiniu config error")
	}
	client := kodo.NewWithoutZone(&kodo.Config{
		AccessKey: config.GetHttpConfig().AccessKey,
		SecretKey: config.GetHttpConfig().SecretKey,
	})
	nBucket := client.Bucket(b)
	return &nBucket, nil
}

func UploadFromUrl(fetchURL string) (string, error) {
	var err error
	if bucket == nil {
		bucket, err = newBucket()
		if err != nil {
			return "", err
		}
	}
	keyFetch := uuid.NewV4().String()
	err = bucket.Fetch(nil, keyFetch, fetchURL)
	if err != nil {
		return "", err
	}
	return config.GetHttpConfig().BucketUrl + keyFetch, nil
}
