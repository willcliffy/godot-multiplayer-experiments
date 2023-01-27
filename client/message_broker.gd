extends Node

onready var localPlayer = $Player
onready var opponentController = $OpponentController

export var websocket_url = "ws://kilnwood-game.com/connect"
var _client = WebSocketClient.new()

var connected = false
var num_connected = false

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
	if err != OK:
		print("failed to connect: ", err)
		return

	connected = true
	print("connected")

func _process(_delta):
	_client.poll()

func _on_data():
	print("on_data")
	if not connected:
		return

	var packet = _client.get_peer(1).get_packet()
	if not packet:
		print("got no data")
		return

	var json = JSON.parse(packet.get_string_from_utf8())
	if json.error:
		print("got error from json parse: ", json.error)
		print(packet.get_string_from_utf8())
		return
	if not json.result or not json.result.has("Type"):
		print("unexpected json from server: ", json.result)
		return

	match json.result.Type:
		"join-response":
			on_local_player_joined_game(json.result)
		"join-broadcast":
			if json.result.PlayerId == localPlayer.get_id(): return
			opponentController.joined(json.result.PlayerId, json.result.Spawn, json.result.Color)
		"tick":
			if len(json.result.Events) == 0: return
			process_tick(json.result.Events)

func process_tick(events):
	for event in events:
		print(event)
		var split_event = event.split(":")
		match split_event[0]:
			"m": on_move_event_received(split_event)
			"a": on_attack_event_received(split_event)
			"d": opponentController.delete(split_event[1])

func on_local_player_joined_game(msg):
	print(msg)
	localPlayer.translation = Vector3(msg.Spawn.X, 0, msg.Spawn.Z)
	localPlayer.set_team("team", msg.Color)
	localPlayer.set_id(msg.PlayerId)
	
	if not msg.Others: return
	for player in msg.Others:
		opponentController.joined(player.PlayerId, player.Spawn, player.Color)

func on_move_event_received(event):
	var sourcePlayerId = event[1]
	if event[1] == localPlayer.get_id():
		localPlayer.set_moving(Vector3(event[2], 0, event[3]))
		return

	opponentController.get(sourcePlayerId).set_moving(Vector3(event[2], 0, event[3]))

func on_attack_event_received(event):
	var sourcePlayerId = event[1]
	if sourcePlayerId == localPlayer.get_id():
		localPlayer.set_attacking(Vector3(event[2], 0, event[3]))
		return

	if not opponentController.has(sourcePlayerId):
		print("NETWORK BADNESS: got message from {s} but wasnt in game!".format({"s": sourcePlayerId}))

	opponentController.get(sourcePlayerId).set_attacking(Vector3(event[2], 0, event[3]))

func player_requested_move(location):
	var msg = "m:{source}:{x}:{z}".format({
		"source": localPlayer.get_id(),
		"x": location.x,
		"z": location.z
	})

	var err = _client.get_peer(1).put_packet(msg.to_utf8())
	if err: print(err)

func player_requested_attack(target_id):
	var msg = "a:{source}:{target}".format({
		"source": localPlayer.get_id(),
		"target": target_id
	})
	var err = _client.get_peer(1).put_packet(msg.to_utf8())
	if err: print(err)
