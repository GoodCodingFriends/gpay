package gcs

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoodCodingFriends/gpay/source"
	"github.com/morikuni/failure"
	"google.golang.org/api/iterator"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type gcsSource struct {
	c           *storage.Client
	bucketNames []string
}

func New(ctx context.Context, bucketNames []string) (source.Source, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	for _, name := range bucketNames {
		if _, err := bucketSize(name); err != nil {
			return nil, failure.Wrap(
				err,
				failure.WithCode(source.InvalidParameterCode),
				failure.Context{"bucket": name},
			)
		}
	}

	return &gcsSource{c: client, bucketNames: bucketNames}, nil
}

func (s *gcsSource) Random(ctx context.Context) (io.ReadCloser, error) {
	bktName := s.bucketNames[rand.Int31n(int32(len(s.bucketNames)))]
	bkt := s.c.Bucket(bktName)
	n, err := bucketSize(bktName)
	if err != nil {
		return nil, failure.Wrap(err)
	}
	iter := bkt.Objects(ctx, nil)
	var (
		i   int
		obj *storage.ObjectAttrs
	)
	for {
		obj, err = iter.Next()
		if err == iterator.Done {
			return nil, failure.Wrap(err, failure.Message("unexpected iterator.Done"))
		}
		if i == n {
			break
		}
		i++
	}

	return bkt.Object(obj.Name).NewReader(ctx)
}

func bucketSize(name string) (int, error) {
	s := fmt.Sprintf("GCS_BUCKET_SIZE_%s", strings.ToUpper(name))
	n, err := strconv.Atoi(os.Getenv(s))
	if err != nil {
		return 0, failure.Wrap(err)
	}
	return n, nil
}
