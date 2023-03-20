package engine

type KillFeed struct {
	queue     *FeedItemList
	length    int
	maxLength int
}

// Object that stores last kills. Should be used as Singleton. Basically a linkedlist
type FeedItemList struct {
	Item *FeedItem
	Next *FeedItemList
}

// A single kill entity
type FeedItem struct {
	killer     string
	killerTeam string
	victim     string
	victimTeam string
	weapon     string
	hs         bool
}

func NewFeed() *KillFeed {
	queue := FeedItemList{
		Item: nil,
		Next: nil,
	}
	return &KillFeed{queue: &queue, length: 0, maxLength: 5}
}

// Push an item to the list of last killFeeds. If length is maxed out, remove the oldest element
func (f *KillFeed) pushFeedItem(fi *FeedItem) {
	lastItem := f.queue
	for f.queue.Next != nil {
		lastItem = f.queue.Next
	}
	lastItem.Next = &FeedItemList{Item: fi, Next: nil}

	if f.length == f.maxLength || f.length == 0 {
		f.queue = f.queue.Next
	}

	if f.length != f.maxLength {
		f.length++
	}
}

func (f *KillFeed) getItems() []FeedItem {
	list := f.queue
	res := make([]FeedItem, 0, f.length)
	for list.Next != nil {
		if list.Item == nil {
			continue
		}

		res = append(res, *list.Item)
	}

	return res
}
