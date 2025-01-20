package opus

import "errors"

var (
	ErrTooLargeLastPacket = errors.New("last packet length is greater than frame size")
	ErrTooLargePacket     = errors.New("packet larger than expecting buffer size")
)
