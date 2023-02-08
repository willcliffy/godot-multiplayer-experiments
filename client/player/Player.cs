using Game;
using Godot;

public class Player : KinematicBody
{
    const int SPEED = 3;
    const float ACCEPTABLE_DIST_TO_TARGET_RANGE = 0.05f;
    const float ATTACK_RANGE = 1 + ACCEPTABLE_DIST_TO_TARGET_RANGE;

    public ulong id { get; set; }

    private Spatial model;
    private NavigationAgent nav;
    private AnimationNodeStateMachinePlayback animations;

    private bool moving;
    private bool attacking;

    private Target target;
    private Vector3 targetLocation;
    private ulong? targetPlayerId;


    public override void _Ready()
    {
        this.model = GetNode<Spatial>("Robot");
        this.nav = GetNode<NavigationAgent>("NavigationAgent");

        var animationNode = GetNode<AnimationTree>("AnimationTree");
        this.animations = (AnimationNodeStateMachinePlayback)animationNode.Get("parameters/playback");

        // TODO - hacky way to check if this is the local player
        if (Visible) this.target = GetNode<Target>("../Target");
    }

    public override void _PhysicsProcess(float delta)
    {
        if (Input.IsActionJustPressed("exit")) GetTree().Quit();
        if (!this.moving) return;

        var distToTarget = (this.targetLocation - this.Translation).Length();
        if (attacking && distToTarget < ATTACK_RANGE)
        {
            GD.Print(this.targetLocation);
            GD.Print(this.Translation);
            GD.Print($"reached target, attacking with dist {distToTarget}");
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
        this.moving = true;
        this.attacking = true;
        this.targetPlayerId = targetPlayerId;
        this.targetLocation = targetLocation;
        this.nav.SetTargetLocation(targetLocation);
        this.animations.Travel("walk");
        this.target?.SetLocation(targetLocation, attacking = true);
    }

    private void setAttackingTargetReached()
    {
        this.moving = false;
        this.attacking = true;
        this.animations.Travel("punch");
        //this.target?.OnArrived();
    }

    private void setIdle()
    {
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

    public bool IsAttacking(ulong? playerId)
    {
        return this.attacking && playerId != null && playerId == this.targetPlayerId;
    }
}
