package sync

import (
	"github.com/juju/errors"

	"github.com/bangwork/import-tools/serve/services/importer/constants"
	"github.com/bangwork/import-tools/serve/services/importer/resolve"
	"github.com/bangwork/import-tools/serve/services/importer/resolve/jira"
	"github.com/bangwork/import-tools/serve/services/importer/types"
)

var resolverFactoryMap = map[int]resolve.ResolverFactory{}

func registerResolverFactory(importType int, resolverFactory resolve.ResolverFactory) {
	resolverFactoryMap[importType] = resolverFactory
}

func getResolverFactory(importType int) resolve.ResolverFactory {
	return resolverFactoryMap[importType]
}

func createResolver(importTask *types.ImportTask) (resolve.ResourceResolver, error) {
	factory := getResolverFactory(importTask.ImportType)
	if factory != nil {
		return factory.CreateResolver(importTask)
	}
	return nil, nil
}

func InitImportFile(importTask *types.ImportTask) (resolve.ResourceResolver, error) {
	factory := getResolverFactory(importTask.ImportType)
	if factory != nil {
		res, err := factory.InitImportFile(importTask)
		if err != nil {
			return nil, errors.Trace(err)
		}
		return res, nil
	}
	return nil, nil
}

func InitResolverFactory() error {
	registerResolverFactory(constants.ImportTypeJira, &jira.JiraResolverFactory{})

	return nil
}
