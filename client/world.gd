extends Spatial

onready var camerabase : Spatial = $Player/CameraBase
onready var camera : Camera = $Player/CameraBase/Camera

var show_path = true

func _unhandled_input(event):
	if not event is InputEventMouseButton or not event.pressed:
		return
	if event.button_index != BUTTON_LEFT:
		return

	var from = camera.project_ray_origin(event.position)
	var to = from + camera.project_ray_normal(event.position) * 500
	var res = get_world().direct_space_state.intersect_ray(from, to, [$Player])
	if not res:
		return

	var pos = res.position
	pos.x = round(res.position.x)
	pos.z = round(res.position.z)
	pos.y = round(res.position.y)

	$Target.set_target(pos)
	$Player.set_moving(pos)
	$MessageBroker.on_player_move(pos)
