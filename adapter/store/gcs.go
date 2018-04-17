package store

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/store"
	"google.golang.org/api/iterator"
)

var (
	urlFormat       = "https://storage.googleapis.com/%s/%s"
	bucketNameEupho = "sound-euphonium"
)

func NewGCSStore(cfg *config.Config) (*store.Store, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	client.Bucket(bucketNameEupho)
	eupho, err := newGCSEuphoStore(cfg, client)
	if err != nil {
		return nil, err
	}
	return &store.Store{
		Eupho: eupho,
	}, nil
}

type gcsEuphoStore struct {
	cfg *config.Config

	bkt *storage.BucketHandle

	objs     []*storage.ObjectAttrs
	cachedAt time.Time
}

func newGCSEuphoStore(cfg *config.Config, client *storage.Client) (*gcsEuphoStore, error) {
	s := &gcsEuphoStore{
		bkt: client.Bucket(bucketNameEupho),
		cfg: cfg,
	}
	return s, s.cacheObjs()
}

func (s *gcsEuphoStore) Get() (*entity.EuphoImage, error) {
	if expired(s.cachedAt, s.cfg.Store.GCS.CacheDuration) {
		if err := s.cacheObjs(); err != nil {
			return nil, err
		}
	}
	n := rand.Intn(len(s.objs))
	return &entity.EuphoImage{
		URL: fmt.Sprintf(urlFormat, bucketNameEupho, s.objs[n].Name),
	}, nil
}

func (s *gcsEuphoStore) cacheObjs() error {
	it := s.bkt.Objects(context.Background(), nil)

	attrs := make([]*storage.ObjectAttrs, 0, 500)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		attrs = append(attrs, objAttrs)
	}
	s.objs = attrs
	rand.Shuffle(len(s.objs), func(i, j int) {
		s.objs[i], s.objs[j] = s.objs[j], s.objs[i]
	})
	s.cachedAt = time.Now()
	return nil
}

func expired(t time.Time, duration time.Duration) bool {
	return time.Now().After(t.Add(duration))
}
