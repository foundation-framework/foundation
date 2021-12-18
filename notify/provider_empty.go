package notify

type emptyProvider struct {
}

func NewEmptyProvider() Provider {
	return &emptyProvider{}
}

func (e *emptyProvider) Send(text string, chatIds []int64, attachments ...Attachment) error {
	// Doing nothing
	return nil
}
