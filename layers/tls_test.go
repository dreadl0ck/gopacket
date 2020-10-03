// Copyright 2020 The GoPacket Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file in the root of the source tree.
package layers

import (
	"reflect"
	"testing"

	"github.com/dreadl0ck/gopacket"
)

// https://github.com/tintinweb/scapy-ssl_tls/blob/master/tests/files/RSA_WITH_AES_128_CBC_SHA.pcap
// WARNING! Tests are specific for each packet. If you change a packet, please review their tests.

// Packet 4 - Client Hello (full packet, from Ethernet to TLS layers)
var testClientHello = []byte{
	0x00, 0x0c, 0x29, 0x1f, 0xab, 0x17, 0x00, 0x50, 0x56, 0xc0, 0x00, 0x08, 0x08, 0x00, 0x45, 0x00,
	0x00, 0xfe, 0x71, 0x42, 0x40, 0x00, 0x80, 0x06, 0x4e, 0xe1, 0xc0, 0xa8, 0xdc, 0x01, 0xc0, 0xa8,
	0xdc, 0x83, 0x2f, 0x0e, 0x01, 0xbb, 0x25, 0x6c, 0xbd, 0x3d, 0xcc, 0xce, 0xe1, 0xf7, 0x50, 0x18,
	0xff, 0xff, 0x7c, 0xaf, 0x00, 0x00, 0x16, 0x03, 0x01, 0x00, 0xd1, 0x01, 0x00, 0x00, 0xcd, 0x03,
	0x01, 0xff, 0xa2, 0x88, 0x97, 0x7c, 0x41, 0xa1, 0x08, 0x34, 0x2c, 0x98, 0xc2, 0x70, 0x04, 0xa0,
	0x5d, 0x5f, 0x39, 0xef, 0xe0, 0x70, 0xd5, 0x12, 0xf1, 0x35, 0x17, 0xb6, 0x0d, 0xc4, 0xd3, 0x09,
	0x85, 0x00, 0x00, 0x5a, 0xc0, 0x14, 0xc0, 0x0a, 0x00, 0x39, 0x00, 0x38, 0x00, 0x88, 0x00, 0x87,
	0xc0, 0x0f, 0xc0, 0x05, 0x00, 0x35, 0x00, 0x84, 0xc0, 0x13, 0xc0, 0x09, 0x00, 0x33, 0x00, 0x32,
	0x00, 0x9a, 0x00, 0x99, 0x00, 0x45, 0x00, 0x44, 0xc0, 0x0e, 0xc0, 0x04, 0x00, 0x2f, 0x00, 0x96,
	0x00, 0x41, 0xc0, 0x11, 0xc0, 0x07, 0xc0, 0x0c, 0xc0, 0x02, 0x00, 0x05, 0x00, 0x04, 0xc0, 0x12,
	0xc0, 0x08, 0x00, 0x16, 0x00, 0x13, 0xc0, 0x0d, 0xc0, 0x03, 0x00, 0x0a, 0x00, 0x15, 0x00, 0x12,
	0x00, 0x09, 0x00, 0x14, 0x00, 0x11, 0x00, 0x08, 0x00, 0x06, 0x00, 0x03, 0x00, 0xff, 0x02, 0x01,
	0x00, 0x00, 0x49, 0x00, 0x0b, 0x00, 0x04, 0x03, 0x00, 0x01, 0x02, 0x00, 0x0a, 0x00, 0x34, 0x00,
	0x32, 0x00, 0x0e, 0x00, 0x0d, 0x00, 0x19, 0x00, 0x0b, 0x00, 0x0c, 0x00, 0x18, 0x00, 0x09, 0x00,
	0x0a, 0x00, 0x16, 0x00, 0x17, 0x00, 0x08, 0x00, 0x06, 0x00, 0x07, 0x00, 0x14, 0x00, 0x15, 0x00,
	0x04, 0x00, 0x05, 0x00, 0x12, 0x00, 0x13, 0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x0f, 0x00,
	0x10, 0x00, 0x11, 0x00, 0x23, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x01, 0x01,
}
var testClientHelloDecoded = &TLS{
	BaseLayer: BaseLayer{
		Contents: testClientHello[54:],
		Payload:  nil,
	},
	ChangeCipherSpec: nil,
	Handshake: []TLSHandshakeRecord{
		{
			TLSRecordHeader{
				ContentType: 22,
				Version:     0x0301,
				Length:      209,
			},
		},
	},
	AppData: nil,
	Alert:   nil,
}

// Packet 6 - Server Hello, Certificate, Server Hello Done
var testServerHello = []byte{
	0x16, 0x03, 0x01, 0x00, 0x3a, 0x02, 0x00, 0x00, 0x36, 0x03, 0x01, 0x55, 0x5c, 0xd6, 0x97, 0xa3,
	0x97, 0xe9, 0xf4, 0x0c, 0xf4, 0x56, 0x14, 0x9f, 0xe4, 0x24, 0xf9, 0xeb, 0x49, 0xd4, 0xd1, 0x5f,
	0xfc, 0x12, 0xb4, 0xfd, 0x45, 0x4e, 0x3d, 0xeb, 0x6a, 0xad, 0xcf, 0x00, 0x00, 0x2f, 0x01, 0x00,
	0x0e, 0xff, 0x01, 0x00, 0x01, 0x00, 0x00, 0x23, 0x00, 0x00, 0x00, 0x0f, 0x00, 0x01, 0x01, 0x16,
	0x03, 0x01, 0x01, 0x90, 0x0b, 0x00, 0x01, 0x8c, 0x00, 0x01, 0x89, 0x00, 0x01, 0x86, 0x30, 0x82,
	0x01, 0x82, 0x30, 0x82, 0x01, 0x2c, 0x02, 0x01, 0x04, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48,
	0x86, 0xf7, 0x0d, 0x01, 0x01, 0x04, 0x05, 0x00, 0x30, 0x38, 0x31, 0x0b, 0x30, 0x09, 0x06, 0x03,
	0x55, 0x04, 0x06, 0x13, 0x02, 0x41, 0x55, 0x31, 0x0c, 0x30, 0x0a, 0x06, 0x03, 0x55, 0x04, 0x08,
	0x13, 0x03, 0x51, 0x4c, 0x44, 0x31, 0x1b, 0x30, 0x19, 0x06, 0x03, 0x55, 0x04, 0x03, 0x13, 0x12,
	0x53, 0x53, 0x4c, 0x65, 0x61, 0x79, 0x2f, 0x72, 0x73, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x20,
	0x43, 0x41, 0x30, 0x1e, 0x17, 0x0d, 0x39, 0x35, 0x31, 0x30, 0x30, 0x39, 0x32, 0x33, 0x33, 0x32,
	0x30, 0x35, 0x5a, 0x17, 0x0d, 0x39, 0x38, 0x30, 0x37, 0x30, 0x35, 0x32, 0x33, 0x33, 0x32, 0x30,
	0x35, 0x5a, 0x30, 0x60, 0x31, 0x0b, 0x30, 0x09, 0x06, 0x03, 0x55, 0x04, 0x06, 0x13, 0x02, 0x41,
	0x55, 0x31, 0x0c, 0x30, 0x0a, 0x06, 0x03, 0x55, 0x04, 0x08, 0x13, 0x03, 0x51, 0x4c, 0x44, 0x31,
	0x19, 0x30, 0x17, 0x06, 0x03, 0x55, 0x04, 0x0a, 0x13, 0x10, 0x4d, 0x69, 0x6e, 0x63, 0x6f, 0x6d,
	0x20, 0x50, 0x74, 0x79, 0x2e, 0x20, 0x4c, 0x74, 0x64, 0x2e, 0x31, 0x0b, 0x30, 0x09, 0x06, 0x03,
	0x55, 0x04, 0x0b, 0x13, 0x02, 0x43, 0x53, 0x31, 0x1b, 0x30, 0x19, 0x06, 0x03, 0x55, 0x04, 0x03,
	0x13, 0x12, 0x53, 0x53, 0x4c, 0x65, 0x61, 0x79, 0x20, 0x64, 0x65, 0x6d, 0x6f, 0x20, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x30, 0x5c, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d,
	0x01, 0x01, 0x01, 0x05, 0x00, 0x03, 0x4b, 0x00, 0x30, 0x48, 0x02, 0x41, 0x00, 0xb7, 0x2c, 0x25,
	0xdc, 0x49, 0xc5, 0xae, 0x6b, 0x43, 0xc5, 0x2e, 0x41, 0xc1, 0x2e, 0x6d, 0x95, 0x7a, 0x3a, 0xa9,
	0x03, 0x51, 0x78, 0x45, 0x0f, 0x2a, 0xd1, 0x58, 0xd1, 0x88, 0xf6, 0x9f, 0x8f, 0x1f, 0xd9, 0xfd,
	0xa5, 0x87, 0xde, 0x2a, 0x5d, 0x31, 0x5b, 0xee, 0x24, 0x66, 0xbf, 0xc0, 0x55, 0xdb, 0xfe, 0x70,
	0xc5, 0x2c, 0x39, 0x5f, 0x5a, 0x9f, 0xa8, 0x08, 0xfc, 0x21, 0x06, 0xd5, 0x4f, 0x02, 0x03, 0x01,
	0x00, 0x01, 0x30, 0x0d, 0x06, 0x09, 0x2a, 0x86, 0x48, 0x86, 0xf7, 0x0d, 0x01, 0x01, 0x04, 0x05,
	0x00, 0x03, 0x41, 0x00, 0x2b, 0x34, 0x5b, 0x22, 0x85, 0x62, 0x23, 0x07, 0x36, 0xf4, 0x0c, 0x2b,
	0x14, 0xd0, 0x1b, 0xcb, 0xd9, 0xbb, 0xd2, 0xc0, 0x9a, 0xcf, 0x12, 0xa1, 0x65, 0x90, 0x3a, 0xb7,
	0x17, 0x83, 0x3a, 0x10, 0x6b, 0xad, 0x2f, 0xd6, 0xb1, 0x11, 0xc0, 0x0d, 0x5a, 0x06, 0xdb, 0x11,
	0xd0, 0x2f, 0x34, 0x90, 0xf5, 0x76, 0x61, 0x26, 0xa1, 0x69, 0xf2, 0xdb, 0xb3, 0xe7, 0x20, 0xcb,
	0x3a, 0x64, 0xe6, 0x41, 0x16, 0x03, 0x01, 0x00, 0x04, 0x0e, 0x00, 0x00, 0x00,
}

// Packet 7 - Client Key Exchange, Change Cipher Spec, Encrypted Handshake Message
var testClientKeyExchange = []byte{
	0x16, 0x03, 0x01, 0x00, 0x46, 0x10, 0x00, 0x00, 0x42, 0x00, 0x40, 0x9e, 0x73, 0xdf, 0xe0, 0xf2,
	0xd0, 0x40, 0x32, 0x44, 0x9a, 0x34, 0x7f, 0x57, 0x86, 0x10, 0xea, 0x3d, 0xc5, 0xe2, 0xf9, 0xa5,
	0x69, 0x43, 0xc9, 0x0b, 0x00, 0x7e, 0x91, 0x31, 0x57, 0xfc, 0xc5, 0x65, 0x18, 0x0d, 0x44, 0xfd,
	0x51, 0xf8, 0xda, 0x8a, 0x7a, 0xab, 0x16, 0x03, 0xeb, 0xac, 0x23, 0x6e, 0x8d, 0xdd, 0xbb, 0xf4,
	0x75, 0xe7, 0xb7, 0xa3, 0xce, 0xdb, 0x67, 0x6b, 0x7d, 0x30, 0x2a, 0x14, 0x03, 0x01, 0x00, 0x01,
	0x01, 0x16, 0x03, 0x01, 0x00, 0x30, 0x15, 0xcb, 0x7a, 0x5b, 0x2d, 0xc0, 0x27, 0x09, 0x28, 0x62,
	0x95, 0x44, 0x9f, 0xa1, 0x1e, 0x4e, 0x6a, 0xfb, 0x49, 0x9d, 0x6a, 0x24, 0x44, 0xc6, 0x8e, 0x26,
	0xbc, 0xc1, 0x28, 0x8c, 0x27, 0xcc, 0xa2, 0xba, 0xec, 0x38, 0x63, 0x6e, 0x64, 0xd8, 0x52, 0x94,
	0x17, 0x96, 0x61, 0xfd, 0x9c, 0x54,
}
var testClientKeyExchangeDecoded = &TLS{
	BaseLayer: BaseLayer{
		Contents: testClientKeyExchange[81:],
		Payload:  nil,
	},
	ChangeCipherSpec: []TLSChangeCipherSpecRecord{
		{
			TLSRecordHeader{
				ContentType: 20,
				Version:     0x0301,
				Length:      1,
			},
			1,
		},
	},
	Handshake: []TLSHandshakeRecord{
		{
			TLSRecordHeader{
				ContentType: 22,
				Version:     0x0301,
				Length:      70,
			},
		},
		{
			TLSRecordHeader{
				ContentType: 22,
				Version:     0x0301,
				Length:      48,
			},
		},
	},
	AppData: nil,
	Alert:   nil,
}

// Packet 9 - New Session Ticket, Change Cipher Spec, Encryption Handshake Message
var testNewSessionTicket = []byte{
	0x16, 0x03, 0x01, 0x00, 0xaa, 0x04, 0x00, 0x00, 0xa6, 0x00, 0x00, 0x1c, 0x20, 0x00, 0xa0, 0xd4,
	0xee, 0xb0, 0x9b, 0xb5, 0xa2, 0xd3, 0x00, 0x57, 0x84, 0x59, 0xec, 0x0d, 0xbf, 0x05, 0x0c, 0xd5,
	0xb9, 0xe2, 0xf8, 0x32, 0xb5, 0xec, 0xce, 0xe2, 0x9c, 0x25, 0x25, 0xd9, 0x3e, 0x4a, 0x94, 0x5b,
	0xca, 0x18, 0x2b, 0x0f, 0x5f, 0xf6, 0x73, 0x38, 0x62, 0xcd, 0xcc, 0xf1, 0x32, 0x39, 0xe4, 0x5e,
	0x30, 0xf3, 0x94, 0xf5, 0xc5, 0x94, 0x3a, 0x8c, 0x8e, 0xe5, 0x12, 0x4a, 0x1e, 0xd8, 0x31, 0xb5,
	0x17, 0x09, 0xa6, 0x4c, 0x69, 0xca, 0xae, 0xfb, 0x04, 0x17, 0x64, 0x54, 0x9e, 0xc2, 0xfa, 0xf3,
	0x6d, 0xe9, 0xa5, 0xed, 0xa6, 0x65, 0xfe, 0x2f, 0xf3, 0xc6, 0xce, 0x78, 0x40, 0xf7, 0x65, 0xe0,
	0x13, 0xd3, 0x77, 0xc7, 0xc5, 0x79, 0x16, 0x56, 0x4c, 0x30, 0x94, 0xcf, 0xb0, 0x3c, 0x00, 0x91,
	0xbd, 0x86, 0x08, 0x9f, 0x2f, 0x05, 0x67, 0x03, 0x6f, 0xa7, 0x3b, 0xb9, 0x36, 0xf2, 0x80, 0x4f,
	0x60, 0x5d, 0x4c, 0xc4, 0x42, 0x5d, 0x02, 0x44, 0xba, 0x31, 0x8f, 0x39, 0x8e, 0x0c, 0x1e, 0xa8,
	0x26, 0x4f, 0x3e, 0x01, 0x96, 0xb3, 0x6f, 0xc6, 0x25, 0xe4, 0x30, 0x03, 0xd6, 0x3a, 0x7d, 0x14,
	0x03, 0x01, 0x00, 0x01, 0x01, 0x16, 0x03, 0x01, 0x00, 0x30, 0x25, 0xb8, 0x58, 0xc1, 0xa6, 0x3f,
	0xf8, 0xbd, 0xe6, 0xae, 0xbd, 0x98, 0xd4, 0x75, 0xa5, 0x45, 0x1b, 0xd8, 0x6a, 0x70, 0x79, 0x86,
	0x29, 0x4e, 0x4f, 0x64, 0xba, 0xe7, 0x1f, 0xca, 0x4b, 0x96, 0x9b, 0xf7, 0x0b, 0x50, 0xf5, 0x4f,
	0xfd, 0xda, 0xda, 0xcd, 0xcd, 0x4b, 0x12, 0x2e, 0xdf, 0xd5,
}

// Packet 13 - Two Application Data Records
var testDoubleAppData = []byte{
	0x17, 0x03, 0x01, 0x00, 0x20, 0x77, 0x3a, 0x94, 0x7d, 0xb4, 0x47, 0x4a, 0x1d, 0xd4, 0x6c, 0x5a,
	0x69, 0x74, 0x03, 0x93, 0x32, 0xca, 0x54, 0x5e, 0xa5, 0x81, 0x99, 0x6a, 0x73, 0x66, 0xbf, 0x06,
	0xa0, 0xdc, 0x6a, 0x9c, 0xb1, 0x17, 0x03, 0x01, 0x00, 0x20, 0x44, 0x64, 0xc8, 0xc2, 0x5a, 0xfc,
	0x4a, 0x82, 0xdd, 0x53, 0x6d, 0x30, 0x82, 0x4d, 0x35, 0x22, 0xf1, 0x5f, 0x3b, 0x96, 0x66, 0x79,
	0x61, 0x9f, 0x51, 0x93, 0x1b, 0xbf, 0x53, 0x3b, 0xf8, 0x26,
}
var testDoubleAppDataDecoded = &TLS{
	BaseLayer: BaseLayer{
		Contents: testDoubleAppData[37:],
		Payload:  nil,
	},
	ChangeCipherSpec: nil,
	Handshake:        nil,
	AppData: []TLSAppDataRecord{
		{
			TLSRecordHeader{
				ContentType: 23,
				Version:     0x0301,
				Length:      32,
			},
			testDoubleAppData[5 : 5+32],
		},
		{
			TLSRecordHeader{
				ContentType: 23,
				Version:     0x0301,
				Length:      32,
			},
			testDoubleAppData[42 : 42+32],
		},
	},
	Alert: nil,
}

var testAlertEncrypted = []byte{
	0x15, 0x03, 0x03, 0x00, 0x20, 0x44, 0xb9, 0x9c, 0x2c, 0x6e, 0xab, 0xa3, 0xdf, 0xb1, 0x77, 0x04,
	0xa2, 0xa4, 0x3a, 0x9a, 0x08, 0x1d, 0xe6, 0x51, 0xac, 0xa0, 0x5f, 0xab, 0x74, 0xa7, 0x96, 0x24,
	0xfe, 0x62, 0xfe, 0xe8, 0x5e,
}
var testAlertEncryptedDecoded = &TLS{
	BaseLayer: BaseLayer{
		Contents: testAlertEncrypted,
		Payload:  nil,
	},
	ChangeCipherSpec: nil,
	Handshake:        nil,
	AppData:          nil,
	Alert: []TLSAlertRecord{
		{
			TLSRecordHeader{
				ContentType: 21,
				Version:     0x0303,
				Length:      32,
			},
			0xFF,
			0xFF,
			testAlertEncrypted[5:],
		},
	},
}

// Malformed TLS records
var testMalformed = []byte{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
	0xd1, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xdb, 0xdc, 0xdd, 0xde, 0xdf,
}

var testTLSDecodeOptions = gopacket.DecodeOptions{
	SkipDecodeRecovery:       true,
	DecodeStreamsAsDatagrams: true,
}

func TestParseTLSClientHello(t *testing.T) {
	p := gopacket.NewPacket(testClientHello, LinkTypeEthernet, testTLSDecodeOptions)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeEthernet, LayerTypeIPv4, LayerTypeTCP, LayerTypeTLS}, t)

	if got, ok := p.Layer(LayerTypeTLS).(*TLS); ok {
		want := testClientHelloDecoded
		if !reflect.DeepEqual(got, want) {
			t.Errorf("TLS ClientHello packet processing failed:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No TLS layer type found in packet")
	}
}

func TestTLSClientHelloDecodeFromBytes(t *testing.T) {
	var got TLS
	want := *testClientKeyExchangeDecoded

	if err := got.DecodeFromBytes(testClientKeyExchange, gopacket.NilDecodeFeedback); err != nil {
		t.Errorf("TLS DecodeFromBytes first decode failed:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TLS DecodeFromBytes first decode doesn't match:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
	}

	if err := got.DecodeFromBytes(testClientKeyExchange, gopacket.NilDecodeFeedback); err != nil {
		t.Errorf("TLS DecodeFromBytes second decode failed:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TLS DecodeFromBytes second decode doesn't match:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
	}
}

func TestParseTLSChangeCipherSpec(t *testing.T) {
	p := gopacket.NewPacket(testClientKeyExchange, LayerTypeTLS, testTLSDecodeOptions)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeTLS}, t)

	if got, ok := p.Layer(LayerTypeTLS).(*TLS); ok {
		want := testClientKeyExchangeDecoded
		if !reflect.DeepEqual(got, want) {
			t.Errorf("TLS ChangeCipherSpec packet processing failed:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No TLS layer type found in packet")
	}
}

func TestParseTLSAppData(t *testing.T) {
	p := gopacket.NewPacket(testDoubleAppData, LayerTypeTLS, testTLSDecodeOptions)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeTLS}, t)

	if got, ok := p.Layer(LayerTypeTLS).(*TLS); ok {
		want := testDoubleAppDataDecoded
		if !reflect.DeepEqual(got, want) {
			t.Errorf("TLS TLSAppData packet processing failed:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No TLS layer type found in packet")
	}
}

func TestSerializeTLSAppData(t *testing.T) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true}
	err := gopacket.SerializeLayers(buf, opts, testDoubleAppDataDecoded)
	if err != nil {
		t.Fatal(err)
	}

	p2 := gopacket.NewPacket(buf.Bytes(), LayerTypeTLS, testTLSDecodeOptions)
	if p2.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p2.ErrorLayer().Error())
	}
	checkLayers(p2, []gopacket.LayerType{LayerTypeTLS}, t)

	if got, ok := p2.Layer(LayerTypeTLS).(*TLS); ok {
		want := testDoubleAppDataDecoded
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Reconstructed TLSAppData packet processing failed:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No TLS layer type found in reconstructed packet")
	}
}

func TestParseTLSMalformed(t *testing.T) {
	p := gopacket.NewPacket(testMalformed, LayerTypeTLS, testTLSDecodeOptions)
	if p.ErrorLayer() == nil {
		t.Error("No Decoding Error when parsing a malformed data")
	}
}

func TestParseTLSTooShort(t *testing.T) {
	p := gopacket.NewPacket(testMalformed[0:2], LayerTypeTLS, testTLSDecodeOptions)
	if p.ErrorLayer() == nil {
		t.Error("No Decoding Error when parsing a malformed data")
	}
}

func TestParseTLSLengthMismatch(t *testing.T) {
	var testLengthMismatch = make([]byte, len(testDoubleAppData))
	copy(testLengthMismatch, testDoubleAppData)
	testLengthMismatch[3] = 0xFF
	testLengthMismatch[4] = 0xFF
	p := gopacket.NewPacket(testLengthMismatch, LayerTypeTLS, testTLSDecodeOptions)
	if p.ErrorLayer() == nil {
		t.Error("No Decoding Error when parsing a malformed data")
	}
}

func TestParseTLSAlertEncrypted(t *testing.T) {
	p := gopacket.NewPacket(testAlertEncrypted, LayerTypeTLS, testTLSDecodeOptions)
	if p.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p.ErrorLayer().Error())
	}
	checkLayers(p, []gopacket.LayerType{LayerTypeTLS}, t)

	if got, ok := p.Layer(LayerTypeTLS).(*TLS); ok {
		want := testAlertEncryptedDecoded
		if !reflect.DeepEqual(got, want) {
			t.Errorf("TLS TLSAlert packet processing failed:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No TLS layer type found in packet")
	}
}

func TestSerializeTLSAlertEncrypted(t *testing.T) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true}
	err := gopacket.SerializeLayers(buf, opts, testAlertEncryptedDecoded)
	if err != nil {
		t.Fatal(err)
	}

	p2 := gopacket.NewPacket(buf.Bytes(), LayerTypeTLS, testTLSDecodeOptions)
	if p2.ErrorLayer() != nil {
		t.Error("Failed to decode packet:", p2.ErrorLayer().Error())
	}
	checkLayers(p2, []gopacket.LayerType{LayerTypeTLS}, t)

	if got, ok := p2.Layer(LayerTypeTLS).(*TLS); ok {
		want := testAlertEncryptedDecoded
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Reconstructed TLSAlertEncrypted packet processing failed:\ngot:\n%#v\n\nwant:\n%#v\n\n", got, want)
		}
	} else {
		t.Error("No TLS layer type found in reconstructed packet")
	}
}
