using Godot;

namespace Game
{
    public partial class Location
    {
        public int x { get; set; }
        public int z { get; set; }

        public Vector3 ToVector3()
        {
            return new Vector3(this.x, 0, this.z);
        }
    }

    public enum ClientActionType
    {
        ACTION_PING = 0,
        ACTION_CONNECT = 1,
        ACTION_DISCONNECT = 2,
        ACTION_MOVE = 3,
        ACTION_ATTACK = 4,
        ACTION_DAMAGE = 5,
        ACTION_DEATH = 6,
        ACTION_RESPAWN = 7,
    }

    public partial class ClientAction
    {
        public int type { get; set; }
        public string payload { get; set; }
    }

    public partial class Connect
    {
        public ulong playerId { get; set; }
        public string color { get; set; }
        public Location spawn { get; set; }
    }

    public partial class Disconnect
    {
        public ulong playerId { get; set; }
    }

    public partial class Move
    {
        public ulong playerId { get; set; }
        public Location target { get; set; }
    }

    public partial class Attack
    {
        public ulong sourcePlayerId { get; set; }
        public Location sourcePlayerLocation { get; set; }
        public ulong targetPlayerId { get; set; }
        public Location targetPlayerLocation { get; set; }
    }

    public partial class Damage
    {
        public ulong sourcePlayerId { get; set; }
        public ulong targetPlayerId { get; set; }
        public int damageDealt { get; set; }
    }


    public partial class Death
    {
        public ulong playerId { get; set; }
        public Location location { get; set; }
    }

    public partial class Respawn
    {
        public ulong playerId { get; set; }
        public Location spawn { get; set; }
    }

    public enum ServerMessageType
    {
        MESSAGE_PING = 0,
        MESSAGE_JOIN = 1,
        MESSAGE_TICK = 2,
    }

    public partial class ServerMessage
    {
        public int type { get; set; }
        public string payload { get; set; }
    }

    public partial class JoinGameResponse
    {
        public ulong playerId { get; set; }
        public string color { get; set; }
        public Location spawn { get; set; }
        public Connect[] others { get; set; }
    }

    public partial class GameTickAction
    {
        public uint type { get; set; }
        public string value { get; set; }
    }

    public partial class GameTick
    {
        public uint tick { get; set; }
        public GameTickAction[] actions { get; set; }
    }
}
