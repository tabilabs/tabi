package types

func NewNode(id, divisionId, owner string) Node {
	return Node{
		Id:         id,
		DivisionId: divisionId,
		Owner:      owner,
	}
}
