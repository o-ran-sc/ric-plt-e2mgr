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

Resource   ../Resource/resource.robot
Resource   ../Resource/Keywords.robot
Library     OperatingSystem
Library     REST      ${url}

*** Test Cases ***
X2 - Setup Test
    Post Request setup node b x-2
    Integer     response status       200

X2 - Get Nodeb
    Get Request node b enb test1
    Integer  response status  200
    String   response body ranName    test1
    String   response body ip    10.0.2.15
    Integer  response body port     5577
    String   response body connectionStatus    CONNECTED
    String   response body nodeType     ENB
    String   response body enb enbType     MACRO_ENB
    Integer  response body enb servedCells 0 pci  99
    String   response body enb servedCells 0 cellId   02f829:0007ab00
    String   response body enb servedCells 0 tac    0102
    String   response body enb servedCells 0 broadcastPlmns 0   "02f829"
    Integer  response body enb servedCells 0 choiceEutraMode fdd ulearFcn    1
    Integer  response body enb servedCells 0 choiceEutraMode fdd dlearFcn    1
    String   response body enb servedCells 0 choiceEutraMode fdd ulTransmissionBandwidth   BW50
    String   response body enb servedCells 0 choiceEutraMode fdd dlTransmissionBandwidth   BW50





