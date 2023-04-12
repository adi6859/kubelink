// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/devtron-labs/kubelink/api/router"
	"github.com/devtron-labs/kubelink/internal/lock"
	"github.com/devtron-labs/kubelink/internal/logger"
	"github.com/devtron-labs/kubelink/pkg/cluster"
	"github.com/devtron-labs/kubelink/pkg/k8sInformer"
	"github.com/devtron-labs/kubelink/pkg/service"
	"github.com/devtron-labs/kubelink/pkg/sql"
	"github.com/devtron-labs/kubelink/pprof"
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
	k8sInformerImpl := k8sInformer.Newk8sInformerImpl(sugaredLogger, clusterRepositoryImpl, helmReleaseConfig)
	serviceHelmReleaseConfig, err := service.GetHelmReleaseConfig()
	if err != nil {
		return nil, err
	}
	helmAppServiceImpl := service.NewHelmAppServiceImpl(sugaredLogger, k8sServiceImpl, k8sInformerImpl, serviceHelmReleaseConfig)
	applicationServiceServerImpl := service.NewApplicationServiceServerImpl(sugaredLogger, chartRepositoryLocker, helmAppServiceImpl)
	pProfRestHandlerImpl := pprof.NewPProfRestHandler(sugaredLogger)
	pProfRouterImpl := pprof.NewPProfRouter(sugaredLogger, pProfRestHandlerImpl)
	routerImpl := router.NewRouter(sugaredLogger, pProfRouterImpl)
	app := NewApp(sugaredLogger, applicationServiceServerImpl, routerImpl, k8sInformerImpl)
	return app, nil
}
