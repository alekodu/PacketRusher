/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package ngap

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/misc"

	"github.com/free5gc/ngap"

	"github.com/free5gc/ngap/ngapType"
	log "github.com/sirupsen/logrus"
)

func Dispatch(amf *context.GNBAmf, gnb *context.GNBContext, message []byte) {

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	if message == nil {
		// TODO return error
		log.WithFields(logFields).Info("NGAP message is nil")
	}

	// decode NGAP message.
	ngapMsg, err := ngap.Decoder(message)
	if err != nil {
		log.WithFields(logFields).Error("Error decoding NGAP message in GNB", ": ", err)
	}

	// check RanUeId and get UE.

	// handle NGAP message.
	switch ngapMsg.Present {

	case ngapType.NGAPPDUPresentInitiatingMessage:

		switch ngapMsg.InitiatingMessage.ProcedureCode.Value {

		case ngapType.ProcedureCodeDownlinkNASTransport:
			// handler NGAP Downlink NAS Transport.
			log.WithFields(logFields).Info("Receive Downlink NAS Transport")
			HandlerDownlinkNasTransport(gnb, ngapMsg)

		case ngapType.ProcedureCodeInitialContextSetup:
			// handler NGAP Initial Context Setup Request.
			log.WithFields(logFields).Info("Receive Initial Context Setup Request")
			HandlerInitialContextSetupRequest(gnb, ngapMsg)

		case ngapType.ProcedureCodePDUSessionResourceSetup:
			// handler NGAP PDU Session Resource Setup Request.
			log.WithFields(logFields).Info("Receive PDU Session Resource Setup Request")
			HandlerPduSessionResourceSetupRequest(gnb, ngapMsg)

		case ngapType.ProcedureCodePDUSessionResourceRelease:
			// handler NGAP PDU Session Resource Release
			log.WithFields(logFields).Info("Receive PDU Session Release Command")
			HandlerPduSessionReleaseCommand(gnb, ngapMsg)

		case ngapType.ProcedureCodeUEContextRelease:
			// handler NGAP UE Context Release
			log.WithFields(logFields).Info("Receive UE Context Release Command")
			HandlerUeContextReleaseCommand(gnb, ngapMsg)

		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			// handler NGAP AMF Configuration Update
			log.WithFields(logFields).Info("Receive AMF Configuration Update")
			HandlerAmfConfigurationUpdate(amf, gnb, ngapMsg)

		case ngapType.ProcedureCodeAMFStatusIndication:
			log.WithFields(logFields).Info("Receive AMF Status Indication")
			HandlerAmfStatusIndication(amf, gnb, ngapMsg)

		case ngapType.ProcedureCodeHandoverResourceAllocation:
			// handler NGAP Handover Request
			log.WithFields(logFields).Info("Receive Handover Request")
			HandlerHandoverRequest(amf, gnb, ngapMsg)

		case ngapType.ProcedureCodePaging:
			// handler NGAP Paging
			log.WithFields(logFields).Info("Receive Paging")
			HandlerPaging(gnb, ngapMsg)

		case ngapType.ProcedureCodeErrorIndication:
			// handler Error Indicator
			log.WithFields(logFields).Error("Receive Error Indication")
			HandlerErrorIndication(gnb, ngapMsg)

		default:
			log.WithFields(logFields).Warnf("Received unknown NGAP message 0x%x", ngapMsg.InitiatingMessage.ProcedureCode.Value)
		}

	case ngapType.NGAPPDUPresentSuccessfulOutcome:

		switch ngapMsg.SuccessfulOutcome.ProcedureCode.Value {

		case ngapType.ProcedureCodeNGSetup:
			// handler NGAP Setup Response.
			log.WithFields(logFields).Info("Receive NG Setup Response")
			HandlerNgSetupResponse(amf, gnb, ngapMsg)

		case ngapType.ProcedureCodePathSwitchRequest:
			// handler PathSwitchRequestAcknowledge
			log.WithFields(logFields).Info("Receive PathSwitchRequestAcknowledge")
			HandlerPathSwitchRequestAcknowledge(gnb, ngapMsg)

		case ngapType.ProcedureCodeHandoverPreparation:
			// handler NGAP AMF Handover Command
			log.WithFields(logFields).Info("Receive Handover Command")
			HandlerHandoverCommand(amf, gnb, ngapMsg)

		default:
			log.WithFields(logFields).Warnf("Received unknown NGAP message 0x%x", ngapMsg.SuccessfulOutcome.ProcedureCode.Value)
		}

	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:

		switch ngapMsg.UnsuccessfulOutcome.ProcedureCode.Value {

		case ngapType.ProcedureCodeNGSetup:
			// handler NGAP Setup Failure.
			log.WithFields(logFields).Info("Receive Ng Setup Failure")
			HandlerNgSetupFailure(amf, gnb, ngapMsg)

		default:
			log.WithFields(logFields).Warnf("Received unknown NGAP message 0x%x", ngapMsg.UnsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}
