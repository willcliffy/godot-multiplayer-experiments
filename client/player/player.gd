extends KinematicBody

onready var agent : NavigationAgent = $NavigationAgent
onready var animation: AnimationTree = $AnimationTree

onready var HEADMESH : MeshInstance = $Robot/RobotArmature/Skeleton/BoneAttachment2/Head

const SPEED = 3
const ACCEPTABLE_DIST_TO_TARGET_RANGE = 0.05
const ATTACK_RANGE = 2 + ACCEPTABLE_DIST_TO_TARGET_RANGE

var id
var moving = false
var attacking = false
var target_translation = null

func set_id(new_id):
	id = new_id

func get_id():
	return id

func set_team(_t, hex):
	var material = HEADMESH.mesh.surface_get_material(0).duplicate()
	material.albedo_color = Color(hex)
	HEADMESH.set_surface_material(0, material)

func _physics_process(delta):
	if Input.is_action_just_pressed("exit"):
		get_tree().quit()

	if not moving or not target_translation:
		return

	var dist_to_target = (target_translation - translation).length()
	if attacking && dist_to_target < ATTACK_RANGE:
		set_attacking_target_reached()
		return

	if dist_to_target <= ACCEPTABLE_DIST_TO_TARGET_RANGE:
		set_idle()
		return

	var next = agent.get_next_location()
	if not next:
		set_idle()
		return

	var direction = (next - translation).normalized()
	var _collision = move_and_collide(direction * delta * SPEED)

	var facing_direction = (translation - direction)
	facing_direction.y = translation.y
	$Robot.look_at(facing_direction, Vector3.UP)

func set_moving(location):
	moving = true
	attacking = false
	agent.set_target_location(location)
	target_translation = location
	animation.get("parameters/playback").travel("walk")

func set_attacking(location):
	print("attack started")
	moving = true
	attacking = true
	agent.set_target_location(location)
	target_translation = location
	animation.get("parameters/playback").travel("walk")

func set_attacking_target_reached():
	moving = false
	attacking = false
	target_translation = null
	animation.get("parameters/playback").travel("attack")

func set_idle():
	moving = false
	attacking = false
	target_translation = null
	animation.get("parameters/playback").travel("idle")
