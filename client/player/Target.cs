using Godot;

public partial class Target : MeshInstance3D
{
    public override void _Process(double delta)
    {
        if (this.Visible)
        {
            this.RotateY((float)delta);
            this.Transparency += (float)delta;
            if (this.Transparency >= 1.0) this.Visible = false;
        }
    }

    public void SetLocation(Vector3 location, bool attacking = false)
    {
        if (attacking)
        {
            this.Position = location + new Vector3(0, 5 / 3f, 0);
            var mat = (StandardMaterial3D)Mesh.SurfaceGetMaterial(0).Duplicate();
            mat.AlbedoColor = new Color(1, 0, 0);
            this.SetSurfaceOverrideMaterial(0, mat);
        }
        else
        {
            this.Position = location + new Vector3(0, 1 / 3f, 0);
            var mat = (StandardMaterial3D)Mesh.SurfaceGetMaterial(0).Duplicate();
            mat.AlbedoColor = new Color(0, 0, 1);
            this.Transparency = 0.0f;
            this.SetSurfaceOverrideMaterial(0, mat);
        }
        this.Visible = true;
    }

    public void OnArrived()
    {
        this.Visible = false;
    }
}
