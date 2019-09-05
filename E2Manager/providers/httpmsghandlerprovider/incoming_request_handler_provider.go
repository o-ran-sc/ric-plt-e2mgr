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
//

package httpmsghandlerprovider

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/handlers/httpmsghandlers"
	"e2mgr/logger"
	"e2mgr/rNibWriter"
	"e2mgr/services"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
)

type IncomingRequest string

const (
	ShutdownRequest IncomingRequest = "Shutdown"
	ResetRequest    IncomingRequest = "Reset"
)

type IncomingRequestHandlerProvider struct {
	requestMap map[IncomingRequest]httpmsghandlers.RequestHandler
	logger     *logger.Logger
}

func NewIncomingRequestHandlerProvider(logger *logger.Logger, rmrService *services.RmrService, config *configuration.Configuration, rNibWriterProvider func() rNibWriter.RNibWriter, rNibReaderProvider func() reader.RNibReader) *IncomingRequestHandlerProvider {

	return &IncomingRequestHandlerProvider{
		requestMap: initRequestHandlerMap(rmrService, config, rNibWriterProvider, rNibReaderProvider),
		logger:     logger,
	}
}

func initRequestHandlerMap(rmrService *services.RmrService, config *configuration.Configuration, rNibWriterProvider func() rNibWriter.RNibWriter,
	rNibReaderProvider func() reader.RNibReader) map[IncomingRequest]httpmsghandlers.RequestHandler {

	return map[IncomingRequest]httpmsghandlers.RequestHandler{
		ShutdownRequest: httpmsghandlers.NewDeleteAllRequestHandler(rmrService, config, rNibWriterProvider, rNibReaderProvider), //TODO change to pointer
		ResetRequest:    httpmsghandlers.NewX2ResetRequestHandler(rmrService, config, rNibWriterProvider, rNibReaderProvider),   //TODO change to pointer
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
