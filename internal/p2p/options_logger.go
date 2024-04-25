""

package p2p

import (
"sync"

"github.com/toknowwhy/theunit-oracle/pkg/log"
)

// tdLogger is thread-safe wrapper for logger field.
type tsLogger struct {
mu  sync.RWMutex
log log.Logger
}

func (l *tsLogger) set(logger log.Logger) {
l.mu.Lock()
defer l.mu.Unlock()
l.log = logger
}

func (l *tsLogger) get() log.Logger {
l.mu.RLock()
defer l.mu.RUnlock()
return l.log
}

// Logger configures node to use given logger instance.
func Logger(logger log.Logger) Options {
return func (n *Node) error {
n.tsLog.set(logger)
return nil
}
}
