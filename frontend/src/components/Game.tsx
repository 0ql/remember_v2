import { For, Match, Switch, createSignal, useContext } from "solid-js"
import { GlobalContext, endGame } from "../socket"

const [state, setState] = useContext(GlobalContext)

export default () => {
	const [index, setIndex] = createSignal(0)
	if (state.game) {
		console.log()
	}
	const answers: number[] = []
	return (
		<Switch fallback={<h1>Warte auf Gegner...</h1>}>
			<Match when={index() < 5}>
				<h1>Game</h1>
				<h3>{state.game?.Questions[index()].Question}</h3>
				<For each={state.game?.Questions[index()].Options}>
					{(q, i) => (
						<button onClick={() => {
							answers[index()] = i()
							if (index() < 4) {
								setIndex(index() + 1)
							} else {
								endGame(answers)
								setIndex(index() + 1)
							}
						}}>{q}</button>
					)}
				</For>
			</Match>
		</Switch>
	)
}
