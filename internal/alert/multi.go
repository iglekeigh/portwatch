package alert

// MultiNotifier fans out an Event to multiple Notifier implementations.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier returns a MultiNotifier wrapping the provided notifiers.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Add appends a notifier to the chain.
func (m *MultiNotifier) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Notify dispatches the event to all registered notifiers, collecting errors.
// It continues even if one notifier fails and returns the last non-nil error.
func (m *MultiNotifier) Notify(e Event) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Notify(e); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
