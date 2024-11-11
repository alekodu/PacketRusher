/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package ngap

import (
	"fmt"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/misc"

	"github.com/ishidawataru/sctp"
	log "github.com/sirupsen/logrus"
)

var ConnCount int

func InitConn(amf *context.GNBAmf, gnb *context.GNBContext) error {

	// check AMF IP and AMF port.
	remote := fmt.Sprintf("%s:%d", amf.GetAmfIp(), amf.GetAmfPort())
	local := fmt.Sprintf("%s:%d", gnb.GetGnbIp(), gnb.GetGnbPort()+ConnCount)
	ConnCount++

	rem, err := sctp.ResolveSCTPAddr("sctp", remote)
	if err != nil {
		return err
	}
	loc, err := sctp.ResolveSCTPAddr("sctp", local)
	if err != nil {
		return err
	}

	// streams := amf.GetTNLAStreams()

	conn, err := sctp.DialSCTPExt(
		"sctp",
		loc,
		rem,
		sctp.InitMsg{NumOstreams: 2, MaxInstreams: 2})
	if err != nil {
		amf.SetSCTPConn(nil)
		return err
	}

	// set streams and other information about TNLA

	// successful established SCTP (TNLA - N2)
	amf.SetSCTPConn(conn)
	gnb.SetN2(conn)

	conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)

	go GnbListen(amf, gnb)

	return nil
}

func GnbListen(amf *context.GNBAmf, gnb *context.GNBContext) {

	buf := make([]byte, 65535)
	conn := amf.GetSCTPConn()

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.SCTP

	/*
		defer func() {
			err := conn.Close()
			if err != nil {
				log.WithFields(logFields).Info("Error in closing SCTP association for %d AMF\n", amf.GetAmfId())
			}
		}()
	*/

	for {

		n, info, err := conn.SCTPRead(buf[:])
		if err != nil {
			break
		}

		log.WithFields(logFields).Info("Receive message in ", info.Stream, " stream\n")

		forwardData := make([]byte, n)
		copy(forwardData, buf[:n])

		// handling NGAP message.
		go Dispatch(amf, gnb, forwardData)

	}

}
