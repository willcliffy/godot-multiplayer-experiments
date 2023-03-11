using Proto;
using Godot;
using System;

public partial class Player : CharacterBody3D
{
    const int MAX_HP = 10;
    const int SPEED = 300;

    public ulong Id { get; set; }

    public int Hp { get; set; }

    public Proto.Location CurrentLocation
    {
        get
        {
            return new Location()
            {
                X = Mathf.RoundToInt(this.Position.X),
                Z = Mathf.RoundToInt(this.Position.Z),
            };
        }
    }

    public NavigationAgent3D Nav;
    private MeshInstance3D healthBar;
    private MeshInstance3D healthBarBase;
    private CollisionShape3D collider;

    private bool moving;
    private bool attacking;
    private bool alive;

    private ulong? targetPlayerId;

    public override void _Ready()
    {
        this.Nav = GetNode<NavigationAgent3D>("NavigationAgent");

        this.alive = true;
        this.Hp = MAX_HP;
        this.healthBar = GetNode<MeshInstance3D>("HealthBar/Health");
        this.healthBarBase = GetNode<MeshInstance3D>("HealthBar/Base");

        this.collider = GetNode<CollisionShape3D>("CollisionShape");
        this.collider.Disabled = true;
    }

    public override void _PhysicsProcess(double delta)
    {
        if (!this.moving) return;
        if (this.Nav.IsTargetReached())
        {
            if (attacking) this.setAttackingTargetReached();
            else if (moving) this.setIdle();
            return;
        }

        var next = this.Nav.GetNextPathPosition();
        var direction = this.GlobalPosition.DirectionTo(next);
        var velocity = direction * (float)delta * SPEED;
        //GD.Print(velocity.Length());
        this.Nav.SetVelocity(velocity);
        this.Velocity = velocity;
        this.MoveAndSlide();
    }

    public void SetTeam(string color)
    {
        // TODO HACKY
        var mesh = GetNode<MeshInstance3D>("Mesh");
        var material = (StandardMaterial3D)mesh.Mesh.SurfaceGetMaterial(0).Duplicate();
        material.AlbedoColor = new Color(color);
        mesh.SetSurfaceOverrideMaterial(0, material);
    }

    public Vector3[] SetMoving(Vector3 targetVec3)
    {
        if (!this.alive) return null;
        this.moving = true;
        this.attacking = false;
        this.targetPlayerId = null;
        this.Nav.TargetPosition = targetVec3;
        this.Nav.GetNextPathPosition();
        return this.Nav.GetCurrentNavigationPath();
    }

    public void SetAttacking(ulong targetPlayerId, Vector3 targetVec3)
    {
        if (!this.alive) return;
        this.moving = true;
        this.attacking = true;
        this.targetPlayerId = targetPlayerId;
        this.Nav.TargetPosition = targetVec3;
    }

    public bool IsAttacking(ulong? playerId)
    {
        return this.attacking && playerId != null && playerId == this.targetPlayerId;
    }

    public void StopAttacking()
    {
        if (!attacking) return;
        this.setIdle();
    }

    private void setAttackingTargetReached()
    {
        if (!this.alive) return;
        this.moving = false;
        this.attacking = false;

        // TODO - Player should not be responsible for this - this should happen on the server
        // The reason it doesnt now is because the golang server doesnt understand pathfinding
        // so it doesn't know when someone has arrived when attacking.
        // if (targetPlayerId.HasValue)
        // {
        //     this.mb?.PlayerRequestedDamage(this.targetPlayerId.GetValueOrDefault());
        // }
    }

    private void setIdle()
    {
        if (!this.alive) return;
        this.moving = false;
        this.attacking = false;
        this.targetPlayerId = null;
        this.Nav.SetVelocity(Vector3.Zero);
        this.Velocity = Vector3.Zero;
    }

    public void ApplyDamage(int amount)
    {
        if (!this.alive) return;
        if (this.Hp <= amount) return;

        this.Hp -= amount;
        var hpMesh = (CapsuleMesh)this.healthBar.Mesh.Duplicate();
        hpMesh.Height = 1.0f * this.Hp / MAX_HP;
        this.healthBar.Mesh = hpMesh;
        this.healthBar.Position -= new Vector3(0.5f * amount / MAX_HP, 0, 0);
    }

    public void Die()
    {
        this.alive = false;
        this.Hp = 0;
        this.healthBar.Visible = false;
        this.healthBar.Position = new Vector3(0, 0, 0); // todo
        this.healthBarBase.Visible = false;
        this.collider.Disabled = true;
    }

    public void Spawn(Location spawn)
    {
        this.Visible = true;
        this.alive = true;
        this.Hp = MAX_HP;
        this.healthBar.Visible = true;
        this.healthBarBase.Visible = true;
        this.collider.Disabled = false;
        this.Position = new Vector3(spawn.X, 0, spawn.Z);
        ((CapsuleMesh)this.healthBar.Mesh).Height = 1.0f;
        this.setIdle();

        // TODO - hacky, might need this to refresh the hp bar meshes
        //this.ApplyDamage(0);
    }
}
