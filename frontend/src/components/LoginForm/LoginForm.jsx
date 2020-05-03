import React, { useState, useEffect } from "react";
import "./LoginForm.scss";
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

const LoginForm = function(props){

  const [userData, setUserData] = useState({ "user": "" });

  let comp= (
    <div className="join-form" >
      <form onSubmit={(e)=>onSubmit(e,props.onLogin)}>
        <label htmlFor="username">Username:</label>
        <input id="username" name="username" type="text" /><br />
        <label htmlFor="password">Password:</label>
        <input id="password" name="password" type="text" /><br />
        <input type="submit" value="Login" /><br />
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

export default LoginForm;