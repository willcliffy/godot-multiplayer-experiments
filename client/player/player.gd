extends KinematicBody

onready var agent : NavigationAgent = $NavigationAgent
onready var animation_tree: AnimationTree = $AnimationTree

const SPEED = 2.5

var moving = false

func _physics_process(delta):
	if Input.is_action_just_pressed("exit"):
		get_tree().quit()

	if not moving:
		return

	$CameraBase/Camera.look_at(translation, Vector3.UP)

	var next = agent.get_next_location()
	if not next or (next - translation).length() < 0.01:
		moving = false
		$AnimationTree.get("parameters/playback").travel("idle")
		return

	var direction = (next - translation).normalized()
	var _collision = move_and_collide(direction * delta * SPEED)
	$Robot.look_at(translation - direction, Vector3.UP)


func set_moving():
	moving = true
	$AnimationTree.get("parameters/playback").travel("walk")
