// let websockt = new WebSocket("ws://localhost:8080/initwebsocket")
// console.log("Attempting websocket connection")
//
// websockt.onopen = () => {
//     console.log("Websocket connection established.")
//     websockt.send("hello from client!!")
// }
//
// websockt.onmessage = (msg) => {
//     addReceivedMsgElement(msg.data, new Date())
//     console.log("Message received from server: ", msg)
// }
//
// websockt.onclose = (event) => {
//     console.log("Websocket connection closed. ", event)
// }
//
// websockt.onerror = (err) => {
//     console.log("Websocket connection error: ", err)
// }
//
// document.getElementById("sendMsgButton").onclick = () => {
//     var msg = document.getElementById("msgTxt").value
//     if (msg === "") {
//         return
//     }
//
//     websockt.send(msg)
//     addMyMsgElement(msg, new Date())
//     document.getElementById("msgTxt").textContent = ""
// }
//
// function addReceivedMsgElement(msg, currentTime) {
//     var chatWindow = document.getElementById("chat-window");
//     chatWindow.innerHTML += `<div class="container" id="received-msg-element">
//     <img src="/w3images/bandmember.jpg" alt="Avatar" style="width:100%;">
//     <p>`+ msg +`</p>
//     <span class="time-right">`+currentTime+`</span>
//     </div>`
// }
//
// function addMyMsgElement(msg, currentTime) {
//     var chatWindow = document.getElementById("chat-window");
//     chatWindow.innerHTML += `<div class="container darker" id="my-msg-element">
//     <img src="/w3images/avatar_g2.jpg" alt="Avatar" class="right" style="width:100%;">
//     <p>`+ msg +`</p>
//     <span class="time-left">`+currentTime+`</span>
// </div>`
// }