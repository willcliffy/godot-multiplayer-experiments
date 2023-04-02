using Godot;

public partial class CameraBase : Node3D
{
	private const float CAMERA_MIN_ZOOM = 5;
	private const float CAMERA_MAX_ZOOM = 50;

	private Vector3 CAMERA_ZOOM_SPEED = new Vector3(0, 2, 1);
	private Camera3D camera;

	public override void _Ready()
	{
		this.camera = this.GetNode<Camera3D>("Camera");
	}

	public override void _Input(InputEvent @event)
	{
		// TODO - allow camera rotation?

		if (!@event.IsPressed() || !(@event is InputEventMouseButton e)) return;
		if (e.ButtonIndex == MouseButton.WheelUp && camera.Position.Y > CAMERA_MIN_ZOOM)
		{
			camera.Position -= CAMERA_ZOOM_SPEED;
		}
		else if (e.ButtonIndex == MouseButton.WheelDown && camera.Position.Y < CAMERA_MAX_ZOOM)
		{
			camera.Position += CAMERA_ZOOM_SPEED;
		}
	}
}
