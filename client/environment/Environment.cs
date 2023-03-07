using Proto;
using Godot;
using System.Collections.Generic;

public partial class Environment : Node
{
    const int CAMERA_RAY_CAST_DISTANCE = 1000;
    private Camera3D camera;
    private MessageBroker mb;

    private Godot.Rid mapRid;

    public override void _Ready()
    {
        this.camera = this.GetParent().GetNode<Camera3D>("CameraBase/Camera");
        this.mb = this.GetParent().GetNode<MessageBroker>("MessageBroker");
        // TODO - Doing this allows us to pre-calculate the route, but breaks pathfinding for characters
        // this.mapRid = NavigationServer3D.MapCreate();
        // NavigationServer3D.MapSetUp(mapRid, Vector3.Up);
        // NavigationServer3D.MapSetActive(mapRid, true);
        // var navRegionRid = this.GetNode<NavigationRegion3D>("NavRegion").GetRegionRid();
        // NavigationServer3D.RegionSetMap(navRegionRid, mapRid);
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

        //var path = this.calculatePath(players.LocalPlayer.Position, targetVec3);
        this.mb.PlayerRequestedMove(new Vector3[] { targetVec3 });
    }

    private Vector3[] calculatePath(Vector3 from, Vector3 to)
    {
        var pathVectors = NavigationServer3D.MapGetPath(this.mapRid, from, to, false);
        var finalPathVectors = new List<Vector3>();
        var last = Vector3.Inf;
        foreach (var vec in pathVectors)
        {
            var vecRounded = vec.Round();
            if (!vecRounded.IsEqualApprox(last))
            {
                finalPathVectors.Add(vecRounded);
            }
        }
        return finalPathVectors.ToArray();
    }
}
