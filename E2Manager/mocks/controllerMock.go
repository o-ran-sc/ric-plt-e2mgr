package mocks

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type ControllerMock struct {
	mock.Mock
}

func (c *ControllerMock) ShutdownHandler(writer http.ResponseWriter, r *http.Request){
	c.Called()
}

func (c *ControllerMock) X2ResetHandler(writer http.ResponseWriter, r *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	ranName := vars["ranName"]

	writer.Write([]byte(ranName))

	c.Called()
}

