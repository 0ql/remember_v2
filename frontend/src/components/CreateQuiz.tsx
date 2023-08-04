import { createStore } from "solid-js/store";

const [state, setState] = createStore({
	Question: "",
	Answer: "",
	Wrong0: "",
	Wrong1: "",
	Wrong2: "",
});

const sendQuestion = () => {
	console.log(JSON.stringify(state));
	fetch("http://localhost:5000/api/v1/createQuestion", {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		credentials: "include",
		body: JSON.stringify(state),
	});
};

export default () => {
	return (
		<>
			<h2>Create Question</h2>
			<hr></hr>
			<fieldset>
				<legend>CREATE QUESTION</legend>
				<p>
					<label>Question</label>
					<input
						type="text"
						placeholder="Enter Question"
						onChange={(e: any) => setState("Question", e.target.value)}
					></input>
				</p>
				<p>
					<label>Answer</label>
					<input
						type="text"
						placeholder="Enter Answer"
						onChange={(e: any) => setState("Answer", e.target.value)}
					></input>
				</p>
				<p>
					<label>1. Wrong Answer</label>
					<input
						type="text"
						placeholder="Wrong Answer"
						onChange={(e: any) => setState("Wrong0", e.target.value)}
					></input>
				</p>
				<p>
					<label>2. Wrong Answer</label>
					<input
						type="text"
						placeholder="Wrong Answer"
						onChange={(e: any) => setState("Wrong1", e.target.value)}
					></input>
				</p>
				<p>
					<label>3. Wrong Answer</label>
					<input
						type="text"
						placeholder="Wrong Answer"
						onChange={(e: any) => setState("Wrong2", e.target.value)}
					></input>
				</p>
				<button onClick={sendQuestion} class="button primary">
					Create Question
				</button>
			</fieldset>
		</>
	);
};
