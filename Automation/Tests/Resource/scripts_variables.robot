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
Documentation    Message types resource file


*** Variables ***
${configurationupdate_message_type}    MType: 10370
${Meid_test1}   Meid: \\"test1\\"
${Meid_test2}   Meid: \\"test2\\"
${configurationupdate_ack_message_type}    MType: 10371
${RAN_CONNECTED_message_type}     MType: 1200
${RAN_RESTARTED_message_type}     MType: 1210
${RIC_X2_RESET_REQ_message_type}    MType: 10070
${RIC_X2_RESET_RESP_message_type}    MType: 10070
${failed_to_retrieve_nodeb_message}     failed to retrieve nodeB entity. RanName: test1.
${first_retry_to_retrieve_from_db}      RnibDataService.retry - retrying 1 GetNodeb
${third_retry_to_retrieve_from_db}      RnibDataService.retry - after 3 attempts of GetNodeb
${RIC_RES_STATUS_REQ_message_type_successfully_sent}     Message type: 10090 - Successfully sent RMR message
${RAN_NAME_test1}    RAN name: test1
${RAN_NAME_test2}    RAN name: test2
${E2ADAPTER_Setup_Resp}    Send dummy ENDCX2SetupResponse to RIC

