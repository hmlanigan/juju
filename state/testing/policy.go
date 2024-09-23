// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	stdcontext "context"
	"fmt"

	"github.com/juju/collections/transform"
	"github.com/juju/errors"
	"github.com/juju/schema"
	"gopkg.in/juju/environschema.v1"

	"github.com/juju/juju/cloud"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/credential"
	domainstorage "github.com/juju/juju/domain/storage"
	storageerrors "github.com/juju/juju/domain/storage/errors"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/environs/envcontext"
	"github.com/juju/juju/internal/storage"
	"github.com/juju/juju/state"
)

type MockPolicy struct {
	GetConfigValidator         func() (config.Validator, error)
	GetConstraintsValidator    func() (constraints.Validator, error)
	GetStorageProviderRegistry func() (storage.ProviderRegistry, error)

	Providers map[string]domainstorage.StoragePoolDetails
}

func (p *MockPolicy) ConfigValidator() (config.Validator, error) {
	if p.GetConfigValidator != nil {
		return p.GetConfigValidator()
	}
	return nil, errors.NotImplementedf("ConfigValidator")
}

func (p *MockPolicy) ConstraintsValidator(ctx envcontext.ProviderCallContext) (constraints.Validator, error) {
	if p.GetConstraintsValidator != nil {
		return p.GetConstraintsValidator()
	}
	return nil, errors.NotImplementedf("ConstraintsValidator")
}

type mockStoragePoolService struct {
	state.StoragePoolGetter
	Providers map[string]domainstorage.StoragePoolDetails
}

func (m *mockStoragePoolService) GetStoragePoolByName(_ stdcontext.Context, name string) (*storage.Config, error) {
	p, ok := m.Providers[name]
	if !ok {
		return nil, fmt.Errorf("storage pool %q not found%w", name, errors.Hide(storageerrors.PoolNotFoundError))
	}
	attrs := transform.Map(p.Attrs, func(k, v string) (string, any) { return k, v })
	return storage.NewConfig(name, storage.ProviderType(p.Provider), attrs)
}

type noopStoragePoolGetter struct{}

func (noopStoragePoolGetter) GetStoragePoolByName(ctx stdcontext.Context, name string) (*storage.Config, error) {
	return nil, fmt.Errorf("storage pool %q not found%w", name, errors.Hide(storageerrors.PoolNotFoundError))
}

func (p *MockPolicy) StorageServices() (state.StoragePoolGetter, storage.ProviderRegistry, error) {
	if p.GetStorageProviderRegistry != nil {
		registry, err := p.GetStorageProviderRegistry()
		return &mockStoragePoolService{Providers: p.Providers}, registry, err
	}
	return noopStoragePoolGetter{}, nil, errors.NotImplementedf("StorageProviderRegistry")
}

type MockConfigSchemaSource struct {
	CloudName string
}

func (m *MockConfigSchemaSource) ConfigSchema() schema.Fields {
	configSchema := environschema.Fields{
		"providerAttr" + m.CloudName: {
			Type: environschema.Tstring,
		},
	}
	fs, _, err := configSchema.ValidationSchema()
	if err != nil {
		panic(err)
	}
	return fs
}

func (m *MockConfigSchemaSource) ConfigDefaults() schema.Defaults {
	return schema.Defaults{
		"providerAttr" + m.CloudName: "vulch",
	}
}

type MockCredentialService struct {
	Credential *cloud.Credential
}

func (m *MockCredentialService) CloudCredential(ctx stdcontext.Context, key credential.Key) (cloud.Credential, error) {
	if m.Credential == nil {
		return cloud.Credential{}, errors.NotFoundf("credential %q", key)
	}
	return *m.Credential, nil
}

type MockCloudService struct {
	CloudInfo *cloud.Cloud
}

func (m *MockCloudService) Cloud(ctx stdcontext.Context, name string) (*cloud.Cloud, error) {
	if m.CloudInfo == nil {
		return nil, errors.NotFoundf("cloud %q", name)
	}
	return m.CloudInfo, nil
}

type MockConfigService struct {
	Config *config.Config
}

func (m *MockConfigService) ModelConfig(_ stdcontext.Context) (*config.Config, error) {
	// name (model name)
	// uuid (model uuid)
	return m.Config, nil
}
