import React from "react";
import "./PlayerList.scss";
import { sendMsg } from "../../api";

function send(msg) {
  sendMsg(msg);
}

const PlayerList = function(props){

  let comp= (
    <div className="player-list" >
      players here
    </div>
  );

  return  comp;
} 

export default PlayerList;