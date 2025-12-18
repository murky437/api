package longreminderpush

import (
	"context"
	"murky_api/internal/app"
	"murky_api/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type mockMessageService struct {
	SentMessages []string
}

func newMockMessageService() *mockMessageService {
	return &mockMessageService{
		SentMessages: []string{},
	}
}

func (ms *mockMessageService) Send(message string) error {
	ms.SentMessages = append(ms.SentMessages, message)
	return nil
}

type mockClock struct{}

func (f mockClock) Now() time.Time {
	return time.Date(2025, 2, 4, 6, 8, 10, 0, time.UTC)
}

func TestHandleSuccess(t *testing.T) {
	c := app.NewTestContainer(t)
	defer c.Close()

	mockMs := newMockMessageService()

	err := Handle(c.Db, c.Config, mockMs, mockClock{})(context.Background(), nil)
	require.NoError(t, err)

	require.Equal(t, []string{
		"Long reminder 1",
		"Long reminder 2",
	}, mockMs.SentMessages)

	longReminders, err := model.GetLongRemindersByUserId(c.Db, 1)
	require.NoError(t, err)

	lastRemindedAtTime := time.Date(2025, 2, 4, 6, 8, 10, 0, time.UTC)

	require.Equal(t, 2, len(longReminders))
	require.Equal(t, &lastRemindedAtTime, longReminders[0].LastRemindedAt)
	require.Equal(t, &lastRemindedAtTime, longReminders[1].LastRemindedAt)
}
