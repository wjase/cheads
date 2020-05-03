import React, { useState, useEffect } from "react";
import "./JoinStartForm.scss";
import { connect, sendMsg } from "../../api";

function send(msg) {
  console.log(msg);
  sendMsg(msg);
}

function onSubmit(event, resultFn) {
  event.preventDefault();
  const formData = new FormData(event.target);
  const obj = {}

  for (let entry of formData.entries()) {
    obj[entry[0]] = entry[1]
  }

  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(obj)
  };
  
  fetch('/signin', requestOptions)
    .then(response => {
      response.text().then(function(data) {
        console.log(data);
        resultFn(data);
        });
    })

  return false;
}

const JoinStartForm = function(show){
  const [userData, setUserData] = useState({ "user": "" });
  let onJoin = ()=>{
  };
  
  let onStartGame = ()=>{
  };

  console.log("show",show);
  let comp= (
    <div className="join-start-form"  >
      <form onSubmit={(e)=>onSubmit(e,onJoin)}>
        <label htmlFor="username">Room Code:</label>
        <input id="roomcode" name="roomcode" type="text" /><br />
        <input type="submit" value="Join Game" /><br />
      </form>
      <form onSubmit={(e)=>onSubmit(e,onStartGame)}>
        <input type="submit" value="Start Game" /><br />
      </form>
    </div>
  );
  
  if(userData.user!=""){
    comp=(
    <div>Totes logged in as{userData.user}</div>
    );
  }

  return  comp;
} 

export default JoinStartForm;