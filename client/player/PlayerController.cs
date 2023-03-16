using Godot;
using System.Collections.Generic;

public partial class PlayerController : Node
{
    public Player LocalPlayer;

    public ulong LocalPlayerId
    {
        get { return LocalPlayer.Id; }
    }

    public Godot.Rid LocalPlayerRid
    {
        get { return LocalPlayer.GetRid(); }
    }

    private Dictionary<ulong, Player> players = new Dictionary<ulong, Player>();

    private Target target;

    public override void _Ready()
    {
        this.LocalPlayer = this.GetNode<Player>("Player");
        this.target = this.GetParent().GetNode<Target>("Target");
    }

    public Player OnLocalPlayerJoined(Proto.JoinGameResponse msg)
    {
        this.LocalPlayer = this.OnPlayerConnected(msg.PlayerId, msg.Spawn, msg.Color);
        this.players[msg.PlayerId] = this.LocalPlayer;
        if (msg.Others != null)
        {
            foreach (var other in msg.Others)
            {
                var connect = Proto.Connect.Parser.ParseFrom(other.Payload);
                this.OnPlayerConnected(other.PlayerId, connect.Spawn, connect.Color);
            }
        }

        return this.LocalPlayer;
    }

    public Player OnPlayerConnected(ulong id, Proto.Location spawn, string color)
    {
        Player p;
        if (this.players.ContainsKey(id))
        {
            p = this.players[id];
        }
        else if (id != this.LocalPlayerId)
        {
            p = this.LocalPlayer;
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
        return p.CurrentLocation;
    }

    public void SetMoving(ulong id, Vector3I target)
    {
        if (!this.players.TryGetValue(id, out var p)) return;
        if (id == this.LocalPlayerId) this.target.SetLocation(target);
        p.SetMoving(target);
    }

    // public void SetAttacking(ulong sourcePlayerId, ulong targetPlayerId, Proto.Location target)
    // {
    //     if (!this.players.TryGetValue(sourcePlayerId, out var p)) return;
    //     var targetVec3 = new Vector3(target.X, 0, target.Z);
    //     p.SetAttacking(targetPlayerId, targetVec3);
    //     if (sourcePlayerId == this.LocalPlayerId) this.target.SetLocation(targetVec3, true);
    // }

    public void StopAttacking(ulong playerId)
    {
        foreach (var kvp in this.players)
        {
            if (kvp.Value.IsAttacking(playerId))
            {
                kvp.Value.StopAttacking();
            }
        }
    }

    public void ApplyDamage(ulong playerId, int amount)
    {
        if (!this.players.TryGetValue(playerId, out var p)) return;
        p.ApplyDamage(amount);
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
