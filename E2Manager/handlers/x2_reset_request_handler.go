//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package handlers
// #cgo CFLAGS: -I../asn1codec/inc/  -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <x2reset_request_wrapper.h>
import "C"
import (
	"e2mgr/configuration"
	"e2mgr/e2managererrors"
	"e2mgr/logger"
	"e2mgr/rNibWriter"
	"e2mgr/rmrCgo"
	"e2mgr/services"
	"e2mgr/sessions"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/common"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/reader"
	"unsafe"

	"e2mgr/models"
)

type X2ResetRequestHandler struct {
	readerProvider func() reader.RNibReader
	writerProvider func() rNibWriter.RNibWriter
	rmrService *services.RmrService
	config         *configuration.Configuration
}

/*

C
*/

type cause struct {
causeGroup uint32
cause      int
}

var knownCauses = map[string] cause  {
"misc:control-processing-overload": {causeGroup:C.Cause_PR_misc,  cause: C.CauseMisc_control_processing_overload},
"misc:hardware-failure": {causeGroup:C.Cause_PR_misc,  cause: C.CauseMisc_hardware_failure},
"misc:om-intervention": {causeGroup:C.Cause_PR_misc,  cause: C.CauseMisc_om_intervention},
"misc:not-enough-user-plane-processing-resources": {causeGroup:C.Cause_PR_misc,  cause: C.CauseMisc_not_enough_user_plane_processing_resources},
"misc:unspecified": {causeGroup:C.Cause_PR_misc,  cause: C.CauseMisc_unspecified},

"protocol:transfer-syntax-error": {causeGroup:C.Cause_PR_protocol,  cause: C.CauseProtocol_transfer_syntax_error},
"protocol:abstract-syntax-error-reject": {causeGroup:C.Cause_PR_protocol,  cause: C.CauseProtocol_abstract_syntax_error_reject},
"protocol:abstract-syntax-error-ignore-and-notify": {causeGroup:C.Cause_PR_protocol,  cause: C.CauseProtocol_abstract_syntax_error_ignore_and_notify},
"protocol:message-not-compatible-with-receiver-state": {causeGroup:C.Cause_PR_protocol,  cause: C.CauseProtocol_message_not_compatible_with_receiver_state},
"protocol:semantic-error": {causeGroup:C.Cause_PR_protocol,  cause: C.CauseProtocol_semantic_error},
"protocol:unspecified": {causeGroup:C.Cause_PR_protocol,  cause: C.CauseProtocol_unspecified},
"protocol:abstract-syntax-error-falsely-constructed-message": {causeGroup:C.Cause_PR_protocol,  cause: C.CauseProtocol_abstract_syntax_error_falsely_constructed_message},

"transport:transport-resource-unavailable": {causeGroup:C.Cause_PR_transport,  cause: C.CauseTransport_transport_resource_unavailable},
"transport:unspecified":{causeGroup:C.Cause_PR_transport,  cause: C.CauseTransport_unspecified},

"radioNetwork:handover-desirable-for-radio-reasons": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_handover_desirable_for_radio_reasons},
"radioNetwork:time-critical-handover": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_time_critical_handover},
"radioNetwork:resource-optimisation-handover": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_resource_optimisation_handover},
"radioNetwork:reduce-load-in-serving-cell": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_reduce_load_in_serving_cell},
"radioNetwork:partial-handover": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_partial_handover},
"radioNetwork:unknown-new-enb-ue-x2ap-id": {causeGroup:C.Cause_PR_radioNetwork,  cause:C.CauseRadioNetwork_unknown_new_eNB_UE_X2AP_ID},
"radioNetwork:unknown-old-enb-ue-x2ap-id": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_unknown_old_eNB_UE_X2AP_ID},
"radioNetwork:unknown-pair-of-ue-x2ap-id": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_unknown_pair_of_UE_X2AP_ID},
"radioNetwork:ho-target-not-allowed": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_ho_target_not_allowed},
"radioNetwork:tx2relocoverall-expiry": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_tx2relocoverall_expiry},
"radioNetwork:trelocprep-expiry": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_trelocprep_expiry},
"radioNetwork:cell-not-available": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_cell_not_available},
"radioNetwork:no-radio-resources-available-in-target-cell": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_no_radio_resources_available_in_target_cell},
"radioNetwork:invalid-mme-groupid": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_invalid_MME_GroupID},
"radioNetwork:unknown-mme-code": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_unknown_MME_Code},
"radioNetwork:encryption-and-or-integrity-protection-algorithms-not-supported": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_encryption_and_or_integrity_protection_algorithms_not_supported},
"radioNetwork:reportcharacteristicsempty": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_reportCharacteristicsEmpty},
"radioNetwork:noreportperiodicity": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_noReportPeriodicity},
"radioNetwork:existingMeasurementID": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_existingMeasurementID},
"radioNetwork:unknown-enb-measurement-id": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_unknown_eNB_Measurement_ID},
"radioNetwork:measurement-temporarily-not-available": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_measurement_temporarily_not_available},
"radioNetwork:unspecified": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_unspecified},
"radioNetwork:load-balancing": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_load_balancing},
"radioNetwork:handover-optimisation": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_handover_optimisation},
"radioNetwork:value-out-of-allowed-range": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_value_out_of_allowed_range},
"radioNetwork:multiple-E-RAB-ID-instances": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_multiple_E_RAB_ID_instances},
"radioNetwork:switch-off-ongoing": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_switch_off_ongoing},
"radioNetwork:not-supported-qci-value": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_not_supported_QCI_value},
"radioNetwork:measurement-not-supported-for-the-object": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_measurement_not_supported_for_the_object},
"radioNetwork:tdcoverall-expiry": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_tDCoverall_expiry},
"radioNetwork:tdcprep-expiry": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_tDCprep_expiry},
"radioNetwork:action-desirable-for-radio-reasons": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_action_desirable_for_radio_reasons},
"radioNetwork:reduce-load": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_reduce_load},
"radioNetwork:resource-optimisation": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_resource_optimisation},
"radioNetwork:time-critical-action": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_time_critical_action},
"radioNetwork:target-not-allowed": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_target_not_allowed},
"radioNetwork:no-radio-resources-available": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_no_radio_resources_available},
"radioNetwork:invalid-qos-combination": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_invalid_QoS_combination},
"radioNetwork:encryption-algorithms-not-aupported": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_encryption_algorithms_not_aupported},
"radioNetwork:procedure-cancelled":{causeGroup:C.Cause_PR_radioNetwork,  cause:C.CauseRadioNetwork_procedure_cancelled},
"radioNetwork:rrm-purpose": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_rRM_purpose},
"radioNetwork:improve-user-bit-rate": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_improve_user_bit_rate},
"radioNetwork:user-inactivity": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_user_inactivity},
"radioNetwork:radio-connection-with-ue-lost": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_radio_connection_with_UE_lost},
"radioNetwork:failure-in-the-radio-interface-procedure": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_failure_in_the_radio_interface_procedure},
"radioNetwork:bearer-option-not-supported": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_bearer_option_not_supported},
"radioNetwork:mcg-mobility": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_mCG_Mobility},
"radioNetwork:scg-mobility": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_sCG_Mobility},
"radioNetwork:count-reaches-max-value": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_count_reaches_max_value},
"radioNetwork:unknown-old-en-gnb-ue-x2ap-id": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_unknown_old_en_gNB_UE_X2AP_ID} ,
"radioNetwork:pdcp-Overload": {causeGroup:C.Cause_PR_radioNetwork,  cause: C.CauseRadioNetwork_pDCP_Overload},
}

func NewX2ResetRequestHandler(rmrService *services.RmrService, config *configuration.Configuration, writerProvider func() rNibWriter.RNibWriter,
	readerProvider func() reader.RNibReader) *X2ResetRequestHandler {
	return &X2ResetRequestHandler{
		readerProvider: readerProvider,
		writerProvider: writerProvider,
		rmrService: rmrService,
		config:         config,
	}
}

func (handler *X2ResetRequestHandler) Handle(logger *logger.Logger, request models.Request) error {
	resetRequest := request.(models.ResetRequest)

	if len(resetRequest.Cause) == 0 {
		resetRequest.Cause = "misc:om-intervention"
	}
	cause, ok:= knownCauses[resetRequest.Cause]
	if !ok {
		logger.Errorf("#reset_request_handler.Handle - Unknown cause (%s)", resetRequest.Cause)
		return e2managererrors.NewRequestValidationError()
	}

	nodeb, err  := handler.readerProvider().GetNodeb(resetRequest.RanName)
	if err != nil {
		logger.Errorf("#reset_request_handler.Handle - failed to get status of RAN: %s from RNIB. Error: %s", resetRequest.RanName,  err.Error())
		if err.GetCode() == common.RESOURCE_NOT_FOUND {
			return e2managererrors.NewResourceNotFoundError()
		}
		return e2managererrors.NewRnibDbError()
	}

	if nodeb.ConnectionStatus != entities.ConnectionStatus_CONNECTED {
		logger.Errorf("#reset_request_handler.Handle - RAN: %s in wrong state (%s)", resetRequest.RanName, entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)])
		return e2managererrors.NewWrongStateError(entities.ConnectionStatus_name[int32(nodeb.ConnectionStatus)])
	}

	var payloadSize  = C.ulong(MaxAsn1PackedBufferSize)
	packedBuffer := [MaxAsn1PackedBufferSize]C.uchar{}
	errorBuffer := [MaxAsn1CodecMessageBufferSize]C.char{}

	if status := C.build_pack_x2reset_request(cause.causeGroup, C.int(cause.cause), &payloadSize, &packedBuffer[0], MaxAsn1CodecMessageBufferSize, &errorBuffer[0]); !status {
		logger.Errorf("#reset_request_handler.Handle - failed to build and pack the reset message %s ", C.GoString(&errorBuffer[0]))
		return  e2managererrors.NewInternalError()
	}
	transactionId := resetRequest.RanName
	handler.rmrService.E2sessions[transactionId] = sessions.E2SessionDetails{SessionStart: resetRequest.StartTime, Request: &models.RequestDetails{RanName: resetRequest.RanName}}
	response := models.NotificationResponse{MgsType: rmrCgo.RIC_X2_RESET, RanName: resetRequest.RanName, Payload: C.GoBytes(unsafe.Pointer(&packedBuffer[0]), C.int(payloadSize))}
	if err:= handler.rmrService.SendRmrMessage(&response); err != nil {
		logger.Errorf("#reset_request_handler.Handle - failed to send reset message to RMR: %s", err)
		return  e2managererrors.NewRmrError()
	}

	logger.Infof("#reset_request_handler.Handle - sent x2 reset to RAN: %s with cause: %s", resetRequest.RanName, resetRequest.Cause)
	return nil
}


