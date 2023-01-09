extends Spatial

var cam_rot = 0.0

const CAMERA_X_ZOOM_RATIO = 0.10
const CAMERA_Y_ZOOM_RATIO = 0.975

func _ready():
	#look_at(translation, Vector3.UP)
	pass

func _unhandled_input(event):
	if event is InputEventMouseButton and event.pressed:
		pass # TODO - fix camera zoom with mouse wheel

	if event is InputEventMouseMotion:
		if event.button_mask & (BUTTON_MASK_MIDDLE + BUTTON_MASK_RIGHT):
			cam_rot += event.relative.x * 0.005
			set_rotation(Vector3(0, cam_rot, 0))
