package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	totalLen int
	first    *ListItem
	last     *ListItem
}

func NewList() List {
	return new(list)
}

func (listP *list) Len() int {
	return listP.totalLen
}

func (listP *list) Front() *ListItem {
	return listP.first
}

func (listP *list) Back() *ListItem {
	return listP.last
}

func (listP *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: listP.first}
	if listP.first != nil {
		listP.first.Prev = item
	}
	listP.first = item
	if listP.last == nil {
		listP.last = item
	}
	listP.totalLen++
	return item
}

func (listP *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{v, nil, listP.last}
	if listP.last != nil {
		listP.last.Next = item
	}
	listP.last = item
	if listP.first == nil {
		listP.first = item
	}
	listP.totalLen++
	return item
}

func (listP *list) Remove(item *ListItem) {
	if listP.first == item {
		listP.first = item.Next
		if item.Next != nil {
			item.Next.Prev = nil
		}
	}
	if listP.last == item {
		listP.last = item.Prev
		if item.Prev != nil {
			item.Prev.Next = nil
		}
	}
	if item.Next != nil && item.Prev != nil {
		item.Prev.Next, item.Next.Prev = item.Next, item.Prev
	}
	listP.totalLen--
}

func (listP *list) MoveToFront(item *ListItem) {
	if listP.first == item {
		return
	}
	listP.Remove(item)
	item.Next = listP.first
	if listP.first != nil {
		listP.first.Prev = item
		item.Prev = nil
	}
	listP.first = item
	listP.totalLen++
}
