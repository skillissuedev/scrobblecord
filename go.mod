module scrobblecord

go 1.25.1

require (
	github.com/rikkuness/discord-rpc v0.0.0-20250917191941-23d4e4280ce7
	github.com/shkh/lastfm-go v0.0.0-20191215035245-89a801c244e0
)

require (
	github.com/google/uuid v1.1.1 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)

replace github.com/rikkuness/discord-rpc v0.0.0-20250917191941-23d4e4280ce7 => ../discord-rpc/
