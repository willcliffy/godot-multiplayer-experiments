using Game;
using Godot;
using System.Collections.Generic;

public class OpponentController : Node
{
    Queue<Player> available = new Queue<Player>();
    Dictionary<ulong, Player> inUse = new Dictionary<ulong, Player>();

    public override void _Ready()
    {
        available.Enqueue(GetNode<Player>("1"));
        available.Enqueue(GetNode<Player>("2"));
        available.Enqueue(GetNode<Player>("3"));
        available.Enqueue(GetNode<Player>("4"));
        available.Enqueue(GetNode<Player>("5"));
        available.Enqueue(GetNode<Player>("6"));
        available.Enqueue(GetNode<Player>("7"));
    }

    public void OnPlayerConnected(ulong id, Location spawn, string color)
    {
        Player p;
        if (inUse.ContainsKey(id))
        {
            p = inUse[id];
        }
        else
        {
            if (available.Count == 0)
            {
                GD.Print("SERVER BADNESS: game overflow");
                return;
            }
            p = available.Dequeue();
            inUse[id] = p;
        }

        p.id = id;
        p.SetTeam(color);
        p.Spawn(spawn);
    }

    public void OnPlayerDisconnected(ulong id)
    {
        if (!inUse.TryGetValue(id, out var p)) return;

        p.Visible = false;
        p.Translation = new Vector3(-10, -10, -10);
        inUse.Remove(id);
        available.Enqueue(p);
    }

    public Location CurrentLocation(ulong id)
    {
        if (!inUse.TryGetValue(id, out var p)) return null;
        return p.CurrentLocation();
    }

    public void SetMoving(ulong id, Vector3 target)
    {
        if (!inUse.TryGetValue(id, out var p)) return;
        p.SetMoving(target);
    }

    public void SetAttacking(ulong sourcePlayerId, ulong targetPlayerId, Vector3 target)
    {
        if (!inUse.TryGetValue(sourcePlayerId, out var p)) return;
        p.SetAttacking(targetPlayerId, target);
    }

    public void StopAttacking(ulong playerId)
    {
        foreach (var kvp in inUse)
        {
            if (kvp.Value.IsAttacking(playerId))
            {
                kvp.Value.StopAttacking();
            }
        }
    }

    public void ApplyDamage(ulong playerId, int amount)
    {
        if (!inUse.TryGetValue(playerId, out var p)) return;
        p.ApplyDamage(amount);
    }

    public void PlayAttackingAnimation(ulong playerId)
    {
        if (!inUse.TryGetValue(playerId, out var p)) return;
        p.PlayAttackingAnimation();
    }

    public void Die(ulong playerId)
    {
        if (!inUse.TryGetValue(playerId, out var p)) return;
        p.Die();
    }

    public void Spawn(ulong playerId, Location spawn)
    {
        if (!inUse.TryGetValue(playerId, out var p)) return;
        p.Spawn(spawn);
    }
}
