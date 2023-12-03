package custom

func (m *Model) SetYOffset(n int) {
	m.viewport.YOffset = n
	m.UpdateViewport()
}

func (m *Model) YOffset() int {
	return m.viewport.YOffset
}
