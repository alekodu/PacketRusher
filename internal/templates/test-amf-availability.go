/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package templates

import (
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/monitoring"
	"my5G-RANTester/misc"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestAvailability(interval int) {

	monitor := monitoring.Monitor{}

	conf := config.GetConfig()

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.TESTER
	logFields[misc.FUNCTION] = misc.CONFIG

	ranPort := 1000
	for y := 1; y <= interval; y++ {

		monitor.InitAvaibility()

		for i := 1; i <= 1; i++ {

			conf.GNodeB.PlmnList.GnbId = gnbIdGenerator(i)

			conf.GNodeB.ControlIF.Port = ranPort

			go gnb.InitGnbForAvaibility(conf, &monitor)

			ranPort++
		}

		time.Sleep(1020 * time.Millisecond)

		if monitor.GetAvailability() {
			log.WithFields(logFields).Warn("AMF Availability:", 1)

		} else {
			log.WithFields(logFields).Warn("AMF Availability:", 0)

		}
	}

	return
}
