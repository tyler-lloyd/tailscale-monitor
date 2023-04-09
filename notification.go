package tswatcher

import "time"

type Notification struct {
	Device    *TailnetDevice
	Timestamp time.Time
}
