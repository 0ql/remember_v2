import { For, Match, Switch, useContext } from "solid-js";
import { GlobalContext, joinWaitingRoom, leaveWaitingRoom, Status } from "../socket";

const [state, setState] = useContext(GlobalContext)

export default () => {
	return (
		<div>
			<Switch fallback={<button class="button" onClick={leaveWaitingRoom}>Warteraum verlassen</button>}>
				<Match when={state.clients.filter(user => user.Username === state.userdata?.Username)[0]?.Status === Status.Idle}>
					<button class="button" onClick={joinWaitingRoom}>In Warteraum eintreten</button>
				</Match>
			</Switch>
			<table>
				<thead>
					<tr>
						<th>Status</th>
						<th>Benutzername</th>
						<th>Rang</th>
						<th>Elo Zahl</th>
					</tr>
				</thead>
				<tbody>
					<For each={state.clients}>
						{(wu) => (
							<tr>
								<td>
									<Switch>
										<Match when={wu.Status === Status.Idle}>Idle</Match>
										<Match when={wu.Status === Status.Waiting}>Waiting</Match>
										<Match when={wu.Status === Status.InGame}>In Game</Match>
									</Switch>
								</td>
								<td>{wu.Username}</td>
								<td>{wu.Rank}</td>
								<td>{wu.Elo}</td>
							</tr>
						)}
					</For>
				</tbody>
			</table>
		</div>
	);
};
