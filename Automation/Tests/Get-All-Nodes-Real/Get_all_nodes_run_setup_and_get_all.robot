##############################################################################
#
#   Copyright (c) 2019 AT&T Intellectual Property.
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
##############################################################################
#
#   This source code is part of the near-RT RIC (RAN Intelligent Controller)
#   platform project (RICP).
#

*** Settings ***
Suite Setup  Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     REST      ${url}


*** Test Cases ***
Run x2 setup
    Post Request setup node b x-2
    Integer     response status       204
    Sleep  2s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    Integer  response body port     5577
    String   response body connectionStatus    CONNECTED

Run endc setup
    Post Request setup node b endc-setup
    Integer     response status       204
    Sleep  2s
    GET      /v1/nodeb/test2
    Integer  response status  200
    String   response body ranName    test2
    String   response body connectionStatus    CONNECTED


Get all node ids
    GET     v1/nodeb/ids
    Sleep  2s
    Integer  response status   200
    String   response body 0 inventoryName  test1
    String   response body 0 globalNbId plmnId   02f829
    String   response body 0 globalNbId nbId     007ab0
    String   response body 1 inventoryName  test2
    String   response body 1 globalNbId plmnId   42f490
    String   response body 1 globalNbId nbId     000004








