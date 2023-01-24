extends Spatial

onready var messageBroker : Node = $MessageBroker
onready var player : KinematicBody = $MessageBroker/Player
onready var camera : Camera = $MessageBroker/Player/CameraBase/Camera

func _input(event):
	if not event is InputEventMouseButton or not event.pressed: return
	if event.button_index != BUTTON_LEFT: return

	var from = camera.project_ray_origin(event.position)
	var to = from + camera.project_ray_normal(event.position) * 500
	var res = get_world().direct_space_state.intersect_ray(from, to, [player])
	if not res: return

	var pos = res.position
	pos.x = round(res.position.x)
	pos.z = round(res.position.z)
	pos.y = round(res.position.y)

	# TODO - I don't like this, but it's really simple
	var attacking = false
	if res.collider.has_method("get_id"):
		attacking = true
		player.set_attacking(pos)
		messageBroker.player_requested_attack(res.collider.get_id()) # TODO - send pos to server for additional validation

	if not attacking:
		messageBroker.player_requested_move(pos)
