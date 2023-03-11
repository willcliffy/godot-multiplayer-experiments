using Proto;
using Godot;
using System.Collections.Generic;

public partial class World : Node3D
{
    const int CAMERA_RAY_CAST_DISTANCE = 1000;
    private Camera3D camera;
    private MessageBroker mb;

    private Godot.Rid mapRid;

    public override void _Ready()
    {
        this.camera = this.GetNode<Camera3D>("CameraBase/Camera");
        this.mb = this.GetNode<MessageBroker>("MessageBroker");
    }

    public override void _Input(InputEvent @event)
    {
        if (camera == null) return;
        if (!(@event is InputEventMouseButton eventKey)) return;
        if (!eventKey.IsPressed() || eventKey.ButtonIndex != MouseButton.Left) return;

        var from = camera.ProjectRayOrigin(eventKey.Position);
        var to = from + camera.ProjectRayNormal(eventKey.Position) * CAMERA_RAY_CAST_DISTANCE;
        var exclude = new Godot.Collections.Array<Rid>();
        var players = this.mb.GetNode<PlayerController>("PlayerController");
        exclude.Add(players.LocalPlayerRid);
        var param = PhysicsRayQueryParameters3D.Create(from, to, exclude: exclude);

        // TODO - spagooti
        var result = this.GetWorld3D().DirectSpaceState.IntersectRay(param);
        if (result == null || !result.ContainsKey("position")) return;

        var targetVec3 = (Vector3)result["position"];
        var targetLocation = new Location()
        {
            X = Mathf.RoundToInt(targetVec3.X),
            Z = Mathf.RoundToInt(targetVec3.Z),
        };

        var targetCollider = result["collider"].AsGodotObject();
        if (targetCollider is Player targetPlayer)
        {
            this.mb.PlayerRequestedAttack(targetPlayer.Id);
            return;
        }
        else if (targetCollider is Resource targetResource)
        {
            // TODO - find closest point to player from resource and make that the target
            GD.Print(!targetResource.IsDepleted);
            targetResource.Collect();
            this.mb.PlayerRequestedCollect(targetVec3, targetResource.Type);
            return;
        }

        this.mb.PlayerRequestedMove(targetVec3);
    }
}
