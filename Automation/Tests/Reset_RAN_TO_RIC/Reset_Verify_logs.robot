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
Library     String
Library     OperatingSystem
Library     Process
Library     ${CURDIR}/Reset_Ran_To_Ric_Verify_logs.py


*** Test Cases ***
Verify logs - Reset Sent by simulator
    ${Reset}=   Grep File  ./gnb.log  ResetRequest has been sent
    #Log to console      ${Reset}
    Should Be Equal     ${Reset}     gnbe2_simu: ResetRequest has been sent

Verify logs - e2mgr logs
   ${result}    Reset_Ran_To_Ric_Verify_logs.verify   ${EXECDIR}
   Should Be Equal As Strings    ${result}      True