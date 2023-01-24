extends Node

export var websocket_url = "ws://kilnwood-game.com/connect"
var _client = WebSocketClient.new()

var connected = false
var id = 0

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
		print("got no data")
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
			on_local_player_joined_game(json.result)
		"join-broadcast":
			if json.result.PlayerId == id: return
			print(packet.get_string_from_utf8())
			on_remote_player_joined_game(json.result)
		"tick":
			if len(json.result.Events) == 0: return
			print(packet.get_string_from_utf8())
			process_tick(json.result.Events)

func process_tick(events):
	for event in events:
		var split_event = event.split(":")
		match split_event[0]:
			"m": on_move_event_received(split_event)
			"a": on_attack_event_received(split_event)
		print(event)

func player_requested_move(location):
	var err = _client.get_peer(1).put_packet("m:{id}:{x}:{z}".format({"id": id, "x": location.x, "z": location.z}).to_utf8())
	if err:
		print(err)

func player_requested_attack(target_id):
	var err = _client.get_peer(1).put_packet("a:{source}:{target}".format({"source": id, "target": target_id}).to_utf8())
	if err:
		print(err)

func on_local_player_joined_game(msg):
	id = msg.PlayerId
	$"../Player".translation = Vector3(msg.Spawn.X, 0, msg.Spawn.Z)
	$"../Player".set_team("team", msg.Color)
	$"../Player".set_id(id)
	
	if not msg.Others: return
	for player in msg.Others:
		print(player)
		on_remote_player_joined_game(player)

func on_remote_player_joined_game(msg):
	#$"../Opponent1".translation = Vector3(msg.Spawn.X, 0, msg.Spawn.Z)
	$"../Opponent1".translation = Vector3(10, 0, 10)
	$"../Opponent1".set_team("team", msg.Color)
	$"../Opponent1".set_id(msg.PlayerId)
	$"../Opponent1".visible = true

func on_move_event_received(event):
	if event[1] == id:
		$"../Player".set_moving(Vector3(event[2], 0, event[3]))
		return
	$"../Opponent1".set_moving(Vector3(event[2], 0, event[3]))

func on_attack_event_received(event):
	print("attack nyi")
