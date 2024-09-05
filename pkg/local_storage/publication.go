package store

import (
	"slices"
	"sync"
	"time"
)

// PubData base storage for sent/delete time
// publicationID main identification for all searching operation, usage in del/sub
// channelID usage only for delete message, adding in struct when message successfully sent to channel
// sentMsgID message in channel which need to delete
type PubData struct {
	PubDate       time.Time
	DelDate       time.Time
	PublicationID int
	ChannelID     int64
	SentMsgID     int
}

type PublicationArray struct {
	publicationDataArray []*PubData
	deleteDataArray      []*PubData

	mu sync.RWMutex
}

func NewSortPublication(capacity int) *PublicationArray {
	return &PublicationArray{
		publicationDataArray: make([]*PubData, 0, capacity),
		deleteDataArray:      make([]*PubData, 0, capacity),
	}
}

func (p *PublicationArray) GetPub() []*PubData {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.publicationDataArray
}

func (p *PublicationArray) AppendPub(publication *PubData) {
	p.mu.Lock()
	p.publicationDataArray = append(p.publicationDataArray, publication)
	p.mu.Unlock()
}

func (p *PublicationArray) LenPub() int {
	return len(p.publicationDataArray)
}

func (p *PublicationArray) CapPub() int {
	return cap(p.publicationDataArray)
}

func (p *PublicationArray) SortPub() {
	p.mu.Lock()
	slices.SortFunc(p.publicationDataArray, func(a *PubData, b *PubData) int {
		return a.PubDate.Compare(b.PubDate)
	})
	p.mu.Unlock()
}

func (p *PublicationArray) RemovePub(arr *PubData) {
	p.mu.Lock()
	for key, value := range p.publicationDataArray {
		if value.PublicationID == arr.PublicationID {
			copy(p.publicationDataArray[key:], p.publicationDataArray[key+1:])
			p.publicationDataArray[len(p.publicationDataArray)-1] = nil
			p.publicationDataArray = p.publicationDataArray[:len(p.publicationDataArray)-1]
			break
		}
	}
	p.mu.Unlock()
}

func (p *PublicationArray) ReplacePub(arr *PubData) {
	p.mu.Lock()
	for key, value := range p.publicationDataArray {
		if value.PublicationID == arr.PublicationID {
			p.publicationDataArray[key] = arr
		}
	}
	p.mu.Unlock()
}

func (p *PublicationArray) GetDel() []*PubData {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.deleteDataArray
}

func (p *PublicationArray) AppendDel(publication *PubData) {
	p.mu.Lock()
	p.deleteDataArray = append(p.deleteDataArray, publication)
	p.mu.Unlock()
}

func (p *PublicationArray) LenDel() int {
	return len(p.deleteDataArray)
}

func (p *PublicationArray) CapDel() int {
	return cap(p.deleteDataArray)
}

func (p *PublicationArray) SortDel() {
	p.mu.Lock()
	slices.SortFunc(p.deleteDataArray, func(a *PubData, b *PubData) int {
		return a.PubDate.Compare(b.PubDate)
	})
	p.mu.Unlock()
}

func (p *PublicationArray) RemoveDel(arr *PubData) {
	p.mu.Lock()
	for key, value := range p.deleteDataArray {
		if value.PublicationID == arr.PublicationID {
			copy(p.deleteDataArray[key:], p.deleteDataArray[key+1:])
			p.deleteDataArray[len(p.deleteDataArray)-1] = nil
			p.deleteDataArray = p.deleteDataArray[:len(p.deleteDataArray)-1]
		}
	}
	p.mu.Unlock()
}

func (p *PublicationArray) ReplaceDel(arr *PubData) {
	p.mu.Lock()
	for key, value := range p.deleteDataArray {
		if value.PublicationID == arr.PublicationID {
			p.deleteDataArray[key] = arr
		}
	}
	p.mu.Unlock()
}
