FROM golang:1.19-alpine AS server-builder

RUN mkdir -p /app
WORKDIR /app
COPY server/ .

RUN COOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/dist/app /app/main.go

FROM willcliffy/kilnwood-webclient-builder:v0.0.0 as client-builder

# init webclient builder
WORKDIR /client/_engine

ENV GODOT_VERSION "3.5.1"

RUN unzip -q Godot_v${GODOT_VERSION}-stable_mono_export_templates.tpz
RUN mkdir -p ~/.local/share/godot/templates/3.5.1.stable.mono
RUN mv templates/* ~/.local/share/godot/templates/${GODOT_VERSION}.stable.mono

RUN unzip -q Godot_v${GODOT_VERSION}-stable_mono_linux_headless_64.zip
RUN mv Godot_v${GODOT_VERSION}-stable_mono_linux_headless_64/GodotSharp /usr/local/bin/
RUN mv Godot_v${GODOT_VERSION}-stable_mono_linux_headless_64/Godot_v${GODOT_VERSION}-stable_mono_linux_headless.64 /usr/local/bin/godot
# end init webclient builder

WORKDIR /client
COPY client/ .

RUN mkdir -p dist
RUN godot --path . --export "HTML5" dist/index.html

FROM alpine:3.17

RUN mkdir -p /app
WORKDIR /app

COPY --from=server-builder /app/dist/app /app/app
COPY --from=client-builder /client/dist/ /app/dist/

EXPOSE 8080
CMD /app/app
