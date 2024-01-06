package mail

import (
	"github.com/ngohoang211020/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	err := util.Config.LoadConfig("../")
	require.NoError(t, err)
	sender := NewGmailSender(util.Config.EmailSenderName, util.Config.EmailSenderAddress, util.Config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="google.com">Google</a></p>
	`
	to := []string{"hoanggg2110@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
