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
	"e2mgr/handlers"
	"e2mgr/logger"
	"e2mgr/mocks"
	"e2mgr/rNibWriter"
	"testing"
)

const x2SetupRequestType = "x2-setup"
const endcSetupRequestType = "endc-setup"

/*
 * Verify consturctor.
 */
func TestNewRequestHandlerProvider(t *testing.T) {

	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}

	provider := NewRequestHandlerProvider(rnibWriterProvider)
	if provider == nil {
		t.Errorf("want: provider, got: nil")
	}
}

/*
 * Verify support for known providers.
 */

func TestGetX2SetupRequestHandler(t *testing.T) {

	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}

	provider := NewRequestHandlerProvider(rnibWriterProvider)

	handler, err := provider.GetHandler(log, x2SetupRequestType)

	if handler == nil {
		t.Errorf("failed to get x2 setup handler")
	}

	_, ok := handler.(*handlers.SetupRequestHandler)

	if !ok {
		t.Errorf("failed to get x2 setup handler")
	}
}

func TestGetEndcSetupRequestHandler(t *testing.T) {

	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}

	provider := NewRequestHandlerProvider(rnibWriterProvider)

	handler, err := provider.GetHandler(log, endcSetupRequestType)

	if handler == nil {
		t.Errorf("failed to get endc setup handler")
	}

	_, ok := handler.(*handlers.EndcSetupRequestHandler)

	if !ok {
		t.Errorf("failed to get endc setup handler")
	}
}

/*
 * Verify handling of a request for an unsupported request.
 */

func TestGetHandlerFailure(t *testing.T) {
	log, err := logger.InitLogger(logger.InfoLevel)
	if err != nil {
		t.Errorf("failed to initialize logger, error: %s", err)
	}

	rnibWriterProvider := func() rNibWriter.RNibWriter {
		return &mocks.RnibWriterMock{}
	}

	provider := NewRequestHandlerProvider(rnibWriterProvider)

	_, err = provider.GetHandler(log, "dummy")

	if err == nil {
		t.Errorf("Provider should had respond with error for dummy request type")
	}
}
