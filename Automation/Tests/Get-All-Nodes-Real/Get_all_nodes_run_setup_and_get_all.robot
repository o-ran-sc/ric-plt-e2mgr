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
Variables  ../Scripts/variables.py
Suite Setup   Prepare Enviorment     ${True}
Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     REST      ${url}

*** Variables ***
${url}  ${e2mgr_address}


*** Test Cases ***

Add eNb Node
    Sleep   2s
    Add eNb Request


Get all node ids
    &{res}=   GET     v1/nodeb/states
    Sleep  2s
    Integer  response status   200

    ${is_enb_first}=    set variable if  '${enb_ran_name}'=='${res.body[0].inventoryName}'   True    False

    run keyword if  ${is_enb_first}    RUN KEYWORDS
...      String   response body 1 inventoryName    ${ranName}
...      AND   String   response body 1 globalNbId plmnId   02F829
...      AND   String   response body 1 globalNbId nbId     001100000011000000110000
...      AND   String   response body 0 inventoryName    ${enb_ran_name}
...      AND   String   response body 0 connectionStatus    DISCONNECTED
...      AND   String   response body 0 globalNbId plmnId   def
...      AND   String   response body 0 globalNbId nbId     abc
...      AND   Log To Console    enb index is 0 - all rans were verified successfully

...  ELSE     RUN KEYWORDS
...      String   response body 0 inventoryName    ${ranName}
...      AND   String   response body 0 globalNbId plmnId   02F829
...      AND   String   response body 0 globalNbId nbId     001100000011000000110000
...      AND   String   response body 1 inventoryName    ${enb_ran_name}
...      AND   String   response body 1 connectionStatus    DISCONNECTED
...      AND   String   response body 1 globalNbId plmnId   def
...      AND   String   response body 1 globalNbId nbId     abc
...      AND   Log To Console    enb index is 1 - all rans were verified successfully

Prepare Logs For Tests
    Remove log files
    Save logs






