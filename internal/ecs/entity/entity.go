package entity

type Component interface {
    // Update is called every frame to update the component.
    Update()

    // Draw is called every frame to draw the component.
    Draw()
}

type Entity struct {
	id uint32
    components []Component
}
