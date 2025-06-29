package mail

// Put your senders here
// Example:
//
//	NotificationSender = "Notification <notification@sender.com>"
//	SupportSender     = "Support <support@sender.com>"

type EmailSender string

const (
	NotificationSender EmailSender = "Gulg <notificacao@gulg.io>"
)

type Template string

const (
	WelcomeTmpl Template = "welcome.html"
)
