/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package handler

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/ngap_control/nas_transport"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/sender"
	"my5G-RANTester/misc"

	log "github.com/sirupsen/logrus"
)

func HandlerUeInitialized(ue *context.GNBUe, message []byte, gnb *context.GNBContext) {

	var logFields log.Fields
	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	// encode NAS message in NGAP.
	ngap, err := nas_transport.SendInitialUeMessage(message, ue, gnb)
	if err != nil {
		log.WithFields(logFields).Errorln("Error making initial UE message: ", err)
	}

	// change state of UE.
	ue.SetStateOngoing()

	// Send Initial UE Message
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngap, conn)
	if err != nil {
		log.WithFields(logFields).Errorln("Error sending initial UE message to AMF: ", err)
	}
}

func HandlerUeOngoing(ue *context.GNBUe, message []byte, gnb *context.GNBContext) {

	var logFields log.Fields
	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	ngap, err := nas_transport.SendUplinkNasTransport(message, ue, gnb)
	if err != nil {
		log.WithFields(logFields).Errorln("Error making Uplink Nas Transport for AMF: ", err)
	}

	// Send Uplink Nas Transport
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngap, conn)
	if err != nil {
		log.WithFields(logFields).Errorln("Error sending Uplink Nas Transport to AMF: ", err)
	}
}

func HandlerUeReady(ue *context.GNBUe, message []byte, gnb *context.GNBContext) {

	var logFields log.Fields
	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NGAP

	ngap, err := nas_transport.SendUplinkNasTransport(message, ue, gnb)
	if err != nil {
		log.WithFields(logFields).Errorln("Error making Uplink Nas Transport for AMF: ", err)
	}

	// Send Uplink Nas Transport
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngap, conn)
	if err != nil {
		log.WithFields(logFields).Errorln("Error sending Uplink Nas Transport to AMF: ", err)
	}
}
