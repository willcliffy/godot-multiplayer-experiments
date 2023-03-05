using Godot;
using System.Text;
using System.Text.Json;
using Game;

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
        var action = JsonSerializer.Deserialize<ServerMessage>(packet);
        switch ((ServerMessageType)action.type)
        {
            case ServerMessageType.MESSAGE_PING:
                GD.Print("ping"); // TODO
                break;
            case ServerMessageType.MESSAGE_JOIN:
                var joinGameRes = JsonSerializer.Deserialize<JoinGameResponse>(action.payload);
                var localPlayer = this.players.OnLocalPlayerJoined(joinGameRes);
                // TODO - spagoot
                var cameraFollowing = this.GetParent().GetNode<RemoteTransform3D>("CameraBase/Following");
                cameraFollowing.RemotePath = localPlayer.GetPath();
                break;
            case ServerMessageType.MESSAGE_TICK:
                var tick = JsonSerializer.Deserialize<GameTick>(action.payload);
                this.processGameTick(tick);
                break;
            default:
                GD.Print($"Unknown server message type: '{action.type}'");
                break;
        }
    }

    private void processGameTick(GameTick tick)
    {
        // TODO - 🍝
        foreach (var action in tick.actions)
        {
            switch (action.type)
            {
                case (uint)ClientActionType.ACTION_CONNECT:
                    var connect = JsonSerializer.Deserialize<Connect>(action.value);
                    if (connect.playerId == this.players.LocalPlayerId) return;
                    this.players.OnPlayerConnected(connect.playerId, connect.spawn, connect.color);
                    break;
                case (uint)ClientActionType.ACTION_DISCONNECT:
                    var disconnect = JsonSerializer.Deserialize<Disconnect>(action.value);
                    this.players.OnPlayerDisconnected(disconnect.playerId);
                    break;
                case (uint)ClientActionType.ACTION_MOVE:
                    var move = JsonSerializer.Deserialize<Move>(action.value);
                    this.players.SetMoving(move.playerId, move.target.ToVector3());
                    this.players.StopAttacking(move.playerId);
                    break;
                case (uint)ClientActionType.ACTION_ATTACK:
                    var attack = JsonSerializer.Deserialize<Attack>(action.value);
                    this.players.SetAttacking(
                        attack.sourcePlayerId,
                        attack.targetPlayerId,
                        attack.targetPlayerLocation.ToVector3());
                    break;
                case (uint)ClientActionType.ACTION_DAMAGE:
                    var damage = JsonSerializer.Deserialize<Damage>(action.value);
                    GD.Print("damage received");
                    this.players.PlayAttackingAnimation(damage.sourcePlayerId);
                    this.players.ApplyDamage(damage.targetPlayerId, damage.damageDealt);
                    break;
                case (uint)ClientActionType.ACTION_DEATH:
                    var death = JsonSerializer.Deserialize<Death>(action.value);
                    this.players.Die(death.playerId);
                    break;
                case (uint)ClientActionType.ACTION_RESPAWN:
                    var respawn = JsonSerializer.Deserialize<Respawn>(action.value);
                    this.players.Spawn(respawn.playerId, respawn.spawn);
                    break;
                default:
                    break;
            }
        }
    }
    #endregion

    #region CLIENT_TO_SERVER
    public void PlayerRequestedMove(Location target)
    {
        var msg = new ClientAction()
        {
            type = (int)ClientActionType.ACTION_MOVE,
            payload = JsonSerializer.Serialize(new Move()
            {
                playerId = this.players.LocalPlayerId,
                target = target,
            })
        };
        var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        Error error = this.client.PutPacket(msgBytes);
        if (error != Error.Ok)
        {
            GD.Print($"Failed to request move: {error}");
            return;
        }
    }

    public void PlayerRequestedAttack(ulong targetPlayerId)
    {
        var msg = new ClientAction()
        {
            type = (int)ClientActionType.ACTION_ATTACK,
            payload = JsonSerializer.Serialize(new Attack()
            {
                sourcePlayerId = this.players.LocalPlayerId,
                sourcePlayerLocation = this.players.CurrentLocation(this.players.LocalPlayerId),
                targetPlayerId = targetPlayerId,
                targetPlayerLocation = this.players.CurrentLocation(targetPlayerId),
            })
        };
        var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        Error error = this.client.PutPacket(msgBytes);
        if (error != Error.Ok)
        {
            GD.Print($"Failed to request move: {error}");
            return;
        }
    }

    public void PlayerRequestedDamage(ulong targetPlayerId)
    {
        var msg = new ClientAction()
        {
            type = (int)ClientActionType.ACTION_DAMAGE,
            payload = JsonSerializer.Serialize(new Damage()
            {
                sourcePlayerId = this.players.LocalPlayerId,
                targetPlayerId = targetPlayerId,
                damageDealt = 1,
            })
        };
        var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        Error error = this.client.PutPacket(msgBytes);
        if (error != Error.Ok)
        {
            GD.Print($"Failed to request move: {error}");
            return;
        }
    }
    #endregion
}
