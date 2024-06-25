package mailer

import (
	"context"
	"cyclic/ent/user"
	"cyclic/pkg/colonel"
	"cyclic/pkg/dispatcher"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/scribe"
	"cyclic/pkg/secretary"
	"errors"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strings"
	"sync"
)

func Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done() // send signal to the wait group that this goroutine is done

	go func() {
		for {
			scribe.Scribe.Debug("start sending email")

			// dequeue email message
			message, err := dispatcher.Dequeue()
			if err != nil {
				scribe.Scribe.Error("failed to get email message", zap.Error(err))
				continue
			}

			// send email
			if err := SendEmail(ctx, message); err != nil {
				scribe.Scribe.Error("failed to send email", zap.Error(err))

				// if failed to send email, re-enqueue the message
				if err := dispatcher.Enqueue(message); err != nil {
					scribe.Scribe.Error("failed to re-enqueue email message", zap.Error(err))
				}

				continue
			}

			scribe.Scribe.Debug("email sent", zap.String("type", message.Type), zap.String("target", message.Target))
		}
	}()

	scribe.Scribe.Info("mailer started")
	<-ctx.Done()
	scribe.Scribe.Info("mailer stopped")
}

func SendEmail(ctx context.Context, message *dispatcher.Message) error {
	switch message.Type {
	case dispatcher.Verify:
		// generate signup token
		m := magistrate.New()

		token, err := m.Issue([]string{"verify"}, message.Target)
		if err != nil {
			return err
		}

		scribe.Scribe.Debug("verify token", zap.String("user", message.Target), zap.String("token", token))

		// get user email
		result, err := secretary.Minute.User.Query().Where(user.ID(uuid.MustParse(message.Target))).Only(ctx)
		if err != nil {
			return err
		}

		// send email
		auth := sasl.NewPlainClient("", colonel.Writ.SMTP.User, colonel.Writ.SMTP.Password)

		to := []string{result.Email}
		msg := strings.NewReader(fmt.Sprintf("Subject: Verify your email\n\n"+
			"Please verify your email address.\n\n"+
			"%s", token))

		if err := smtp.SendMail(fmt.Sprintf("%s:%d", colonel.Writ.SMTP.Host, colonel.Writ.SMTP.Port), auth, colonel.Writ.SMTP.User, to, msg); err != nil {
			return err
		}
	case dispatcher.Notify:
		// send notify email
		// TODO: implement notify email
	default:
		return errors.New("unknown email type")
	}

	return nil
}
