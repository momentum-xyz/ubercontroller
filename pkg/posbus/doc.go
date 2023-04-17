// The posbus package implements the custom protocol for communication between the backend and frontends.
//
// Originally it implemented messaging for position updates of objects in a 3D space, hence the PosBus name.
// But has been expanded to include more types of information.
//
// Messages are in a compact binary format and usually send through a websocket connection.
package posbus
