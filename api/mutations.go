package api

import (
	"context"

	"github.com/arrow2nd/anct/gen"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// UpdateWorkState : 作品の視聴ステータスを更新
func (a *API) UpdateWorkState(id string, state gen.StatusState) error {
	ctx := context.Background()
	if _, err := a.client.UpdateWorkState(ctx, id, state); err != nil {
		return err
	}

	return nil
}

// CreateEpisodeRecords : エピソードの視聴記録を作成
func (a *API) CreateEpisodeRecords(episodeIDs []string, rating gen.RatingState, comment string) error {
	const (
		resourceMax = int64(2)
		weight      = int64(1)
	)

	// 同時実行数を制限
	sem := semaphore.NewWeighted(resourceMax)
	eg, ctx := errgroup.WithContext(context.Background())

	for _, ID := range episodeIDs {
		ID := ID

		sem.Acquire(ctx, weight)

		eg.Go(func() error {
			if _, err := a.client.CreateEpisodeRecord(ctx, ID, rating, &comment); err != nil {
				return err
			}

			sem.Release(weight)
			return nil
		})
	}

	return eg.Wait()
}
