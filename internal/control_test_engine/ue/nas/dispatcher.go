/**
 * SPDX-License-Identifier: Apache-2.0
 * © Copyright 2023 Hewlett Packard Enterprise Development LP
 */
package nas

import (
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/handler"
	"my5G-RANTester/misc"
	"reflect"

	"github.com/free5gc/nas"
	"github.com/free5gc/nas/nasMessage"
	"github.com/free5gc/nas/nasType"
	"github.com/free5gc/nas/security"
	log "github.com/sirupsen/logrus"
)

func DispatchNas(ue *context.UEContext, message []byte) {

	var cph bool

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.UE
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NAS
	logFields[misc.UE_PR_ID] = ue.GetPrUeId()
	logFields[misc.UE_MSIN] = ue.GetMsin()
	logFields[misc.PROCEDURE] = ue.GetProcedureType()
	logFields[misc.STAGE] = ue.GetProcedureStage()

	// check if message is null.
	if message == nil {
		// TODO return error
		log.WithFields(logFields).Fatal("NAS message is nil")
	}

	// decode NAS message.
	m := new(nas.Message)
	m.SecurityHeaderType = nas.GetSecurityHeaderType(message) & 0x0f

	payload := make([]byte, len(message))
	copy(payload, message)

	newSecurityContext := false

	// check if NAS is security protected
	if m.SecurityHeaderType != nas.SecurityHeaderTypePlainNas {

		log.WithFields(logFields).Info("Message with security header")

		// information to check integrity and ciphered.

		// sequence number.
		sequenceNumber := message[6]

		// mac verification.
		macReceived := message[2:6]

		// remove security Header
		payload := payload[7:]

		// check security header type.
		cph = false
		switch m.SecurityHeaderType {

		case nas.SecurityHeaderTypeIntegrityProtected:
			log.WithFields(logFields).Info("Message with integrity")

		case nas.SecurityHeaderTypeIntegrityProtectedAndCiphered:
			log.WithFields(logFields).Info("Message with integrity and ciphered")
			cph = true

		case nas.SecurityHeaderTypeIntegrityProtectedWithNew5gNasSecurityContext:
			log.WithFields(logFields).Info("Message with integrity and with NEW 5G NAS SECURITY CONTEXT")
			newSecurityContext = true

		case nas.SecurityHeaderTypeIntegrityProtectedAndCipheredWithNew5gNasSecurityContext:
			log.WithFields(logFields).Error("Received message with security header \"Integrity protected and ciphered with new 5G NAS security context\", this is reserved for a SECURITY MODE COMPLETE and UE should not receive this code")
			return
		}

		// check security header(Downlink data).
		if ue.UeSecurity.DLCount.SQN() > sequenceNumber {
			ue.UeSecurity.DLCount.SetOverflow(ue.UeSecurity.DLCount.Overflow() + 1)
		}
		ue.UeSecurity.DLCount.SetSQN(sequenceNumber)

		// check ciphering.
		if cph {
			if err := security.NASEncrypt(ue.UeSecurity.CipheringAlg, ue.UeSecurity.KnasEnc, ue.UeSecurity.DLCount.Get(), security.Bearer3GPP,
				security.DirectionDownlink, payload); err != nil {
				log.WithFields(logFields).Error("error in encrypt algorithm")
				return
			} else {
				log.WithFields(logFields).Info("successful NAS CIPHERING")
			}
		}

		// decode NAS message.
		err := m.PlainNasDecode(&payload)
		if err != nil {
			log.WithFields(logFields).Error("Decode NAS error", err)
		}

		if newSecurityContext {
			if m.GmmHeader.GetMessageType() == nas.MsgTypeSecurityModeCommand {
				ue.UeSecurity.DLCount.Set(0, 0)
				ue.UeSecurity.CipheringAlg = m.SecurityModeCommand.SelectedNASSecurityAlgorithms.GetTypeOfCipheringAlgorithm()
				ue.UeSecurity.IntegrityAlg = m.SecurityModeCommand.SelectedNASSecurityAlgorithms.GetTypeOfIntegrityProtectionAlgorithm()
				ue.DerivateAlgKey()
			} else {
				log.WithFields(logFields).Error("Received message with security header \"Integrity protected with new 5G NAS security context\", but message type is not SECURITY MODE COMMAND")
				return
			}
		}

		mac32, err := security.NASMacCalculate(ue.UeSecurity.IntegrityAlg,
			ue.UeSecurity.KnasInt,
			ue.UeSecurity.DLCount.Get(),
			security.Bearer3GPP,
			security.DirectionDownlink, message[6:])
		if err != nil {
			log.WithFields(logFields).Error("NAS MAC error", err)
			return
		}

		// check integrity
		if !reflect.DeepEqual(mac32, macReceived) {
			log.WithFields(logFields).Error("NAS MAC verification failed(received:", macReceived, "expected:", mac32)
			return
		} else {
			log.WithFields(logFields).Info("successful NAS MAC verification")
		}

	} else {

		log.WithFields(logFields).Info("Message without security header")

		// decode NAS message.
		err := m.PlainNasDecode(&payload)
		if err != nil {
			// TODO return error
			log.WithFields(logFields).Info("Decode NAS error", err)
		}
	}

	switch m.GmmHeader.GetMessageType() {

	case nas.MsgTypeAuthenticationRequest:
		// handler authentication request.
		log.WithFields(logFields).Info("Receive Authentication Request")
		handler.HandlerAuthenticationRequest(ue, m)

	case nas.MsgTypeAuthenticationReject:
		// handler authentication reject.
		log.WithFields(logFields).Info("Receive Authentication Reject")
		handler.HandlerAuthenticationReject(ue, m)

	case nas.MsgTypeIdentityRequest:
		log.WithFields(logFields).Info("Receive Identify Request")
		// handler identity request.
		handler.HandlerIdentityRequest(ue, m)

	case nas.MsgTypeSecurityModeCommand:
		// handler security mode command.
		log.WithFields(logFields).Info("Receive Security Mode Command")
		if !newSecurityContext {
			log.WithFields(logFields).Warn("Received Security Mode Command with security header different from \"Integrity protected with new 5G NAS security context\" ")
		}
		handler.HandlerSecurityModeCommand(ue, m)

	case nas.MsgTypeRegistrationAccept:
		// handler registration accept.
		log.WithFields(logFields).Info("Receive Registration Accept")
		handler.HandlerRegistrationAccept(ue, m)

	case nas.MsgTypeConfigurationUpdateCommand:
		log.WithFields(logFields).Info("Receive Configuration Update Command")
		handler.HandlerConfigurationUpdateCommand(ue, m)

	case nas.MsgTypeDLNASTransport:
		// handler DL NAS Transport.
		log.WithFields(logFields).Info("Receive DL NAS Transport")
		handleCause5GMM(m.DLNASTransport.Cause5GMM)
		handler.HandlerDlNasTransportPduaccept(ue, m)

	case nas.MsgTypeServiceAccept:
		// handler service reject
		log.WithFields(logFields).Info("Receive Service Accept")
		handler.HandlerServiceAccept(ue, m)

	case nas.MsgTypeServiceReject:
		// handler service reject
		log.WithFields(logFields).Error("Receive Service Reject")
		handleCause5GMM(&m.ServiceReject.Cause5GMM)

	case nas.MsgTypeRegistrationReject:
		// handler registration reject
		log.WithFields(logFields).Error("Receive Registration Reject")
		handleCause5GMM(&m.RegistrationReject.Cause5GMM)

	case nas.MsgTypeStatus5GMM:
		log.WithFields(logFields).Error("Receive Status 5GMM")
		handleCause5GMM(&m.Status5GMM.Cause5GMM)

	case nas.MsgTypeStatus5GSM:
		log.WithFields(logFields).Error("Receive Status 5GSM")
		handleCause5GSM(&m.Status5GSM.Cause5GSM)

	default:
		log.WithFields(logFields).Warnf("Received unknown NAS message 0x%x", m.GmmHeader.GetMessageType())
	}

}

func handleCause5GSM(cause5SMM *nasType.Cause5GSM) {

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.UE
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NAS

	if cause5SMM != nil {
		log.WithFields(logFields).Error("UE received a 5GSM Failure, cause: ", cause5GMMToString(cause5SMM.Octet))
	}
}

func handleCause5GMM(cause5GMM *nasType.Cause5GMM) {

	logFields := make(log.Fields)

	logFields[misc.NODE] = misc.UE
	logFields[misc.FUNCTION] = misc.MESSAG
	logFields[misc.PROTOCOL] = misc.NAS

	if cause5GMM != nil {
		log.WithFields(logFields).Error("UE received a 5GMM Failure, cause: ", cause5GMMToString(cause5GMM.Octet))
	}
}

func cause5GMMToString(cause5GMM uint8) string {
	switch cause5GMM {
	case nasMessage.Cause5GMMIllegalUE:
		return "Illegal UE"
	case nasMessage.Cause5GMMPEINotAccepted:
		return "PEI not accepted"
	case nasMessage.Cause5GMMIllegalME:
		return "5GS services not allowed"
	case nasMessage.Cause5GMM5GSServicesNotAllowed:
		return "5GS services not allowed"
	case nasMessage.Cause5GMMUEIdentityCannotBeDerivedByTheNetwork:
		return "UE identity cannot be derived by the network"
	case nasMessage.Cause5GMMImplicitlyDeregistered:
		return "Implicitly de-registered"
	case nasMessage.Cause5GMMPLMNNotAllowed:
		return "PLMN not allowed"
	case nasMessage.Cause5GMMTrackingAreaNotAllowed:
		return "Tracking area not allowed"
	case nasMessage.Cause5GMMRoamingNotAllowedInThisTrackingArea:
		return "Roaming not allowed in this tracking area"
	case nasMessage.Cause5GMMNoSuitableCellsInTrackingArea:
		return "No suitable cells in tracking area"
	case nasMessage.Cause5GMMMACFailure:
		return "MAC failure"
	case nasMessage.Cause5GMMSynchFailure:
		return "Synch failure"
	case nasMessage.Cause5GMMCongestion:
		return "Congestion"
	case nasMessage.Cause5GMMUESecurityCapabilitiesMismatch:
		return "UE security capabilities mismatch"
	case nasMessage.Cause5GMMSecurityModeRejectedUnspecified:
		return "Security mode rejected, unspecified"
	case nasMessage.Cause5GMMNon5GAuthenticationUnacceptable:
		return "Non-5G authentication unacceptable"
	case nasMessage.Cause5GMMN1ModeNotAllowed:
		return "N1 mode not allowed"
	case nasMessage.Cause5GMMRestrictedServiceArea:
		return "Restricted service area"
	case nasMessage.Cause5GMMLADNNotAvailable:
		return "LADN not available"
	case nasMessage.Cause5GMMMaximumNumberOfPDUSessionsReached:
		return "Maximum number of PDU sessions reached"
	case nasMessage.Cause5GMMInsufficientResourcesForSpecificSliceAndDNN:
		return "Insufficient resources for specific slice and DNN"
	case nasMessage.Cause5GMMInsufficientResourcesForSpecificSlice:
		return "Insufficient resources for specific slice"
	case nasMessage.Cause5GMMngKSIAlreadyInUse:
		return "ngKSI already in use"
	case nasMessage.Cause5GMMNon3GPPAccessTo5GCNNotAllowed:
		return "Non-3GPP access to 5GCN not allowed"
	case nasMessage.Cause5GMMServingNetworkNotAuthorized:
		return "Serving network not authorized"
	case nasMessage.Cause5GMMPayloadWasNotForwarded:
		return "Payload was not forwarded"
	case nasMessage.Cause5GMMDNNNotSupportedOrNotSubscribedInTheSlice:
		return "DNN not supported or not subscribed in the slice"
	case nasMessage.Cause5GMMInsufficientUserPlaneResourcesForThePDUSession:
		return "Insufficient user-plane resources for the PDU session"
	case nasMessage.Cause5GMMSemanticallyIncorrectMessage:
		return "Semantically incorrect message"
	case nasMessage.Cause5GMMInvalidMandatoryInformation:
		return "Invalid mandatory information"
	case nasMessage.Cause5GMMMessageTypeNonExistentOrNotImplemented:
		return "Message type non-existent or not implementedE"
	case nasMessage.Cause5GMMMessageTypeNotCompatibleWithTheProtocolState:
		return "Message type not compatible with the protocol state"
	case nasMessage.Cause5GMMInformationElementNonExistentOrNotImplemented:
		return "Information element non-existent or not implemented"
	case nasMessage.Cause5GMMConditionalIEError:
		return "Conditional IE error"
	case nasMessage.Cause5GMMMessageNotCompatibleWithTheProtocolState:
		return "Message not compatible with the protocol state"
	case nasMessage.Cause5GMMProtocolErrorUnspecified:
		return "Protocol error, unspecified. Please share the pcap with packetrusher@hpe.com."
	default:
		return "Protocol error, unspecified. Please share the pcap with packetrusher@hpe.com."
	}
}
