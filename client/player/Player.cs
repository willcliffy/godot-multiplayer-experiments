using Game;
using Godot;

public class Player : KinematicBody
{
    const int MAX_HP = 10;
    const int SPEED = 3;
    const float ACCEPTABLE_DIST_TO_TARGET_RANGE = 0.05f;
    const float ATTACK_RANGE = 1 + 2 * ACCEPTABLE_DIST_TO_TARGET_RANGE;

    public ulong id { get; set; }
    public int hp { get; set; }

    private Spatial model;
    private NavigationAgent nav;
    private MessageBroker mb;
    private AnimationNodeStateMachinePlayback animations;
    private MeshInstance healthBar;
    private MeshInstance healthBarBase;
    private CollisionShape collider;

    private bool moving;
    private bool attacking;
    private bool alive;

    private Target target;
    private Vector3 targetLocation;
    private ulong? targetPlayerId;


    public override void _Ready()
    {
        this.model = GetNode<Spatial>("Robot");
        this.nav = GetNode<NavigationAgent>("NavigationAgent");

        var animationNode = GetNode<AnimationTree>("AnimationTree");
        this.animations = (AnimationNodeStateMachinePlayback)animationNode.Get("parameters/playback");

        this.alive = true;
        this.hp = MAX_HP;
        this.healthBar = GetNode<MeshInstance>("Robot/HealthBar/Health");
        this.healthBarBase = GetNode<MeshInstance>("Robot/HealthBar/Base");

        this.collider = GetNode<CollisionShape>("CollisionShape");

        this.collider.Disabled = true;

        // TODO - hacky way to check if this is the local player
        if (Visible)
        {
            this.mb = GetParent<MessageBroker>();
            this.target = mb.GetNode<Target>("Target");
        }
    }

    public override void _PhysicsProcess(float delta)
    {
        if (Input.IsActionJustPressed("exit")) GetTree().Quit();
        if (!this.moving) return;

        var distToTarget = (this.targetLocation - this.Translation).Length();
        if (attacking && distToTarget < ATTACK_RANGE)
        {
            this.setAttackingTargetReached();
            return;
        }

        if (distToTarget <= ACCEPTABLE_DIST_TO_TARGET_RANGE)
        {
            this.setIdle();
            return;
        }

        var next = this.nav.GetNextLocation();
        if (next == null)
        {
            this.setIdle();
            return;
        }

        var direction = (next - Translation).Normalized();
        var _collision = MoveAndCollide(direction * delta * SPEED);
        var facingDirection = Translation - direction;
        facingDirection.y = Translation.y;
        this.model.LookAt(facingDirection, Vector3.Up);
    }

    public void SetTeam(string color)
    {
        // TODO HACKY
        var head = GetNode<MeshInstance>("Robot/RobotArmature/Skeleton/BoneAttachment2/Head");
        var material = (SpatialMaterial)head.Mesh.SurfaceGetMaterial(0).Duplicate();
        material.AlbedoColor = new Color(color);
        head.SetSurfaceMaterial(0, material);
    }

    public void SetMoving(Vector3 target)
    {
        if (!this.alive) return;
        this.moving = true;
        this.attacking = false;
        this.targetPlayerId = null;
        this.targetLocation = target;
        this.nav.SetTargetLocation(target);
        this.animations.Travel("walk");
        this.target?.SetLocation(target);
    }

    public void SetAttacking(ulong targetPlayerId, Vector3 targetLocation)
    {
        if (!this.alive) return;
        var currentLocation = this.Translation;
        this.nav.SetTargetLocation(targetLocation);
        if (!this.nav.IsTargetReachable())
        {
            GD.Print($"FAILED TO SET ATTACKING! cannot reach {targetLocation} from {currentLocation}");
            this.nav.SetTargetLocation(currentLocation);
            return;
        }

        this.moving = true;
        this.attacking = true;
        this.targetPlayerId = targetPlayerId;
        this.targetLocation = targetLocation;
        this.animations.Travel("walk");
        this.target?.SetLocation(targetLocation, attacking = true);
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
        this.target?.OnArrived();

        if (targetPlayerId.HasValue)
        {
            this.mb?.PlayerRequestedDamage(this.targetPlayerId.GetValueOrDefault());
        }
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
        this.animations.Travel("idle");
        this.target?.OnArrived();
    }

    public Location CurrentLocation()
    {
        return new Location()
        {
            x = (uint)Mathf.RoundToInt(this.Translation.x),
            z = (uint)Mathf.RoundToInt(this.Translation.z),
        };
    }

    public void ApplyDamage(int amount)
    {
        if (!this.alive) return;
        if (this.hp <= amount) return;

        this.hp -= amount;
        var hpMesh = (CapsuleMesh)this.healthBar.Mesh.Duplicate();
        hpMesh.MidHeight = 1.0f * this.hp / MAX_HP;
        this.healthBar.Mesh = hpMesh;
        this.healthBar.Translation -= new Vector3(0.5f * amount / MAX_HP, 0, 0);
    }

    public void Die()
    {
        this.alive = false;
        this.hp = 0;
        this.healthBar.Visible = false;
        this.healthBar.Translation = new Vector3(0, 1.75f, 0);
        this.healthBarBase.Visible = false;
        this.collider.Disabled = true;
        this.animations.Travel("death");
    }

    public void Spawn(Location spawn)
    {
        this.Visible = true;
        this.alive = true;
        this.hp = MAX_HP;
        this.healthBar.Visible = true;
        this.healthBarBase.Visible = true;
        this.collider.Disabled = false;
        this.Translation = new Vector3(spawn.x, 0, spawn.z);
        ((CapsuleMesh)this.healthBar.Mesh).MidHeight = 1.0f;
        this.setIdle();

        // TODO - hacky, might need this to refresh the hp bar meshes
        //this.ApplyDamage(0);
    }
}
