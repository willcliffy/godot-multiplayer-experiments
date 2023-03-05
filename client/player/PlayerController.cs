using Game;
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
            //p.Scale = new Vector3(0.25f, 0.25f, 0.25f); // TODO - do this in blender instead for better visual effect
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
