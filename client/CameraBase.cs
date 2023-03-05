using Godot;

public partial class CameraBase : Node3D
{
    private const float CAMERA_MIN_ZOOM = 2;
    private const float CAMERA_MAX_ZOOM = 50;

    private Vector3 CAMERA_ZOOM_SPEED = new Vector3(0, 0.5f, 0.5f);
    private Camera3D camera;

    public override void _Ready()
    {
        this.camera = this.GetNode<Camera3D>("Camera");
    }

    public override void _Input(InputEvent @event)
    {
        // TODO - allow camera rotation?

        if (!@event.IsPressed() || !(@event is InputEventMouseButton eventKey)) return;
        if (eventKey.ButtonIndex == MouseButton.WheelUp)
        {
            if (camera.Position.Y > CAMERA_MIN_ZOOM)
            {
                camera.Position -= CAMERA_ZOOM_SPEED;
            }
            return;
        }
        else if (eventKey.ButtonIndex == MouseButton.WheelDown)
        {
            if (camera.Position.Y < CAMERA_MAX_ZOOM)
            {
                camera.Position += CAMERA_ZOOM_SPEED;
            }
            return;
        }
    }
}
