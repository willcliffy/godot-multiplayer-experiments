extends KinematicBody

onready var agent : NavigationAgent = $NavigationAgent
onready var animation_tree: AnimationTree = $AnimationTree

onready var HEADMESH : MeshInstance = $Robot/RobotArmature/Skeleton/BoneAttachment2/Head

const SPEED = 3
const ACCEPTABLE_DIST_TO_TARGET_RANGE = 0.05

var id
var team
var moving = false
var moving_to_attack = false

func set_id(new_id):
	id = new_id
	
func get_id():
	return id

func set_team(t):
	if t == 1: # Red team
		var material = HEADMESH.mesh.surface_get_material(0).duplicate()
		material.albedo_color = Color(1, 0.25, 0.25)
		HEADMESH.set_surface_material(0, material)
		team = t
	elif t == 2: # Blue team
		var material = HEADMESH.mesh.surface_get_material(0).duplicate()
		material.albedo_color = Color(0.25, 0.25, 1)
		HEADMESH.set_surface_material(0, material)
		team = t
	else: #rainbow party bitch
		var material = HEADMESH.mesh.surface_get_material(0).duplicate()
		material.albedo_color = Color(randf(), randf(), randf())
		HEADMESH.set_surface_material(0, material)
		team = t

func _physics_process(delta):
	if Input.is_action_just_pressed("exit"):
		get_tree().quit()

	if not moving:
		return

	var next = agent.get_next_location()
	if not next:
		set_idle()
		return

	var dist_to_target = (next - translation).length()
	if dist_to_target <= ACCEPTABLE_DIST_TO_TARGET_RANGE:
		set_idle()
		return
	
	if moving_to_attack:
		print("nyi should be attacking")
		set_idle()
		return

	var direction = (next - translation).normalized()
	var _collision = move_and_collide(direction * delta * SPEED)

	var facing_direction = (translation - direction)
	facing_direction.y = translation.y
	$Robot.look_at(facing_direction, Vector3.UP)

func set_moving(location):
	moving = true
	agent.set_target_location(location)
	$AnimationTree.get("parameters/playback").travel("walk")

func set_moving_to_attack():
	moving = true
	moving_to_attack = true
	
func set_attacking():
	moving = false
	moving_to_attack = false
	$AnimationTree.get("parameters/playback").travel("punch")

func set_idle():
	moving = false
	moving_to_attack = false
	$AnimationTree.get("parameters/playback").travel("idle")
	$"../Target".on_arrived()
