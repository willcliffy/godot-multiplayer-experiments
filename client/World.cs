using Godot;
using Game;

public class World : Spatial
{
    const int RAY_TRACE_DISTANCE = 500;

    MessageBroker mb;
    Player player;
    Camera camera;

    public override void _Ready()
    {
        mb = GetNode<MessageBroker>("MessageBroker");
        player = GetNode<Player>("MessageBroker/Player");
        camera = GetNode<Camera>("MessageBroker/Player/CameraBase/Camera");
    }


    public override void _Input(InputEvent @event)
    {
        if (!(@event is InputEventMouseButton eventKey) || !@event.IsPressed()) return;
        if (eventKey.ButtonIndex != (int)ButtonList.Left) return;

        var from = camera.ProjectRayOrigin(eventKey.Position);
        var to = from + camera.ProjectRayNormal(eventKey.Position) * RAY_TRACE_DISTANCE;
        var exclude = new Godot.Collections.Array();
        exclude.Add(player);

        var result = GetWorld().DirectSpaceState.IntersectRay(from, to, exclude: exclude);
        if (result == null) return;

        var targetVec3 = (Vector3)result["position"];
        var targetLocation = new Location()
        {
            x = (uint)Mathf.RoundToInt(targetVec3.x),
            z = (uint)Mathf.RoundToInt(targetVec3.z),
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
