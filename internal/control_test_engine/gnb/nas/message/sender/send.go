/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package sender

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/misc"

	log "github.com/sirupsen/logrus"
)

func SendToUe(ue *context.GNBUe, gnb *context.GNBContext, message []byte) {
	ue.Lock()

	logFields := make(log.Fields)
	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NAS

	gnbTx := ue.GetGnbTx()
	if gnbTx == nil {
		log.WithFields(logFields).Warn("Do not send NAS messages to UE as channel is closed")
	} else {
		gnbTx <- context.UEMessage{IsNas: true, Nas: message}
	}
	ue.Unlock()
}

func SendMessageToUe(ue *context.GNBUe, gnb *context.GNBContext, message context.UEMessage) {
	ue.Lock()

	logFields := make(log.Fields)
	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NAS

	gnbTx := ue.GetGnbTx()
	if gnbTx == nil {
		log.WithFields(logFields).Warn("Do not send NAS messages to UE as channel is closed")
	} else {
		gnbTx <- message
	}
	ue.Unlock()
}
