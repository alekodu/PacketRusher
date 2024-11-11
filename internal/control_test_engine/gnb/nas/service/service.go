/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 * © Copyright 2023 Valentin D'Emmanuele
 */
package service

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/control_test_engine/gnb/nas"
	"my5G-RANTester/internal/control_test_engine/gnb/nas/message/sender"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/trigger"
	"my5G-RANTester/misc"

	log "github.com/sirupsen/logrus"
)

func InitServer(gnb *context.GNBContext) {
	go gnbListen(gnb)
}

func gnbListen(gnb *context.GNBContext) {
	ln := gnb.GetInboundChannel()

	logFields := make(log.Fields)

	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.PROTOCOL] = misc.NAS

	for {
		message := <-ln
		logFields[misc.UE_PR_ID] = message.PrUeId
		logFields[misc.UE_TMSI] = message.Tmsi

		if message.FetchPagedUEs {
			if message.GNBTx != nil {
				message.GNBTx <- context.UEMessage{PagedUEs: gnb.GetPagedUEs()}
				close(message.GNBTx)
			} else {
				logFields[misc.FUNCTION] = misc.MESSAG
				log.WithFields(logFields).Info("Unable to give PagedUEs to UE, GNBTx is nill")
			}
			continue
		}

		// TODO this region of the code may induces race condition.

		// new instance GNB UE context
		// store UE in UE Pool
		// store UE connection
		// select AMF and get sctp association
		// make a tun interface
		ue, _ := gnb.GetGnbUeByPrUeId(message.PrUeId)
		if ue != nil && message.IsHandover {
			// We already have a context for this UE since it was sent to us by the AMF from a NGAP Handover
			// Notify the AMF that the UE has succesfully been handed over to US
			ue.SetGnbRx(message.GNBRx)
			ue.SetGnbTx(message.GNBTx)

			// We enable the new PDU Session handed over to us
			msg := context.UEMessage{GNBPduSessions: ue.GetPduSessions(), GnbIp: gnb.GetN3GnbIp()}
			sender.SendMessageToUe(ue, gnb, msg)

			ue.SetStateReady()

			trigger.SendHandoverNotify(gnb, ue)
		} else {
			var err error
			ue, err = gnb.NewGnBUe(message.GNBTx, message.GNBRx, message.PrUeId, message.Tmsi)

			if ue == nil && err != nil {
				logFields[misc.FUNCTION] = misc.SETUP
				log.WithFields(logFields).Errorf("UE was not created succesfully: %s. Closing connection with UE.", err)
				close(message.GNBTx)
				continue
			}
			if message.UEContext != nil && message.IsHandover {
				// Xn Handover
				ue.SetProcedureType(string(context.XN_HANDOVER))
				ue.SetProcedureStage(string(context.INITIATED))

				logFields[misc.PROCEDURE] = ue.GetProcedureType()
				logFields[misc.STAGE] = ue.GetProcedureStage()
				logFields[misc.FUNCTION] = misc.MESSAG

				log.WithFields(logFields).Info("Received incoming handover for UE from another gNodeB")
				ue.SetStateReady()
				ue.CopyFromPreviousContext(message.UEContext)
				trigger.SendPathSwitchRequest(gnb, ue)

			} else {
				// Usual first UE connection to a gNodeB

				logFields[misc.PROCEDURE] = ue.GetProcedureType()
				logFields[misc.STAGE] = ue.GetProcedureStage()
				logFields[misc.FUNCTION] = misc.MESSAG

				log.WithFields(logFields).Info("Received incoming connection from new UE")
				mcc, mnc := gnb.GetMccAndMnc()
				message.GNBTx <- context.UEMessage{Mcc: mcc, Mnc: mnc}
				ue.SetPduSessions(message.GNBPduSessions)
			}
		}

		if ue == nil {
			log.WithFields(logFields).Errorf("UE has not been created")
			continue
		}

		// accept and handle connection.
		go processingConn(ue, gnb)
	}
}

func processingConn(ue *context.GNBUe, gnb *context.GNBContext) {
	logFields := make(log.Fields)
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.PROTOCOL] = misc.NAS
	logFields[misc.FUNCTION] = misc.MESSAG

	rx := ue.GetGnbRx()
	for {
		message, done := <-rx
		gnbUeContext, err := gnb.GetGnbUe(ue.GetRanUeId())
		if (gnbUeContext == nil || err != nil) && done {
			log.WithFields(logFields).Error("Ignoring message from UE ", ue.GetRanUeId(), " as UE Context was cleaned as requested by AMF.")
			break
		}
		if !done {
			if gnbUeContext != nil {
				gnbUeContext.SetStateDown()
			}
			break
		}

		// send to dispatch.
		if message.ConnectionClosed {
			logFields[misc.FUNCTION] = misc.SETUP
			log.WithFields(logFields).Info("Cleaning up context on current gNb")
			gnbUeContext.SetStateDown()
			if gnbUeContext.GetHandoverGnodeB() == nil {
				// We do not clean the context if it's a NGAP Handover, as AMF will request the context clean-up
				// Otherwise, we do clean the context
				gnb.DeleteGnBUe(ue)
			}
		} else if message.IsNas {
			nas.Dispatch(ue, message.Nas, gnb)
		} else if message.Idle {
			trigger.SendUeContextReleaseRequest(ue, gnb.GetGnbId())
		} else {
			log.Error("Received unknown message from UE")
		}
	}
}
