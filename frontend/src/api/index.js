// api/index.js
var url = "ws://"+window.location.hostname+":8000/ws";
var socket = null;

let connect = (token,callbackFn) => {
  if(socket==null){
    console.log("Attempting Connection...");
    socket = new WebSocket(url+"?token="+token);
  }

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = callbackFn;

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = msg => {
  console.log("sending msg: ", msg);
  socket.send(msg);
};

export { connect, sendMsg };
