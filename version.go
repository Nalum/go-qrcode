// go-qrcode
// Copyright 2014 Tom Harwood

package qrcode

import (
	"log"

	bitset "github.com/nalum/go-qrcode/bitset"
)

// Error detection/recovery capacity.
//
// There are several levels of error detection/recovery capacity. Higher levels
// of error recovery are able to correct more errors, with the trade-off of
// increased symbol size.
type RecoveryLevel int

const (
	// Level L: 7% error recovery.
	Low RecoveryLevel = iota

	// Level M: 15% error recovery. Good default choice.
	Medium

	// Level Q: 25% error recovery.
	High

	// Level H: 30% error recovery.
	Highest
)

// qrCodeVersion describes the data length and encoding order of a single QR
// Code version. There are 40 versions numbers x 4 recovery levels == 160
// possible qrCodeVersion structures.
type qrCodeVersion struct {
	// Version number (1-40 inclusive).
	version int

	// Error recovery level.
	level RecoveryLevel

	dataEncoderType dataEncoderType

	// Encoded data can be split into multiple blocks. Each block contains data
	// and error recovery bytes.
	//
	// Larger QR Codes contain more blocks.
	block []block

	// Number of bits required to pad the combined data & error correction bit
	// stream up to the symbol's full capacity.
	numRemainderBits int
}

type block struct {
	numBlocks int

	// Total codewords (numCodewords == numErrorCodewords+numDataCodewords).
	numCodewords int

	// Number of data codewords.
	numDataCodewords int
}

var (
	versions = []qrCodeVersion{
		{
			1,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					1,
					26,
					19,
				},
			},
			0,
		},
		{
			1,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					1,
					26,
					16,
				},
			},
			0,
		},
		{
			1,
			High,
			dataEncoderType1To9,
			[]block{
				{
					1,
					26,
					13,
				},
			},
			0,
		},
		{
			1,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					1,
					26,
					9,
				},
			},
			0,
		},
		{
			2,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					1,
					44,
					34,
				},
			},
			7,
		},
		{
			2,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					1,
					44,
					28,
				},
			},
			7,
		},
		{
			2,
			High,
			dataEncoderType1To9,
			[]block{
				{
					1,
					44,
					22,
				},
			},
			7,
		},
		{
			2,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					1,
					44,
					16,
				},
			},
			7,
		},
		{
			3,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					1,
					70,
					55,
				},
			},
			7,
		},
		{
			3,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					1,
					70,
					44,
				},
			},
			7,
		},
		{
			3,
			High,
			dataEncoderType1To9,
			[]block{
				{
					2,
					35,
					17,
				},
			},
			7,
		},
		{
			3,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					2,
					35,
					13,
				},
			},
			7,
		},
		{
			4,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					1,
					100,
					80,
				},
			},
			7,
		},
		{
			4,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					2,
					50,
					32,
				},
			},
			7,
		},
		{
			4,
			High,
			dataEncoderType1To9,
			[]block{
				{
					2,
					50,
					24,
				},
			},
			7,
		},
		{
			4,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					4,
					25,
					9,
				},
			},
			7,
		},
		{
			5,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					1,
					134,
					108,
				},
			},
			7,
		},
		{
			5,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					2,
					67,
					43,
				},
			},
			7,
		},
		{
			5,
			High,
			dataEncoderType1To9,
			[]block{
				{
					2,
					33,
					15,
				},
				{
					2,
					34,
					16,
				},
			},
			7,
		},
		{
			5,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					2,
					33,
					11,
				},
				{
					2,
					34,
					12,
				},
			},
			7,
		},
		{
			6,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					2,
					86,
					68,
				},
			},
			7,
		},
		{
			6,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					4,
					43,
					27,
				},
			},
			7,
		},
		{
			6,
			High,
			dataEncoderType1To9,
			[]block{
				{
					4,
					43,
					19,
				},
			},
			7,
		},
		{
			6,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					4,
					43,
					15,
				},
			},
			7,
		},
		{
			7,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					2,
					98,
					78,
				},
			},
			0,
		},
		{
			7,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					4,
					49,
					31,
				},
			},
			0,
		},
		{
			7,
			High,
			dataEncoderType1To9,
			[]block{
				{
					2,
					32,
					14,
				},
				{
					4,
					33,
					15,
				},
			},
			0,
		},
		{
			7,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					4,
					39,
					13,
				},
				{
					1,
					40,
					14,
				},
			},
			0,
		},
		{
			8,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					2,
					121,
					97,
				},
			},
			0,
		},
		{
			8,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					2,
					60,
					38,
				},
				{
					2,
					61,
					39,
				},
			},
			0,
		},
		{
			8,
			High,
			dataEncoderType1To9,
			[]block{
				{
					4,
					40,
					18,
				},
				{
					2,
					41,
					19,
				},
			},
			0,
		},
		{
			8,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					4,
					40,
					14,
				},
				{
					2,
					41,
					15,
				},
			},
			0,
		},
		{
			9,
			Low,
			dataEncoderType1To9,
			[]block{
				{
					2,
					146,
					116,
				},
			},
			0,
		},
		{
			9,
			Medium,
			dataEncoderType1To9,
			[]block{
				{
					3,
					58,
					36,
				},
				{
					2,
					59,
					37,
				},
			},
			0,
		},
		{
			9,
			High,
			dataEncoderType1To9,
			[]block{
				{
					4,
					36,
					16,
				},
				{
					4,
					37,
					17,
				},
			},
			0,
		},
		{
			9,
			Highest,
			dataEncoderType1To9,
			[]block{
				{
					4,
					36,
					12,
				},
				{
					4,
					37,
					13,
				},
			},
			0,
		},
		{
			10,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					2,
					86,
					68,
				},
				{
					2,
					87,
					69,
				},
			},
			0,
		},
		{
			10,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					4,
					69,
					43,
				},
				{
					1,
					70,
					44,
				},
			},
			0,
		},
		{
			10,
			High,
			dataEncoderType10To26,
			[]block{
				{
					6,
					43,
					19,
				},
				{
					2,
					44,
					20,
				},
			},
			0,
		},
		{
			10,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					6,
					43,
					15,
				},
				{
					2,
					44,
					16,
				},
			},
			0,
		},
		{
			11,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					4,
					101,
					81,
				},
			},
			0,
		},
		{
			11,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					1,
					80,
					50,
				},
				{
					4,
					81,
					51,
				},
			},
			0,
		},
		{
			11,
			High,
			dataEncoderType10To26,
			[]block{
				{
					4,
					50,
					22,
				},
				{
					4,
					51,
					23,
				},
			},
			0,
		},
		{
			11,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					3,
					36,
					12,
				},
				{
					8,
					37,
					13,
				},
			},
			0,
		},
		{
			12,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					2,
					116,
					92,
				},
				{
					2,
					117,
					93,
				},
			},
			0,
		},
		{
			12,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					6,
					58,
					36,
				},
				{
					2,
					59,
					37,
				},
			},
			0,
		},
		{
			12,
			High,
			dataEncoderType10To26,
			[]block{
				{
					4,
					46,
					20,
				},
				{
					6,
					47,
					21,
				},
			},
			0,
		},
		{
			12,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					7,
					42,
					14,
				},
				{
					4,
					43,
					15,
				},
			},
			0,
		},
		{
			13,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					4,
					133,
					107,
				},
			},
			0,
		},
		{
			13,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					8,
					59,
					37,
				},
				{
					1,
					60,
					38,
				},
			},
			0,
		},
		{
			13,
			High,
			dataEncoderType10To26,
			[]block{
				{
					8,
					44,
					20,
				},
				{
					4,
					45,
					21,
				},
			},
			0,
		},
		{
			13,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					12,
					33,
					11,
				},
				{
					4,
					34,
					12,
				},
			},
			0,
		},
		{
			14,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					3,
					145,
					115,
				},
				{
					1,
					146,
					116,
				},
			},
			3,
		},
		{
			14,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					4,
					64,
					40,
				},
				{
					5,
					65,
					41,
				},
			},
			3,
		},
		{
			14,
			High,
			dataEncoderType10To26,
			[]block{
				{
					11,
					36,
					16,
				},
				{
					5,
					37,
					17,
				},
			},
			3,
		},
		{
			14,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					11,
					36,
					12,
				},
				{
					5,
					37,
					13,
				},
			},
			3,
		},
		{
			15,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					5,
					109,
					87,
				},
				{
					1,
					110,
					88,
				},
			},
			3,
		},
		{
			15,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					5,
					65,
					41,
				},
				{
					5,
					66,
					42,
				},
			},
			3,
		},
		{
			15,
			High,
			dataEncoderType10To26,
			[]block{
				{
					5,
					54,
					24,
				},
				{
					7,
					55,
					25,
				},
			},
			3,
		},
		{
			15,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					11,
					36,
					12,
				},
				{
					7,
					37,
					13,
				},
			},
			3,
		},
		{
			16,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					5,
					122,
					98,
				},
				{
					1,
					123,
					99,
				},
			},
			3,
		},
		{
			16,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					7,
					73,
					45,
				},
				{
					3,
					74,
					46,
				},
			},
			3,
		},
		{
			16,
			High,
			dataEncoderType10To26,
			[]block{
				{
					15,
					43,
					19,
				},
				{
					2,
					44,
					20,
				},
			},
			3,
		},
		{
			16,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					3,
					45,
					15,
				},
				{
					13,
					46,
					16,
				},
			},
			3,
		},
		{
			17,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					1,
					135,
					107,
				},
				{
					5,
					136,
					108,
				},
			},
			3,
		},
		{
			17,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					10,
					74,
					46,
				},
				{
					1,
					75,
					47,
				},
			},
			3,
		},
		{
			17,
			High,
			dataEncoderType10To26,
			[]block{
				{
					1,
					50,
					22,
				},
				{
					15,
					51,
					23,
				},
			},
			3,
		},
		{
			17,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					2,
					42,
					14,
				},
				{
					17,
					43,
					15,
				},
			},
			3,
		},
		{
			18,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					5,
					150,
					120,
				},
				{
					1,
					151,
					121,
				},
			},
			3,
		},
		{
			18,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					9,
					69,
					43,
				},
				{
					4,
					70,
					44,
				},
			},
			3,
		},
		{
			18,
			High,
			dataEncoderType10To26,
			[]block{
				{
					17,
					50,
					22,
				},
				{
					1,
					51,
					23,
				},
			},
			3,
		},
		{
			18,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					2,
					42,
					14,
				},
				{
					19,
					43,
					15,
				},
			},
			3,
		},
		{
			19,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					3,
					141,
					113,
				},
				{
					4,
					142,
					114,
				},
			},
			3,
		},
		{
			19,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					3,
					70,
					44,
				},
				{
					11,
					71,
					45,
				},
			},
			3,
		},
		{
			19,
			High,
			dataEncoderType10To26,
			[]block{
				{
					17,
					47,
					21,
				},
				{
					4,
					48,
					22,
				},
			},
			3,
		},
		{
			19,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					9,
					39,
					13,
				},
				{
					16,
					40,
					14,
				},
			},
			3,
		},
		{
			20,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					3,
					135,
					107,
				},
				{
					5,
					136,
					108,
				},
			},
			3,
		},
		{
			20,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					3,
					67,
					41,
				},
				{
					13,
					68,
					42,
				},
			},
			3,
		},
		{
			20,
			High,
			dataEncoderType10To26,
			[]block{
				{
					15,
					54,
					24,
				},
				{
					5,
					55,
					25,
				},
			},
			3,
		},
		{
			20,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					15,
					43,
					15,
				},
				{
					10,
					44,
					16,
				},
			},
			3,
		},
		{
			21,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					4,
					144,
					116,
				},
				{
					4,
					145,
					117,
				},
			},
			4,
		},
		{
			21,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					17,
					68,
					42,
				},
			},
			4,
		},
		{
			21,
			High,
			dataEncoderType10To26,
			[]block{
				{
					17,
					50,
					22,
				},
				{
					6,
					51,
					23,
				},
			},
			4,
		},
		{
			21,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					19,
					46,
					16,
				},
				{
					6,
					47,
					17,
				},
			},
			4,
		},
		{
			22,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					2,
					139,
					111,
				},
				{
					7,
					140,
					112,
				},
			},
			4,
		},
		{
			22,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					17,
					74,
					46,
				},
			},
			4,
		},
		{
			22,
			High,
			dataEncoderType10To26,
			[]block{
				{
					7,
					54,
					24,
				},
				{
					16,
					55,
					25,
				},
			},
			4,
		},
		{
			22,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					34,
					37,
					13,
				},
			},
			4,
		},
		{
			23,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					4,
					151,
					121,
				},
				{
					5,
					152,
					122,
				},
			},
			4,
		},
		{
			23,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					4,
					75,
					47,
				},
				{
					14,
					76,
					48,
				},
			},
			4,
		},
		{
			23,
			High,
			dataEncoderType10To26,
			[]block{
				{
					11,
					54,
					24,
				},
				{
					14,
					55,
					25,
				},
			},
			4,
		},
		{
			23,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					16,
					45,
					15,
				},
				{
					14,
					46,
					16,
				},
			},
			4,
		},
		{
			24,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					6,
					147,
					117,
				},
				{
					4,
					148,
					118,
				},
			},
			4,
		},
		{
			24,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					6,
					73,
					45,
				},
				{
					14,
					74,
					46,
				},
			},
			4,
		},
		{
			24,
			High,
			dataEncoderType10To26,
			[]block{
				{
					11,
					54,
					24,
				},
				{
					16,
					55,
					25,
				},
			},
			4,
		},
		{
			24,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					30,
					46,
					16,
				},
				{
					2,
					47,
					17,
				},
			},
			4,
		},
		{
			25,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					8,
					132,
					106,
				},
				{
					4,
					133,
					107,
				},
			},
			4,
		},
		{
			25,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					8,
					75,
					47,
				},
				{
					13,
					76,
					48,
				},
			},
			4,
		},
		{
			25,
			High,
			dataEncoderType10To26,
			[]block{
				{
					7,
					54,
					24,
				},
				{
					22,
					55,
					25,
				},
			},
			4,
		},
		{
			25,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					22,
					45,
					15,
				},
				{
					13,
					46,
					16,
				},
			},
			4,
		},
		{
			26,
			Low,
			dataEncoderType10To26,
			[]block{
				{
					10,
					142,
					114,
				},
				{
					2,
					143,
					115,
				},
			},
			4,
		},
		{
			26,
			Medium,
			dataEncoderType10To26,
			[]block{
				{
					19,
					74,
					46,
				},
				{
					4,
					75,
					47,
				},
			},
			4,
		},
		{
			26,
			High,
			dataEncoderType10To26,
			[]block{
				{
					28,
					50,
					22,
				},
				{
					6,
					51,
					23,
				},
			},
			4,
		},
		{
			26,
			Highest,
			dataEncoderType10To26,
			[]block{
				{
					33,
					46,
					16,
				},
				{
					4,
					47,
					17,
				},
			},
			4,
		},
		{
			27,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					8,
					152,
					122,
				},
				{
					4,
					153,
					123,
				},
			},
			4,
		},
		{
			27,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					22,
					73,
					45,
				},
				{
					3,
					74,
					46,
				},
			},
			4,
		},
		{
			27,
			High,
			dataEncoderType27To40,
			[]block{
				{
					8,
					53,
					23,
				},
				{
					26,
					54,
					24,
				},
			},
			4,
		},
		{
			27,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					12,
					45,
					15,
				},
				{
					28,
					46,
					16,
				},
			},
			4,
		},
		{
			28,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					3,
					147,
					117,
				},
				{
					10,
					148,
					118,
				},
			},
			3,
		},
		{
			28,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					3,
					73,
					45,
				},
				{
					23,
					74,
					46,
				},
			},
			3,
		},
		{
			28,
			High,
			dataEncoderType27To40,
			[]block{
				{
					4,
					54,
					24,
				},
				{
					31,
					55,
					25,
				},
			},
			3,
		},
		{
			28,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					11,
					45,
					15,
				},
				{
					31,
					46,
					16,
				},
			},
			3,
		},
		{
			29,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					7,
					146,
					116,
				},
				{
					7,
					147,
					117,
				},
			},
			3,
		},
		{
			29,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					21,
					73,
					45,
				},
				{
					7,
					74,
					46,
				},
			},
			3,
		},
		{
			29,
			High,
			dataEncoderType27To40,
			[]block{
				{
					1,
					53,
					23,
				},
				{
					37,
					54,
					24,
				},
			},
			3,
		},
		{
			29,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					19,
					45,
					15,
				},
				{
					26,
					46,
					16,
				},
			},
			3,
		},
		{
			30,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					5,
					145,
					115,
				},
				{
					10,
					146,
					116,
				},
			},
			3,
		},
		{
			30,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					19,
					75,
					47,
				},
				{
					10,
					76,
					48,
				},
			},
			3,
		},
		{
			30,
			High,
			dataEncoderType27To40,
			[]block{
				{
					15,
					54,
					24,
				},
				{
					25,
					55,
					25,
				},
			},
			3,
		},
		{
			30,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					23,
					45,
					15,
				},
				{
					25,
					46,
					16,
				},
			},
			3,
		},
		{
			31,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					13,
					145,
					115,
				},
				{
					3,
					146,
					116,
				},
			},
			3,
		},
		{
			31,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					2,
					74,
					46,
				},
				{
					29,
					75,
					47,
				},
			},
			3,
		},
		{
			31,
			High,
			dataEncoderType27To40,
			[]block{
				{
					42,
					54,
					24,
				},
				{
					1,
					55,
					25,
				},
			},
			3,
		},
		{
			31,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					23,
					45,
					15,
				},
				{
					28,
					46,
					16,
				},
			},
			3,
		},
		{
			32,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					17,
					145,
					115,
				},
			},
			3,
		},
		{
			32,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					10,
					74,
					46,
				},
				{
					23,
					75,
					47,
				},
			},
			3,
		},
		{
			32,
			High,
			dataEncoderType27To40,
			[]block{
				{
					10,
					54,
					24,
				},
				{
					35,
					55,
					25,
				},
			},
			3,
		},
		{
			32,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					19,
					45,
					15,
				},
				{
					35,
					46,
					16,
				},
			},
			3,
		},
		{
			33,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					17,
					145,
					115,
				},
				{
					1,
					146,
					116,
				},
			},
			3,
		},
		{
			33,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					14,
					74,
					46,
				},
				{
					21,
					75,
					47,
				},
			},
			3,
		},
		{
			33,
			High,
			dataEncoderType27To40,
			[]block{
				{
					29,
					54,
					24,
				},
				{
					19,
					55,
					25,
				},
			},
			3,
		},
		{
			33,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					11,
					45,
					15,
				},
				{
					46,
					46,
					16,
				},
			},
			3,
		},
		{
			34,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					13,
					145,
					115,
				},
				{
					6,
					146,
					116,
				},
			},
			3,
		},
		{
			34,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					14,
					74,
					46,
				},
				{
					23,
					75,
					47,
				},
			},
			3,
		},
		{
			34,
			High,
			dataEncoderType27To40,
			[]block{
				{
					44,
					54,
					24,
				},
				{
					7,
					55,
					25,
				},
			},
			3,
		},
		{
			34,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					59,
					46,
					16,
				},
				{
					1,
					47,
					17,
				},
			},
			3,
		},
		{
			35,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					12,
					151,
					121,
				},
				{
					7,
					152,
					122,
				},
			},
			0,
		},
		{
			35,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					12,
					75,
					47,
				},
				{
					26,
					76,
					48,
				},
			},
			0,
		},
		{
			35,
			High,
			dataEncoderType27To40,
			[]block{
				{
					39,
					54,
					24,
				},
				{
					14,
					55,
					25,
				},
			},
			0,
		},
		{
			35,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					22,
					45,
					15,
				},
				{
					41,
					46,
					16,
				},
			},
			0,
		},
		{
			36,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					6,
					151,
					121,
				},
				{
					14,
					152,
					122,
				},
			},
			0,
		},
		{
			36,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					6,
					75,
					47,
				},
				{
					34,
					76,
					48,
				},
			},
			0,
		},
		{
			36,
			High,
			dataEncoderType27To40,
			[]block{
				{
					46,
					54,
					24,
				},
				{
					10,
					55,
					25,
				},
			},
			0,
		},
		{
			36,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					2,
					45,
					15,
				},
				{
					64,
					46,
					16,
				},
			},
			0,
		},
		{
			37,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					17,
					152,
					122,
				},
				{
					4,
					153,
					123,
				},
			},
			0,
		},
		{
			37,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					29,
					74,
					46,
				},
				{
					14,
					75,
					47,
				},
			},
			0,
		},
		{
			37,
			High,
			dataEncoderType27To40,
			[]block{
				{
					49,
					54,
					24,
				},
				{
					10,
					55,
					25,
				},
			},
			0,
		},
		{
			37,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					24,
					45,
					15,
				},
				{
					46,
					46,
					16,
				},
			},
			0,
		},
		{
			38,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					4,
					152,
					122,
				},
				{
					18,
					153,
					123,
				},
			},
			0,
		},
		{
			38,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					13,
					74,
					46,
				},
				{
					32,
					75,
					47,
				},
			},
			0,
		},
		{
			38,
			High,
			dataEncoderType27To40,
			[]block{
				{
					48,
					54,
					24,
				},
				{
					14,
					55,
					25,
				},
			},
			0,
		},
		{
			38,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					42,
					45,
					15,
				},
				{
					32,
					46,
					16,
				},
			},
			0,
		},
		{
			39,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					20,
					147,
					117,
				},
				{
					4,
					148,
					118,
				},
			},
			0,
		},
		{
			39,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					40,
					75,
					47,
				},
				{
					7,
					76,
					48,
				},
			},
			0,
		},
		{
			39,
			High,
			dataEncoderType27To40,
			[]block{
				{
					43,
					54,
					24,
				},
				{
					22,
					55,
					25,
				},
			},
			0,
		},
		{
			39,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					10,
					45,
					15,
				},
				{
					67,
					46,
					16,
				},
			},
			0,
		},
		{
			40,
			Low,
			dataEncoderType27To40,
			[]block{
				{
					19,
					148,
					118,
				},
				{
					6,
					149,
					119,
				},
			},
			0,
		},
		{
			40,
			Medium,
			dataEncoderType27To40,
			[]block{
				{
					18,
					75,
					47,
				},
				{
					31,
					76,
					48,
				},
			},
			0,
		},
		{
			40,
			High,
			dataEncoderType27To40,
			[]block{
				{
					34,
					54,
					24,
				},
				{
					34,
					55,
					25,
				},
			},
			0,
		},
		{
			40,
			Highest,
			dataEncoderType27To40,
			[]block{
				{
					20,
					45,
					15,
				},
				{
					61,
					46,
					16,
				},
			},
			0,
		},
	}
)

var (
	// Each QR Code contains a 15-bit Format Information value.  The 15 bits
	// consist of 5 data bits concatenated with 10 error correction bits.
	//
	// The 5 data bits consist of:
	// - 2 bits for the error correction level (L=01, M=00, G=11, H=10).
	// - 3 bits for the data mask pattern identifier.
	//
	// formatBitSequence is a mapping from the 5 data bits to the completed 15-bit
	// Format Information value.
	//
	// For example, a QR Code using error correction level L, and data mask
	// pattern identifier 001:
	//
	// 01 | 001 = 01001 = 0x9
	// formatBitSequence[0x9].qrCode = 0x72f3 = 111001011110011
	formatBitSequence = []struct {
		regular uint32
		micro   uint32
	}{
		{0x5412, 0x4445},
		{0x5125, 0x4172},
		{0x5e7c, 0x4e2b},
		{0x5b4b, 0x4b1c},
		{0x45f9, 0x55ae},
		{0x40ce, 0x5099},
		{0x4f97, 0x5fc0},
		{0x4aa0, 0x5af7},
		{0x77c4, 0x6793},
		{0x72f3, 0x62a4},
		{0x7daa, 0x6dfd},
		{0x789d, 0x68ca},
		{0x662f, 0x7678},
		{0x6318, 0x734f},
		{0x6c41, 0x7c16},
		{0x6976, 0x7921},
		{0x1689, 0x06de},
		{0x13be, 0x03e9},
		{0x1ce7, 0x0cb0},
		{0x19d0, 0x0987},
		{0x0762, 0x1735},
		{0x0255, 0x1202},
		{0x0d0c, 0x1d5b},
		{0x083b, 0x186c},
		{0x355f, 0x2508},
		{0x3068, 0x203f},
		{0x3f31, 0x2f66},
		{0x3a06, 0x2a51},
		{0x24b4, 0x34e3},
		{0x2183, 0x31d4},
		{0x2eda, 0x3e8d},
		{0x2bed, 0x3bba},
	}

	// QR Codes version 7 and higher contain an 18-bit Version Information value,
	// consisting of a 6 data bits and 12 error correction bits.
	//
	// versionBitSequence is a mapping from QR Code version to the completed
	// 18-bit Version Information value.
	//
	// For example, a QR code of version 7:
	// versionBitSequence[0x7] = 0x07c94 = 000111110010010100
	versionBitSequence = []uint32{
		0x00000,
		0x00000,
		0x00000,
		0x00000,
		0x00000,
		0x00000,
		0x00000,
		0x07c94,
		0x085bc,
		0x09a99,
		0x0a4d3,
		0x0bbf6,
		0x0c762,
		0x0d847,
		0x0e60d,
		0x0f928,
		0x10b78,
		0x1145d,
		0x12a17,
		0x13532,
		0x149a6,
		0x15683,
		0x168c9,
		0x177ec,
		0x18ec4,
		0x191e1,
		0x1afab,
		0x1b08e,
		0x1cc1a,
		0x1d33f,
		0x1ed75,
		0x1f250,
		0x209d5,
		0x216f0,
		0x228ba,
		0x2379f,
		0x24b0b,
		0x2542e,
		0x26a64,
		0x27541,
		0x28c69,
	}
)

const (
	formatInfoLengthBits  = 15
	versionInfoLengthBits = 18
)

// formatInfo returns the 15-bit Format Information value for a QR
// code.
func (v qrCodeVersion) formatInfo(maskPattern int) *bitset.Bitset {
	formatID := 0

	switch v.level {
	case Low:
		formatID = 0x08 // 0b01000
	case Medium:
		formatID = 0x00 // 0b00000
	case High:
		formatID = 0x18 // 0b11000
	case Highest:
		formatID = 0x10 // 0b10000
	default:
		log.Panicf("Invalid level %d", v.level)
	}

	if maskPattern < 0 || maskPattern > 7 {
		log.Panicf("Invalid maskPattern %d", maskPattern)
	}

	formatID |= maskPattern & 0x7

	result := bitset.New()

	result.AppendUint32(formatBitSequence[formatID].regular, formatInfoLengthBits)

	return result
}

// versionInfo returns the 18-bit Version Information value for a QR Code.
//
// Version Information is applicable only to QR Codes versions 7-40 inclusive.
// nil is returned if Version Information is not required.
func (v qrCodeVersion) versionInfo() *bitset.Bitset {
	if v.version < 7 {
		return nil
	}

	result := bitset.New()
	result.AppendUint32(versionBitSequence[v.version], 18)

	return result
}

// numDataBits returns the data capacity in bits.
func (v qrCodeVersion) numDataBits() int {
	numDataBits := 0
	for _, b := range v.block {
		numDataBits += 8 * b.numBlocks * b.numDataCodewords // 8 bits in a byte
	}

	return numDataBits
}

// chooseQRCodeVersion chooses the most suitable QR Code version for a stated
// data length in bits, the error recovery level required, and the data encoder
// used.
//
// The chosen QR Code version is the smallest version able to fit numDataBits
// and the optional terminator bits required by the specified encoder.
//
// On success the chosen QR Code version is returned.
func chooseQRCodeVersion(level RecoveryLevel, encoder *dataEncoder, numDataBits int) *qrCodeVersion {
	var chosenVersion *qrCodeVersion

	for _, v := range versions {
		if v.level != level {
			continue
		} else if v.version < encoder.minVersion {
			continue
		} else if v.version > encoder.maxVersion {
			break
		}

		numFreeBits := v.numDataBits() - numDataBits

		if numFreeBits >= 0 {
			chosenVersion = &v
			break
		}
	}

	return chosenVersion
}

func (v qrCodeVersion) numTerminatorBitsRequired(numDataBits int) int {
	numFreeBits := v.numDataBits() - numDataBits

	var numTerminatorBits int

	switch {
	case numFreeBits >= 4:
		numTerminatorBits = 4
	default:
		numTerminatorBits = numFreeBits
	}

	return numTerminatorBits
}

// numBlocks returns the number of blocks.
func (v qrCodeVersion) numBlocks() int {
	numBlocks := 0

	for _, b := range v.block {
		numBlocks += b.numBlocks
	}

	return numBlocks
}

// numBitsToPadToCodeword returns the number of bits required to pad data of
// length numDataBits upto the nearest codeword size.
func (v qrCodeVersion) numBitsToPadToCodeword(numDataBits int) int {
	if numDataBits == v.numDataBits() {
		return 0
	}

	return (8 - numDataBits%8) % 8
}

// symbolSize returns the size of the QR Code symbol in number of modules (which
// is both the width and height, since QR codes are square). The QR Code has
// size symbolSize() x symbolSize() pixels. This does not include the quiet
// zone.
func (v qrCodeVersion) symbolSize() int {
	return 21 + (v.version-1)*4
}

// quietZoneSize returns the number of pixels of border space on each side of
// the QR Code. The quiet space assists with decoding.
func (v qrCodeVersion) quietZoneSize() int {
	return 4
}

// getQRCodeVersion returns the QR Code version by version number and recovery
// level. Returns nil if the requested combination is not defined.
func getQRCodeVersion(level RecoveryLevel, version int) *qrCodeVersion {
	for _, v := range versions {
		if v.level == level && v.version == version {
			return &v
		}
	}

	return nil
}
