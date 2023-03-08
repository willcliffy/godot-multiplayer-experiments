using Godot;
using System;

public partial class Resource : StaticBody3D
{
    public Proto.ResourceType Type;

    public bool IsDepleted { get { return remaining == 0; } }
    private int remaining = 5;
    private MeshInstance3D mesh;

    public override void _Process(double delta)
    {
        this.mesh = this.GetNode<MeshInstance3D>("RedResource");
    }

    public void Collect()
    {
        this.remaining -= 1;
        if (this.remaining == 0)
        {
            var mat = (StandardMaterial3D)this.mesh.Mesh.SurfaceGetMaterial(0).Duplicate();
            mat.AlbedoColor = new Color(0.5f, 0.1f, 0.1f, 0.9f);
            this.mesh.SetSurfaceOverrideMaterial(0, mat);
        }
    }
}
