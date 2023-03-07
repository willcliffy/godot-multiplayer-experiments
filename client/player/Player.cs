using Proto;
using Godot;

public partial class Player : CharacterBody3D
{
    const int MAX_HP = 10;
    const int SPEED = 1000;

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

    private Node3D model;
    public NavigationAgent3D Nav;
    private AnimationNodeStateMachinePlayback animations;
    private MeshInstance3D healthBar;
    private MeshInstance3D healthBarBase;
    private CollisionShape3D collider;

    private bool moving;
    private bool attacking;
    private bool alive;

    private ulong? targetPlayerId;

    public override void _Ready()
    {
        this.model = GetNode<Node3D>("Robot");
        this.Nav = GetNode<NavigationAgent3D>("NavigationAgent");

        var animationNode = GetNode<AnimationTree>("AnimationTree");
        this.animations = (AnimationNodeStateMachinePlayback)animationNode.Get("parameters/playback");

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
            GD.Print("t reached");
            if (attacking) this.setAttackingTargetReached();
            else if (moving) this.setIdle();
            return;
        }

        // TODO - why is this coming up zero when nav.Target is getting set properly
        var next = this.Nav.GetNextPathPosition();

        var direction = this.GlobalPosition.DirectionTo(next);
        var velocity = direction * (float)delta * SPEED;
        this.Nav.SetVelocity(velocity);
        this.Velocity = velocity;
        this.MoveAndSlide();
        GD.Print($"moving towards {this.Nav.GetCurrentNavigationPath().Length} with velocity. {this.Nav.TargetPosition}");
    }

    public void SetTeam(string color)
    {
        // TODO HACKY
        var head = GetNode<MeshInstance3D>("Robot/RobotArmature/Skeleton3D/Head_2/Head_2");
        var material = (StandardMaterial3D)head.Mesh.SurfaceGetMaterial(0).Duplicate();
        material.AlbedoColor = new Color(color);
        head.SetSurfaceOverrideMaterial(0, material);
    }

    public void SetMoving(Vector3 targetVec3)
    {
        if (!this.alive) return;
        this.moving = true;
        this.attacking = false;
        this.targetPlayerId = null;
        this.Nav.TargetPosition = targetVec3;
        this.animations.Travel("run");
        GD.Print($"set moving towards {targetVec3}");
    }

    public void SetAttacking(ulong targetPlayerId, Vector3 targetVec3)
    {
        if (!this.alive) return;
        this.moving = true;
        this.attacking = true;
        this.targetPlayerId = targetPlayerId;
        this.Nav.TargetPosition = targetVec3;
        this.animations.Travel("run");
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

    public void PlayAttackingAnimation()
    {
        this.animations.Travel("punch");
    }

    private void setIdle()
    {
        if (!this.alive) return;
        this.moving = false;
        this.attacking = false;
        this.targetPlayerId = null;
        this.Nav.SetVelocity(Vector3.Zero);
        this.Velocity = Vector3.Zero;
        this.animations.Travel("idle");
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
        this.animations.Travel("death");
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
