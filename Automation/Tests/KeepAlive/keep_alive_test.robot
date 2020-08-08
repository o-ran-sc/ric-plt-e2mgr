robot##############################################################################
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
Variables  ../Scripts/variables.py
Suite Setup   Prepare Enviorment    ${True}    ${False}
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library    ../Scripts/find_error_script.py
Library     ../Scripts/find_rmr_message.py
Library     ../Scripts/e2mdbscripts.py
Library    OperatingSystem
Library    Collections
Library     REST      ${url}

*** Variables ***
${url}  ${e2mgr_address}

*** Test Cases ***

Get request gnb
    Sleep    2s
    Get Request nodeb
    Integer  response status  200
    String   response body ranName    ${ranname}
    String   response body connectionStatus    CONNECTED
    String   response body nodeType     GNB
    String   response body associatedE2tInstanceAddress  ${e2t_alpha_address}
    Integer  response body gnb ranFunctions 0 ranFunctionId  1
    Integer  response body gnb ranFunctions 0 ranFunctionRevision  1
    Integer  response body gnb ranFunctions 1 ranFunctionId  2
    Integer  response body gnb ranFunctions 1 ranFunctionRevision  1
    Integer  response body gnb ranFunctions 2 ranFunctionId  3
    Integer  response body gnb ranFunctions 2 ranFunctionRevision  1


Verify RAN is associated with E2T instance
   ${result}    e2mdbscripts.verify_ran_is_associated_with_e2t_instance      ${ranname}    ${e2t_alpha_address}
   Should Be True    ${result}

Stop E2T
    Stop E2

Prepare logs
    Remove log files
    Save logs

Verify RAN is not associated with E2T instance
    Sleep  8m
    Get Request nodeb
    Integer  response status  200
    String   response body ranName    ${ranname}
    Missing  response body associatedE2tInstanceAddress
    String   response body connectionStatus    DISCONNECTED

Verify E2T instance removed from db
    ${result}    e2mdbscripts.verify_e2t_instance_key_exists     ${e2t_alpha_address}
    Should Be True    ${result} == False

    ${result}    e2mdbscripts.verify_e2t_instance_exists_in_addresses     ${e2t_alpha_address}
    Should Be True    ${result} == False


[Teardown]    Run Keywords
              Start E2
              AND wait until keyword succeeds  2 min    10 sec    Validate Required Dockers