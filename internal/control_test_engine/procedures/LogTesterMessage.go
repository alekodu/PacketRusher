/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package procedures

// Log Fields
const (
	PROCEDURE = "procedure"
	STAGE     = "stage"
	UE_PR_ID  = "uePrId"
	GNB_ID    = "gnbId"
	NODE      = "node"
	PROTOCOL  = "protocol"
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
	CONFIG = "CONFIG"
	SCTP   = "SCTP"
	NGAP   = "NGAP"
	GTP    = "GTP"
	NAS    = "NAS"
)
