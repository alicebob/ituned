package main

import (
	"crypto/aes"
	"crypto/cipher"
	"net"

	"github.com/alicebob/alac"
	"github.com/mesilliac/pulse-simple"
)

func handleStream(conn *net.UDPConn, aesiv, aeskey []byte, fmtp []int) {
	dec, err := alac.New()
	if err != nil {
		panic(err)
	}

	ss := pulse.SampleSpec{pulse.SAMPLE_S16LE, 44100, 2}
	stream, err := pulse.Playback("ituned", "my stream", &ss)
	if err != nil {
		panic(err)
	}
	defer stream.Free()
	defer stream.Drain()

	block, err := aes.NewCipher(aeskey)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024*16)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			break
		}
		audio := buf[12:n]
		AESDec := cipher.NewCBCDecrypter(block, aesiv)
		// for some reason todec isn't always a multiple of aes.BlockSize.
		rounded := aes.BlockSize * (len(audio) / aes.BlockSize)
		AESDec.CryptBlocks(audio[:rounded], audio[:rounded])

		// fmt.Printf("encoded\n%x\n", audio)
		decoded := dec.Decode(audio)
		stream.Write(decoded)
	}
}
