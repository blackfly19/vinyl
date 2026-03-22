package qwe

type Commit struct {
	CommitID          string
	Datetime          string
	Checkpoint        bool
	Message           string
	parentNodeAddress *Commit
	ChildNodeAddress  []*Commit
}

func (c *Commit) GetParentNodeAddress() *Commit {
	return c.parentNodeAddress
}
