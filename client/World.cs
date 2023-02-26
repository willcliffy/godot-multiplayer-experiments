using Godot;
using Game;

public class World : Spatial
{
	const int RAY_TRACE_DISTANCE = 1000;
	const float CAMERA_MIN_ZOOM = 2;
	const float CAMERA_MAX_ZOOM = 50;
	const float CAMERA_ROTATION_SPEED = 1.5f;
	Vector3 CAMERA_ZOOM_SPEED = new Vector3(0, 0.5f, 0.5f);

	MessageBroker mb;
	Player player;
	Camera camera;
	Spatial cameraBase;
	float cameraRotation;

	public override void _Ready()
	{
		mb = GetNode<MessageBroker>("MessageBroker");
		player = GetNode<Player>("MessageBroker/Player");
		camera = GetNode<Camera>("MessageBroker/Player/CameraBase/Camera");
		cameraBase = GetNode<Spatial>("MessageBroker/Player/CameraBase");
	}

	public override void _Process(float delta)
	{
		cameraBase.RotateY(cameraRotation * delta * CAMERA_ROTATION_SPEED);
	}

	public override void _Input(InputEvent @event)
	{
		// rotate camera
		if (@event.IsActionReleased("move_right") || @event.IsActionReleased("move_left"))
		{
			cameraRotation = 0;
		}
		if (!@event.IsPressed()) return;

		cameraRotation = @event.GetActionStrength("move_right") - @event.GetActionStrength("move_left");

		// Zoom Camera
		if (!(@event is InputEventMouseButton eventKey)) return;
		if (eventKey.ButtonIndex == (int)ButtonList.WheelUp)
		{
			if (camera.Translation.y > CAMERA_MIN_ZOOM)
			{
				camera.Translation -= CAMERA_ZOOM_SPEED;
			}
			return;
		}
		else if (eventKey.ButtonIndex == (int)ButtonList.WheelDown)
		{
			if (camera.Translation.y < CAMERA_MAX_ZOOM)
			{
				camera.Translation += CAMERA_ZOOM_SPEED;
			}
			return;
		}

		// Move character
		if (eventKey.ButtonIndex != (int)ButtonList.Left) return;

		var from = camera.ProjectRayOrigin(eventKey.Position);
		var to = from + camera.ProjectRayNormal(eventKey.Position) * RAY_TRACE_DISTANCE;
		var exclude = new Godot.Collections.Array();
		exclude.Add(player);

		var result = GetWorld().DirectSpaceState.IntersectRay(from, to, exclude: exclude);
		if (result == null || !result.Contains("position"))
		{
			GD.Print($"result null");
			GD.Print(result);
			return;
		}

		var targetVec3 = (Vector3)result["position"];
		var targetLocation = new Location()
		{
			x = Mathf.RoundToInt(targetVec3.x),
			z = Mathf.RoundToInt(targetVec3.z),
		};

		var targetPlayer = result["collider"] as Player;
		if (targetPlayer != null)
		{
			this.mb.PlayerRequestedAttack(targetPlayer.id);
			return;
		}

		this.mb.PlayerRequestedMove(targetLocation);
	}
}
