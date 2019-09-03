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
Library     ${CURDIR}/scripts.py
Resource   ../Resource/resource.robot
Library     OperatingSystem


*** Keywords ***
Post Request setup node b x-2
    Set Headers     ${header}
    POST        /v1/nodeb/x2-setup    ${json}


Get Request node b enb test1
    Sleep    1s
    GET      /v1/nodeb/test1


Get Request node b enb test2
    Sleep    1s
    GET      /v1/nodeb/test2

Post Request setup node b endc-setup
    Set Headers     ${header}
    POST        /v1/nodeb/endc-setup    ${endcjson}


Prepare Simulator For Load Information
     Run And Return Rc And Output    ${stop_simu}
     Run And Return Rc And Output    ${docker_Remove}
     ${flush}  scripts.flush
     Should Be Equal As Strings  ${flush}  True
     Run And Return Rc And Output    ${run_simu_load}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    5

Prepare Enviorment
     ${flush}  scripts.flush
     Should Be Equal As Strings  ${flush}  True
     Run And Return Rc And Output    ${stop_simu}
     Run And Return Rc And Output    ${docker_Remove}
     Run And Return Rc And Output    ${run_simu_regular}
     Run And Return Rc And Output    ${restart_e2adapter}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    5

Start E2
     Run And Return Rc And Output    ${start_e2}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    5
     Sleep  2s

Start Redis
     Run And Return Rc And Output    ${redis_remove}
     Run And Return Rc And Output    ${start_redis}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    5





