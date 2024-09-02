package protocol

func (k TeamKind) ToString() string {
	switch k {
	case BlueTeam:
		return "Blue Team"
	case RedTeam:
		return "Red Team"
	default:
		panic("invalid team kind")
	}
}
