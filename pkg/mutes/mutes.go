package mutesmanager

import (
	"time"

	"main/pkg/utils"
)

type Mute struct {
	Chain      string
	ProposalID string
	Expires    time.Time
	Comment    string
}

type Mutes struct {
	Mutes []*Mute
}

func (m *Mute) IsExpired() bool {
	return m.Expires.Before(time.Now())
}

func (m *Mute) Matches(chain string, proposalID string) bool {
	match := true

	if m.Chain != "" {
		match = match && chain == m.Chain
	}

	if m.ProposalID != "" {
		match = match && proposalID == m.ProposalID
	}

	return match
}

func (m *Mutes) IsMuted(chain string, proposalID string) bool {
	for _, mute := range m.Mutes {
		if mute.IsExpired() {
			continue
		}

		if mute.Matches(chain, proposalID) {
			return true
		}
	}

	return false
}

func (m *Mute) GetExpirationTime() string {
	return m.Expires.Format(time.RFC822)
}

func (m *Mutes) AddMute(mute *Mute) {
	m.Mutes = append(m.Mutes, mute)
	m.Mutes = utils.Filter(m.Mutes, func(m *Mute) bool {
		return !m.IsExpired()
	})
}

func (m *Mutes) DeleteMute(chain string, proposalID string) bool {
	for index, mute := range m.Mutes {
		if mute.Chain == chain && mute.ProposalID == proposalID {
			m.Mutes = append(m.Mutes[:index], m.Mutes[index+1:]...)
			m.Mutes = utils.Filter(m.Mutes, func(m *Mute) bool {
				return !m.IsExpired()
			})
			return true
		}
	}

	m.Mutes = utils.Filter(m.Mutes, func(m *Mute) bool {
		return !m.IsExpired()
	})

	return false
}
