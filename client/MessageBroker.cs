using Godot;
using System.Text;
using System.Text.Json;
using Game;

public class MessageBroker : Node
{
    [Export]
    string webSocketURL = System.Environment.GetEnvironmentVariable("GAMESERVER_WEBSOCKET_URL") ?? "ws://kilnwood-game.com/ws/v1/connect";

    WebSocketClient client = null;
    Player player;
    OpponentController opponents;

    #region GODOT
    public override void _Ready()
    {
        player = GetNode<Player>("Player");
        opponents = GetNode<OpponentController>("OpponentController");

        client = new WebSocketClient();

        client.Connect("connection_established", this, nameof(onConnectionEstablished));
        client.Connect("data_received", this, nameof(onDataReceived));
        client.Connect("server_close_request", this, nameof(onServerCloseRequest));
        client.Connect("connection_closed", this, nameof(onConnectionClosed));

        Error error = client.ConnectToUrl(webSocketURL);
        if (error != Error.Ok)
        {
            client.GetPeer(1).Close();
            GD.Print("Error connect to " + webSocketURL);
            return;
        }

        GD.Print("Starting socket connetion to " + webSocketURL);
    }

    public override void _Process(float delta)
    {
        if (client.GetConnectionStatus() == NetworkedMultiplayerPeer.ConnectionStatus.Connected ||
            client.GetConnectionStatus() == NetworkedMultiplayerPeer.ConnectionStatus.Connecting)
        {
            client.Poll();
        }
    }
    #endregion

    #region WEBSOCKET
    private void onConnectionEstablished(string protocol)
    {
        var msg = new ClientAction();
        msg.type = (int)ClientActionType.ACTION_CONNECT;
        var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        Error error = client.GetPeer(1).PutPacket(msgBytes);
        if (error != Error.Ok)
        {
            GD.Print($"Failed to establish connection: {error}");
            return;
        }
        GD.Print("Connection established");
    }

    private void onServerCloseRequest(int code, string reason)
    {
        GD.Print("Close request, reason: " + reason);
    }

    private void onConnectionClosed(bool wasCleanClose)
    {
        GD.Print("Connection closed. was clean close: " + wasCleanClose.ToString());
    }
    #endregion

    #region SERVER_TO_CLIENT
    private void onDataReceived()
    {
        var packet = client.GetPeer(1).GetPacket();
        GD.Print(Encoding.UTF8.GetString(packet));
        var action = JsonSerializer.Deserialize<ServerMessage>(packet);
        switch ((ServerMessageType)action.type)
        {
            case ServerMessageType.MESSAGE_PING:
                GD.Print("ping"); // TODO
                break;
            case ServerMessageType.MESSAGE_JOIN:
                var joinGameRes = JsonSerializer.Deserialize<JoinGameResponse>(action.payload);
                this.onLocalPlayerJoinedGame(joinGameRes);
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
        if (tick.connects != null)
        {
            foreach (var connect in tick.connects)
            {
                if (connect.playerId == player.id) continue;
                opponents.OnPlayerConnected(connect.playerId, connect.spawn, connect.color);
            }
        }

        if (tick.disconnects != null)
        {
            foreach (var disconnect in tick.disconnects)
            {
                this.onDisconnectEventReceived(disconnect);
            }
        }

        if (tick.moves != null)
        {
            foreach (var move in tick.moves)
            {
                this.onMoveEventReceived(move);
            }
        }

        if (tick.attacks != null)
        {
            foreach (var attack in tick.attacks)
            {
                this.onAttackEventReceived(attack);
            }
        }

        if (tick.damage != null)
        {
            GD.Print("damage not none");
            foreach (var damage in tick.damage)
            {
                this.onDamageEventReceived(damage);
            }
        }
    }

    private void onLocalPlayerJoinedGame(JoinGameResponse msg)
    {
        player.id = msg.playerId;
        player.Translation = new Vector3(msg.spawn.x, 0, msg.spawn.z);
        player.SetTeam(msg.color);

        if (msg.others != null)
        {
            foreach (var other in msg.others)
            {
                opponents.OnPlayerConnected(other.playerId, other.spawn, other.color);
            }
        }
    }

    private void onMoveEventReceived(Move move)
    {
        if (move.playerId == player.id)
        {
            player.SetMoving(move.target.ToVector3());
            opponents.StopAttacking(player.id);
            return;
        }

        opponents.SetMoving(move.playerId, move.target.ToVector3());
        player.StopAttacking();
    }

    private void onDisconnectEventReceived(Disconnect disconnect)
    {
        if (disconnect.playerId == player.id)
        {
            this.client.DisconnectFromHost();
            // TODO - handle this in UI
            return;
        }

        this.opponents.OnPlayerDisconnected(disconnect.playerId);
    }

    private void onAttackEventReceived(Attack attack)
    {
        if (attack.sourcePlayerId == player.id)
        {
            this.player.SetAttacking(attack.targetPlayerId, attack.targetPlayerLocation.ToVector3());
            return;
        }

        this.opponents.SetAttacking(attack.sourcePlayerId, attack.targetPlayerId, attack.targetPlayerLocation.ToVector3());
    }

    private void onDamageEventReceived(Damage damage)
    {
        GD.Print("damage received");
        if (damage.targetPlayerId == player.id)
        {
            this.player.ApplyDamage(damage.damageDealt);
            return;
        }

        this.opponents.ApplyDamage(damage.targetPlayerId, damage.damageDealt);
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
                playerId = player.id,
                target = target,
            })
        };
        var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        Error error = client.GetPeer(1).PutPacket(msgBytes);
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
                sourcePlayerId = player.id,
                sourcePlayerLocation = player.CurrentLocation(),
                targetPlayerId = targetPlayerId,
                targetPlayerLocation = opponents.CurrentLocation(targetPlayerId),
            })
        };
        var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        Error error = client.GetPeer(1).PutPacket(msgBytes);
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
                sourcePlayerId = player.id,
                targetPlayerId = targetPlayerId,
                damageDealt = 1,
            })
        };
        var msgBytes = JsonSerializer.SerializeToUtf8Bytes(msg);
        Error error = client.GetPeer(1).PutPacket(msgBytes);
        if (error != Error.Ok)
        {
            GD.Print($"Failed to request move: {error}");
            return;
        }
    }
    #endregion
}
