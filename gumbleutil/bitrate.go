package gumbleutil

import (
	"time"

	"layeh.com/gumble/gumble"
)

var autoBitrate = &Listener{
	Connect: func(e *gumble.ConnectEvent) {
		if e.MaximumBitrate != nil {
			const safety = 5
			interval := e.Client.Config.AudioInterval
			dataBytes := (*e.MaximumBitrate / (8 * (int(time.Second/interval) + safety))) - 32 - 10

			// If we use default config and specify maximum bitrate of <35000, then we get a negative number
			// that is then passed into a make function in opus_nonshared that then panics.
			//
			// As a hack - let's set a minimum of 6 bytes for AudioDataBytes 
			// (lowest value that worked in manual testing)
			if dataBytes < 6 {
				dataBytes = 6
			}
			
			e.Client.Config.AudioDataBytes = dataBytes
		}
	},
}

// AutoBitrate is a gumble.EventListener that automatically sets the client's
// AudioDataBytes to suitable value, based on the server's bitrate.
var AutoBitrate gumble.EventListener

func init() {
	AutoBitrate = autoBitrate
}
