""

package p2p

import (
"time"

"github.com/libp2p/go-libp2p-core/peer"
pubsub "github.com/libp2p/go-libp2p-pubsub"

"github.com/toknowwhy/theunit-oracle/internal/p2p/sets"
"github.com/toknowwhy/theunit-oracle/pkg/log"
)

// PeerScoring configures peer scoring parameters used in a pubsub system.
func PeerScoring(
params *pubsub.PeerScoreParams,
thresholds *pubsub.PeerScoreThresholds,
topicScoreParams func (topic string) *pubsub.TopicScoreParams) Options {

return func (n *Node) error {
n.pubsubOpts = append(
n.pubsubOpts,
pubsub.WithPeerScore(params, thresholds),
pubsub.WithPeerScoreInspect(func (m map[peer.ID]*pubsub.PeerScoreSnapshot) {
for id, ps := range m {
n.tsLog.get().
WithField("peerID", id).
WithField("score", log.Format(ps)).
Debug("Peer score")
}
}, time.Minute),
)

n.AddNodeEventHandler(sets.NodeEventHandlerFunc(func (event interface{}) {
if e, ok := event.(sets.NodeTopicSubscribedEvent); ok {
var err error
defer func () {
if err != nil {
n.tsLog.get().
WithError(err).
WithField("topic", e.Topic).
Warn("Unable to set topic score parameters")
}
}()
sub, err := n.Subscription(e.Topic)
if err != nil {
return
}
if sp := topicScoreParams(e.Topic); sp != nil {
n.tsLog.get().
WithField("topic", e.Topic).
WithField("params", log.Format(sp)).
Info("Topic score params")
err = sub.topic.SetScoreParams(sp)
if err != nil {
return
}
}
}
}))
return nil
}
}
