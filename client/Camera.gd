extends Camera

const ray_length = 1000

func _input(event):
	if event is InputEventMouseButton and event.pressed and event.button_index == 1:
		var from = project_ray_origin(event.position)
		var to = from + project_ray_normal(event.position) * ray_length
		print(from, to)
	
