package externalsort

type ValAndFileReader struct {
	Val    string
	Reader LineReader
}
type MyHeapWittReaders []ValAndFileReader

func (h MyHeapWittReaders) Len() int { return len(h) }
func (h MyHeapWittReaders) Less(i, j int) bool {
	return h[i].Val < h[j].Val
}
func (h MyHeapWittReaders) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *MyHeapWittReaders) Push(x interface{}) {
	*h = append(*h, x.(ValAndFileReader))
}
func (h *MyHeapWittReaders) Pop() interface{} {
	old := *h
	oldLen := len(old)
	x := old[oldLen-1]
	*h = old[0 : oldLen-1]
	return x
}

type MyHeap []string

func (h MyHeap) Len() int { return len(h) }
func (h MyHeap) Less(i, j int) bool {
	return h[i] < h[j]
}
func (h MyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *MyHeap) Push(x interface{}) {
	*h = append(*h, x.(string))
}

func (h *MyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
