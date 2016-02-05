package main

import (
	"crypto/aes"
	"crypto/cipher"
	"net"

	"github.com/alicebob/alac"
	"github.com/mesilliac/pulse-simple"
)

func writeUdp(aesiv, aeskey []byte, fmtp []int) {
	udpaddr, err := net.ResolveUDPAddr("udp", ":6000")
	if err != nil {
		panic(err)
	}
	udpconn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		panic(err)
	}
	// Never closes zomg

	dec, err := alac.New()
	if err != nil {
		panic(err)
	}

	// packetchan := make(chan []byte, 1000)
	// go CreateALACPlayer(fmtp, packetchan)

	ss := pulse.SampleSpec{pulse.SAMPLE_S16LE, 44100, 2}
	stream, _ := pulse.Playback("ituned", "my stream", &ss)
	defer stream.Free()
	defer stream.Drain()

	buf := make([]byte, 1024*16)
	for {
		n, _, err := udpconn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}
		packet := buf[:n]
		audio := packet[12:]
		todec := audio
		block, err := aes.NewCipher(aeskey)
		if err != nil {
			panic(err)
		}
		AESDec := cipher.NewCBCDecrypter(block, aesiv)
		for len(todec) >= aes.BlockSize {
			AESDec.CryptBlocks(todec[:aes.BlockSize], todec[:aes.BlockSize])
			todec = todec[aes.BlockSize:]
		}

		decoded := dec.Decode(audio)
		// fmt.Printf("audio packet %d->%d\n", len(audio), len(decoded))
		// send := make([]byte, len(audio))
		// copy(send, audio)
		// packetchan <- send
		stream.Write(decoded)
	}
}
