using Godot;

public class Player : KinematicBody
{
    const int SPEED = 3;
    const float ACCEPTABLE_DIST_TO_TARGET_RANGE = 0.05f;
    const float ATTACK_RANGE = 2 + ACCEPTABLE_DIST_TO_TARGET_RANGE;

    public ulong id { get; set; }


    private bool moving;
    private bool attacking;
    private Vector3 targetLocation;

    private Spatial model;
    private NavigationAgent nav;
    private AnimationNodeStateMachinePlayback animations;
    private Target target;

    // Called when the node enters the scene tree for the first time.
    public override void _Ready()
    {
        this.model = GetNode<Spatial>("Robot");
        this.nav = GetNode<NavigationAgent>("NavigationAgent");
        var animationNode = GetNode<AnimationTree>("AnimationTree");
        this.animations = (AnimationNodeStateMachinePlayback)animationNode.Get("parameters/playback");

        if (Visible) // TODO - hacky way to check if you are the local player
        {
            this.target = GetNode<Target>("../Target");
        }
    }

    public override void _PhysicsProcess(float delta)
    {
        if (Input.IsActionJustPressed("exit")) GetTree().Quit();
        if (!this.moving) return;

        var distToTarget = (this.targetLocation - Translation).Length();
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
        }

        var direction = (next - Translation).Normalized();
        var _collision = MoveAndCollide(direction * delta * SPEED);
        var facingDirection = Translation - direction;
        facingDirection.y = Translation.y;
        this.model.LookAt(facingDirection, Vector3.Up);
    }

    public void SetTeam(string color)
    {
        // HACKY
        var head = GetNode<MeshInstance>("Robot/RobotArmature/Skeleton/BoneAttachment2/Head");
        var material = (SpatialMaterial)head.Mesh.SurfaceGetMaterial(0).Duplicate();
        material.AlbedoColor = new Color(color);
        head.SetSurfaceMaterial(0, material);
    }

    public void SetMoving(Vector3 target)
    {
        this.moving = true;
        this.attacking = false;
        this.targetLocation = target;
        this.nav.SetTargetLocation(target);
        this.animations.Travel("walk");

        if (this.target != null)
        {
            this.target.SetLocation(target);
        }
    }

    public void SetAttacking(Vector3 target)
    {
        this.moving = true;
        this.attacking = true;
        this.target.SetLocation(target);
        this.nav.SetTargetLocation(target);
        this.animations.Travel("walk");

        if (this.target != null)
        {
            this.target.SetLocation(target);
        }
    }

    private void setAttackingTargetReached()
    {
        this.moving = false;
        this.attacking = true;
        this.animations.Travel("attack");

        if (this.target != null)
        {
            this.target.OnArrived();
        }
    }

    private void setIdle()
    {
        this.moving = false;
        this.attacking = false;
        this.animations.Travel("idle");

        if (this.target != null)
        {
            this.target.OnArrived();
        }
    }
}