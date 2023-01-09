extends Node

var udp = PacketPeerUDP.new()
var connected = false
var id = 0


func _ready():
	var err = udp.connect_to_host("kilnwood-production.up.railway.app", 4444)
	if err:
		print("error connecting to host: ", err)
		return

	err = udp.put_packet("J:::f".to_utf8())
	if err:
		print("error joining game: ", err)
		return

	connected = true


func _process(_delta):
	if not connected: return

	var packet = udp.get_packet()
	if not packet: return

	var json = JSON.parse(packet.get_string_from_utf8())
	if json.error:
		print("got error from json parse: ", json.error)
		return
	if not json.result or not json.result.has("Type"):
		print("unexpected json from server: ", json.result)
		return
	
	match json.result.Type:
		"join-response":
			id = json.result.PlayerId
			$"../Player".translation = Vector3(json.result.Spawn.X, 0, json.result.Spawn.Z)
		"join-broadcast":
			if json.result.PlayerId == id: return
			$"../Opponent".translation = Vector3(json.result.Spawn.X, 0, json.result.Spawn.Z)
			$"../Opponent".visible = true
		"tick":
			process_tick(json.result.Events)


func process_tick(events):
	for event in events:
		var split_event = event.split(":")
		match split_event[0]:
			"m":
				if split_event[1] == id: return
				$"../Opponent".set_moving(Vector3(split_event[2], 0, split_event[3]))
			"a": print("attack nyi")
		print(event)


func on_player_move(position):
	print("player moving to", position)
	var err = udp.put_packet("m:{id}:{x}:{y}".format({"id": id, "x": position.x, "y": position.y}).to_utf8())
	if err:
		print(err)
