package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/toknowwhy/theunit-oracle/internal/p2p"
	"github.com/toknowwhy/theunit-oracle/pkg/ethereum"
	"github.com/toknowwhy/theunit-oracle/pkg/log"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/messages"
	"github.com/toknowwhy/theunit-oracle/pkg/transport/p2p/crypto/ethkey"
)

// oracle adds a validator for price messages. The validator checks if the
// author of the message is allowed to send price messages, the price
// message is valid, and if the price is not older than 5 min.
func oracle(feeders []ethereum.Address, signer ethereum.Signer, logger log.Logger) p2p.Options {
	return func(n *p2p.Node) error {
		n.AddValidator(func(ctx context.Context, topic string, id peer.ID, psMsg *pubsub.Message) pubsub.ValidationResult {
			priceMsg, ok := psMsg.ValidatorData.(*messages.Price)
			if !ok {
				return pubsub.ValidationAccept
			}
			// Check is a message signature is valid and extract author's address:
			priceFrom, err := priceMsg.Price.From(signer)
			wat := priceMsg.Price.Wat
			age := priceMsg.Price.Age.UTC().Format(time.RFC3339)
			val := priceMsg.Price.Val.String()
			if err != nil {
				logger.
					WithError(err).
					WithField("peerID", psMsg.GetFrom().String()).
					WithField("wat", wat).
					WithField("age", age).
					WithField("val", val).
					Warn("The price message was rejected, invalid signature")
				return pubsub.ValidationReject
			}
			// The libp2p message should be created by the same person who signs the price message:
			if ethkey.AddressToPeerID(*priceFrom) != psMsg.GetFrom() {
				logger.
					WithField("peerID", psMsg.GetFrom().String()).
					WithField("from", priceFrom.String()).
					WithField("wat", wat).
					WithField("age", age).
					WithField("val", val).
					Warn("The price message was rejected, the message author and price signature don't match")
				return pubsub.ValidationReject
			}
			// Check if an author is allowed to send price messages:
			feedAllowed := false
			for _, addr := range feeders {
				if addr == *priceFrom {
					feedAllowed = true
					break
				}
			}
			if !feedAllowed {
				logger.
					WithField("peerID", psMsg.GetFrom().String()).
					WithField("from", priceFrom.String()).
					WithField("wat", wat).
					WithField("age", age).
					WithField("val", val).
					Warn("The price message was ignored, the feeder is not allowed to send price messages")
				return pubsub.ValidationIgnore
			}
			// Check when message was created, ignore if older than 5 min, reject if older than 10 min:
			if time.Since(priceMsg.Price.Age) > 5*time.Minute {
				logger.
					WithField("peerID", psMsg.GetFrom().String()).
					WithField("from", priceFrom.String()).
					WithField("wat", wat).
					WithField("age", age).
					WithField("val", val).
					Warn("The price message was rejected, the message is older than 5 min")
				if time.Since(priceMsg.Price.Age) > 10*time.Minute {
					return pubsub.ValidationReject
				}
				return pubsub.ValidationIgnore
			}

			return pubsub.ValidationAccept
		})
		return nil
	}
}
