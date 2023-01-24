extends KinematicBody

onready var agent : NavigationAgent = $NavigationAgent
onready var animation_tree: AnimationTree = $AnimationTree

onready var HEADMESH : MeshInstance = $Robot/RobotArmature/Skeleton/BoneAttachment2/Head

const SPEED = 3
const ATTACK_RANGE = 0.05 # todo - why doesnt this work as I expect it to? attack range should be 1 but that makes the character stop short
const ACCEPTABLE_DIST_TO_TARGET_RANGE = 0.05

var id
var moving = false
var attacking = false
var target = null

func set_id(new_id):
	id = new_id

func get_id():
	return id

func set_team(t, hex):
	var material = HEADMESH.mesh.surface_get_material(0).duplicate()
	material.albedo_color = Color(hex)
	HEADMESH.set_surface_material(0, material)

func _physics_process(delta):
	if Input.is_action_just_pressed("exit"):
		get_tree().quit()

	if not moving:
		return

	var dist_to_target = (target - translation).length()
	if attacking && dist_to_target < ATTACK_RANGE:
		print(dist_to_target)
		print(target)
		print(translation)
		print("nyi should be attacking")
		set_idle()
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
	target = location
	print(target)
	print(translation)
	$AnimationTree.get("parameters/playback").travel("walk")

func set_attacking(location):
	moving = true
	attacking = true

func set_attacking_target_reached():
	moving = false
	attacking = false
	$AnimationTree.get("parameters/playback").travel("punch")

func set_idle():
	moving = false
	attacking = false
	$AnimationTree.get("parameters/playback").travel("idle")
	$"../Target".on_arrived()
