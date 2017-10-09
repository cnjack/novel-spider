package qiniu

import (
	"git.oschina.net/cnjack/novel-spider/config"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/satori/go.uuid"
)

var m *storage.BucketManager

func newBucket() (*storage.BucketManager, error) {
	mac := qbox.NewMac(config.GetHttpConfig().AccessKey, config.GetHttpConfig().SecretKey)
	zone, err := storage.GetZone(config.GetHttpConfig().AccessKey, config.GetHttpConfig().BucketName)
	if err != nil {
		return nil, err
	}
	cfg := &storage.Config{
		Zone:          zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	return storage.NewBucketManager(mac, cfg), nil
}

func UploadFromUrl(fetchURL string) (string, error) {
	var err error
	if m == nil {
		m, err = newBucket()
		if err != nil {
			return "", err
		}
	}

	keyFetch := "spider/" + uuid.NewV4().String()

	_, err = m.Fetch(fetchURL, config.GetHttpConfig().BucketName, keyFetch)
	if err != nil {
		return "", err
	}

	return config.GetHttpConfig().BucketUrl + keyFetch, nil
}
