# Kilnwood Monorepo

## Quirks

- Going to use protobuf to define types, but actually going to just use json serialization. The chain of reasoning was:
  - I don't want to swap from Godot, and I don't want to lose the web preview yet as it's crazy useful for quick testing
  - The web preview requires the networking to be over websocket. Websockets (or at least the golang impl Gorilla) does not support non-UTF8 encodings, which includes protobuf's raw binary
  - I don't want to permanently shoehorn myself to that decision, so I'll continue to define the proto schema and use generated types where reasonable
