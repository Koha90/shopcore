package telegram

import "github.com/koha90/shopcore/internal/flow"

func (r *Runner) rememberActiveMessage(key flow.SessionKey, messageID int) {
	if r == nil || messageID == 0 {
		return
	}

	r.activeMessageMu.Lock()
	defer r.activeMessageMu.Unlock()

	if r.activeMessageID == nil {
		r.activeMessageID = make(map[flow.SessionKey]int)
	}
	r.activeMessageID[key] = messageID
}

func (r *Runner) activeMessageFor(key flow.SessionKey) (int, bool) {
	if r == nil {
		return 0, false
	}

	r.activeMessageMu.RLock()
	defer r.activeMessageMu.RUnlock()

	id, ok := r.activeMessageID[key]
	return id, ok
}
