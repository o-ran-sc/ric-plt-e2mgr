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
Resource    ../Resource/scripts_variables.robot
Library     String
Library     OperatingSystem
Library     Process
Library     ../Scripts/find_rmr_message.py



*** Test Cases ***
Verify logs - Confiugration update - Begin Tag Get
    ${Configuration}=   Grep File  ./gnb.log  <ENDCConfigurationUpdate>
    ${ConfigurationAfterStrip}=     Strip String    ${Configuration}
    Should Be Equal     ${ConfigurationAfterStrip}        <ENDCConfigurationUpdate>

Verify logs - Confiugration update - End Tag Get
    ${ConfigurationEnd}=   Grep File  ./gnb.log  </ENDCConfigurationUpdate>
    ${ConfigurationEndAfterStrip}=     Strip String    ${ConfigurationEnd}
    Should Be Equal     ${ConfigurationEndAfterStrip}        </ENDCConfigurationUpdate>

Verify logs - Confiugration update - Ack Tag Begin
    ${ConfigurationAck}=   Grep File  ./gnb.log   <ENDCConfigurationUpdateAcknowledge>
    ${ConfigurationAckAfter}=     Strip String    ${ConfigurationAck}
    Should Be Equal     ${ConfigurationAckAfter}        <ENDCConfigurationUpdateAcknowledge>

Verify logs - Confiugration update - Ack Tag End
    ${ConfigurationAckEnd}=   Grep File  ./gnb.log  </ENDCConfigurationUpdateAcknowledge>
    ${ConfigurationAckEndAfterStrip}=     Strip String    ${ConfigurationAckEnd}
    Should Be Equal     ${ConfigurationAckEndAfterStrip}        </ENDCConfigurationUpdateAcknowledge>

Verify logs - find RIC_ENDC_CONF_UPDATE
   ${result}   find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${configurationupdate_message_type}  ${Meid_test1}
   Should Be Equal As Strings    ${result}      True
Verify logs - find RIC_ENDC_CONF_UPDATE_ACK
   ${result1}  find_rmr_message.verify_logs  ${EXECDIR}  ${e2mgr_log_filename}  ${configurationupdate_ack_message_type}  ${Meid_test1}
   Should Be Equal As Strings    ${result1}      True
