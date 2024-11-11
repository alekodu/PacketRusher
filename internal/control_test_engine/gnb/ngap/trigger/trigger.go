/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package trigger

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	ueSender "my5G-RANTester/internal/control_test_engine/gnb/nas/message/sender"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/ngap_control/interface_management"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/ngap_control/pdu_session_management"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/ngap_control/ue_context_management"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/ngap_control/ue_mobility_management"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/sender"
	"my5G-RANTester/misc"

	"github.com/free5gc/ngap/ngapType"
	log "github.com/sirupsen/logrus"
)

func SendPduSessionResourceSetupResponse(pduSessions []*context.GnbPDUSession, ue *context.GNBUe, gnb *context.GNBContext) {

	logFields := make(log.Fields)

	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating PDU Session Resource Setup Response")

	// send PDU Session Resource Setup Response.
	ngapMsg, err := pdu_session_management.PDUSessionResourceSetupResponse(pduSessions, ue, gnb)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending PDU Session Resource Setup Response: ", err)
	}

	ue.SetStateReady()

	// Send PDU Session Resource Setup Response.
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending PDU Session Resource Setup Response: ", err)
	}
}

func SendPduSessionReleaseResponse(pduSessionIds []ngapType.PDUSessionID, ue *context.GNBUe, gnbId string) {

	logFields := make(log.Fields)

	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnbId
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating PDU Session Release Response")

	if len(pduSessionIds) == 0 {
		log.WithFields(logFields).Fatal("Trying to send a PDU Session Release Reponse for no PDU Session")
	}

	ngapMsg, err := pdu_session_management.PDUSessionReleaseResponse(pduSessionIds, ue)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending PDU Session Release Response: ", err)
	}

	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending PDU Session Release Response: ", err)
	}
}

func SendInitialContextSetupResponse(ue *context.GNBUe, gnb *context.GNBContext) {
	logFields := make(log.Fields)

	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating Initial Context Setup Response")

	// send Initial Context Setup Response.
	ngapMsg, err := ue_context_management.InitialContextSetupResponse(ue, gnb)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending Initial Context Setup Response: ", err)
	}

	// Send Initial Context Setup Response.
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending Initial Context Setup Response: ", err)
	}
}

func SendUeContextReleaseRequest(ue *context.GNBUe, gnbId string) {
	logFields := make(log.Fields)

	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnbId
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating UE Context Release Request")

	// send UE Context Release Complete
	ngapMsg, err := ue_context_management.UeContextReleaseRequest(ue)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending UE Context Release Request: ", err)
	}

	// Send UE Context Release Complete
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending UE Context Release Request: ", err)
	}
}

func SendUeContextReleaseComplete(ue *context.GNBUe, gnbId string) {
	logFields := make(log.Fields)

	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnbId
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating UE Context Complete")

	// send UE Context Release Complete
	ngapMsg, err := ue_context_management.UeContextReleaseComplete(ue)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending UE Context Complete: ", err)
	}

	// Send UE Context Release Complete
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending UE Context Complete: ", err)
	}
}

func SendAmfConfigurationUpdateAcknowledge(amf *context.GNBAmf) {

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.GNB
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating AMF Configuration Update Acknowledge")

	// send AMF Configure Update Acknowledge
	ngapMsg, err := interface_management.AmfConfigurationUpdateAcknowledge()
	if err != nil {
		log.WithFields(logFields).Warn("Error sending AMF Configuration Update Acknowledge: ", err)
	}

	// Send AMF Configure Update Acknowledge
	conn := amf.GetSCTPConn()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Warn("Error sending AMF Configuration Update Acknowledge: ", err)
	}
}

func SendNgSetupRequest(gnb *context.GNBContext, amf *context.GNBAmf) {
	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating NG Setup Request ", gnb.GetGnbIp(), ":", gnb.GetGnbPort())

	// send NG setup response.
	ngapMsg, err := interface_management.NGSetupRequest(gnb, "PacketRusher")
	if err != nil {
		log.WithFields(logFields).Info("Error sending NG Setup Request: ", err)
	}

	conn := amf.GetSCTPConn()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Info("Error sending NG Setup Request: ", err)
	}

}

func SendPathSwitchRequest(gnb *context.GNBContext, ue *context.GNBUe) {

	logFields := make(log.Fields)

	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating Path Switch Request")

	// send NG setup response.
	ngapMsg, err := ue_mobility_management.PathSwitchRequest(gnb, ue)
	if err != nil {
		log.WithFields(logFields).Info("Error sending Path Switch Request: ", err)
	}

	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending Path Switch Request: ", err)
	}
}

func SendHandoverRequestAcknowledge(gnb *context.GNBContext, ue *context.GNBUe) {

	logFields := make(log.Fields)

	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating Handover Request Acknowledge")

	// send NG setup response.
	ngapMsg, err := ue_mobility_management.HandoverRequestAcknowledge(gnb, ue)
	if err != nil {
		log.WithFields(logFields).Info("Error sending Handover Request Acknowledge: ", err)
	}

	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending Handover Request Acknowledge for UE: ", err)
	}
}

func SendHandoverNotify(gnb *context.GNBContext, ue *context.GNBUe) {

	logFields := make(log.Fields)

	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating Handover Notify")

	// send NG setup response.
	ngapMsg, err := ue_mobility_management.HandoverNotify(gnb, ue)
	if err != nil {
		log.WithFields(logFields).Info("Error sending Handover Notify: ", err)
	}

	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending Handover Notify for UE: ", err)
	}
}

func TriggerXnHandover(oldGnb *context.GNBContext, newGnb *context.GNBContext, prUeId int64) {

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = oldGnb.GetGnbId()
	logFields[misc.UE_PR_ID] = prUeId
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating Xn Handover")

	gnbUeContext, err := oldGnb.GetGnbUeByPrUeId(prUeId)
	if err != nil {
		log.WithFields(logFields).Fatal("Error getting UE from PR UE ID: ", err)
	}

	newGnbRx := make(chan context.UEMessage, 1)
	newGnbTx := make(chan context.UEMessage, 1)
	newGnb.GetInboundChannel() <- context.UEMessage{GNBRx: newGnbRx, GNBTx: newGnbTx, PrUeId: gnbUeContext.GetPrUeId(), UEContext: gnbUeContext, IsHandover: true}

	msg := context.UEMessage{GNBRx: newGnbRx, GNBTx: newGnbTx, GNBInboundChannel: newGnb.GetInboundChannel()}

	ueSender.SendMessageToUe(gnbUeContext, newGnb, msg)
}

func TriggerNgapHandover(oldGnb *context.GNBContext, newGnb *context.GNBContext, prUeId int64) {

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = oldGnb.GetGnbId()
	logFields[misc.UE_PR_ID] = prUeId
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	log.WithFields(logFields).Info("Initiating NGAP Handover")

	gnbUeContext, err := oldGnb.GetGnbUeByPrUeId(prUeId)
	if err != nil {
		log.WithFields(logFields).Fatal("Error getting UE from PR UE ID: ", err)
	}

	gnbUeContext.SetHandoverGnodeB(newGnb)

	// send NG setup response.
	ngapMsg, err := ue_mobility_management.HandoverRequired(oldGnb, newGnb, gnbUeContext)
	if err != nil {
		log.WithFields(logFields).Info("Error sending Handover Required: ", err)
	}

	conn := gnbUeContext.GetSCTP()
	err = sender.SendToAmF(ngapMsg, conn)
	if err != nil {
		log.WithFields(logFields).Fatal("Error sending Handover Required: ", err)
	}
}
