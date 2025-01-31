// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/devtron-labs/authenticator/client"
	"github.com/devtron-labs/common-lib/utils/k8s"
	"github.com/devtron-labs/kubelink/api/router"
	"github.com/devtron-labs/kubelink/internal/lock"
	"github.com/devtron-labs/kubelink/internal/logger"
	"github.com/devtron-labs/kubelink/pkg/cache"
	"github.com/devtron-labs/kubelink/pkg/cluster"
	"github.com/devtron-labs/kubelink/pkg/k8sInformer"
	"github.com/devtron-labs/kubelink/pkg/service"
	"github.com/devtron-labs/kubelink/pkg/sql"
	"github.com/devtron-labs/kubelink/pprof"
	"github.com/devtron-labs/kubelink/statsViz"
)

// Injectors from Wire.go:

func InitializeApp() (*App, error) {
	sugaredLogger := logger.NewSugaredLogger()
	chartRepositoryLocker := lock.NewChartRepositoryLocker(sugaredLogger)
	k8sServiceImpl := service.NewK8sServiceImpl(sugaredLogger)
	config, err := sql.GetConfig()
	if err != nil {
		return nil, err
	}
	db, err := sql.NewDbConnection(config, sugaredLogger)
	if err != nil {
		return nil, err
	}
	clusterRepositoryImpl := repository.NewClusterRepositoryImpl(db, sugaredLogger)
	helmReleaseConfig, err := k8sInformer.GetHelmReleaseConfig()
	if err != nil {
		return nil, err
	}
	runtimeConfig, err := client.GetRuntimeConfig()
	if err != nil {
		return nil, err
	}
	k8sUtil := k8s.NewK8sUtil(sugaredLogger, runtimeConfig)
	k8sInformerImpl := k8sInformer.Newk8sInformerImpl(sugaredLogger, clusterRepositoryImpl, helmReleaseConfig, k8sUtil)
	serviceHelmReleaseConfig, err := service.GetHelmReleaseConfig()
	if err != nil {
		return nil, err
	}
	clusterCacheConfig, err := cache.GetClusterCacheConfig()
	if err != nil {
		return nil, err
	}
	clusterCacheImpl := cache.NewClusterCacheImpl(sugaredLogger, clusterCacheConfig, clusterRepositoryImpl, k8sUtil, k8sInformerImpl)
	helmAppServiceImpl := service.NewHelmAppServiceImpl(sugaredLogger, k8sServiceImpl, k8sInformerImpl, serviceHelmReleaseConfig, k8sUtil, clusterRepositoryImpl, clusterCacheImpl)
	applicationServiceServerImpl := service.NewApplicationServiceServerImpl(sugaredLogger, chartRepositoryLocker, helmAppServiceImpl)
	pProfRestHandlerImpl := pprof.NewPProfRestHandler(sugaredLogger)
	pProfRouterImpl := pprof.NewPProfRouter(sugaredLogger, pProfRestHandlerImpl)
	statVizConfig, err := statsViz.GetStatsVizConfig()
	if err != nil {
		return nil, err
	}
	statsVizRouterImpl := statsViz.NewStatsVizRouter(sugaredLogger, statVizConfig)
	routerImpl := router.NewRouter(sugaredLogger, pProfRouterImpl, statsVizRouterImpl)
	app := NewApp(sugaredLogger, applicationServiceServerImpl, routerImpl, k8sInformerImpl)
	return app, nil
}
