package Utils

import "ty/csi/ws/SensorsWS/Global"

//type node Global.ItemInfo
type qt Global.QueueItem

func (q *qt) Push(n *Global.ItemInfo) {
	*q = append(*q, n)
}

func (q *qt) Pop() (n *Global.ItemInfo) {
	n = (*q)[0]
	*q = (*q)[1:]
	return
}

func (q *qt) Len() int {
	return len(*q)
}

type Stack []*Global.ItemInfo

func (q *Stack) Push(n *Global.ItemInfo) {
	*q = append(*q, n)
}

func (q *Stack) Pop() (n *Global.ItemInfo) {
	x := q.Len() - 1
	n = (*q)[x]
	*q = (*q)[:x]
	return
}
func (q *Stack) Len() int {
	return len(*q)
}
