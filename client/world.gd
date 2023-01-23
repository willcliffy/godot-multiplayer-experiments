extends Spatial

onready var camera : Camera = $Player/CameraBase/Camera

func _input(event):
	if not event is InputEventMouseButton or not event.pressed: return
	if event.button_index != BUTTON_LEFT: return

	var from = camera.project_ray_origin(event.position)
	var to = from + camera.project_ray_normal(event.position) * 500
	var res = get_world().direct_space_state.intersect_ray(from, to, [$Player])
	if not res: return

	var pos = res.position
	pos.x = round(res.position.x)
	pos.z = round(res.position.z)
	pos.y = round(res.position.y)

	var attacking = false
	# TODO - I don't like this, but it's really simple
	if res.collider.has_method("get_id"):
		attacking = true
		$Player.set_moving_to_attack()
		$MessageBroker.player_requested_attack(res.collider.get_id()) # TODO - send pos to server for additional validation

	$Target.set_target(pos)
	if not attacking:
		$MessageBroker.player_requested_move(pos)
