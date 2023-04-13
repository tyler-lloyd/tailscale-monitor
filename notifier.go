package tsmon

import (
	"context"
	"time"

	"golang.org/x/exp/slog"
)

type NotificationService interface {
	Send(ctx context.Context, notification Notification)
}

type Notifier struct {
	logger              *slog.Logger
	notificationQueue   <-chan Notification
	notificationService NotificationService
	lastNotification    time.Time
}

type NotifierOption func(n *Notifier)

func NewNotifier(opts ...NotifierOption) *Notifier {
	noti := &Notifier{
		logger: slog.Default(),
	}
	for _, o := range opts {
		o(noti)
	}
	return noti
}

func WithQueue(ch <-chan Notification) NotifierOption {
	return func(n *Notifier) {
		n.notificationQueue = ch
	}
}

func WithNotificationService(svc NotificationService) NotifierOption {
	return func(n *Notifier) {
		n.notificationService = svc
	}
}

func (n *Notifier) Start() {
	for {
		event := <-n.notificationQueue
		go n.HandleNotification(event)
	}
}

func (n *Notifier) HandleNotification(notification Notification) {
	n.logger.Info("event received!", "time", notification.Timestamp)
	if time.Since(n.lastNotification) > time.Hour {
		n.logger.Info("sending notification to service")
		n.lastNotification = notification.Timestamp
		n.notificationService.Send(context.TODO(), notification)
	}
}
