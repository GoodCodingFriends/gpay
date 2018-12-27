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
	urlFormat = "https://storage.googleapis.com/%s/%s"

	bucketNameEupho   = "eupho"
	bucketNameHamachi = "hamachi"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	hamachi, err := newGCSHamachiStore(cfg, client)
	if err != nil {
		return nil, err
	}
	return &store.Store{
		Eupho:   eupho,
		Hamachi: hamachi,
	}, nil
}

type gcsEuphoStore struct {
	*baseStore
}

func newGCSEuphoStore(cfg *config.Config, client *storage.Client) (*gcsEuphoStore, error) {
	base, err := newBaseStore(cfg, client.Bucket(bucketNameEupho))
	if err != nil {
		return nil, err
	}
	return &gcsEuphoStore{base}, nil
}

func (s *gcsEuphoStore) Get() (*entity.EuphoImage, error) {
	obj, err := s.baseStore.get()
	if err != nil {
		return nil, err
	}
	return &entity.EuphoImage{
		URL: fmt.Sprintf(urlFormat, bucketNameEupho, obj.Name),
	}, nil
}

func (s *gcsEuphoStore) Close() error {
	return nil
}

type gcsHamachiStore struct {
	*baseStore
}

func newGCSHamachiStore(cfg *config.Config, client *storage.Client) (*gcsHamachiStore, error) {
	base, err := newBaseStore(cfg, client.Bucket(bucketNameHamachi))
	if err != nil {
		return nil, err
	}
	return &gcsHamachiStore{base}, nil
}

func (s *gcsHamachiStore) Get() (*entity.HamachiImage, error) {
	obj, err := s.baseStore.get()
	if err != nil {
		return nil, err
	}
	return &entity.HamachiImage{
		URL: fmt.Sprintf(urlFormat, bucketNameHamachi, obj.Name),
	}, nil
}

func (s *gcsHamachiStore) Close() error {
	return nil
}

type baseStore struct {
	cfg *config.Config

	bkt *storage.BucketHandle

	objs     []*storage.ObjectAttrs
	cachedAt time.Time
}

func newBaseStore(cfg *config.Config, bkt *storage.BucketHandle) (*baseStore, error) {
	s := &baseStore{
		cfg: cfg,
		bkt: bkt,
	}
	return s, s.cacheObjs()
}

func (s *baseStore) cacheObjs() error {
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

func (s *baseStore) get() (*storage.ObjectAttrs, error) {
	if expired(s.cachedAt, s.cfg.Store.GCS.CacheDuration) {
		if err := s.cacheObjs(); err != nil {
			return nil, err
		}
	}
	n := rand.Intn(len(s.objs))
	return s.objs[n], nil
}

func expired(t time.Time, duration time.Duration) bool {
	return time.Now().After(t.Add(duration))
}
