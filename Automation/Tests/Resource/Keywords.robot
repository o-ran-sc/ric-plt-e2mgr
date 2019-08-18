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
Resource   ../Resource/resource.robot
Library     OperatingSystem





*** Keywords ***
Post Request setup node b x-2
    Set Headers     ${header}
    POST        /v1/nodeb/x2-setup    ${json}



Get Request node b enb
    Sleep    1s
    GET      /v1/nodeb/test1



Post Request setup node b endc-setup
    Set Headers     ${header}
    POST        /v1/nodeb/endc-setup    ${endcjson}




Start dockers
     Run And Return Rc And Output    ${run_script}
     ${result}=  Run And Return Rc And Output     ${docker_command}
     Should Be Equal As Integers    ${result[1]}    5

