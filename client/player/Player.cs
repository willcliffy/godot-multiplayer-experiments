using Game;
using Godot;

public class Player : KinematicBody
{
    const int MAX_HP = 20;
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
        this.healthBar = GetNodeOrNull<MeshInstance>("HealthBar/Health");

        // TODO - get health bar to change size when taking damage
        this.ApplyDamage(1);

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
        targetLocation.y = Translation.y;
        this.model.LookAt(-targetLocation, Vector3.Up);
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
        this.animations.Travel("punch");
        this.target?.OnArrived();

        if (targetPlayerId.HasValue)
        {
            this.mb?.PlayerRequestedDamage(this.targetPlayerId.GetValueOrDefault());
        }
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
        if (this.hp <= amount)
        {
            this.alive = false;
            this.hp = 0;
            if (this.healthBar != null)
            {
                this.healthBar.Visible = false;
            }

            this.animations.Travel("death");
            // Play death animation. Restrict movement. Respawn soon?
        }
        this.hp -= amount;
        if (this.healthBar != null)
        {
            GD.Print(((CapsuleMesh)this.healthBar.Mesh).MidHeight);
            ((CapsuleMesh)this.healthBar.Mesh).MidHeight = 1.0f * this.hp / MAX_HP;
            this.healthBar.Translation -= new Vector3(0.5f * amount / MAX_HP, 0, 0);
            GD.Print(((CapsuleMesh)this.healthBar.Mesh).MidHeight);
        }
    }
}
