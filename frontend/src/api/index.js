// api/index.js
var url = "ws://"+window.location.hostname+":8000/ws";
var socket = null;

let start = (token,callbackFn) => {
  join(token,callbackFn)
}

let join = (token,callbackFn,gameId="") => {
  if(socket==null){
    console.log("Attempting Connection...");
    socket = new WebSocket(url+"?token="+token+(gameId!==""?"&gameId="+gameId:""));
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

export { join, sendMsg };
