window.onload = () => {
  const lobbyMessageDialog = document.getElementById("lobby-message-dialog");
  if (lobbyMessageDialog) lobbyMessageDialog.showModal();

  let wsProtocol = "wss://";
  if (document.location.protocol === "http:") {
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

  lobbyChatForm.onsubmit = (event) => {
    event.preventDefault();
    if (!lobbyChatInput.value) return;
    conn.send(lobbyChatInput.value);
    lobbyChatInput.value = "";
  };

  conn.onmessage = (event) => {
    if (event.data === "refresh") {
      htmx.ajax(
        "GET",
        "/api" + document.location.pathname + "/game-interface",
        { target: "#lobby-grid-interface" }
      );
      return;
    }

    if (event.data === "kick") {
      document.location.href = "/lobbies";
      return;
    }

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

let lobbyPlayerDataScrollTop = 0;
let lobbyGameBoardScrollTop = 0;
let lobbyGameStatsScrollTop = 0;

document.addEventListener("htmx:beforeSwap", function () {
  const lobbyPlayerData = document.getElementById("lobby-player-data");
  if (lobbyPlayerData) lobbyPlayerDataScrollTop = lobbyPlayerData.scrollTop;
  const lobbyGameBoard = document.getElementById("lobby-game-board");
  if (lobbyGameBoard) lobbyGameBoardScrollTop = lobbyGameBoard.scrollTop;
  const lobbyGameStats = document.getElementById("lobby-game-stats");
  if (lobbyGameStats) lobbyGameStatsScrollTop = lobbyGameStats.scrollTop;
});

document.addEventListener("htmx:afterSwap", function () {
  const lobbyPlayerData = document.getElementById("lobby-player-data");
  if (lobbyPlayerData) lobbyPlayerData.scrollTop = lobbyPlayerDataScrollTop;
  const lobbyGameBoard = document.getElementById("lobby-game-board");
  if (lobbyGameBoard) lobbyGameBoard.scrollTop = lobbyGameBoardScrollTop;
  const lobbyGameStats = document.getElementById("lobby-game-stats");
  if (lobbyGameStats) lobbyGameStats.scrollTop = lobbyGameStatsScrollTop;
});
