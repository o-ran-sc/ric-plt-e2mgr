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
	"e2mgr/providers/httpmsghandlerprovider"
	"e2mgr/rNibWriter"
	"e2mgr/services"
	"e2mgr/sessions"
	"e2mgr/utils"
	"encoding/json"
	"errors"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"
	"time"
)

const (
	parseErrorCode            int = 401
	validationErrorCode       int = 402
	notFoundErrorCode         int = 404
	internalErrorCode         int = 501
	requiredInputErrorMessage     = "Mandatory fields are missing"
	validationFailedMessage       = "Validation failed"
	parseErrorMessage             = "Parse failure"
	notFoundErrorMessage          = "Resource not found"
	internalErrorMessage          = "Internal Server Error. Please try again later"
	sendMessageErrorMessage       = "Failed to send message. For more information please check logs"
)

var E2Sessions = make(sessions.E2Sessions)

var messageChannel chan *models.E2RequestMessage
var errorChannel chan error

type INodebController interface {
	HandleRequest(writer http.ResponseWriter, request *http.Request)
	GetNodebIdList (writer http.ResponseWriter, request *http.Request)
	GetNodeb(writer http.ResponseWriter, request *http.Request)
	HandleHealthCheckRequest(writer http.ResponseWriter, request *http.Request)
}

type NodebController struct {
	rmrService         *services.RmrService
	Logger             *logger.Logger
	rnibReaderProvider func() reader.RNibReader
	rnibWriterProvider func() rNibWriter.RNibWriter
}

func NewNodebController(logger *logger.Logger, rmrService *services.RmrService, rnibReaderProvider func() reader.RNibReader,
	rnibWriterProvider func() rNibWriter.RNibWriter) *NodebController {
	messageChannel = make(chan *models.E2RequestMessage)
	errorChannel = make(chan error)
	return &NodebController{
		rmrService:         rmrService,
		Logger:             logger,
		rnibReaderProvider: rnibReaderProvider,
		rnibWriterProvider: rnibWriterProvider,
	}
}

func prettifyRequest(request *http.Request) string {
	dump, _ := httputil.DumpRequest(request, true)
	requestPrettyPrint := strings.Replace(string(dump), "\r\n", " ", -1)
	return strings.Replace(requestPrettyPrint, "\n", "", -1)
}

func (rc NodebController) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	rc.Logger.Infof("[Client -> E2 Manager] #nodeb_controller.HandleRequest - request: %v", prettifyRequest(request))

	vars := mux.Vars(request)
	messageTypeParam := vars["messageType"]
	requestHandlerProvider := httpmsghandlerprovider.NewRequestHandlerProvider(rc.rnibWriterProvider)
	handler, err := requestHandlerProvider.GetHandler(rc.Logger, messageTypeParam)

	if err != nil {
		handleErrorResponse(rc.Logger, http.StatusNotFound, notFoundErrorCode, notFoundErrorMessage, writer, startTime)
		return
	}

	requestDetails, err := parseJson(rc.Logger, request)

	if err != nil {
		handleErrorResponse(rc.Logger, http.StatusBadRequest, parseErrorCode, parseErrorMessage, writer, startTime)
		return
	}

	rc.Logger.Infof("#nodeb_controller.HandleRequest - request: %+v", requestDetails)

	if err := validateRequestDetails(rc.Logger, requestDetails); err != nil {
		handleErrorResponse(rc.Logger, http.StatusBadRequest, validationErrorCode, requiredInputErrorMessage, writer, startTime)
		return
	}

	err = handler.PreHandle(rc.Logger, &requestDetails)

	if err != nil {
		handleErrorResponse(rc.Logger, http.StatusInternalServerError, internalErrorCode, err.Error(), writer, startTime)
		return
	}

	rc.Logger.Infof("[E2 Manager -> Client] #nodeb_controller.HandleRequest - http status: 200")
	writer.WriteHeader(http.StatusOK)

	var wg sync.WaitGroup

	go handler.CreateMessage(rc.Logger, &requestDetails, messageChannel, E2Sessions, startTime, wg)

	go rc.rmrService.SendMessage(handler.GetMessageType(), messageChannel, errorChannel, wg)

	wg.Wait()

	err = <-errorChannel

	if err != nil {
		handleErrorResponse(rc.Logger, http.StatusInternalServerError, internalErrorCode, sendMessageErrorMessage, writer, startTime)
		return
	}

	printHandlingRequestElapsedTimeInMs(rc.Logger, startTime)
}

func (rc NodebController) GetNodebIdList (writer http.ResponseWriter, request *http.Request) {
	startTime := time.Now()
	rnibReaderService := services.NewRnibReaderService(rc.rnibReaderProvider)
	nodebIdList, rnibError := rnibReaderService.GetNodebIdList()

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
	// WAS: respondingNode, rnibError := reader.GetRNibReader().GetNodeb(ranName)
	rnibReaderService := services.NewRnibReaderService(rc.rnibReaderProvider);
	respondingNode, rnibError := rnibReaderService.GetNodeb(ranName)
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

func parseJson(logger *logger.Logger, request *http.Request) (models.RequestDetails, error) {
	var requestDetails models.RequestDetails
	if err := json.NewDecoder(request.Body).Decode(&requestDetails); err != nil {
		logger.Errorf("#nodeb_controller.parseJson - cannot deserialize incoming request. request: %v, error: %v", request, err)
		return requestDetails, err
	}
	return requestDetails, nil
}

func validateRequestDetails(logger *logger.Logger, requestDetails models.RequestDetails) error {

	if requestDetails.RanPort == 0 {
		logger.Errorf("#nodeb_controller.validateRequestDetails - validation failure: port cannot be zero")
		return errors.New("port: cannot be blank")
	}
	err := validation.ValidateStruct(&requestDetails,
		validation.Field(&requestDetails.RanIp, validation.Required, is.IP),
		validation.Field(&requestDetails.RanName, validation.Required),
	)
	if err != nil {
		logger.Errorf("#nodeb_controller.validateRequestDetails - validation failure, error: %v", err)
	}

	return err
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
