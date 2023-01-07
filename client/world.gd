extends Spatial

onready var player_nav_agent : NavigationAgent = $Player/NavigationAgent
onready var camerabase : Spatial = $Player/CameraBase
onready var camera : Camera = $Player/CameraBase/Camera

var show_path = true
var cam_rot = 0.0

const CAMERA_X_ZOOM_RATIO = 0.10
const CAMERA_Y_ZOOM_RATIO = 0.975

func _unhandled_input(event):
	if event is InputEventMouseButton and event.pressed:
		if event.button_index == BUTTON_LEFT:
			var from = camera.project_ray_origin(event.position)
			var to = from + camera.project_ray_normal(event.position) * 500
			var res = get_world().direct_space_state.intersect_ray(from, to, [$Player])
			if res: 
				player_nav_agent.set_target_location(res.position)
				$NavMesh/Target.translation = res.position + Vector3(0, 0.25, 0)
				$Player.set_moving()
		elif event.button_index == BUTTON_WHEEL_UP:
			camera.translation.x /= CAMERA_X_ZOOM_RATIO
			camera.translation.y *= CAMERA_Y_ZOOM_RATIO
		elif event.button_index == BUTTON_WHEEL_DOWN:
			camera.translation.x *= CAMERA_X_ZOOM_RATIO
			camera.translation.y /= CAMERA_Y_ZOOM_RATIO

	if event is InputEventMouseMotion:
		if event.button_mask & (BUTTON_MASK_MIDDLE + BUTTON_MASK_RIGHT):
			cam_rot += event.relative.x * 0.005
			camerabase.set_rotation(Vector3(0, cam_rot, 0))
