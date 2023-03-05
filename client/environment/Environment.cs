using Proto;
using Godot;

public partial class Environment : Node
{
    const int CAMERA_RAY_CAST_DISTANCE = 1000;
    private Camera3D camera;
    private MessageBroker mb;

    public override void _Ready()
    {
        this.camera = this.GetParent().GetNode<Camera3D>("CameraBase/Camera");
        this.mb = this.GetParent().GetNode<MessageBroker>("MessageBroker");
    }

    public override void _Input(InputEvent @event)
    {
        if (camera == null) return;
        if (!(@event is InputEventMouseButton eventKey)) return;
        // Move character
        if (!eventKey.IsPressed() || eventKey.ButtonIndex != MouseButton.Left) return;

        var from = camera.ProjectRayOrigin(eventKey.Position);
        var to = from + camera.ProjectRayNormal(eventKey.Position) * CAMERA_RAY_CAST_DISTANCE;
        var exclude = new Godot.Collections.Array<Rid>();
        var players = this.mb.GetNode<PlayerController>("PlayerController");
        exclude.Add(players.LocalPlayerRid);
        var param = PhysicsRayQueryParameters3D.Create(from, to, exclude: exclude);

        // TODO - spagooti
        var result = this.GetParent<Node3D>().GetWorld3D().DirectSpaceState.IntersectRay(param);
        if (result == null || !result.ContainsKey("position")) return;

        var targetVec3 = (Vector3)result["position"];
        var targetLocation = new Location()
        {
            X = Mathf.RoundToInt(targetVec3.X),
            Z = Mathf.RoundToInt(targetVec3.Z),
        };


        var targetPlayer = result["collider"].AsGodotObject();
        if (targetPlayer is Player)
        {
            this.mb.PlayerRequestedAttack(((Player)targetPlayer).Id);
            return;
        }

        this.mb.PlayerRequestedMove(targetLocation);
    }
}
