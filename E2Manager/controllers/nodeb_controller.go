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
	"e2mgr/logger"
	"e2mgr/models"
	"e2mgr/services"
	"e2mgr/utils"
	"encoding/json"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	validationErrorCode       int = 402
	notFoundErrorCode         int = 404
	internalErrorCode         int = 501
	validationFailedMessage       = "Validation failed"
	notFoundErrorMessage          = "Resource not found"
	internalErrorMessage          = "Internal Server Error. Please try again later"
)

var messageChannel chan *models.E2RequestMessage
var errorChannel chan error

type INodebController interface {
	GetNodebIdList (writer http.ResponseWriter, request *http.Request)
	GetNodeb(writer http.ResponseWriter, request *http.Request)
	HandleHealthCheckRequest(writer http.ResponseWriter, request *http.Request)
}

type NodebController struct {
	rmrService         *services.RmrService
	Logger             *logger.Logger
	rnibDataService services.RNibDataService
}

func NewNodebController(logger *logger.Logger, rmrService *services.RmrService, rnibDataService services.RNibDataService) *NodebController {
	messageChannel = make(chan *models.E2RequestMessage)
	errorChannel = make(chan error)
	return &NodebController{
		rmrService:         rmrService,
		Logger:             logger,
		rnibDataService: rnibDataService,
	}
}

func (rc NodebController) GetNodebIdList (writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	nodebIdList, rnibError := rc.rnibDataService.GetListNodebIds()

	if rnibError != nil {
		rc.Logger.Errorf("%v", rnibError);
		httpStatusCode, errorCode, errorMessage := rnibErrorToHttpError(rnibError)
		handleErrorResponse(rc.Logger, httpStatusCode, errorCode, errorMessage, writer, startTime)
		return;
	}

	pmList := utils.ConvertNodebIdListToProtoMessageList(nodebIdList)
	result, err := utils.MarshalProtoMessageListToJsonArray(pmList)

	if err != nil {
		rc.Logger.Errorf("%v", err);
		handleErrorResponse(rc.Logger, http.StatusInternalServerError, internalErrorCode, internalErrorMessage, writer, startTime)
		return;
	}

	writer.Header().Set("Content-Type", "application/json")
	rc.Logger.Infof("[E2 Manager -> Client] #nodeb_controller.GetNodebIdList - response: %s", result)
	writer.Write([]byte(result))
}

func (rc NodebController) GetNodeb(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	vars := mux.Vars(request)
	ranName := vars["ranName"]
	respondingNode, rnibError := rc.rnibDataService.GetNodeb(ranName)
	if rnibError != nil {
		rc.Logger.Errorf("%v", rnibError)
		httpStatusCode, errorCode, errorMessage := rnibErrorToHttpError(rnibError)
		handleErrorResponse(rc.Logger, httpStatusCode, errorCode, errorMessage, writer, startTime)
		return
	}

	m := jsonpb.Marshaler{}
	result, err := m.MarshalToString(respondingNode)

	if err != nil {
		rc.Logger.Errorf("%v", err)
		handleErrorResponse(rc.Logger, http.StatusInternalServerError, internalErrorCode, internalErrorMessage, writer, startTime)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	rc.Logger.Infof("[E2 Manager -> Client] #nodeb_controller.GetNodeb - response: %s", result)
	writer.Write([]byte(result))
}

func (rc NodebController) HandleHealthCheckRequest(writer http.ResponseWriter, request *http.Request) {
	//fmt.Println("[X-APP -> Client] #HandleHealthCheckRequest - http status: 200")
	writer.WriteHeader(http.StatusOK)
}

func handleErrorResponse(logger *logger.Logger, httpStatus int, errorCode int, errorMessage string, writer http.ResponseWriter, startTime time.Time) {
	errorResponseDetails := models.ErrorResponse{errorCode, errorMessage}
	errorResponse, _ := json.Marshal(errorResponseDetails)
	printHandlingRequestElapsedTimeInMs(logger, startTime)
	logger.Infof("[E2 Manager -> Client] #nodeb_controller.handleErrorResponse - http status: %d, error response: %+v", httpStatus, errorResponseDetails)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(httpStatus)
	_, err := writer.Write(errorResponse)

	if err != nil {
		logger.Errorf("#nodeb_controller.handleErrorResponse - Cannot send response. writer:%v", writer)
	}
}

func printHandlingRequestElapsedTimeInMs(logger *logger.Logger, startTime time.Time) {
	logger.Infof("Summary: #nodeb_controller.printElapsedTimeInMs - Elapsed time for handling request from client to E2 termination: %f ms",
		float64(time.Since(startTime))/float64(time.Millisecond))
}

func rnibErrorToHttpError(rnibError error) (int, int, string) {
	switch rnibError.(type) {
	case *common.ResourceNotFoundError:
		return http.StatusNotFound, notFoundErrorCode, notFoundErrorMessage
	case *common.InternalError:
		return http.StatusInternalServerError, internalErrorCode, internalErrorMessage
	case *common.ValidationError:
		return http.StatusBadRequest, validationErrorCode, validationFailedMessage
	default:
		return http.StatusInternalServerError, internalErrorCode, internalErrorMessage
	}
}
