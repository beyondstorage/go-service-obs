package obs

import (
	"context"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"

	ps "github.com/beyondstorage/go-storage/v4/pairs"
	. "github.com/beyondstorage/go-storage/v4/types"
)

func (s *Service) create(ctx context.Context, name string, opt pairServiceCreate) (store Storager, err error) {
	pairs := append(opt.pairs, ps.WithName(name))

	st, err := s.newStorage(pairs...)
	if err != nil {
		return nil, err
	}

	input := &obs.CreateBucketInput{
		Bucket: name,
	}

	_, err = s.service.CreateBucket(input)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (s *Service) delete(ctx context.Context, name string, opt pairServiceDelete) (err error) {
	_, err = s.service.DeleteBucket(name)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) get(ctx context.Context, name string, opt pairServiceGet) (store Storager, err error) {
	pairs := append(opt.pairs, ps.WithName(name))

	st, err := s.newStorage(pairs...)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (s *Service) list(ctx context.Context, opt pairServiceList) (sti *StoragerIterator, err error) {
	input := &storagePageStatus{}

	return NewStoragerIterator(ctx, s.nextStoragePage, input), nil
}

func (s *Service) nextStoragePage(ctx context.Context, page *StoragerPage) error {
	input := &obs.ListBucketsInput{
		QueryLocation: true,
	}

	output, err := s.service.ListBuckets(input)
	if err != nil {
		return err
	}

	for _, v := range output.Buckets {
		store, err := s.newStorage(ps.WithName(v.Name))
		if err != nil {
			return err
		}

		page.Data = append(page.Data, store)
	}

	return IterateDone
}
