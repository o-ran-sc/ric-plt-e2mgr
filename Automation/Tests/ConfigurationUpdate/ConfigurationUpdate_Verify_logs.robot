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

*** Variables ***
${verifylogsscript}     ${CURDIR}/verifylogs.py


*** Test Cases ***
Verify logs - Confiugration update - Begin Tag Get
    ${Configuration}=   Grep File  ./gnb.log  <ENDCConfigurationUpdate>
    #Log to console      ${Configuration}
    ${ConfigurationAfterStrip}=     Strip String    ${Configuration}
    Should Be Equal     ${ConfigurationAfterStrip}        <ENDCConfigurationUpdate>

Verify logs - Confiugration update - End Tag Get
    ${ConfigurationEnd}=   Grep File  ./gnb.log  </ENDCConfigurationUpdate>
    #Log to console      ${ConfigurationEnd}
    ${ConfigurationEndAfterStrip}=     Strip String    ${ConfigurationEnd}
    Should Be Equal     ${ConfigurationEndAfterStrip}        </ENDCConfigurationUpdate>

Verify logs - Confiugration update - Ack Tag Begin
    ${ConfigurationAck}=   Grep File  ./gnb.log   <ENDCConfigurationUpdateAcknowledge>
    #Log to console      ${ConfigurationEnd}
    ${ConfigurationAckAfter}=     Strip String    ${ConfigurationAck}
    Should Be Equal     ${ConfigurationAckAfter}        <ENDCConfigurationUpdateAcknowledge>

Verify logs - Confiugration update - Ack Tag End
    ${ConfigurationAckEnd}=   Grep File  ./gnb.log  </ENDCConfigurationUpdateAcknowledge>
    #Log to console      ${ConfigurationEnd}
    ${ConfigurationAckEndAfterStrip}=     Strip String    ${ConfigurationAckEnd}
    Should Be Equal     ${ConfigurationAckEndAfterStrip}        </ENDCConfigurationUpdateAcknowledge>

Verify logs - e2mgr logs
   ${result}=  Run Process    python3.6    ${verifylogsscript}  ${EXECDIR}
   Should Be Equal As Strings    ${result.stdout}       Found All Configuration Update logs