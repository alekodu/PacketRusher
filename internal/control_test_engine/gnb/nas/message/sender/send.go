/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package sender

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/utils"

	log "github.com/sirupsen/logrus"
)

func SendToUe(ue *context.GNBUe, gnb *context.GNBContext, message []byte) {
	ue.Lock()
	gnbTx := ue.GetGnbTx()
	if gnbTx == nil {
		log.WithFields(log.Fields{
			utils.PROCEDURE: ue.GetProcedureType(),
			utils.STAGE:     ue.GetProcedureStage(),
			utils.UE_PR_ID:  ue.GetPrUeId(),
			utils.GNB_ID:    gnb.GetGnbId(),
			utils.NODE:      utils.GNB,
			utils.PROTOCOL:  utils.NAS,
		}).Warn("Do not send NAS messages to UE as channel is closed")
	} else {
		gnbTx <- context.UEMessage{IsNas: true, Nas: message}
	}
	ue.Unlock()
}

func SendMessageToUe(ue *context.GNBUe, gnb *context.GNBContext, message context.UEMessage) {
	ue.Lock()
	gnbTx := ue.GetGnbTx()
	if gnbTx == nil {
		log.WithFields(log.Fields{
			utils.PROCEDURE: ue.GetProcedureType(),
			utils.STAGE:     ue.GetProcedureStage(),
			utils.UE_PR_ID:  ue.GetPrUeId(),
			utils.GNB_ID:    gnb.GetGnbId(),
			utils.NODE:      utils.GNB,
			utils.PROTOCOL:  utils.NAS,
		}).Warn("Do not send NAS messages to UE as channel is closed")
	} else {
		gnbTx <- message
	}
	ue.Unlock()
}
