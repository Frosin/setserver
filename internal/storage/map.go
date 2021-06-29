package storage

import (
	"context"
	"sync"
)

type MapStorage struct {
	storage map[string]string
	mtx     sync.RWMutex
}

func NewMapStorage() Storage {
	return &MapStorage{
		storage: make(map[string]string),
		mtx:     sync.RWMutex{},
	}
}

func (m *MapStorage) Add(ctx context.Context, key, value string) error {
	errChan := make(chan error)
	go func() {
		m.mtx.Lock()
		m.storage[key] = value
		m.mtx.Unlock()
		errChan <- nil
	}()

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errChan:
	}
	return err
}

func (m *MapStorage) Delete(ctx context.Context, key string) (string, error) {
	errChan := make(chan error)
	value := ""
	go func() {
		ok := false
		err := error(nil)
		m.mtx.RLock()
		value, ok = m.storage[key]
		m.mtx.RUnlock()
		if ok {
			m.mtx.Lock()
			delete(m.storage, key)
			m.mtx.Unlock()
		}
		if !ok {
			err = ErrNotFound
		}
		errChan <- err
	}()

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errChan:
	}
	return value, err
}
