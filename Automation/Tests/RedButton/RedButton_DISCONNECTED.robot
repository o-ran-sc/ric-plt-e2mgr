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
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Resource    red_button_keywords.robot
Library     OperatingSystem
Library    Collections
Library     REST      ${url}

*** Variables ***
${stop_docker_e2}      docker stop e2adapter

*** Test Cases ***

Pre Condition for Connecting - no E2ADAPTER
    Run And Return Rc And Output    ${stop_docker_e2}
    ${result}=  Run And Return Rc And Output     ${docker_command}
    Should Be Equal As Integers    ${result[1]}    ${docker_number-1}

Prepare Ran in Connecting connectionStatus
    Post Request setup node b endc-setup
    Integer     response status       204
    Sleep  10s
    GET      /v1/nodeb/test2
    Integer  response status  200
    String   response body ranName    test2
    String   response body connectionStatus    DISCONNECTED

Execute Shutdown
    Execute Shutdown

Verify nodeb's connection status is SHUT_DOWN and it's NOT associated to an e2t instance
    Verify shutdown for gnb