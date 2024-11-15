package resource

import (
	"context"
	"fmt"
	"io"

	"github.com/juju/juju/core/objectstore"
	"github.com/juju/juju/internal/charm/resource"
)

type ResourceStore interface {
	// Get returns an io.ReadCloser for data at path, namespaced to the
	// model.
	Get(context.Context, string) (io.ReadCloser, int64, error)
}

type OCIResourceStore struct {
}

func (f OCIResourceStore) Get(ctx context.Context, path string) (io.ReadCloser, int64, error) {
	return nil, -1, nil
}

type FileResourceStore struct {
	objectStore objectstore.ObjectStore
}

func (f FileResourceStore) Get(ctx context.Context, path string) (io.ReadCloser, int64, error) {
	return f.objectStore.Get(ctx, path)
}

type ResourceStoreFactory struct {
	objectStore objectstore.ModelObjectStoreGetter
}

func NewResourceStoreFactory(objectStore objectstore.ModelObjectStoreGetter) ResourceStoreFactory {
	return ResourceStoreFactory{objectStore: objectStore}
}

func (f ResourceStoreFactory) GetResourceStore(ctx context.Context, t resource.Type) (ResourceStore, error) {
	switch t {
	case resource.TypeContainerImage:
		return OCIResourceStore{}, nil
	case resource.TypeFile:
		objectstore, err := f.objectStore.GetObjectStore(ctx)
		if err != nil {
			return nil, err
		}

		return FileResourceStore{objectStore: objectstore}, nil
	default:
		return nil, fmt.Errorf("unknown resource type %q", t)
	}
}
