/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package ngap

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/control_test_engine/procedures"

	"github.com/free5gc/ngap"

	"github.com/free5gc/ngap/ngapType"
	log "github.com/sirupsen/logrus"
)

func Dispatch(amf *context.GNBAmf, gnb *context.GNBContext, message []byte) {

	if message == nil {
		// TODO return error
		log.Info("[GNB][NGAP][](", gnb.GetGnbId(), ")()() NGAP message is nil")
	}

	// decode NGAP message.
	ngapMsg, err := ngap.Decoder(message)
	if err != nil {
		log.Error("[GNB][NGAP][](", gnb.GetGnbId(), ")()() Error decoding NGAP message in GNB", ": ", err)
	}

	// check RanUeId and get UE.

	// handle NGAP message.
	switch ngapMsg.Present {

	case ngapType.NGAPPDUPresentInitiatingMessage:

		switch ngapMsg.InitiatingMessage.ProcedureCode.Value {

		case ngapType.ProcedureCodeDownlinkNASTransport:
			// handler NGAP Downlink NAS Transport.
			log.Info("<", procedures.Registration, "><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive Downlink NAS Transport")
			HandlerDownlinkNasTransport(gnb, ngapMsg)

		case ngapType.ProcedureCodeInitialContextSetup:
			// handler NGAP Initial Context Setup Request.
			log.Info("<", procedures.Registration, "><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive Initial Context Setup Request")
			HandlerInitialContextSetupRequest(gnb, ngapMsg)

		case ngapType.ProcedureCodePDUSessionResourceSetup:
			// handler NGAP PDU Session Resource Setup Request.
			log.Info("<", procedures.NewPDUSession, "><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive PDU Session Resource Setup Request")
			HandlerPduSessionResourceSetupRequest(gnb, ngapMsg)

		case ngapType.ProcedureCodePDUSessionResourceRelease:
			// handler NGAP PDU Session Resource Release
			log.Info("<", procedures.DestroyPDUSession, "><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive PDU Session Release Command")
			HandlerPduSessionReleaseCommand(gnb, ngapMsg)

		case ngapType.ProcedureCodeUEContextRelease:
			// handler NGAP UE Context Release
			log.Info("<", procedures.Deregistration, "><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive UE Context Release Command")
			HandlerUeContextReleaseCommand(gnb, ngapMsg)

		case ngapType.ProcedureCodeAMFConfigurationUpdate:
			// handler NGAP AMF Configuration Update
			log.Info("<><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive AMF Configuration Update")
			HandlerAmfConfigurationUpdate(amf, gnb, ngapMsg)

		case ngapType.ProcedureCodeAMFStatusIndication:
			log.Info("<><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive AMF Status Indication")
			HandlerAmfStatusIndication(amf, gnb, ngapMsg)

		case ngapType.ProcedureCodeHandoverResourceAllocation:
			// handler NGAP Handover Request
			log.Info("<", procedures.NgHandover, "><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive Handover Request")
			HandlerHandoverRequest(amf, gnb, ngapMsg)

		case ngapType.ProcedureCodePaging:
			// handler NGAP Paging
			log.Info("<><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive Paging")
			HandlerPaging(gnb, ngapMsg)

		case ngapType.ProcedureCodeErrorIndication:
			// handler Error Indicator
			log.Error("<><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive Error Indication")
			HandlerErrorIndication(gnb, ngapMsg)

		default:
			log.Warnf("<><>[GNB][NGAP][](", gnb.GetGnbId(), ")()() Received unknown NGAP message 0x%x", ngapMsg.InitiatingMessage.ProcedureCode.Value)
		}

	case ngapType.NGAPPDUPresentSuccessfulOutcome:

		switch ngapMsg.SuccessfulOutcome.ProcedureCode.Value {

		case ngapType.ProcedureCodeNGSetup:
			// handler NGAP Setup Response.
			log.Info("[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive NG Setup Response")
			HandlerNgSetupResponse(amf, gnb, ngapMsg)

		case ngapType.ProcedureCodePathSwitchRequest:
			// handler PathSwitchRequestAcknowledge
			log.Info("[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive PathSwitchRequestAcknowledge")
			HandlerPathSwitchRequestAcknowledge(gnb, ngapMsg)

		case ngapType.ProcedureCodeHandoverPreparation:
			// handler NGAP AMF Handover Command
			log.Info("[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive Handover Command")
			HandlerHandoverCommand(amf, gnb, ngapMsg)

		default:
			log.Warnf("[GNB][NGAP][](", gnb.GetGnbId(), ")()() Received unknown NGAP message 0x%x", ngapMsg.SuccessfulOutcome.ProcedureCode.Value)
		}

	case ngapType.NGAPPDUPresentUnsuccessfulOutcome:

		switch ngapMsg.UnsuccessfulOutcome.ProcedureCode.Value {

		case ngapType.ProcedureCodeNGSetup:
			// handler NGAP Setup Failure.
			log.Info("[GNB][NGAP][](", gnb.GetGnbId(), ")()() Receive Ng Setup Failure")
			HandlerNgSetupFailure(amf, gnb, ngapMsg)

		default:
			log.Warnf("[GNB][NGAP][](", gnb.GetGnbId(), ")()() Received unknown NGAP message 0x%x", ngapMsg.UnsuccessfulOutcome.ProcedureCode.Value)
		}
	}
}
