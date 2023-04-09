package main

import (
	"flag"
	"strings"
	"tsmon"
)

func main() {
	var devices string
	flag.StringVar(&devices, "devices", "", "list of devices to watch")
	flag.Parse()

	deviceList := strings.Split(devices, ",")
	notifyCh := make(chan tsmon.Notification)
	watcher := tsmon.NewWatchProcess(tsmon.NewLocalClient(),
		tsmon.WithNodesToWatch(deviceList),
		tsmon.WithNotificationChan(notifyCh),
	)

	notifier := tsmon.NewNotifier(
		tsmon.WithQueue(notifyCh),
	)
	go notifier.Start()

	watcher.Run()
}
