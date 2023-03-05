using Godot;
using System.IO;
using System.Text;

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
            GD.Print("Error connect to " + this.wsUrl);
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

    #region SERVER_TO_CLIENT
    private void onDataReceived()
    {
        var packet = this.client.GetPacket();
        GD.Print(Encoding.UTF8.GetString(packet));
        // var action = JsonSerializer.Deserialize<Proto.ServerMessage>(packet);


        var action = Proto.ServerMessage.Parser.ParseFrom(packet);
        GD.Print(action);
        switch (action.Type)
        {
            case Proto.ServerMessageType.MessagePing:
                GD.Print("ping"); // TODO
                break;
            case Proto.ServerMessageType.MessageJoin:
                // var joinGameRes = JsonSerializer.Deserialize<Proto.JoinGameResponse>(action.Payload);
                // var localPlayer = this.players.OnLocalPlayerJoined(joinGameRes);
                // // TODO - spagoot
                // var cameraFollowing = this.GetParent().GetNode<RemoteTransform3D>("CameraBase/Following");
                // cameraFollowing.RemotePath = localPlayer.GetPath();
                break;
            case Proto.ServerMessageType.MessageTick:
                // var tick = JsonSerializer.Deserialize<Proto.GameTick>(action.Payload);
                // this.processGameTick(tick);
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
                    // var connect = JsonSerializer.Deserialize<Proto.Connect>(action.Value);
                    // if (connect.PlayerId == this.players.LocalPlayerId) return;
                    // this.players.OnPlayerConnected(connect.PlayerId, connect.Spawn, connect.Color);
                    break;
                case Proto.ClientActionType.ActionDisconnect:
                    // var disconnect = JsonSerializer.Deserialize<Proto.Disconnect>(action.Value);
                    // this.players.OnPlayerDisconnected(disconnect.PlayerId);
                    break;
                case Proto.ClientActionType.ActionMove:
                    // var move = JsonSerializer.Deserialize<Proto.Move>(action.Value);
                    // this.players.SetMoving(move.PlayerId, move.Target);
                    // this.players.StopAttacking(move.PlayerId);
                    break;
                case Proto.ClientActionType.ActionAttack:
                    // var attack = JsonSerializer.Deserialize<Proto.Attack>(action.Value);
                    // this.players.SetAttacking(
                    //     attack.SourcePlayerId,
                    //     attack.TargetPlayerId,
                    //     attack.TargetPlayerLocation);
                    break;
                case Proto.ClientActionType.ActionDamage:
                    // var damage = JsonSerializer.Deserialize<Proto.Damage>(action.Value);
                    // GD.Print("damage received");
                    // this.players.PlayAttackingAnimation(damage.SourcePlayerId);
                    // this.players.ApplyDamage(damage.TargetPlayerId, damage.DamageDealt);
                    break;
                case Proto.ClientActionType.ActionDeath:
                    // var death = JsonSerializer.Deserialize<Proto.Death>(action.Value);
                    // this.players.Die(death.PlayerId);
                    break;
                case Proto.ClientActionType.ActionRespawn:
                    // var respawn = JsonSerializer.Deserialize<Proto.Respawn>(action.Value);
                    // this.players.Spawn(respawn.PlayerId, respawn.Spawn);
                    break;
                default:
                    GD.Print($"Got unexpected client action from server: {action.Type}");
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

        Error error = this.client.PutPacket(actionStreamMem.ToArray());
        if (error != Error.Ok)
        {
            GD.Print($"Failed to write packet: {error}");
            actionStream.Dispose();
            return;
        }

        actionStream.Dispose();
    }

    public void PlayerRequestedMove(Proto.Location target)
    {
        var moveStreamMem = new MemoryStream();
        var moveStream = new Google.Protobuf.CodedOutputStream(new MemoryStream(), false);

        new Proto.Move()
        {
            PlayerId = this.players.LocalPlayerId,
            Target = target,
        }.WriteTo(moveStream);

        this.writeClientActionToServer(
            Proto.ClientActionType.ActionMove,
            moveStreamMem.ToArray());

        moveStream.Dispose();
    }

    public void PlayerRequestedAttack(ulong targetPlayerId)
    {
        // var msg = new Proto.ClientAction()
        // {
        //     Type = Proto.ClientActionType.ActionAttack,
        //     Payload = JsonSerializer.Serialize(new Proto.Attack()
        //     {
        //         SourcePlayerId = this.players.LocalPlayerId,
        //         SourcePlayerLocation = this.players.CurrentLocation(this.players.LocalPlayerId),
        //         TargetPlayerId = targetPlayerId,
        //         TargetPlayerLocation = this.players.CurrentLocation(targetPlayerId),
        //     })
        // };
        // var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        // Error error = this.client.PutPacket(msgBytes);
        // if (error != Error.Ok)
        // {
        //     GD.Print($"Failed to request move: {error}");
        //     return;
        // }
    }

    public void PlayerRequestedDamage(ulong targetPlayerId)
    {
        // var msg = new Proto.ClientAction()
        // {
        //     Type = Proto.ClientActionType.ActionDamage,
        //     Payload = JsonSerializer.Serialize(new Proto.Damage()
        //     {
        //         SourcePlayerId = this.players.LocalPlayerId,
        //         TargetPlayerId = targetPlayerId,
        //         DamageDealt = 1,
        //     })
        // };
        // var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        // Error error = this.client.PutPacket(msgBytes);
        // if (error != Error.Ok)
        // {
        //     GD.Print($"Failed to request move: {error}");
        //     return;
        // }
    }
    #endregion
}
