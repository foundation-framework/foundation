package notify

type Provider interface {
	Send(text string, chatIds []int64, attachments ...Attachment) error
}

func CombineProviders(providers ...Provider) Provider {
	return &providersGroup{
		providers: providers,
	}
}

type providersGroup struct {
	providers []Provider
}

func (p *providersGroup) Send(text string, chatIds []int64, attachments ...Attachment) error {
	for _, provider := range p.providers {
		if err := provider.Send(text, chatIds, attachments...); err != nil {
			return err
		}
	}

	return nil
}
