package ifstat

func (i InterfaceNotExists) Error() string {
	return "Interface " + string(i) + " does not exists"
}
