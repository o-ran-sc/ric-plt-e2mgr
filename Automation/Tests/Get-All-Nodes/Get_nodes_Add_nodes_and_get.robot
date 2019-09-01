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
Library     ${CURDIR}/scripts.py
Library     REST      ${url}




*** Test Cases ***
Add nodes to redis db
    ${result}   scripts.add
    Should Be Equal As Strings  ${result}  True


Get all node ids
    GET     v1/nodeb/ids
    Integer  response status   200
    String   response body 0 inventoryName  test1
    String   response body 0 globalNbId plmnId   02f829
    String   response body 0 globalNbId nbId     007a80
    String   response body 1 inventoryName  test2
    String   response body 1 globalNbId plmnId   03f829
    String   response body 1 globalNbId nbId     001234



*** Keywords ***
Flush Redis
    ${result}   scripts.flush
    Should Be Equal As Strings  ${result}  True
