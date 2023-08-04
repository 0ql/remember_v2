import {
	Component,
	createEffect,
	createResource,
	createSignal,
} from "solid-js";

const fetchQuestions = async () =>
	(
		await fetch("http://localhost:5000/api/v1/getAllQuestions", {
			credentials: "include",
		})
	).json();

const deleteQuestion = async (id: number) =>
	await fetch("http://localhost:5000/api/v1/deleteQuestionById/" + id, {
		method: "DELETE",
		credentials: "include",
	});

export default () => {
	const [questionList] = createResource(fetchQuestions, {
		initialValue: "loading...",
	});

	const [getTable, setTable] = createSignal(["Loading"]);

	createEffect(() => {
		let arr: any[] = [];
		for (let question of questionList()) {
			arr.push(
				<tr>
					<td>{question.question}</td>
					<td>{question.answer}</td>
					<td>{question.wrong0}</td>
					<td>{question.wrong1}</td>
					<td>{question.wrong2}</td>
					<td>
						<button
							onClick={() => deleteQuestion(question.ID)}
							class="button error"
						>
							X
						</button>
					</td>
				</tr>
			);
		}

		setTable(arr);
	});

	return (
		<table class="striped">
			<thead>
				<tr>
					<th>Question</th>
					<th>Answer</th>
					<th>Wrong 1</th>
					<th>Wrong 2</th>
					<th>Wrong 3</th>
				</tr>
			</thead>
			<tbody>{getTable()}</tbody>
		</table>
	);
};
