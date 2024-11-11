/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package sender

import (
	context2 "my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/misc"

	log "github.com/sirupsen/logrus"
)

func SendToGnb(ue *context.UEContext, message []byte) {
	SendToGnbMsg(ue, context2.UEMessage{IsNas: true, Nas: message})
}

func SendToGnbMsg(ue *context.UEContext, message context2.UEMessage) {
	ue.Lock()
	gnbRx := ue.GetGnbRx()

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.UE
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NAS
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_MSIN] = ue.GetMsin()

	if gnbRx == nil {
		log.WithFields(logFields).Warn("Do not send NAS messages to gNB as channel is closed")
	} else {
		gnbRx <- message
	}
	ue.Unlock()
}
