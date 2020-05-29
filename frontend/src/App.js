// App.js
import React, { useState } from "react";
import Header from "./components/Header/Header"
import LoginForm from "./components/LoginForm/LoginForm"
import JoinStartForm from "./components/JoinStartForm/JoinStartForm"
import PlayerList from "./components/PlayerList/PlayerList"
import "./App.css";
import { join, sendMsg } from "./api";
import { Button,Box,Grommet } from 'grommet';

function send(msg) {
  console.log(msg);
  sendMsg(msg);
}

function App(props) {
  // const [message, setMessage] = useState({ "data": "" });
  const [name, setName] = useState("");
  const [token, setToken] = useState( "");
  const [roomCode, setRoomCode] = useState( "");
  const [players, setPlayers] = useState( []);

  const onLogin = function(data){
    setToken(data)
    join(data,(d)=>console.dir(d));
  };

  const onJoin = function(roomCode){
    setRoomCode(roomCode)
  };

  let msgText = message.data;

  const theme = {
    global: {
      font: {
        family: 'sans-serif',
        size: '14px',
        height: '20px',
      },
    },
  };

  let hasToken = token != "";
  let hasRoomCode = roomCode != "";
  return (
    <Grommet theme={theme}>
      <div className="App">
        <Header />
        {!hasToken?<LoginForm onLogin={onLogin} />:null} 
        {hasToken&&!hasRoomCode?<JoinStartForm onJoin={onJoin}/>:null}
        {hasRoomCode?<PlayerList />:null}
      
        {/* <Button onClick={() => { send("hello") }}>Press me</Button> */}
      </div>
    </Grommet>
  );
}

export default App;
