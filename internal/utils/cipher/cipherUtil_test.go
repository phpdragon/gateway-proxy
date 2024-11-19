package cipher

import (
	"fmt"
	"testing"
)

var privateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC6BSDlbRplhMMNZBKHX4xe8AwESpzHVfAcHHsX9FFSMuF91W3c
xgT/g5n+qlLLFzCE3hWG/yX5NMAxR4mS3MlhyXKwko3tK9Ua691afod1lxORR3Ia
Z8nV7v5Bv8y4JDe4E3/f/bQIGzroWiJ0sXTcO41GqvOw3G9leClSvjVnSwIDAQAB
AoGAagooYYCbTomq4wRL562ZCDmQsBWUX7Fmia/Wn6YfgVsN3bx/vx2Gld2AOIMB
ZVJXzzYGUYk7LV9bu/vKudRwWvtPIYzxPEOBoaFGPrEPTAfDVwOklhzTz3zDmCHi
hLSpjQXYCG7boZ6G96NfKyTRSGPqgovHY3+KJvd/gAoZqOkCQQDt9YXfanpuwZx1
7MYjMEvoh5UY1vSsV72ts3/gTVEUkWdfXHj0XhyRhOL9Qu99TGZcoEwecygoUPLm
dySyiXl/AkEAyB+JP8W7oTG/HTHc5QGDfwlPjIH1o5YG2I8Tp02i3G7OTJc/b9/+
SPxcKsT78xrgox5ModPKBX50F783Y2DANQJBAKY4Mjp882b4gWVybnlYHD4ir0h5
ptHYPFvgnfu9plx6sT3Qp4DzWHth2vlUT1w0CPC83E8M28lFula4dP7tvtsCQQCt
a2asdNVbophS3FrnuKAS/iaJRDVxRRk5oQMPACAZlYwAozC96gWZidb02S7cRHZV
5HPT6IwwppxD19hPrg/hAkBdnl+BGF/eRw80OHJtZBkLnY4Hx9YGft9AyIp3SB1z
QDhaWHgFA+8lED+bWTvJxt/IMmzyYBaGnZbEjCECKZLD
-----END RSA PRIVATE KEY-----
`

var publicKey = `
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC6BSDlbRplhMMNZBKHX4xe8AwE
SpzHVfAcHHsX9FFSMuF91W3cxgT/g5n+qlLLFzCE3hWG/yX5NMAxR4mS3MlhyXKw
ko3tK9Ua691afod1lxORR3IaZ8nV7v5Bv8y4JDe4E3/f/bQIGzroWiJ0sXTcO41G
qvOw3G9leClSvjVnSwIDAQAB
-----END PUBLIC KEY-----
`

func TestRsa(t *testing.T) {
	cipherTxt, _ := RsaBase64Encrypt(publicKey, "i like golang")
	fmt.Println(fmt.Sprintf("cipherTxt: %s", cipherTxt))

	golang, _ := RsaBase64Decrypt(privateKey, cipherTxt)
	fmt.Println(fmt.Sprintf("txt: %s", golang))
}

func TestRsaBase64Decrypt(t *testing.T) {
	cipherTxt := "raoQQsfN0KBfPAMRWnxr9kFPvJ6BgQ7PRBCMnz0nWsH03sD4IdlMvKpj78BHe7V7Ga1HZHyDxuJhVaJ0T5qKl8qHXzvKquzNtdMru7G4X9o8ylzkGxJLg-HYCWOrsZ77ZMaKoV9p-TCf-yMI21OpL_5JGot-XNfVVPkmg0z9FW0"
	java, _ := RsaBase64Decrypt(privateKey, cipherTxt)
	fmt.Println(fmt.Sprintf("txt: %s", java))

	cipherTxt = "WopnO2LnolZ7XpOwA_ktOhfkkaQQJQgkJudk3ZH_-ob36GQFv968nE1UBXxNekA9pIHBcvcl0ZWfwFhk-kyOF2FmQvpPY9LkqiCV0T32vhJet0n93ti2PBoFILxvChjzdOgSG9M0flH78Vm696Q4mHo7VMt_XMoHDTd3Rbagvt8"
	php, _ := RsaBase64Decrypt(privateKey, cipherTxt)
	fmt.Println(fmt.Sprintf("txt: %s", php))

	cipherTxt = "T2LFtY3dF_b6OBO07BN-3LtMSEBZqDukovDZ4HGCff8wosvlowf6IFJ3U7LFBIeHfiHBKiFuAV8-pFltCfTXtA4AwgVUnwbBMBWBfIJiLDi02ev30V-5BcYEuSF-cEdnSUd7WecrX4rHhzYLueGuj8H6c7RRbSbrJ6_3EFfU-K0"
	js, _ := RsaBase64Decrypt(privateKey, cipherTxt)
	fmt.Println(fmt.Sprintf("txt: %s", js))
}
