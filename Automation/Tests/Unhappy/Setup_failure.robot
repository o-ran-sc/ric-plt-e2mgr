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
Suite Setup     Start dockers
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST        ${url}




*** Test Cases ***
Post Request setup node b x2-setup - setup failure
    Set Headers     ${header}
    POST        /v1/nodeb/x2-setup   ${json}
    Sleep    1s
    POST        /v1/nodeb/x2-setup   ${json}
    Sleep    1s
    GET      /v1/nodeb/test1
    Integer    response status       200
    String     response body connectionStatus     CONNECTED_SETUP_FAILED
    String     response body failureType     X2_SETUP_FAILURE
    String     response body setupFailure networkLayerCause       HO_TARGET_NOT_ALLOWED




*** Keywords ***
Start dockers
     Run And Return Rc And Output    ${run_script}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    5