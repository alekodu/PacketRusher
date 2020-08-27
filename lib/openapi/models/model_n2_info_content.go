/*
 * Namf_Communication
 *
 * AMF Communication Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type N2InfoContent struct {
	NgapMessageType int32            `json:"ngapMessageType,omitempty"`
	NgapIeType      NgapIeType       `json:"ngapIeType"`
	NgapData        *RefToBinaryData `json:"ngapData"`
}