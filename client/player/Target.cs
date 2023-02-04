using Godot;
using System;

public class Target : CSGMesh
{
    public override void _Process(float delta)
    {
        if (this.Visible)
        {
            RotateY(delta);
        }

    }

    public void SetLocation(Vector3 location)
    {
        this.Translation = location + new Vector3(0, 1 / 3f, 0);
        this.Visible = true;
    }

    public void OnArrived()
    {
        this.Visible = false;
    }
}
