package main

import (
	"encoding/json"
	"github.com/googollee/go-socket.io"
	"github.com/satori/go.uuid"
	"go-signal-server/component"
	"log"
	"net/http"
)

type ExtendedServer struct {
	server *socketio.Server
	result
}

type handshakeDetail struct {
	from     string
	to       string
	sid      string
	roomType string
	payload
	prefix string
}

type payload struct {
	candidate
}

type candidate struct {
	candidate     string
	sdpMid        string
	sdpMLineIndex int
}

func main() {
	eServer := &ExtendedServer{}
	server, err := socketio.NewServer(nil)
	eServer.server = server
	if err != nil {
		log.Fatal(err)
	}

	defaultAdaptor := make(socketio.Broadcast)

	roomCtrl := RoomControl.NewRoomControl(defaultAdaptor)
	server.SetAdaptor(roomCtrl)

	log.Printf("roomCtrl is %+v", roomCtrl)

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so = so.SetResources(false, true, false)

		log.Printf("%+v", so)

		so.On("message", func(msg []byte) {

			var err error
			var m []byte

			hsDetail := &handshakeDetail{}

			err = json.Unmarshal(msg, hsDetail)

			if err != nil {
				log.Fatal(err)
			}

			// {
			// 	to: '/#PVVNhym5878Ldu0mAAAC',
			// 	sid: '1452784289306',
			// 	roomType: 'video',
			// 	type: 'candidate',
			// 	payload:
			// 	{ candidate:
			// 	  { candidate: 'candidate:3545837919 2 udp 2122129150 192.168.42.44 53907 typ host generation 0',
			// 	    sdpMid: 'video',
			// 	    sdpMLineIndex: 1 } },
			// 	prefix: 'webkit'
			// }

			hsDetail.from = so.Id()
			m, err = json.Marshal(hsDetail)

			if err != nil {
				log.Fatal(err)
			}

			log.Println("emit:", so.Emit("message", m))
			server.BroadcastTo(hsDetail.to, "message", m)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})

		so.On("create", func(name string) {
			log.Printf("socket name is %s", name)

			// if (arguments.length == 2) {
			//     cb = (typeof cb == 'function') ? cb : function() {};
			//     name = name || uuid();
			// } else {
			//     cb = name;
			//     name = uuid();
			// }
			log.Println(server)
			log.Printf("uuid is %s", uuid.NewV4())

			so.Emit("room factory message back", name)

			roomCtrl.Join(name, so)

			// check if exists
			// var room = io.nsps['/'].adapter.rooms[name];
			// if (room && room.length) {
			//     safeCb(cb)('taken');
			// } else {
			//     join(name);
			//     safeCb(cb)(null, name);
			// }
		})

		so.On("join", func(name string) {

			joinFeedBack := roomCtrl.DescribeRoom(name)
			so.Emit("join feed back", joinFeedBack)
			roomCtrl.Join(name, so)

		})

	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8888...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}

// function clientsInRoom(name) {
//     return io.sockets.clients(name).length;
// }

// io.sockets.on('connection', function(client) {
//         client.resources = {
//             screen: false,
//             video: true,
//             audio: false
//         };

/**
 * pass a message to another id
 * offer/answer information
 */

//client.on('join', join);

// function removeFeed(type) {
//     if (client.room) {
//         io.sockets.in(client.room).emit('remove', {
//             id: client.id,
//             type: type
//         });
//         if (!type) {
//             client.leave(client.room);
//             client.room = undefined;
//         }
//     }
// }

//         function join(name, cb) {
//             // sanity check
//             if (typeof name !== 'string') return;
//             // check if maximum number of clients reached
//             if (config.rooms && config.rooms.maxClients > 0 &&
//                 clientsInRoom(name) >= config.rooms.maxClients) {
//                 safeCb(cb)('full');
//                 return;
//             }
//             // leave any existing rooms
//             removeFeed();
//             safeCb(cb)(null, describeRoom(name));
//             client.join(name);
//             client.room = name;
//         }

//         // we don't want to pass "leave" directly because the
//         // event type string of "socket end" gets passed too.
//         client.on('disconnect', function() {
//             removeFeed();
//         });
//         client.on('leave', function() {
//             removeFeed();
//         });

//         client.on('create', function(name, cb) {
//             console.log('socket name is %s', name);
//             if (arguments.length == 2) {
//                 cb = (typeof cb == 'function') ? cb : function() {};
//                 name = name || uuid();
//             } else {
//                 cb = name;
//                 name = uuid();
//             }
//             console.log('uuid is %s', uuid());
//             // check if exists
//             var room = io.nsps['/'].adapter.rooms[name];
//             if (room && room.length) {
//                 safeCb(cb)('taken');
//             } else {
//                 join(name);
//                 safeCb(cb)(null, name);
//             }
//         });

// };

// // function safeCb(cb) {
// //     if (typeof cb === 'function') {
// //         return cb;
// //     } else {
// //         return function() {};
// //     }
// // }
