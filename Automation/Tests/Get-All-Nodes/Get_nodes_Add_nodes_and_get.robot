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
Suite Setup  Flush Redis
Library      Process
Resource   ../Resource/resource.robot
Library     OperatingSystem
Library     REST      ${url}


*** Variables ***
${file}     ${CURDIR}/addtoredis.py
${flush_file}     ${CURDIR}/flush.py

*** Test Cases ***
Add nodes to redis db
    ${result}=  Run Process    python3.6    ${file}
    Should Be Equal As Strings      ${result.stdout}        Insert successfully to Redis


Get all node ids
    GET     v1/nodeb-ids
    #Output
    Integer  response status   200
    String   response body 0 inventoryName  test1
    String   response body 0 globalNbId plmnId   02f829
    String   response body 0 globalNbId nbId     007a80
    String   response body 1 inventoryName  test2
    String   response body 1 globalNbId plmnId   03f829
    String   response body 1 globalNbId nbId     001234



*** Keywords ***
Flush Redis
    ${result}=  Run Process    python3.6    ${flush_file}
    Should Be Equal As Strings      ${result.stdout}        Flush Success







