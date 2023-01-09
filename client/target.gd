extends CSGMesh


func _process(delta):
	rotate_y(delta)

func set_target(position):
	position.y = 0.33
	translation = position
	visible = true

func on_arrived():
	visible = false
