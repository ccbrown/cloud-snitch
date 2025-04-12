package apptest

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/protocol/webauthncbor"
	"github.com/go-webauthn/webauthn/protocol/webauthncose"

	"github.com/ccbrown/cloud-snitch/backend/model"
)

type Passkey struct {
	key        *rsa.PrivateKey
	id         model.Id
	userHandle []byte
}

func (p *Passkey) getRawAuthData(rpId string) ([]byte, error) {
	publicKeyData := webauthncose.RSAPublicKeyData{
		PublicKeyData: webauthncose.PublicKeyData{
			KeyType:   int64(webauthncose.RSAKey),
			Algorithm: int64(webauthncose.AlgRS256),
		},
		Modulus:  p.key.PublicKey.N.Bytes(),
		Exponent: binary.LittleEndian.AppendUint32(nil, uint32(p.key.PublicKey.E))[:3],
	}
	publicKeyBytes, err := webauthncbor.Marshal(publicKeyData)
	if err != nil {
		return nil, err
	}

	credentialIdBytes := []byte(p.id)

	rpIdHash := sha256.Sum256([]byte(rpId))

	ret := []byte{}
	ret = append(ret, rpIdHash[:]...)
	ret = append(ret, byte(protocol.FlagUserPresent|protocol.FlagAttestedCredentialData))
	ret = binary.BigEndian.AppendUint32(ret, 0)
	ret = append(ret, make([]byte, 16)...)
	ret = binary.BigEndian.AppendUint16(ret, uint16(len(credentialIdBytes)))
	ret = append(ret, credentialIdBytes...)
	ret = append(ret, publicKeyBytes...)
	return ret, nil
}

func NewPasskey(options *protocol.CredentialCreation) (*Passkey, *protocol.CredentialCreationResponse, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	userHandleString := options.Response.User.ID.(string)
	userHandleBytes, err := base64.RawURLEncoding.DecodeString(userHandleString)
	if err != nil {
		return nil, nil, err
	}

	p := &Passkey{
		key:        key,
		id:         model.NewId("test"),
		userHandle: userHandleBytes,
	}

	clientData := protocol.CollectedClientData{
		Type:         protocol.CreateCeremony,
		Challenge:    options.Response.Challenge.String(),
		Origin:       testFrontendURL,
		TokenBinding: nil,
		Hint:         "",
	}
	clientDataJSON, err := json.Marshal(clientData)
	if err != nil {
		return nil, nil, err
	}

	rawAuthData, err := p.getRawAuthData(options.Response.RelyingParty.ID)
	if err != nil {
		return nil, nil, err
	}

	attestationObject := protocol.AttestationObject{
		AuthData:     protocol.AuthenticatorData{},
		RawAuthData:  rawAuthData,
		Format:       "none",
		AttStatement: nil,
	}
	attestationObjectBytes, err := webauthncbor.Marshal(attestationObject)
	if err != nil {
		return nil, nil, err
	}

	return p, &protocol.CredentialCreationResponse{
		PublicKeyCredential: protocol.PublicKeyCredential{
			Credential: protocol.Credential{
				Type: "public-key",
				ID:   p.id.String(),
			},
			RawID:                   []byte(p.id),
			ClientExtensionResults:  nil,
			AuthenticatorAttachment: string(protocol.Platform),
		},
		AttestationResponse: protocol.AuthenticatorAttestationResponse{
			AuthenticatorResponse: protocol.AuthenticatorResponse{
				ClientDataJSON: clientDataJSON,
			},
			AttestationObject: attestationObjectBytes,
			Transports:        []string{"fake"},
		},
	}, nil
}

func (p *Passkey) Get(options *protocol.CredentialAssertion) (*protocol.CredentialAssertionResponse, error) {
	clientData := protocol.CollectedClientData{
		Type:         protocol.AssertCeremony,
		Challenge:    options.Response.Challenge.String(),
		Origin:       testFrontendURL,
		TokenBinding: nil,
		Hint:         "",
	}
	clientDataJSON, err := json.Marshal(clientData)
	if err != nil {
		return nil, err
	}
	clientDataHash := sha256.Sum256(clientDataJSON)

	rawAuthData, err := p.getRawAuthData(options.Response.RelyingPartyID)
	if err != nil {
		return nil, err
	}

	sigData := []byte{}
	sigData = append(sigData, rawAuthData...)
	sigData = append(sigData, clientDataHash[:]...)
	sigDataHash := sha256.Sum256(sigData)
	signature, err := rsa.SignPKCS1v15(rand.Reader, p.key, crypto.SHA256, sigDataHash[:])
	if err != nil {
		return nil, err
	}

	return &protocol.CredentialAssertionResponse{
		PublicKeyCredential: protocol.PublicKeyCredential{
			Credential: protocol.Credential{
				Type: "public-key",
				ID:   p.id.String(),
			},
			RawID:                   []byte(p.id),
			ClientExtensionResults:  nil,
			AuthenticatorAttachment: string(protocol.Platform),
		},
		AssertionResponse: protocol.AuthenticatorAssertionResponse{
			AuthenticatorResponse: protocol.AuthenticatorResponse{
				ClientDataJSON: clientDataJSON,
			},
			AuthenticatorData: rawAuthData,
			Signature:         signature,
			UserHandle:        p.userHandle,
		},
	}, nil
}
