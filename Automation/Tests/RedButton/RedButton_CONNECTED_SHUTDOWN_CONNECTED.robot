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
Suite Setup   Prepare Enviorment
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     Collections
Library     REST      ${url}


*** Variables ***
${restart_docker_sim}      docker restart gnbe2_simu


*** Test Cases ***

Prepare Ran in Connected connectionStatus
    Post Request setup node b x-2
    Integer     response status       204
    Sleep  1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body connectionStatus    CONNECTED


Disconnect Ran
   PUT    /v1/nodeb/shutdown
   Integer   response status   204



Verfiy Shutdown ConnectionStatus
    Sleep    1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body connectionStatus    SHUT_DOWN

Restart simualtor

    Run And Return Rc And Output    ${restart_docker_sim}
    ${result}=  Run And Return Rc And Output     ${docker_command}
    Should Be Equal As Integers    ${result[1]}    ${docker_number}


repare Ran in Connected connectionStatus
    Post Request setup node b x-2
    Integer     response status       204
    Sleep  1s
    GET      /v1/nodeb/test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body connectionStatus    CONNECTED





