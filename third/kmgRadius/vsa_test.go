package kmgRadius

import (
	"testing"

	. "github.com/bronze1man/kmg/kmgTest"
)

func TestVSAFromFreeRadius(ot *testing.T) {
	//这个数据目前解不开.需要继续修改DecodePacket里面的信息
	inByte := []byte{
		0x02, 0x73, 0x00, 0xaa, 0xd9, 0x05, 0xde, 0x06,
		0x87, 0xae, 0xa9, 0x95, 0x2a, 0x5f, 0x0a, 0x2c,
		0x59, 0x0a, 0xbe, 0x0b, 0x1a, 0x0c, 0x00, 0x00,
		0x01, 0x37, 0x07, 0x06, 0x00, 0x00, 0x00, 0x01,
		0x1a, 0x0c, 0x00, 0x00, 0x01, 0x37, 0x08, 0x06,
		0x00, 0x00, 0x00, 0x06, 0x1a, 0x2a, 0x00, 0x00,
		0x01, 0x37, 0x10, 0x24, 0x92, 0xc3, 0xf4, 0x53,
		0x23, 0x8a, 0x1b, 0x31, 0x84, 0x16, 0xc0, 0x67,
		0xe2, 0x77, 0x29, 0x1b, 0x03, 0x00, 0xf6, 0x9f,
		0x36, 0x9d, 0x65, 0x6e, 0xdb, 0xd6, 0xfe, 0xe6,
		0x43, 0x9b, 0xe9, 0x2c, 0x29, 0x46, 0x1a, 0x2a,
		0x00, 0x00, 0x01, 0x37, 0x11, 0x24, 0x9e, 0x2b,
		0xf1, 0xf0, 0x6e, 0xf0, 0x20, 0x55, 0x5d, 0x5a,
		0xef, 0x36, 0x00, 0x08, 0x58, 0xce, 0x98, 0x9a,
		0x50, 0x80, 0x1b, 0x4d, 0xd5, 0xea, 0x17, 0xb2,
		0x08, 0xe6, 0xed, 0x0a, 0x21, 0xbb, 0x61, 0x0a,
		0x4f, 0x06, 0x03, 0x03, 0x00, 0x04, 0x50, 0x12,
		0x32, 0xaa, 0x90, 0x86, 0x7b, 0x31, 0xb9, 0xc0,
		0x55, 0x43, 0x64, 0x28, 0xef, 0xe7, 0x1c, 0x25,
		0x01, 0x12, 0x79, 0x64, 0x76, 0x62, 0x49, 0x77,
		0x30, 0x63, 0x41, 0x49, 0x34, 0x37, 0x45, 0x4d,
		0x51, 0x57}
	pac, err := DecodeResponsePacket([]byte("sEcReT"), inByte, [16]byte{0x65, 0x4d, 0x3c, 0x73,
		0x87, 0x8c, 0xfa, 0x28, 0xb6, 0xfd, 0x87, 0x96,
		0xba, 0x96, 0xd2, 0xe7})
	Equal(err, nil)
	Equal(pac.Code, CodeAccessAccept)

	Equal(pac.GetVsa(VendorTypeMSMPPESendKey).(*MSMPPESendOrRecvKeyVSA).Salt, [2]byte{0x92, 0xc3})
	Equal(pac.GetVsa(VendorTypeMSMPPESendKey).(*MSMPPESendOrRecvKeyVSA).Key, []byte{0x34, 0x29, 0xe7, 0x78, 0xe5, 0xad, 0x12, 0x14, 0xbf, 0x82, 0x6f, 0x2e, 0x3d, 0xe7, 0x6a, 0x77})
	Equal(pac.GetVsa(VendorTypeMSMPPERecvKey).(*MSMPPESendOrRecvKeyVSA).Salt, [2]byte{0x9e, 0x2b})
	Equal(pac.GetVsa(VendorTypeMSMPPERecvKey).(*MSMPPESendOrRecvKeyVSA).Key, []byte{0x3e, 0x24, 0x79, 0x82, 0xcb, 0x8, 0x1, 0xc7, 0x59, 0x6d, 0x2, 0x94, 0x83, 0xf3, 0x39, 0x1a})

}

func TestMsMPPEKeyEncodeDecode(ot *testing.T) {
	//send
	salt := [2]byte{0x90, 0xde}
	inData := []byte{0xe, 0x4e, 0x2f, 0xd3, 0xe7, 0x6e, 0x52, 0x43, 0xd7, 0xae, 0xd4, 0x7, 0x3, 0x5f, 0x8c, 0xa6}
	outData := []byte{0x90, 0xde, 0x61, 0x10, 0xf8, 0x3a, 0x72, 0x7d, 0x3e, 0x75, 0x2c, 0xb7, 0x28, 0xda, 0xb, 0x5d,
		0xcc, 0xd1, 0x19, 0x7b, 0x6a, 0x8f, 0x12, 0x6a, 0x32, 0x5c, 0x1e, 0x59, 0xe3, 0x4e, 0x2, 0x58, 0x14, 0xb0}
	pac := &Packet{
		Secret:        []byte("sEcReT"),
		Authenticator: [16]byte{0x14, 0xbb, 0xb9, 0xd9, 0xf8, 0xeb, 0x5d, 0xd3, 0xd9, 0x00, 0xb3, 0xeb, 0x8c, 0x84, 0xf5, 0x2c},
	}
	outDataGet, err := msMPPEKeyEncode(pac, salt, inData)
	Equal(err, nil)
	Equal(outDataGet, outData)
	inDataGet, saltGet, err := msMPPEKeyDecode(pac, outDataGet)
	Equal(err, nil)
	Equal(inDataGet, inData)
	Equal(saltGet, salt)

	//recv
	salt = [2]byte{0x9d, 0x23}
	inData = []byte{0x29, 0xc7, 0xf0, 0x63, 0xd, 0x70, 0x5c, 0xcf, 0x54, 0x1f, 0xa2, 0xc2, 0xa, 0xb, 0x2f, 0x41}
	outData = []byte{0x9d, 0x23,
		0x1b, 0x87, 0xfa, 0xaf, 0x7d, 0xe5, 0x38, 0xbe,
		0x5f, 0x4d, 0x40, 0xe2, 0x95, 0x19, 0xa7, 0xd5,
		0xf3, 0x93, 0xb7, 0x87, 0x2b, 0xc9, 0x4a, 0x8b,
		0x48, 0xac, 0x31, 0x05, 0x30, 0xd8, 0xe1, 0x4c}
	pac = &Packet{
		Secret:        []byte("sEcReT"),
		Authenticator: [16]byte{0x14, 0xbb, 0xb9, 0xd9, 0xf8, 0xeb, 0x5d, 0xd3, 0xd9, 0x00, 0xb3, 0xeb, 0x8c, 0x84, 0xf5, 0x2c},
	}
	inDataGet, saltGet, err = msMPPEKeyDecode(pac, outData)
	Equal(err, nil)
	Equal(inDataGet, inData)
	Equal(saltGet, salt)

	outDataGet, err = msMPPEKeyEncode(pac, salt, inData)
	Equal(err, nil)
	Equal(outDataGet, outData)

	/*
		outData2 := []byte{0x90, 0xde, 0x61, 0x10, 0xf8, 0x3a, 0x72, 0x7d, 0x3e, 0x75, 0x2c, 0xb7, 0x28, 0xda, 0xb, 0x5d,
			0xcc, 0xd1, 0x19, 0x7b, 0x6a, 0x8f, 0x12, 0x6a, 0x32, 0x5c, 0x1e, 0x59, 0xe3, 0x4e, 0x2, 0x58, 0x14, 0xb0}
		inDataGet, saltGet, err = msMPPEKeyDecode(pac, outData2)
		Equal(err, nil)
	*/
}
