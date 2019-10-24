package mocks

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"net/http"
)

type ControllerMock struct {
	mock.Mock
}

func (c *ControllerMock) GetNodeb(writer http.ResponseWriter, r *http.Request){
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		vars := mux.Vars(r)
		ranName := vars["ranName"]

		writer.Write([]byte(ranName))

		c.Called()
}

func (c *ControllerMock) GetNodebIdList(writer http.ResponseWriter, r *http.Request){
	c.Called()
}


func (c *ControllerMock) Shutdown(writer http.ResponseWriter, r *http.Request){
	c.Called()
}

func (c *ControllerMock) X2Reset(writer http.ResponseWriter, r *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	ranName := vars["ranName"]

	writer.Write([]byte(ranName))

	c.Called()
}

func (c *ControllerMock) X2Setup(writer http.ResponseWriter, r *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	c.Called()
}

func (c *ControllerMock) EndcSetup(writer http.ResponseWriter, r *http.Request){
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	c.Called()
}