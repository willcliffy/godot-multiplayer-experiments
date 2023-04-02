using Godot;
using System.Collections.Generic;
using System.Linq;

public partial class Environment : Node3D
{
    const float WAYPOINT_REACHED_MARGIN = 1.0f;
    const float TARGET_REACHED_MARGIN = 0.1f;

    private Vector3 targetPosition;

    private Vector3I[] currentPath;
    private Vector3I currentWaypoint;
    private int currentWaypointIdx;

    private AStar3D navAlgo = new AStar3D();
    private Dictionary<string, long> points = new Dictionary<string, long>();
    private GridMap navGrid;

    public override void _Ready()
    {
        this.navGrid = this.GetNode<GridMap>("NavGrid");
        var cells = this.navGrid.GetUsedCells();
        foreach (var cell in cells)
        {
            var idx = this.navAlgo.GetAvailablePointId();
            this.navAlgo.AddPoint(idx, this.navGrid.MapToLocal(cell));
            this.points[Vector3ToString(cell)] = idx;

            for (int x = -1; x <= 1; x++)
                for (int z = -1; z <= 1; z++)
                {
                    var v3 = new Vector3I(x, 0, z);
                    if (v3 == Vector3I.Zero) continue;
                    if (this.points.ContainsKey(Vector3ToString(cell + v3)))
                    {
                        var idx1 = this.points[Vector3ToString(cell)];
                        var idx2 = this.points[Vector3ToString(cell + v3)];
                        if (!this.navAlgo.ArePointsConnected(idx1, idx2))
                            this.navAlgo.ConnectPoints(idx1, idx2);
                    }
                }
        }
    }

    private static string Vector3ToString(Vector3 v3)
    {
        return $"{v3.X}.{v3.Z}";
    }

    public Vector3I[] SetTargetPosition(Vector3 agentPosition, Vector3I target)
    {
        this.targetPosition = target;
        this.currentPath = this.CalculatePath((Vector3I)agentPosition.Round(), target);
        this.currentWaypointIdx = 0;
        this.currentWaypoint = this.currentPath[this.currentWaypointIdx];
        GD.Print($"Starting at {agentPosition} new waypoint {this.currentWaypoint}");
        foreach (var item in currentPath)
        {
            GD.Print($"\t{item}");
        }
        return this.currentPath;
    }

    public Vector3I[] CalculatePath(Vector3I start, Vector3I end)
    {
        GD.Print($"start: {start}");
        var gmStart = Vector3ToString(start);
        long startId = 0;
        if (this.points.ContainsKey(gmStart))
        {
            startId = points[gmStart];
        }
        else
        {
            GD.Print($"No cached point found for {gmStart}");
            startId = this.navAlgo.GetClosestPoint(start);
        }

        var gmEnd = Vector3ToString(this.navGrid.MapToLocal(end));
        GD.Print($"end: {end}");
        long endId = 0;
        if (this.points.ContainsKey(gmEnd))
        {
            endId = points[gmEnd];
        }
        else
        {
            endId = this.navAlgo.GetClosestPoint(end);
        }

        return this.navAlgo.GetPointPath(startId, endId)
            .Select(point => (Vector3I)point.Round())
            .Append(end)
            .Distinct()
            .Skip(1)
            .ToArray();
    }

    public bool IsTargetReached(Vector3 agentPosition)
    {
        var agentPosition2d = new Vector3(agentPosition.X, 0, agentPosition.Z);
        return agentPosition2d.DistanceTo(this.targetPosition) <= TARGET_REACHED_MARGIN;
    }

    public Vector3I GetNextPathPosition(Vector3 agentPosition)
    {
        var agentPosition2d = new Vector3(agentPosition.X, 0, agentPosition.Z);
        if (agentPosition2d.DistanceTo(this.currentWaypoint) <= WAYPOINT_REACHED_MARGIN
            && this.currentWaypointIdx < this.currentPath.Length)
        {
            this.currentWaypoint = this.currentPath[this.currentWaypointIdx];
            this.currentWaypointIdx++;
        }

        return this.currentWaypoint;
    }
}
