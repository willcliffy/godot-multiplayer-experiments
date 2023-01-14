extends KinematicBody

onready var agent : NavigationAgent = $NavigationAgent
onready var animation_tree: AnimationTree = $AnimationTree

onready var HEADMESH : MeshInstance = $Robot/RobotArmature/Skeleton/BoneAttachment2/Head

const SPEED = 3

var material : SpatialMaterial = null;

var moving = false

func set_team(team):
	if team == 1: # Red team
		material = HEADMESH.get_surface_material(1)#.duplicate()
		print(material.albedo_color)
		HEADMESH.set_surface_material(1, material)
	else: # blue team
		pass

func _physics_process(delta):
	if Input.is_action_just_pressed("exit"):
		get_tree().quit()

	if not moving:
		return

	var next = agent.get_next_location()
	if not next or (next - translation).length() < 0.05:
		moving = false
		$AnimationTree.get("parameters/playback").travel("idle")
		$"../Target".on_arrived()
		return

	var direction = (next - translation).normalized()
	var _collision = move_and_collide(direction * delta * SPEED)

	var facing_direction = (translation - direction)
	facing_direction.y = translation.y
	$Robot.look_at(facing_direction, Vector3.UP)

func set_moving(position):
	moving = true
	agent.set_target_location(position)
	$AnimationTree.get("parameters/playback").travel("walk")
