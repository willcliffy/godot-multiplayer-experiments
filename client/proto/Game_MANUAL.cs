namespace Game {
    public class Location {
        public uint x { get; set; }
        public uint z { get; set; }
    }

    public enum ClientActionType {
        ACTION_PING = 0,
        ACTION_CONNECT = 1,
        ACTION_DISCONNECT = 2,
        ACTION_MOVE = 3,
        ACTION_ATTACK = 4,
    }

    public class ClientAction {
        public int type { get; set; }
        public string payload { get; set; }
    }

    public class Connect {
        public uint64 playerId;
        public string color;
        public Location spawn;
    }

    public class Disconnect {
        public ulong playerId { get; set; }
    }

    public class Move {
        public ulong playerId { get; set; }
        public Location target { get; set; }
    }

    public class Attack {
        public ulong sourcePlayerId { get; set; }
        public ulong targetPlayerId { get; set; }
    }

    public class Damage {
        public ulong sourcePlayerId { get; set; }
        public ulong targetPlayerId { get; set; }
    }

    public enum ServerMessageType {
        MESSAGE_PING = 0,
        MESSAGE_JOIN = 1,
        MESSAGE_TICK = 2,
    }

    public class ServerMessage {
        public int type { get; set; }
        public string payload { get; set; }
    }

    public class JoinGameResponse {
        public ulong playerId { get; set; }
        public string color { get; set; }
        public Location spawn { get; set; }

        public Connect others { get; set; }
    }

    public class GameTick {
        public uint tick { get; set; }
        public Disconnect[] disconnects { get; set; }
        public Move[] moves { get; set; }
        public Attack[] attacks { get; set; }
    }
}
