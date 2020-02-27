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
Library     ../Scripts/e2mdbscripts.py
Library     OperatingSystem
Library    Collections
Library     REST      ${url}


*** Test Cases ***

Pre Condition for Connecting - no simu
    Run And Return Rc And Output    ${stop_simu}
    ${result}=  Run And Return Rc And Output     ${docker_command}
    Should Be Equal As Integers    ${result[1]}    ${docker_number-1}


Setup Ran
    Sleep  1s
    Post Request setup node b x-2
    Integer     response status       204

Verify connection status is DISCONNECTED and RAN is not associated with E2T instance
    Sleep    5s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    Missing  response body associatedE2tInstanceAddress
    String   response body connectionStatus    DISCONNECTED

Verify E2T instance is NOT associated with RAN
   ${result}    e2mdbscripts.verify_ran_is_associated_with_e2t_instance     test1    e2t.att.com:38000
   Should Be True    ${result} == False
