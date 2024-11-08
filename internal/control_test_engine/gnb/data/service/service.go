/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package service

import (
	"fmt"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/misc"
	"net"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

func InitGatewayGnb(gnb *context.GNBContext) error {

	// get ip for GNB gateway for data plane.
	ipGateway := gnb.GetGatewayGnbIp()

	conn, err := net.ListenPacket("ip4:4", ipGateway)
	if err != nil {
		return fmt.Errorf("[GNB][DATA] Error setting listen gateway GNB", err)
	}

	dataPlaneConn, err := ipv4.NewRawConn(conn)
	if err != nil {
		return fmt.Errorf("[GNB][DATA] Error setting data plane communication with UEs", err)
	}

	// successful established GNB/UE tunnel.
	gnb.SetUePlane(dataPlaneConn)

	go gatewayListen(gnb)

	return nil
}

func gatewayListen(gnb *context.GNBContext) {

	buffer := make([]byte, 65535)
	conn := gnb.GetUePlane()

	var logFields log.Fields
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.FUNCTION] = misc.DATA

	defer func() {
		err := conn.Close()
		if err != nil {
			log.WithFields(logFields).Info("[Error in closing GNB/UE tunnel\n")
		}
	}()

	for {

		ipHeader, payload, _, err := conn.ReadFrom(buffer)
		// log.WithFields(logFields).Info("Read %d bytes in GNB/UE tunnel", len(payload))
		if err != nil {
			log.WithFields(logFields).Info("Error in reading from GNB/UE tunnel: %+v", err)
			return
		}

		forwardData := make([]byte, len(payload[:]))
		copy(forwardData, payload[:])

		// find owner of  the Data Plane.
		ue, err := gnb.GetGnbUeByIp(ipHeader.Src.String())
		if err != nil || ue == nil {
			log.WithFields(logFields).Info("Invalid GNB UE IP. UE is not found in GNB UE IP Pool")
			return
		}

		go processingData(ue, gnb, forwardData)
	}
}

func processingData(ue *context.GNBUe, gnb *context.GNBContext, packet []byte) {

	var logFields log.Fields
	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()
	logFields[misc.NODE] = misc.GNB
	logFields[misc.GNB_ID] = gnb.GetGnbId()
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_TMSI] = ue.GetTMSI()
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.GTP

	// get GTP/UDP connection.
	conn := gnb.GetN3Plane()
	if conn == nil {
		log.WithFields(logFields).Info("N3 GTP/UDP is not setting")
		return
	}

	// send Data plane with GTP header.
	teidUplink := ue.GetTeidUplink()

	remote := fmt.Sprintf("%s:%d", gnb.GetUpfIp(), gnb.GetUpfPort())
	upfAddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		log.WithFields(logFields).Info("Error resolving UPF address for GTP/UDP tunnel", err)
		return
	}

	// send Data plane with GTP header.
	_, err = conn.WriteToGTP(teidUplink, packet, upfAddr)
	if err != nil {
		log.WithFields(logFields).Info("Error sending data plane in GTP/UDP tunnel")
	}

	//log.WithFields(logFields).Info("[GNB][GTP] Send %d bytes in GNB->UPF tunnel\n", n)
}
