package receivers

import (
	"errors"
	"testing"
	"time"

	"github.com/slack-go/slack"
)

type slackClient interface {
	postMessage(string, slack.MsgOption)
}

type MockSlackClient struct{}
type MockSlackReceiver struct{}

func (m *MockSlackClient) postMessage(ch string, opt slack.MsgOption) (string, time.Time, error) {
	return ch, time.Now(), errors.New("Couldn't send a message")
}

func (mr *MockSlackReceiver) postMessage(mc *MockSlackClient, channel string, attachment *slack.Attachment) error {
	_, _, err := mc.postMessage(channel, slack.MsgOptionAttachments(*attachment))

	return err
}

func (mr *MockSlackReceiver) handleEvent(e ReceiverEvent, c chan error) {
	client := MockSlackClient{}
	text := string(e.EventName) + "" + e.Message

	attach := slack.Attachment{
		Text: text,
	}

	err := mr.postMessage(&client, "mockChannel", &attach)

	if err == nil {
		c <- errors.New("Problem handling event - should receive error from MockSlackClient")
	}
}

func TestPostMessage(t *testing.T) {
	receiver := MockSlackReceiver{}
	client := MockSlackClient{}
	attach := slack.Attachment{}

	err := receiver.postMessage(&client, "mockChannel", &attach)

	if err == nil {
		t.Errorf("Posting message test has failed - should receive an error from slackReceiver.postMessage method: %s \n", err)
	}
}

func TestHandleEvent(t *testing.T) {
	receiver := MockSlackReceiver{}
	event := ReceiverEvent{EventName: "mockEvent", Message: "mockMessage", AdditionalInfo: make(map[string]interface{})}
	channel := make(chan error)

	receiver.handleEvent(event, channel)

	select {
	case err := <-channel:
		t.Errorf("Handling event test has failed with error: %s \n", err)
	default:
	}
}
