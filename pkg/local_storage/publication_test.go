package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSortPublication(t *testing.T) {
	pubArr := NewSortPublication(0)

	test := &Data{
		Date:          time.Now().Add(2 * time.Hour),
		PublicationID: 12498498,
	}

	pubArr.AppendPub(test)
	pubArr.AppendPub(&Data{
		Date:          time.Now().AddDate(0, 0, -2),
		PublicationID: 3423,
	})
	pubArr.AppendPub(&Data{
		Date:          time.Now(),
		PublicationID: 5,
	})
	pubArr.AppendPub(&Data{
		Date:          time.Now().AddDate(0, 0, -1),
		PublicationID: 1,
	})
	pubArr.AppendPub(&Data{
		Date:          time.Now().AddDate(0, 0, -2),
		PublicationID: 3423,
	})
	assert.Equal(t, pubArr.LenPub(), 5)

	pubArr.SortPub()
	t.Logf("%v", pubArr.GetPub())
	for key, value := range pubArr.GetPub() {
		t.Logf("%d - %v", key, value)
	}

	t.Logf("%d", pubArr.CapPub())
	pubArr.RemovePub(test)

	t.Logf("%v", pubArr.GetPub())
	for key, value := range pubArr.GetPub() {
		t.Logf("%d - %v", key, value)
	}

	t.Logf("%d", pubArr.CapPub())
	t.Logf("%d", pubArr.LenPub())
}
