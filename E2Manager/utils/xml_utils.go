//
// Copyright (c) 2020 Samsung Electronics Co., Ltd. All Rights Reserved.
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

package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func CleanXML(payload []byte) []byte {
	xmlStr := string(payload)
	normalized := strings.NewReplacer("\r","","\n",""," ","").Replace(xmlStr)

	return []byte(normalized)
}

func ReadXmlFile(t *testing.T, xmlPath string) []byte {
	path, err := filepath.Abs(xmlPath)
	if err != nil {
		t.Fatal(err)
	}
	xmlAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	return xmlAsBytes
}

func ReadXmlFileNoTest(xmlPath string) ([]byte, error) {
	path, err := filepath.Abs(xmlPath)
	if err != nil {
		return nil,err
	}
	xmlAsBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil,err
	}

	return xmlAsBytes,nil
}

func NormalizeXml(payload []byte) []byte {
	xmlStr := string(payload)
	normalized := strings.NewReplacer("&lt;", "<", "&gt;", ">").Replace(xmlStr)

	return []byte(normalized)
}

func ReplaceEmptyTagsWithSelfClosing(responsePayload []byte, emptyTagsToReplace []string) []byte {
	emptyTagVsSelfClosingTagPairs := make([]string, len(emptyTagsToReplace)*2)
	j := 0

	for i := 0; i < len(emptyTagsToReplace); i++ {
		emptyTagVsSelfClosingTagPairs[j] = fmt.Sprintf("<%[1]s></%[1]s>", emptyTagsToReplace[i])
		emptyTagVsSelfClosingTagPairs[j+1] = fmt.Sprintf("<%s/>", emptyTagsToReplace[i])
		j += 2
	}
	responseString := strings.NewReplacer(emptyTagVsSelfClosingTagPairs...).Replace(string(responsePayload))
	return []byte(responseString)
}
