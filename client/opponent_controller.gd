extends Node

onready var available = [$"1", $"2", $"3", $"4", $"5", $"6", $"7"]

var opponentDict = {}

func has(id):
	return opponentDict.has(id)

func get(id):
	if not opponentDict.has(id): return
	return opponentDict[id]

func delete(id):
	if not opponentDict.has(id): return
	available.append(opponentDict[id])
	opponentDict.delete(id)

func joined(id, spawn, color):
	print("joined: ", id, spawn, color)
	var playerObject
	if opponentDict.has(id):
		playerObject = opponentDict[id]
		print("player obj from dict")
	else:
		if len(available) == 0:
			print("SERVER BADNESS: game overflow")
			return
		playerObject = available.pop_front()
		opponentDict[id] = playerObject
		print("player obj from available")
		print(available)

	playerObject.translation = Vector3(spawn.X, 0, spawn.Z)
	playerObject.set_team("team", color)
	playerObject.set_id(id)
	playerObject.visible = true