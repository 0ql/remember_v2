import { For, Show, useContext } from "solid-js"
import { GlobalContext } from "../socket"

const [state, setState] = useContext(GlobalContext)

export default () => {
	return (
		<>
			<h1>Results</h1>
			<For each={state.game?.Questions}>{(question, i) =>
				<>
					<h2>{question.Question}</h2>
					<Show when={state.answers[i()] === state.results.Answers[i()]} fallback={<h3>Answered Wrongly</h3>}>
						<h3>Answered Correctly</h3>
					</Show>
					<div>Your answer: {question.Options[state.answers[i()]]}</div>
					<div>Correct answer: {question.Options[state.results.Answers[i()]]}</div>
				</>
			}</For>
		</>
	)
}
