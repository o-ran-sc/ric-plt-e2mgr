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

package controllers

import (
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/providers/httpmsghandlerprovider"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

const (
	ParamRanName = "ranName"
	LimitRequest = 2000
)
const ApplicationJson = "application/json"
const ContentType = "Content-Type"

type INodebController interface {
	Shutdown(writer http.ResponseWriter, r *http.Request)
	X2Reset(writer http.ResponseWriter, r *http.Request)
	GetNodeb(writer http.ResponseWriter, r *http.Request)
	UpdateGnb(writer http.ResponseWriter, r *http.Request)
	UpdateEnb(writer http.ResponseWriter, r *http.Request)
	GetNodebIdList(writer http.ResponseWriter, r *http.Request)
	SetGeneralConfiguration(writer http.ResponseWriter, r *http.Request)
	AddEnb(writer http.ResponseWriter, r *http.Request)
	DeleteEnb(writer http.ResponseWriter, r *http.Request)
}

type NodebController struct {
	logger          *logger.Logger
	handlerProvider *httpmsghandlerprovider.IncomingRequestHandlerProvider
}

func NewNodebController(logger *logger.Logger, handlerProvider *httpmsghandlerprovider.IncomingRequestHandlerProvider) *NodebController {
	return &NodebController{
		logger:          logger,
		handlerProvider: handlerProvider,
	}
}

func (c *NodebController) GetNodebIdList(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.GetNodebIdList - request: %v", c.prettifyRequest(r))

	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.GetNodebIdListRequest, nil, false, http.StatusOK)
}

func (c *NodebController) GetNodeb(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.GetNodeb - request: %v", c.prettifyRequest(r))
	vars := mux.Vars(r)
	ranName := vars["ranName"]
	request := models.GetNodebRequest{RanName: ranName}
	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.GetNodebRequest, request, false, http.StatusOK)
}

func (c *NodebController) UpdateGnb(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.UpdateGnb - request: %v", c.prettifyRequest(r))
	vars := mux.Vars(r)
	ranName := vars[ParamRanName]

	request := models.UpdateGnbRequest{}

	gnb := entities.Gnb{}

	if !c.extractRequestBodyToProto(r, &gnb, writer) {
		return
	}

	request.Gnb = &gnb
	request.RanName = ranName
	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.UpdateGnbRequest, request, true, http.StatusOK)
}

func (c *NodebController) UpdateEnb(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.UpdateEnb - request: %v", c.prettifyRequest(r))
	vars := mux.Vars(r)
	ranName := vars[ParamRanName]

	request := models.UpdateNodebRequest{}

	enb := entities.Enb{}

	if !c.extractRequestBodyToProto(r, &enb, writer) {
		return
	}

	request.Enb = &enb
	request.RanName = ranName
	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.UpdateEnbRequest, request, true, http.StatusOK)
}

func (c *NodebController) AddEnb(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.AddEnb - request: %v", c.prettifyRequest(r))

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		c.logger.Errorf("[Client -> E2 Manager] #NodebController.AddEnb - unable to read request body - error: %s", err)
		c.handleErrorResponse(e2managererrors.NewInvalidJsonError(), writer)
		return
	}

	addEnbRequest := models.AddEnbRequest{}
	err = json.Unmarshal(body, &addEnbRequest)

	if err != nil {
		c.logger.Errorf("[Client -> E2 Manager] #NodebController.AddEnb - unable to unmarshal json - error: %s", err)
		c.handleErrorResponse(e2managererrors.NewInvalidJsonError(), writer)
		return
	}

	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.AddEnbRequest, &addEnbRequest, true, http.StatusCreated)
}

func (c *NodebController) DeleteEnb(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.DeleteEnb - request: %v", c.prettifyRequest(r))
	vars := mux.Vars(r)
	ranName := vars["ranName"]
	request := models.DeleteEnbRequest{RanName: ranName}
	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.DeleteEnbRequest, request, true, http.StatusNoContent)
}

func (c *NodebController) SetGeneralConfiguration(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.SetGeneralConfiguration - request: %v", c.prettifyRequest(r))

	request := models.GeneralConfigurationRequest{}

	if !c.extractJsonBodyDisallowUnknownFields(r, &request, writer) {
		return
	}
	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.SetGeneralConfigurationRequest, request, false, http.StatusOK)
}

func (c *NodebController) Shutdown(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.Shutdown - request: %v", c.prettifyRequest(r))
	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.ShutdownRequest, nil, false, http.StatusNoContent)
}

func (c *NodebController) X2Reset(writer http.ResponseWriter, r *http.Request) {
	c.logger.Infof("[Client -> E2 Manager] #NodebController.X2Reset - request: %v", c.prettifyRequest(r))
	request := models.ResetRequest{}
	vars := mux.Vars(r)
	ranName := vars[ParamRanName]

	if err := c.extractJsonBody(r, &request); err != nil {
		c.handleErrorResponse(err, writer)
		return
	}
	request.RanName = ranName
	c.handleRequest(writer, &r.Header, httpmsghandlerprovider.ResetRequest, request, false, http.StatusNoContent)
}

func (c *NodebController) extractRequestBodyToProto(r *http.Request, pb proto.Message, writer http.ResponseWriter) bool {
	defer r.Body.Close()

	err := jsonpb.Unmarshal(r.Body, pb)

	if err != nil {
		c.logger.Errorf("[Client -> E2 Manager] #NodebController.extractJsonBody - unable to extract json body - error: %s", err)
		c.handleErrorResponse(e2managererrors.NewInvalidJsonError(), writer)
		return false
	}

	return true
}

func (c *NodebController) extractJsonBodyDisallowUnknownFields(r *http.Request, request models.Request, writer http.ResponseWriter) bool {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		c.logger.Errorf("[Client -> E2 Manager] #NodebController.extractJsonBody - unable to extract json body - error: %s", err)
		c.handleErrorResponse(e2managererrors.NewInvalidJsonError(), writer)
		return false
	}

	return true
}

func (c *NodebController) extractJsonBody(r *http.Request, request models.Request) error {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, LimitRequest))

	if err != nil {
		c.logger.Errorf("[Client -> E2 Manager] #NodebController.extractJsonBody - unable to extract json body - error: %s", err)
		return e2managererrors.NewInvalidJsonError()
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		c.logger.Errorf("[Client -> E2 Manager] #NodebController.extractJsonBody - unable to extract json body - error: %s", err)
		return e2managererrors.NewInvalidJsonError()
	}

	return nil
}

func (c *NodebController) handleRequest(writer http.ResponseWriter, header *http.Header, requestName httpmsghandlerprovider.IncomingRequest, request models.Request, validateRequestHeaders bool, successStatusCode int) {

	if validateRequestHeaders {

		err := c.validateRequestHeader(header)
		if err != nil {
			c.handleErrorResponse(err, writer)
			return
		}
	}

	handler, err := c.handlerProvider.GetHandler(requestName)

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	response, err := handler.Handle(request)

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	if successStatusCode == http.StatusNoContent {
		writer.WriteHeader(successStatusCode)
		c.logger.Infof("[E2 Manager -> Client] #NodebController.handleRequest - status response: %v", http.StatusNoContent)
		return
	}

	result, err := response.Marshal()

	if err != nil {
		c.handleErrorResponse(err, writer)
		return
	}

	c.logger.Infof("[E2 Manager -> Client] #NodebController.handleRequest - response: %s", result)
	writer.Header().Set(ContentType, ApplicationJson)
	writer.WriteHeader(successStatusCode)
	writer.Write(result)
}

func (c *NodebController) validateRequestHeader(header *http.Header) error {

	if header.Get(ContentType) != ApplicationJson {
		c.logger.Errorf("#NodebController.validateRequestHeader - validation failure, incorrect content type")

		return e2managererrors.NewHeaderValidationError()
	}
	return nil
}

func (c *NodebController) handleErrorResponse(err error, writer http.ResponseWriter) {

	var errorResponseDetails models.ErrorResponse
	var httpError int

	if err != nil {
		switch err.(type) {
		case *e2managererrors.RnibDbError:
			e2Error, _ := err.(*e2managererrors.RnibDbError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusInternalServerError
		case *e2managererrors.CommandAlreadyInProgressError:
			e2Error, _ := err.(*e2managererrors.CommandAlreadyInProgressError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusMethodNotAllowed
		case *e2managererrors.HeaderValidationError:
			e2Error, _ := err.(*e2managererrors.HeaderValidationError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusUnsupportedMediaType
		case *e2managererrors.WrongStateError:
			e2Error, _ := err.(*e2managererrors.WrongStateError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusBadRequest
		case *e2managererrors.RequestValidationError:
			e2Error, _ := err.(*e2managererrors.RequestValidationError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusBadRequest
		case *e2managererrors.InvalidJsonError:
			e2Error, _ := err.(*e2managererrors.InvalidJsonError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusBadRequest
		case *e2managererrors.RmrError:
			e2Error, _ := err.(*e2managererrors.RmrError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusInternalServerError
		case *e2managererrors.ResourceNotFoundError:
			e2Error, _ := err.(*e2managererrors.ResourceNotFoundError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusNotFound
		case *e2managererrors.E2TInstanceAbsenceError:
			e2Error, _ := err.(*e2managererrors.E2TInstanceAbsenceError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusServiceUnavailable
		case *e2managererrors.RoutingManagerError:
			e2Error, _ := err.(*e2managererrors.RoutingManagerError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusServiceUnavailable
		case *e2managererrors.NodebExistsError:
			e2Error, _ := err.(*e2managererrors.NodebExistsError)
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusBadRequest
		default:
			e2Error := e2managererrors.NewInternalError()
			errorResponseDetails = models.ErrorResponse{Code: e2Error.Code, Message: e2Error.Message}
			httpError = http.StatusInternalServerError
		}
	}
	errorResponse, _ := json.Marshal(errorResponseDetails)

	c.logger.Errorf("[E2 Manager -> Client] #NodebController.handleErrorResponse - http status: %d, error response: %+v", httpError, errorResponseDetails)

	writer.Header().Set(ContentType, ApplicationJson)
	writer.WriteHeader(httpError)
	_, err = writer.Write(errorResponse)
}

func (c *NodebController) prettifyRequest(request *http.Request) string {
	dump, _ := httputil.DumpRequest(request, true)
	requestPrettyPrint := strings.Replace(string(dump), "\r\n", " ", -1)
	return strings.Replace(requestPrettyPrint, "\n", "", -1)
}
