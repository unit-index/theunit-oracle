package sets

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Validator func(ctx context.Context, topic string, id peer.ID, msg *pubsub.Message) pubsub.ValidationResult

// ValidatorSet stores multiple instances of validators that implements
// the pubsub.ValidatorEx functions. Validators are groped by topic.
type ValidatorSet struct {
	validators []Validator
}

// NewValidatorSet creates new instance of the ValidatorSet.
func NewValidatorSet() *ValidatorSet {
	return &ValidatorSet{}
}

// Add adds new pubsub.ValidatorEx to the set.
func (n *ValidatorSet) Add(validator ...Validator) {
	n.validators = append(n.validators, validator...)
}

// Validator returns function that implements pubsub.ValidatorEx. That function
// will invoke all registered validators for given topic.
func (n *ValidatorSet) Validator(topic string) pubsub.ValidatorEx {
	return func(ctx context.Context, id peer.ID, psMsg *pubsub.Message) pubsub.ValidationResult {
		for _, validator := range n.validators {
			if result := validator(ctx, topic, id, psMsg); result != pubsub.ValidationAccept {
				return result
			}
		}
		return pubsub.ValidationAccept
	}
}
