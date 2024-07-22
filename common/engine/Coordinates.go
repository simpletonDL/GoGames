package engine

type CoordinatesMapper struct {
	worldWidth   float64
	worldHeight  float64
	screenWidth  float64
	screenHeight float64
}

func NewCoordinatesMapper(worldWidth, worldHeight, screenWidth, screenHeight float64) CoordinatesMapper {
	return CoordinatesMapper{worldWidth, worldHeight, screenWidth, screenHeight}
}

func (cm *CoordinatesMapper) WorldToScreenX(worldX float64) float64 {
	return (worldX / cm.worldWidth) * cm.screenWidth
}

func (cm *CoordinatesMapper) WorldToScreenY(worldY float64) float64 {
	return cm.screenHeight - (worldY/cm.worldHeight)*cm.screenHeight
}

func (cm *CoordinatesMapper) ScreenToWorldX(screenX float64) float64 {
	return (screenX / cm.screenWidth) * cm.worldWidth
}

func (cm *CoordinatesMapper) ScreenToWorldY(screenY float64) float64 {
	return cm.worldHeight - (screenY/cm.screenHeight)*cm.worldHeight
}

func (cm *CoordinatesMapper) ScreenToWorld(screenX, screenY float64) (float64, float64) {
	worldX := cm.ScreenToWorldX(screenX)
	worldY := cm.ScreenToWorldY(screenY)
	return worldX, worldY
}
