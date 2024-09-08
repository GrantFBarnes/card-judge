window.onload = function () {
  let wsProtocol = "wss://";
  if (document.location.protocol == "http:") {
    wsProtocol = "ws://";
  }

  const conn = new WebSocket(
    wsProtocol + document.location.host + "/ws" + document.location.pathname
  );

  conn.onopen = function (evt) {
    conn.send("connection made");
  };

  conn.onclose = function (evt) {
    console.log("disconnected");
  };

  conn.onmessage = function (evt) {
    console.log("message:", evt.data);
  };
};
