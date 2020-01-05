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

*** Settings ***
Documentation   Keywords file
Resource    ../Resource/resource.robot
Resource    ../Resource/Keywords.robot
Library     ../Scripts/e2mdbscripts.py
Library    Collections
Library    OperatingSystem
Library    json
Library    REST      ${url}

*** Keywords ***
Execute setup and verify connected and associated
    Post Request setup node b x-2
    Integer  response status       204
    Get Request node b enb test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body connectionStatus    CONNECTED
    String   response body associatedE2tInstanceAddress     e2t.att.com:38000

Execute setup twice and verify connected setup failed
    Post Request setup node b x-2
    Sleep    1s
    Post Request setup node b x-2
    Get Request node b enb test1
    String     response body connectionStatus     CONNECTED_SETUP_FAILED


Verify shutdown for enb
    Get Request node b enb test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body connectionStatus    SHUT_DOWN
    Missing  response body associatedE2tInstanceAddress


Verify shutdown for gnb
    Get Request node b gnb test2
    Integer  response status  200
    String   response body ranName    test2
    String   response body connectionStatus    SHUT_DOWN
    Missing  response body associatedE2tInstanceAddress

Verify E2T instance has no associated RANs
    ${result}    e2mdbscripts.verify_e2t_instance_has_no_associated_rans     e2t.att.com:38000
    Should Be True    ${result}

Execute Shutdown
   PUT       /v1/nodeb/shutdown
   Integer   response status   204





