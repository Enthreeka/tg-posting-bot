package store

import (
	"sync"
	"testing"
	"time"
)

func TestNewSortPublication(t *testing.T) {
	pubArr := NewSortPublication(5)

	pubArr.AppendPub(&PubData{PublicationID: 1, PubDate: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)})
	pubArr.AppendPub(&PubData{PublicationID: 3, PubDate: time.Date(1, 1, 1, 1, 0, 0, 0, time.UTC)})
	pubArr.AppendPub(&PubData{PublicationID: 4, PubDate: time.Date(1, 1, 1, 1, 1, 0, 0, time.UTC)})
	pubArr.AppendPub(&PubData{PublicationID: 5, PubDate: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)})

	var copyPubData []*PubData
	copyPubData = append(copyPubData, pubArr.GetPub()...)

	for _, value := range copyPubData {
		if value == nil {
			continue
		}

		if value.PubDate == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
			pubArr.RemovePub(value)
			t.Log(value)
		}

	}
}

func TestNewDelPublication(t *testing.T) {
	pubArr := NewSortPublication(5)

	pubArr.AppendDel(&PubData{PublicationID: 1, DelDate: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)})
	pubArr.AppendDel(&PubData{PublicationID: 3, DelDate: time.Date(1, 1, 1, 1, 0, 0, 0, time.UTC)})
	pubArr.AppendDel(&PubData{PublicationID: 4, DelDate: time.Date(1, 1, 1, 1, 1, 0, 0, time.UTC)})
	pubArr.AppendDel(&PubData{PublicationID: 5, DelDate: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)})

	var copyPubData []*PubData
	copyPubData = append(copyPubData, pubArr.GetDel()...)

	var wg sync.WaitGroup

	for _, value := range copyPubData {
		if value == nil {
			continue
		}

		if value.DelDate == time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				pubArr.RemoveDel(value)
			}()
			t.Log(value)
		}
	}

	wg.Wait()
}
