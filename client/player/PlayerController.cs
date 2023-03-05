using Game;
using Godot;
using System.Collections.Generic;

public partial class PlayerController : Node
{
    const int RAY_TRACE_DISTANCE = 1000;

    private Player localPlayer;

    public ulong LocalPlayerId
    {
        get { return localPlayer.Id; }
    }

    private Dictionary<ulong, Player> players = new Dictionary<ulong, Player>();

    private Target target;
    private MessageBroker mb;
    private Camera3D camera;

    public override void _Ready()
    {
        this.mb = this.GetParent<MessageBroker>();
        this.target = this.mb.GetNode<Target>("Target");
    }

    public override void _Input(InputEvent @event)
    {
        if (camera == null) return;
        if (!(@event is InputEventMouseButton eventKey)) return;
        // Move character
        if (eventKey.ButtonIndex != MouseButton.Left) return;

        var from = camera.ProjectRayOrigin(eventKey.Position);
        var to = from + camera.ProjectRayNormal(eventKey.Position) * RAY_TRACE_DISTANCE;
        var exclude = new Godot.Collections.Array<Rid>();
        exclude.Add(this.players[this.LocalPlayerId].GetRid());
        var param = PhysicsRayQueryParameters3D.Create(from, to, exclude: exclude);

        // TODO - spagooti
        var result = this.mb.GetParent<Node3D>().GetWorld3D().DirectSpaceState.IntersectRay(param);
        if (result == null || !result.ContainsKey("position")) return;

        var targetVec3 = (Vector3)result["position"];
        var targetLocation = new Location()
        {
            x = Mathf.RoundToInt(targetVec3.X),
            z = Mathf.RoundToInt(targetVec3.Z),
        };

        var targetPlayer = result["collider"].As<Player>();
        if (targetPlayer != null)
        {
            this.mb.PlayerRequestedAttack(targetPlayer.Id);
            return;
        }

        this.mb.PlayerRequestedMove(targetLocation);
    }

    public Player OnLocalPlayerJoined(JoinGameResponse msg)
    {
        this.localPlayer = this.OnPlayerConnected(msg.playerId, msg.spawn, msg.color);
        if (msg.others != null)
        {
            foreach (var other in msg.others)
            {
                this.OnPlayerConnected(other.playerId, other.spawn, other.color);
            }
        }

        return this.localPlayer;
    }

    public Player OnPlayerConnected(ulong id, Location spawn, string color)
    {
        Player p;
        if (this.players.ContainsKey(id))
        {
            p = this.players[id];
        }
        else
        {
            var scene = ResourceLoader.Load<PackedScene>("res://player/player.tscn");
            p = scene.Instantiate() as Player;
            p.Visible = true;
            p.Scale = new Vector3(0.25f, 0.25f, 0.25f); // TODO - do this in blender instead for better visual effect
            this.AddChild(p);
            this.players[id] = p;
        }

        p.Id = id;
        p.SetTeam(color);
        p.Spawn(spawn);

        return p;
    }

    public void OnPlayerDisconnected(ulong id)
    {
        if (!this.players.TryGetValue(id, out var p)) return;

        p.Visible = false;
        p.Position = new Vector3(-10, -10, -10);
        this.players.Remove(id);
    }

    public Location CurrentLocation(ulong id)
    {
        if (!this.players.TryGetValue(id, out var p)) return null;
        return p.CurrentLocation();
    }

    public void SetMoving(ulong id, Vector3 target)
    {
        if (!this.players.TryGetValue(id, out var p)) return;
        p.SetMoving(target);
        if (id == this.LocalPlayerId) this.target.SetLocation(target);
    }

    public void SetAttacking(ulong sourcePlayerId, ulong targetPlayerId, Vector3 target)
    {
        if (!this.players.TryGetValue(sourcePlayerId, out var p)) return;
        p.SetAttacking(targetPlayerId, target);
        if (sourcePlayerId == this.LocalPlayerId) this.target.SetLocation(target, true);
    }

    public void StopAttacking(ulong playerId)
    {
        foreach (var kvp in this.players)
        {
            if (kvp.Value.IsAttacking(playerId))
            {
                kvp.Value.StopAttacking();
                if (playerId == this.LocalPlayerId) this.target.OnArrived();
            }
        }
    }

    public void ApplyDamage(ulong playerId, int amount)
    {
        if (!this.players.TryGetValue(playerId, out var p)) return;
        p.ApplyDamage(amount);
    }

    public void PlayAttackingAnimation(ulong playerId)
    {
        if (!this.players.TryGetValue(playerId, out var p)) return;
        p.PlayAttackingAnimation();
    }

    public void Die(ulong playerId)
    {
        if (!this.players.TryGetValue(playerId, out var p)) return;
        p.Die();
    }

    public void Spawn(ulong playerId, Location spawn)
    {
        if (!this.players.TryGetValue(playerId, out var p)) return;
        p.Spawn(spawn);
    }
}
