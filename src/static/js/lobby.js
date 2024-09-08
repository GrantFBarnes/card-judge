window.onload = () => {
  let wsProtocol = "wss://";
  if (document.location.protocol == "http:") {
    wsProtocol = "ws://";
  }

  const conn = new WebSocket(
    wsProtocol + document.location.host + "/ws" + document.location.pathname
  );

  if (!conn) {
    alert("Failed to make connection.");
    document.location.href = "/lobbies";
  }

  conn.onclose = () => {
    alert("Connection Lost");
    document.location.href = "/lobbies";
  };

  const lobbyChatForm = document.getElementById("lobby-chat-form");
  const lobbyChatMessages = document.getElementById("lobby-chat-messages");
  const lobbyChatInput = document.getElementById("lobby-chat-input");

  lobbyChatForm.onsubmit = () => {
    if (!lobbyChatInput.value) return;
    conn.send(lobbyChatInput.value);
    lobbyChatInput.value = "";
  };

  conn.onmessage = (event) => {
    const message = document.createElement("div");
    message.innerText = event.data;
    lobbyChatMessages.appendChild(message);

    while (lobbyChatMessages.childNodes.length > 100) {
      lobbyChatMessages.removeChild(lobbyChatMessages.childNodes[0]);
    }

    lobbyChatMessages.scrollTop =
      lobbyChatMessages.scrollHeight - lobbyChatMessages.clientHeight;
  };
};
