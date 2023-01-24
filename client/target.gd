extends CSGMesh

func _process(delta):
	rotate_y(delta)

func set_target_location(location, _attacking):
	translation = location + Vector3(0, 0.333, 0)
	visible = true

func on_arrived():
	visible = false
