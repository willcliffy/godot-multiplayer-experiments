using Godot;
using System;

public class Target : MeshInstance
{
    public override void _Process(float delta)
    {
        if (this.Visible)
        {
            RotateY(delta);
        }

    }

    public void SetLocation(Vector3 location, bool attacking = false)
    {
        if (attacking)
        {
            Translation = location + new Vector3(0, 2 / 3f, 0);
            var mat = (SpatialMaterial)Mesh.SurfaceGetMaterial(0).Duplicate();
            mat.AlbedoColor = new Color(1, 0, 0);
            SetSurfaceMaterial(0, mat);
        }
        else
        {
            Translation = location + new Vector3(0, 1 / 3f, 0);
            var mat = (SpatialMaterial)Mesh.SurfaceGetMaterial(0).Duplicate();
            mat.AlbedoColor = new Color(0, 0, 1);
            SetSurfaceMaterial(0, mat);
        }
        this.Visible = true;
    }

    public void OnArrived()
    {
        this.Visible = false;
    }
}
