package mulitiersorttable

type headersClicksAndIndex struct {
	clicksCount int
	index       int
}

type Pairs []*headersClicksAndIndex

func (p *Pairs) Len() int {
	return len(*p)
}

func (p *Pairs) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *Pairs) Less(i, j int) bool {
	return (*p)[i].clicksCount < (*p)[j].clicksCount
}
