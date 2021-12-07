# EncryptedChat
This repo contains both the server package and client package.

### Local usage
1. Setting up the server
	- clone this repository and navigate to the `server/` directory.
	- run `go build` to compile an executable
	- to start the server run `./server`
	- the server is hosted by default on port `:8000`
2. Connecting via a client on the same network
	- navigate to the `client/` directory and run `go build` to create a client executable
	- run `./client {host:port}` (ie if running locally use `./client localhost:8000`)
	- if multiple clients are required, be sure to run each instance in a different directory as each will generate its own `private.key`

### Commands (upon successful connection)
	- **/login {username}** either creates a new user if both the public key and username are unique, or attempts to authenticate if the public key exists in the database of the server
	- **/rooms** - lists all rooms which are availible to connect to
	- **/join {room name}** - attempts to join the specified room, if the room is full (contains two people) the user will not connect, otherwise it creates a new room named `{room name}`
	- **/name {username}** - changes the user's username to `{username}`
	- **/help** - displays all possible commands
	- **/quit** - disconnects from the server

### Messaging
All messages and commands sent through the server are encrypted using public and private keys. Upon entering a room, public keys are distributed to the users of that room. When a message is sent, the message is encrypted using sha256 and the public key of the recipient. When the message is delivered, the recipient decrypts the message using their own private key.

### Key storage
Upon running the `client` file, a `private.key` will be created in the same directory as the `client`, which contains both the public and private key, saved using a `gob`.
