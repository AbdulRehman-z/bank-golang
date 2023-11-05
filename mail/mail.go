package mail

type Mail interface {
	SendEmail(to []string, subject string, body string) error
}

type MailService struct {
	Host     string
	Port     int
	email    string
	Password string
}

func (m *MailService) SendEmail(to []string, subject string, body string) error {
	return nil
}
