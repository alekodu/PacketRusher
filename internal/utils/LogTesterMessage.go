/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package utils

// Log Fields
const (
	PROCEDURE = "procedure"
	STAGE     = "stage"
	UE_PR_ID  = "uePrId"
	UE_MSIN   = "ueMsin"
	GNB_ID    = "gnbId"
	NODE      = "node"
	PROTOCOL  = "protocol"
	FUNCTION  = "functionality"
)

// Node Types
const (
	TESTER = "TESTER"
	GNB    = "GNB"
	UE     = "UE"
	AMF    = "AMF"
)

// Protocol Types
const (
	SCTP = "SCTP"
	NGAP = "NGAP"
	GTP  = "GTP"
	NAS  = "NAS"
	XN   = "XN"
)

// Tester Functionalities
const (
	CONFIG = "CONFIG"
	SIMUL  = "SIMUL"
	DATA   = "DATA"
)
