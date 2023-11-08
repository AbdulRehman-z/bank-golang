package mail

import (
	"testing"

	"github.com/AbdulRehman-z/bank-golang/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {

	config, err := util.LoadConfig("..")
	require.NoError(t, err)
	senderEmail := "yousafbhaikhan10@gmail.com"
	password := config.APP_PASSWORD
	receiverEmail := []string{"bushrakhan2045@gmail.com"}

	mailSender := NewGmailSender("yousafbhaikhan10", senderEmail, password)
	err = mailSender.SendEmail(receiverEmail, "Test Email", "This is a test email from GoLang Bank App using Gmail SMTP Server. So, if you are reading this email, it means that the email sending functionality is working fine.")
	require.NoError(t, err)
}
