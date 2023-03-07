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
    private Proto.Location vector3ToLocation(Vector3 vector3)
    {
        return new Proto.Location()
        {
            X = (int)vector3.X,
            Z = (int)vector3.Z,
        };
    }

    private Vector3 locationToVector3d(Proto.Location location)
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
                    var connect = Proto.Connect.Parser.ParseFrom(action.Value);
                    if (connect.PlayerId == this.players.LocalPlayerId) return;
                    this.players.OnPlayerConnected(connect.PlayerId, connect.Spawn, connect.Color);
                    break;
                case Proto.ClientActionType.ActionDisconnect:
                    var disconnect = Proto.Disconnect.Parser.ParseFrom(action.Value);
                    this.players.OnPlayerDisconnected(disconnect.PlayerId);
                    break;
                case Proto.ClientActionType.ActionMove:
                    var move = Proto.Move.Parser.ParseFrom(action.Value);
                    GD.Print($"{DateTime.Now.Second}.{DateTime.Now.Millisecond} got move");
                    // this.players.SetMoving(move.PlayerId, move.Path[0]);
                    // this.players.StopAttacking(move.PlayerId);
                    break;
                case Proto.ClientActionType.ActionAttack:
                    var attack = Proto.Attack.Parser.ParseFrom(action.Value);
                    this.players.SetAttacking(
                        attack.SourcePlayerId,
                        attack.TargetPlayerId,
                        attack.TargetPlayerLocation);
                    break;
                case Proto.ClientActionType.ActionDamage:
                    var damage = Proto.Damage.Parser.ParseFrom(action.Value);
                    this.players.ApplyDamage(damage.TargetPlayerId, damage.DamageDealt);
                    break;
                case Proto.ClientActionType.ActionDeath:
                    var death = Proto.Death.Parser.ParseFrom(action.Value);
                    this.players.Die(death.PlayerId);
                    break;
                case Proto.ClientActionType.ActionRespawn:
                    var respawn = Proto.Respawn.Parser.ParseFrom(action.Value);
                    this.players.Spawn(respawn.PlayerId, respawn.Spawn);
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

    public void PlayerRequestedMove(Vector3 target)
    {
        var path = this.players.SetMoving(
            this.players.LocalPlayerId,
            vector3ToLocation(target));

        var moveAction = new Proto.Move()
        {
            PlayerId = this.players.LocalPlayerId,
        };
        foreach (var vector3 in path)
        {
            moveAction.Path.Add(vector3ToLocation(vector3));
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
            SourcePlayerId = this.players.LocalPlayerId,
            SourcePlayerLocation = this.players.CurrentLocation(this.players.LocalPlayerId),
            TargetPlayerId = targetPlayerId,
            TargetPlayerLocation = this.players.CurrentLocation(targetPlayerId),
        };

        var attackStreamMem = new MemoryStream();
        var attackStream = new Google.Protobuf.CodedOutputStream(attackStreamMem, false);
        attackAction.WriteTo(attackStream);
        attackStream.Flush();

        this.writeClientActionToServer(
            Proto.ClientActionType.ActionAttack,
            attackStreamMem.ToArray());

        attackStream.Dispose();
    }

    public void PlayerRequestedDamage(ulong targetPlayerId)
    {
        var damageAction = new Proto.Damage()
        {
            SourcePlayerId = this.players.LocalPlayerId,
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
