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
                // TODO - spagoot
                var cameraFollowing = this.GetParent().GetNode<RemoteTransform3D>("CameraBase/Following");
                cameraFollowing.RemotePath = localPlayer.GetPath();
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
                    this.players.SetMoving(move.PlayerId, move.Target);
                    this.players.StopAttacking(move.PlayerId);
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
                    this.players.PlayAttackingAnimation(damage.SourcePlayerId);
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
        actionStreamMem.Flush();

        Error error = this.client.Send(actionStreamMem.ToArray());
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
        var moveStream = new Google.Protobuf.CodedOutputStream(moveStreamMem, false);

        new Proto.Move()
        {
            PlayerId = this.players.LocalPlayerId,
            Target = target,
        }.WriteTo(moveStream);

        moveStream.Flush();
        moveStreamMem.Flush();

        this.writeClientActionToServer(
            Proto.ClientActionType.ActionMove,
            moveStreamMem.ToArray());

        moveStream.Dispose();
    }

    public void PlayerRequestedAttack(ulong targetPlayerId)
    {
        var attackStreamMem = new MemoryStream();
        var attackStream = new Google.Protobuf.CodedOutputStream(attackStreamMem, false);

        new Proto.Attack()
        {
            SourcePlayerId = this.players.LocalPlayerId,
            SourcePlayerLocation = this.players.CurrentLocation(this.players.LocalPlayerId),
            TargetPlayerId = targetPlayerId,
            TargetPlayerLocation = this.players.CurrentLocation(targetPlayerId),
        }.WriteTo(attackStream);

        attackStream.Flush();
        attackStreamMem.Flush();

        this.writeClientActionToServer(
            Proto.ClientActionType.ActionAttack,
            attackStreamMem.ToArray());

        attackStream.Dispose();
    }

    public void PlayerRequestedDamage(ulong targetPlayerId)
    {
        var dmgStreamMem = new MemoryStream();
        var dmgStream = new Google.Protobuf.CodedOutputStream(dmgStreamMem, false);

        new Proto.Damage()
        {
            SourcePlayerId = this.players.LocalPlayerId,
            TargetPlayerId = targetPlayerId,
            DamageDealt = 1,
        }.WriteTo(dmgStream);

        dmgStream.Flush();
        dmgStreamMem.Flush();

        this.writeClientActionToServer(
            Proto.ClientActionType.ActionDamage,
            dmgStreamMem.ToArray());

        dmgStream.Dispose();
    }
    #endregion
}
