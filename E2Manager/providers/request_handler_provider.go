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

package providers

import (
	"e2mgr/handlers"
	"e2mgr/logger"
	"e2mgr/rNibWriter"
	"errors"
	"fmt"
)

var requestMap map[string]handlers.Handler

type RequestHandlerProvider struct{}

func NewRequestHandlerProvider(rnibWriterProvider func() rNibWriter.RNibWriter) *RequestHandlerProvider {
	requestMap = initRequestMap(rnibWriterProvider)
	return &RequestHandlerProvider{}
}

func initRequestMap(rnibWriterProvider func() rNibWriter.RNibWriter) map[string]handlers.Handler {
	return map[string]handlers.Handler{
		"x2-setup":   handlers.NewSetupRequestHandler(rnibWriterProvider),
		"endc-setup": handlers.NewEndcSetupRequestHandler(rnibWriterProvider),
	}
}

func (provider RequestHandlerProvider) GetHandler(logger *logger.Logger, requestType string) (handlers.Handler, error) {
	handler, ok := requestMap[requestType]

	if !ok {
		errorMessage := fmt.Sprintf("#request_handler_provider.GetHandler - Cannot find handler for request type: %s", requestType)
		logger.Errorf(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return handler, nil
}
