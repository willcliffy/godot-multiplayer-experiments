extends Node

export var websocket_url = "ws://kilnwood-game.com/connect"
var _client = WebSocketClient.new()

var connected = false
var id = 0
var team = 0

func _ready():
	_client.connect("connection_closed", self, "_closed")
	_client.connect("connection_error", self, "_closed")
	_client.connect("connection_established", self, "_connected")
	_client.connect("data_received", self, "_on_data")

	var err = _client.connect_to_url(websocket_url)
	if err != OK:
		print("Unable to connect")
		set_process(false)

func _closed(was_clean = false):
	print("Closed, clean: ", was_clean)
	set_process(false)

func _connected(_proto = ""):
	var err = _client.get_peer(1).put_packet("J:::f".to_utf8())
	if err:
		print("failed to connect: ", err)
		return

	connected = true

func _process(_delta):
	_client.poll()

func _on_data():
	if not connected:
		return

	var packet = _client.get_peer(1).get_packet()
	if not packet:
		return

	var json = JSON.parse(packet.get_string_from_utf8())
	if json.error:
		print("got error from json parse: ", json.error)
		return
	if not json.result or not json.result.has("Type"):
		print("unexpected json from server: ", json.result)
		return

	match json.result.Type:
		"join-response":
			print(packet.get_string_from_utf8())
			id = json.result.PlayerId
			$"../Player".translation = Vector3(json.result.Spawn.X, 0, json.result.Spawn.Z)
			$"../Player".set_team(json.result.Team)
			if not json.result.Others: return
			for player in json.result.Others:
				print(player)
				$"../Opponent1".translation = Vector3(player.Location.X, 0, player.Location.Z)
				$"../Opponent1".set_team(player.Team)
				$"../Opponent1".visible = true
		"join-broadcast":
			if json.result.PlayerId == id: return
			print(packet.get_string_from_utf8())
			$"../Opponent1".translation = Vector3(json.result.Spawn.X, 0, json.result.Spawn.Z)
			$"../Opponent1".set_team(json.result.Team)
			$"../Opponent1".visible = true
		"tick":
			if len(json.result.Events) == 0: return
			print(packet.get_string_from_utf8())
			process_tick(json.result.Events)

func process_tick(events):
	for event in events:
		var split_event = event.split(":")
		match split_event[0]:
			"m":
				if split_event[1] == id: return
				$"../Opponent1".set_moving(Vector3(split_event[2], 0, split_event[3]))
			"a": print("attack nyi")
		print(event)

func on_player_move(location):
	var err = _client.get_peer(1).put_packet("m:{id}:{x}:{z}".format({"id": id, "x": location.x, "z": location.z}).to_utf8())
	if err:
		print(err)
