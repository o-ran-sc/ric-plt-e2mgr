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

package controllers

import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/providers"
	"e2mgr/rNibWriter"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Controller struct {
	logger         *logger.Logger
	handlerProvider *providers.IncomingRequestHandlerProvider
	rmrResponseChannel chan<- *models.NotificationResponse
}

func NewController(logger *logger.Logger, rNibReaderProvider func() reader.RNibReader, rNibWriterProvider func() rNibWriter.RNibWriter,
	config *configuration.Configuration, rmrResponseChannel chan<- *models.NotificationResponse) *Controller {

	provider := providers.NewIncomingRequestHandlerProvider(logger, config, rNibWriterProvider, rNibReaderProvider)
	return &Controller{
		logger: logger,
		handlerProvider: provider,
		rmrResponseChannel: rmrResponseChannel,
	}
}

func (c *Controller)ShutdownHandler(writer http.ResponseWriter, r *http.Request, params httprouter.Params){

	c.handleRequest(writer, &r.Header, providers.ShutdownRequest,nil, false, http.StatusNoContent)
}

//TODO create struct to input parameters
func (c *Controller) handleRequest(writer http.ResponseWriter, header *http.Header, requestName providers.IncomingRequest,
	request *models.Request, validateHeader bool, httpStatusResponse int) {

	c.logger.Infof("[Client -> E2 Manager] #controller.handleRequest - request: %v", requestName) //TODO print request if exist

	if validateHeader {

		err := c.validateRequestHeader(header)
		if err != nil {
			c.handleErrorResponse(err, writer)
			return
		}
	}

	handler,err := c.handlerProvider.GetHandler(requestName)
	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	err = handler.Handle(c.logger, request, c.rmrResponseChannel)

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	writer.WriteHeader(httpStatusResponse)
	c.logger.Infof("[E2 Manager -> Client] #controller.handleRequest - status response: %v", httpStatusResponse)
}

func (c *Controller) validateRequestHeader( header *http.Header) error {

	if header.Get("Content-Type") != "application/json"{
		c.logger.Errorf("#controller.validateRequestHeader - validation failure, incorrect content type")

		return  e2managererrors.NewHeaderValidationError()
	}
	return nil
}

func (c *Controller) handleErrorResponse(err error, writer http.ResponseWriter){

	var errorResponseDetails models.ErrorResponse
	var httpError int

	if err != nil {
		switch err.(type) {
		case *e2managererrors.RnibDbError:
			e2Error, _ := err.(*e2managererrors.RnibDbError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Err.Code, Message: e2Error.Err.Message}
			httpError = http.StatusInternalServerError
		case *e2managererrors.CommandAlreadyInProgressError:
			e2Error, _ := err.(*e2managererrors.CommandAlreadyInProgressError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Err.Code, Message: e2Error.Err.Message}
			httpError = http.StatusMethodNotAllowed
		case *e2managererrors.HeaderValidationError:
			e2Error, _ := err.(*e2managererrors.HeaderValidationError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Err.Code, Message: e2Error.Err.Message}
			httpError = http.StatusUnsupportedMediaType
		default:
			e2Error, _ := err.(*e2managererrors.InternalError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Err.Code, Message: e2Error.Err.Message}
			httpError = http.StatusInternalServerError
		}
	}
	errorResponse, _ := json.Marshal(errorResponseDetails)

	c.logger.Errorf("[E2 Manager -> Client] #controller.handleErrorResponse - http status: %d, error response: %+v", httpError, errorResponseDetails)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpError)
	_, err = writer.Write(errorResponse)

	if err != nil {
		c.logger.Errorf("#controller.handleErrorResponse - Cannot send response. writer:%v", writer)
	}
}