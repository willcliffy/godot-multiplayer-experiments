using Godot;
using System;
using System.IO;
using Game;
using Google.Protobuf;

public class MessageBroker : Node
{
    // onready var localPlayer = $Player
    // onready var opponentController = $OpponentController

    [Export]
    string webSocketURL = "ws://localhost:8080/connect";

    WebSocketClient client = null;

    public override void _Ready()
    {
        client = new WebSocketClient();

        client.Connect("connection_established", this, nameof(OnConnectionEstablished));
        client.Connect("data_received", this, nameof(OnDataRecived));
        client.Connect("server_close_request", this, nameof(OnServerCloseRequest));
        client.Connect("connection_closed", this, nameof(OnConnectionClosed));
        client.Connect("connection_error", this, nameof(OnConnectionError));

        Error error = client.ConnectToUrl(webSocketURL, null, false);
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

    public void OnConnectionEstablished(string protocol)
    {

        Game.Action msg = new Game.Action();
        msg.Type = ActionType.ActionConnect;
        MemoryStream stream = new MemoryStream();
        msg.WriteTo(stream);
        Error error = client.GetPeer(1).PutPacket(stream.ToArray());
        if (error != Error.Ok)
        {
            GD.Print("Failed to establish connection: " + error);
            return;
        }
        GD.Print($"Connection established with proto {protocol}");
    }

    public void OnDataRecived()
    {
        // Example only: Handle reciving text from server.
        var packet = client.GetPeer(1).GetPacket();
        GD.Print(packet);
    }

    public void OnServerCloseRequest(int code, string reason)
    {
        GD.Print("Close request, reason: " + reason);
    }

    public void OnConnectionClosed(bool wasCleanClose)
    {
        GD.Print("Connection closed. was clean close: " + wasCleanClose.ToString());
    }

    public void OnConnectionError()
    {
        GD.Print("Connection error.");
    }
}
