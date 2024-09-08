window.onload = () => {
  let wsProtocol = "wss://";
  if (document.location.protocol == "http:") {
    wsProtocol = "ws://";
  }

  const conn = new WebSocket(
    wsProtocol + document.location.host + "/ws" + document.location.pathname
  );

  conn.onclose = () => {
    alert("Connection Lost");
    document.location.href = "/lobbies";
  };

  const chatbox = document.getElementById("chatbox");
  conn.onmessage = (event) => {
    const message = document.createElement("p");
    message.innerHTML = event.data;
    chatbox.appendChild(message);
  };
};
