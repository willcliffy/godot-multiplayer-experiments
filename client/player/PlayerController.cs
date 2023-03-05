using Godot;
using System.Collections.Generic;

public partial class PlayerController : Node
{
    private Player localPlayer;

    public ulong LocalPlayerId
    {
        get { return localPlayer.Id; }
    }

    public Godot.Rid LocalPlayerRid
    {
        get { return this.players[this.LocalPlayerId].GetRid(); }
    }

    private Dictionary<ulong, Player> players = new Dictionary<ulong, Player>();

    private Target target;

    public override void _Ready()
    {
        this.target = this.GetParent().GetNode<Target>("Target");
    }

    public Player OnLocalPlayerJoined(Proto.JoinGameResponse msg)
    {
        this.localPlayer = this.OnPlayerConnected(msg.PlayerId, msg.Spawn, msg.Color);
        if (msg.Others != null)
        {
            foreach (var other in msg.Others)
            {
                this.OnPlayerConnected(other.PlayerId, other.Spawn, other.Color);
            }
        }

        return this.localPlayer;
    }

    public Player OnPlayerConnected(ulong id, Proto.Location spawn, string color)
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

    public Proto.Location CurrentLocation(ulong id)
    {
        if (!this.players.TryGetValue(id, out var p)) return null;
        return p.CurrentLocation();
    }

    public void SetMoving(ulong id, Proto.Location target)
    {
        if (!this.players.TryGetValue(id, out var p)) return;
        var targetVec3 = new Vector3(target.X, 0, target.Z);
        p.SetMoving(targetVec3);
        if (id == this.LocalPlayerId) this.target.SetLocation(targetVec3);
    }

    public void SetAttacking(ulong sourcePlayerId, ulong targetPlayerId, Proto.Location target)
    {
        if (!this.players.TryGetValue(sourcePlayerId, out var p)) return;
        var targetVec3 = new Vector3(target.X, 0, target.Z);
        p.SetAttacking(targetPlayerId, targetVec3);
        if (sourcePlayerId == this.LocalPlayerId) this.target.SetLocation(targetVec3, true);
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

    public void Spawn(ulong playerId, Proto.Location spawn)
    {
        if (!this.players.TryGetValue(playerId, out var p)) return;
        p.Spawn(spawn);
    }
}
