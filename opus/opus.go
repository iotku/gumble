package opus

import (
	"log"

	gopus "gopkg.in/hraban/opus.v2"
	"layeh.com/gumble/gumble"
)

var Codec gumble.AudioCodec

const ID = 4

func init() {
	Codec = &generator{}
	gumble.RegisterAudioCodec(4, Codec)
}

// generator

type generator struct {
}

func (*generator) ID() int {
	return ID
}

func (*generator) NewEncoder() gumble.AudioEncoder {
	log.Println("Created new hraban/opus Encoder")
	e, _ := gopus.NewEncoder(gumble.AudioSampleRate, gumble.AudioChannels, gopus.AppAudio)
	e.SetBitrate(192000) // TODO: Don't hardcode this?
	return &Encoder{
		e,
	}
}

func (*generator) NewDecoder() gumble.AudioDecoder {
	d, _ := gopus.NewDecoder(gumble.AudioSampleRate, gumble.AudioChannels)
	return &Decoder{
		d,
	}
}

// encoder

type Encoder struct {
	*gopus.Encoder
}

func (*Encoder) ID() int {
	return ID
}

func (e *Encoder) Encode(pcm []int16, mframeSize, maxDataBytes int) ([]byte, error) {
	targetBuffer := make([]byte, maxDataBytes)
	n, err := e.Encoder.Encode(pcm, targetBuffer)
	if err != nil {
		return nil, err
	}
	return targetBuffer[:n], nil
}

func (e *Encoder) Reset() {
	e.Encoder.Reset()
}

// decoder

type Decoder struct {
	*gopus.Decoder
}

func (*Decoder) ID() int {
	return 4
}

func (d *Decoder) Decode(data []byte, frameSize int) ([]int16, error) {
	targetBuffer := make([]int16, frameSize*gumble.AudioChannels)
	n, err := d.Decoder.Decode(data, targetBuffer)
	if err != nil {
		return nil, err
	}
	return targetBuffer[:n], nil
}

func (d *Decoder) Reset() {
	log.Fatalln("Tried to use decoder reset, NOT IMPLEMENTED") // TODO: I don't think this is ever used, but we're going to find out!!!!
	// d.Decoder.Reset() // TODO: This method doesn't exist in hraban
}
