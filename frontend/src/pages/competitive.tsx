import { createSignal, For } from "solid-js";

type WaitingUser = {
  username: string;
  rank: string;
  elo: number;
};

const [getWaitingUsers, setWaitingUsers] = createSignal<WaitingUser[]>([]);

const websocket = () => {
  let socket = new WebSocket("ws://localhost:5000/ws/liveWaitingList");

  socket.onopen = (e) => {
    console.log("Websocket wurde erstellt");
  };

  socket.onmessage = (e: any) => {
    setWaitingUsers(JSON.parse(e.data));
  };

  socket.onclose = (e) => {
    if (e.wasClean) console.log("Ferbindung wurde geschlossen");
    else console.log("Verbindung wurde unterbrochen");
  };

  socket.onerror = (e) => {
    console.error("Etwas ist schiefgelaufen: " + e);
  };
};

export default () => {
  websocket();

  return (
    <div>
      <button class="button">In Warteraum eintreten</button>
      <table>
        <thead>
          <tr>
            <th>Benutzername</th>
            <th>Rang</th>
            <th>Elo Zahl</th>
          </tr>
        </thead>
        <tbody>
          <For each={getWaitingUsers()}>
            {(wu) => (
              <tr>
                <td>{wu.username}</td>
                <td>{wu.rank}</td>
                <td>{wu.elo}</td>
              </tr>
            )}
          </For>
        </tbody>
      </table>
    </div>
  );
};
