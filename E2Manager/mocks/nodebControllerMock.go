package mocks

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type NodebControllerMock struct {
	mock.Mock
}

func (rc *NodebControllerMock) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	rc.Called()
}

func (rc *NodebControllerMock) GetNodebIdList (writer http.ResponseWriter, request *http.Request) {
	rc.Called()
}

func (rc *NodebControllerMock) GetNodeb(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	vars := mux.Vars(request)
	ranName := vars["ranName"]

	writer.Write([]byte(ranName))

	rc.Called()
}

func (rc *NodebControllerMock) HandleHealthCheckRequest(writer http.ResponseWriter, request *http.Request) {
	rc.Called()
}
