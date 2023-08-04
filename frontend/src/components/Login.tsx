import { SetStoreFunction } from "solid-js/store";
import { getCookie, type Credentials } from "../App";

const login = async (credentials: Credentials, setCredential: SetStoreFunction<Credentials>) => {
	const res = await fetch("http://localhost:5000/api/v1/login", {
		method: "POST",
		credentials: "include",
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(credentials),
	});
	if (res.status === 200) {
		console.log("Logged in")
		setCredential("session_id", getCookie("session_id"))
	}
};

export default (props: { credentials: Credentials, setCredential: SetStoreFunction<Credentials> }) => {
	return (
		<>
			<input
				placeholder="Email Eingeben"
				onChange={(e: any) => props.setCredential("email", e.target.value)}
			/>
			<input
				placeholder="Passwort Eingeben"
				onChange={(e: any) => props.setCredential("password", e.target.value)}
			/>
			<button class="button" onClick={() => login(props.credentials, props.setCredential)}>
				Einloggen
			</button>
		</>
	);
};
