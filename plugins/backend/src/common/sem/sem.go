// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//	Hyper-Cluster Common
//
//
//
//																										2020.08.01
//																										DAESEOB.JEONG
//
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
package sem

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"
)

type TSemaphore struct {
	ctx               context.Context
	semaphoreWeighted *semaphore.Weighted
	i64ResourceCount  int64
}

func (sem *TSemaphore) Init(_i64NumberOfResouce int64) {
	sem.i64ResourceCount = _i64NumberOfResouce
	sem.ctx = context.TODO()
	sem.semaphoreWeighted = semaphore.NewWeighted(sem.i64ResourceCount)
}

func (sem *TSemaphore) Acquire() error {
	if sem.semaphoreWeighted != nil {
		return sem.semaphoreWeighted.Acquire(sem.ctx, 1)
	}
	return fmt.Errorf("semaphoreWeighted is nil.")
}

func (sem *TSemaphore) Release() {
	if sem.semaphoreWeighted != nil {
		sem.semaphoreWeighted.Release(1)
	}
}

func (sem *TSemaphore) WaitForFinish() error {
	if sem.semaphoreWeighted != nil {
		if err := sem.semaphoreWeighted.Acquire(sem.ctx, sem.i64ResourceCount); err != nil {
			return err
		}
		sem.semaphoreWeighted.Release(sem.i64ResourceCount)
		sem.semaphoreWeighted = nil
	}
	return nil
}
