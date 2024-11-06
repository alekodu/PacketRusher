/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package procedures

type LogFieldType string

const (
	UeImsi      LogFieldType = "ueImsi"
	UePrId      LogFieldType = "uePrId"
	GnbId       LogFieldType = "gnbId"
	ProcedureId LogFieldType = "procedureId"
	UeStatus    LogFieldType = "ueStatus"
	PduStatus   LogFieldType = "pduStatus"
	Node        LogFieldType = "node"
	Protocol    LogFieldType = "protocol"
)

type NodeType string

const (
	GnB NodeType = "GNB"
	UE  NodeType = "UE"
	AMF NodeType = "AMF"
)

type ProtocolType string

const (
	SCTP ProtocolType = "SCTP"
	NGAP ProtocolType = "NGAP"
	GTP  ProtocolType = "GTP"
	NAS  ProtocolType = "NAS"
)
