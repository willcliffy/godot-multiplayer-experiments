using Godot;
using System;
using System.IO;

public partial class MessageBroker : Node
{
    [Export]
    string wsUrl = "ws://localhost:8080/ws/v1/connect";

    private WebSocketPeer client = new WebSocketPeer();
    private WebSocketPeer.State websocketState = WebSocketPeer.State.Closed;
    private PlayerController players;

    public override void _Ready()
    {
        this.players = GetNode<PlayerController>("PlayerController");

        if (OS.HasEnvironment("GAMESERVER_WEBSOCKET_URL"))
        {
            this.wsUrl = OS.GetEnvironment("GAMESERVER_WEBSOCKET_URL");
        }
        else
        {
            GD.Print($"falling back to default WS URL: {wsUrl}");
        }

        Error error = this.client.ConnectToUrl(this.wsUrl);
        if (error != Error.Ok)
        {
            this.client.Close();
            GD.Print("Error connecting to " + this.wsUrl);
            return;
        }

        GD.Print("Starting socket connetion to " + this.wsUrl);
    }

    public override void _Process(double delta)
    {
        this.client.Poll();
        this.websocketState = this.client.GetReadyState();
        if (this.websocketState != WebSocketPeer.State.Open) return;
        if (this.client.GetAvailablePacketCount() > 0) this.onDataReceived();
    }

    #region UTILS
    public static Proto.Location Vector3ToLocation(Vector3 vector3)
    {
        return new Proto.Location()
        {
            X = (int)vector3.X,
            Z = (int)vector3.Z,
        };
    }

    public static Vector3 LocationToVector3d(Proto.Location location)
    {
        return new Vector3(location.X, 0, location.Z);
    }
    #endregion

    #region SERVER_TO_CLIENT
    private void onDataReceived()
    {
        var action = Proto.ServerMessage.Parser.ParseFrom(this.client.GetPacket());
        switch (action.Type)
        {
            case Proto.ServerMessageType.MessagePing:
                GD.Print("ping"); // TODO
                break;
            case Proto.ServerMessageType.MessageJoin:
                var joinGameRes = Proto.JoinGameResponse.Parser.ParseFrom(action.Payload);
                var localPlayer = this.players.OnLocalPlayerJoined(joinGameRes);
                break;
            case Proto.ServerMessageType.MessageTick:
                var tick = Proto.GameTick.Parser.ParseFrom(action.Payload);
                this.processGameTick(tick);
                break;
            default:
                GD.Print($"Unknown server message type: '{action.Type}'");
                break;
        }
    }

    private void processGameTick(Proto.GameTick tick)
    {
        // TODO - 🍝
        foreach (var action in tick.Actions)
        {
            switch (action.Type)
            {
                case Proto.ClientActionType.ActionConnect:
                    if (action.PlayerId == this.players.LocalPlayerId) return;
                    var connect = Proto.Connect.Parser.ParseFrom(action.Payload);
                    this.players.OnPlayerConnected(action.PlayerId, connect.Spawn, connect.Color);
                    break;
                case Proto.ClientActionType.ActionDisconnect:
                    this.players.OnPlayerDisconnected(action.PlayerId);
                    break;
                case Proto.ClientActionType.ActionMove:
                    var move = Proto.Move.Parser.ParseFrom(action.Payload);
                    var player = players.LocalPlayer.Position;
                    GD.Print($"{player.DistanceTo(LocationToVector3d(move.Path[0]))}");
                    break;
                case Proto.ClientActionType.ActionAttack:
                    var attack = Proto.Attack.Parser.ParseFrom(action.Payload);
                    this.players.SetAttacking(
                        action.PlayerId,
                        attack.TargetPlayerId,
                        attack.TargetPlayerLocation);
                    break;
                case Proto.ClientActionType.ActionDamage:
                    var damage = Proto.Damage.Parser.ParseFrom(action.Payload);
                    this.players.ApplyDamage(damage.TargetPlayerId, damage.DamageDealt);
                    break;
                case Proto.ClientActionType.ActionDeath:
                    this.players.Die(action.PlayerId);
                    break;
                case Proto.ClientActionType.ActionRespawn:
                    var respawn = Proto.Respawn.Parser.ParseFrom(action.Payload);
                    this.players.Spawn(action.PlayerId, respawn.Spawn);
                    break;
                case Proto.ClientActionType.ActionCollect:
                    GD.Print("collect nyi");
                    break;
                case Proto.ClientActionType.ActionBuild:
                    GD.Print("build nyi");
                    break;
                default:
                    GD.Print($"Got unexpected client action: {action.Type}");
                    break;
            }
        }
    }
    #endregion

    #region CLIENT_TO_SERVER
    private void writeClientActionToServer(Proto.ClientActionType type, byte[] bytes)
    {
        var msg = new Proto.ClientAction()
        {
            Type = type,
            PlayerId = this.players.LocalPlayerId,
            Payload = Google.Protobuf.ByteString.CopyFrom(bytes),
        };

        var actionStreamMem = new MemoryStream();
        var actionStream = new Google.Protobuf.CodedOutputStream(actionStreamMem, false);
        msg.WriteTo(actionStream);
        actionStream.Flush();

        Error error = this.client.Send(actionStreamMem.ToArray());
        if (error != Error.Ok)
        {
            GD.Print($"Failed to write packet: {error}");
        }

        actionStream.Dispose();
    }

    public void PlayerRequestedMove(Vector3 target, Proto.ClientAction queued = null)
    {
        var path = this.players.SetMoving(
            this.players.LocalPlayerId,
            Vector3ToLocation(target));

        var moveAction = new Proto.Move()
        {
            Queued = queued,
        };
        foreach (var vector3 in path)
        {
            moveAction.Path.Add(Vector3ToLocation(vector3));
        }

        var moveStreamMem = new MemoryStream();
        var moveStream = new Google.Protobuf.CodedOutputStream(moveStreamMem, false);
        moveAction.WriteTo(moveStream);

        moveStream.Flush();

        this.writeClientActionToServer(
            Proto.ClientActionType.ActionMove,
            moveStreamMem.ToArray());

        moveStream.Dispose();
    }

    public void PlayerRequestedAttack(ulong targetPlayerId)
    {
        var attackAction = new Proto.Attack()
        {
            SourcePlayerLocation = this.players.CurrentLocation(this.players.LocalPlayerId),
            TargetPlayerId = targetPlayerId,
            TargetPlayerLocation = this.players.CurrentLocation(targetPlayerId),
        };

        var attackStreamMem = new MemoryStream();
        var attackStream = new Google.Protobuf.CodedOutputStream(attackStreamMem, false);
        attackAction.WriteTo(attackStream);
        attackStream.Flush();

        this.PlayerRequestedMove(
            LocationToVector3d(this.players.CurrentLocation(targetPlayerId)),
            new Proto.ClientAction()
            {
                Type = Proto.ClientActionType.ActionAttack,
                PlayerId = this.players.LocalPlayerId,
                Payload = Google.Protobuf.ByteString.CopyFrom(attackStreamMem.ToArray()),
            });

        attackStream.Dispose();
    }

    public void PlayerRequestedCollect(Vector3 target, Proto.ResourceType type)
    {
        var collectAction = new Proto.Collect()
        {
            Location = Vector3ToLocation(target),
            Type = type,
        };

        var collectStreamMem = new MemoryStream();
        var collectStream = new Google.Protobuf.CodedOutputStream(collectStreamMem, false);
        collectAction.WriteTo(collectStream);
        collectStream.Flush();

        this.PlayerRequestedMove(
            target,
            new Proto.ClientAction()
            {
                Type = Proto.ClientActionType.ActionCollect,
                PlayerId = this.players.LocalPlayerId,
                Payload = Google.Protobuf.ByteString.CopyFrom(collectStreamMem.ToArray()),
            });

        collectStream.Dispose();
    }

    public void PlayerRequestedBuild(Vector3 target)
    {
        var buildAction = new Proto.Build()
        {
            Location = Vector3ToLocation(target),
        };

        var buildStreamMem = new MemoryStream();
        var buildStream = new Google.Protobuf.CodedOutputStream(buildStreamMem, false);
        buildAction.WriteTo(buildStream);
        buildStream.Flush();

        this.writeClientActionToServer(
            Proto.ClientActionType.ActionMove,
            buildStreamMem.ToArray());

        buildStream.Dispose();
    }

    public void PlayerRequestedDamage(ulong targetPlayerId)
    {
        var damageAction = new Proto.Damage()
        {
            TargetPlayerId = targetPlayerId,
            DamageDealt = 1,
        };

        var dmgStreamMem = new MemoryStream();
        var dmgStream = new Google.Protobuf.CodedOutputStream(dmgStreamMem, false);
        damageAction.WriteTo(dmgStream);
        dmgStream.Flush();

        this.writeClientActionToServer(
            Proto.ClientActionType.ActionDamage,
            dmgStreamMem.ToArray());

        dmgStream.Dispose();
    }
    #endregion
}
