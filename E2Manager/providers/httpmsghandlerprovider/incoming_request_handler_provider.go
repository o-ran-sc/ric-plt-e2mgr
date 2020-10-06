//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).

package httpmsghandlerprovider

import (
	"e2mgr/clients"
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/handlers/httpmsghandlers"
	"e2mgr/logger"
	"e2mgr/managers"
	"e2mgr/services"
	"e2mgr/services/rmrsender"
)

type IncomingRequest string

const (
	SetGeneralConfigurationRequest IncomingRequest = "SetGeneralConfiguration"
	ShutdownRequest                IncomingRequest = "Shutdown"
	ResetRequest                   IncomingRequest = "Reset"
	GetNodebRequest                IncomingRequest = "GetNodebRequest"
	GetNodebIdListRequest          IncomingRequest = "GetNodebIdListRequest"
	GetNodebIdRequest          	   IncomingRequest = "GetNodebIdRequest"
	GetE2TInstancesRequest         IncomingRequest = "GetE2TInstancesRequest"
	UpdateGnbRequest               IncomingRequest = "UpdateGnbRequest"
	UpdateEnbRequest               IncomingRequest = "UpdateEnbRequest"
	AddEnbRequest                  IncomingRequest = "AddEnbRequest"
	DeleteEnbRequest               IncomingRequest = "DeleteEnbRequest"
	HealthCheckRequest             IncomingRequest = "HealthCheckRequest"
)

type IncomingRequestHandlerProvider struct {
	requestMap                    map[IncomingRequest]httpmsghandlers.RequestHandler
	logger                        *logger.Logger
	ranConnectStatusChangeManager managers.IRanConnectStatusChangeManager
}

func NewIncomingRequestHandlerProvider(logger *logger.Logger, rmrSender *rmrsender.RmrSender, config *configuration.Configuration, rNibDataService services.RNibDataService, e2tInstancesManager managers.IE2TInstancesManager, rmClient clients.IRoutingManagerClient, ranConnectStatusChangeManager managers.IRanConnectStatusChangeManager, nodebValidator *managers.NodebValidator, updateEnbManager managers.IUpdateNodebManager, updateGnbManager managers.IUpdateNodebManager, ranListManager managers.RanListManager) *IncomingRequestHandlerProvider {

	return &IncomingRequestHandlerProvider{
		requestMap:                    initRequestHandlerMap(logger, rmrSender, config, rNibDataService, e2tInstancesManager, rmClient, ranConnectStatusChangeManager, nodebValidator, updateEnbManager, updateGnbManager, ranListManager),
		logger:                        logger,
		ranConnectStatusChangeManager: ranConnectStatusChangeManager,
	}
}

func initRequestHandlerMap(logger *logger.Logger, rmrSender *rmrsender.RmrSender, config *configuration.Configuration, rNibDataService services.RNibDataService, e2tInstancesManager managers.IE2TInstancesManager, rmClient clients.IRoutingManagerClient, ranConnectStatusChangeManager managers.IRanConnectStatusChangeManager, nodebValidator *managers.NodebValidator, updateEnbManager managers.IUpdateNodebManager, updateGnbManager managers.IUpdateNodebManager, ranListManager managers.RanListManager) map[IncomingRequest]httpmsghandlers.RequestHandler {

	return map[IncomingRequest]httpmsghandlers.RequestHandler{
		ShutdownRequest:                httpmsghandlers.NewDeleteAllRequestHandler(logger, rmrSender, config, rNibDataService, e2tInstancesManager, rmClient, ranConnectStatusChangeManager, ranListManager),
		ResetRequest:                   httpmsghandlers.NewX2ResetRequestHandler(logger, rmrSender, rNibDataService),
		SetGeneralConfigurationRequest: httpmsghandlers.NewSetGeneralConfigurationHandler(logger, rNibDataService),
		GetNodebRequest:                httpmsghandlers.NewGetNodebRequestHandler(logger, rNibDataService),
		GetNodebIdListRequest:          httpmsghandlers.NewGetNodebIdListRequestHandler(logger, rNibDataService, ranListManager),
		GetNodebIdRequest:          	httpmsghandlers.NewGetNodebIdRequestHandler(logger, ranListManager),
		GetE2TInstancesRequest:         httpmsghandlers.NewGetE2TInstancesRequestHandler(logger, e2tInstancesManager),
		UpdateGnbRequest:               httpmsghandlers.NewUpdateNodebRequestHandler(logger, rNibDataService, updateGnbManager),
		UpdateEnbRequest:               httpmsghandlers.NewUpdateNodebRequestHandler(logger, rNibDataService, updateEnbManager),
		AddEnbRequest:                  httpmsghandlers.NewAddEnbRequestHandler(logger, rNibDataService, nodebValidator, ranListManager),
		DeleteEnbRequest:               httpmsghandlers.NewDeleteEnbRequestHandler(logger, rNibDataService, ranListManager),
		HealthCheckRequest:             httpmsghandlers.NewHealthCheckRequestHandler(logger, rNibDataService, ranListManager, rmrSender),
	}
}

func (provider IncomingRequestHandlerProvider) GetHandler(requestType IncomingRequest) (httpmsghandlers.RequestHandler, error) {
	handler, ok := provider.requestMap[requestType]

	if !ok {
		provider.logger.Errorf("#incoming_request_handler_provider.GetHandler - Cannot find handler for request type: %s", requestType)
		return nil, e2managererrors.NewInternalError()
	}

	return handler, nil
}
