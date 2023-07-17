import { createContext } from "solid-js";
import "./App.css";
import { A, Route, Router, Routes, useNavigate } from "@solidjs/router";
import Login from "./components/Login";
import Dashboard from "./components/Dashboard";
import Competitive from "./components/Competitive";
import CreateQuiz from "./components/CreateQuiz";
import QuestionList from "./components/QuestionList";
import User from "./components/User";
import Game from "./components/Game";
import { GlobalContextProvider, createWebSocket } from "./socket";

export const globalContext = createContext({ isLoggedIn: false });

const getCookie = (name: string) => {
	var dc, prefix, begin, end;

	dc = document.cookie;
	prefix = name + "=";
	begin = dc.indexOf("; " + prefix);
	end = dc.length;

	if (begin !== -1) {
		begin += 2;
	} else {
		begin = dc.indexOf(prefix);
		if (begin === -1 || begin !== 0) return null;
	}
	if (dc.indexOf(";", begin) !== -1) {
		end = dc.indexOf(";", begin);
	}

	return decodeURI(dc.substring(begin + prefix.length, end)).replace(/\"/g, "");
};

const Component = () => {
	const navigate = useNavigate()
	createWebSocket(navigate)

	return (
		<Routes>
			<Route path="/dashboard" component={Dashboard} />
			<Route path="/competitive" component={Competitive} />
			<Route path="/create" component={CreateQuiz} />
			<Route path="/questions" component={QuestionList} />
			<Route path="/user" component={User} />
			<Route path="/game" component={Game} />
		</Routes>
	)
}

export default () => {
	if (getCookie("session_id") === null) {
		return <Login></Login>
	}

	return (
		<GlobalContextProvider>
			<Router>
				<nav
					class="tabs is-full"
					style="background-color: var(--bg-color);position: sticky; top: 0;"
				>
					<A href="/dashboard">Dashboard</A>
					<A href="/competitive">Competitive</A>
					<A href="/create">Create</A>
					<A href="/questions">Questions</A>
					<A href="/user">User</A>
				</nav>
				<Component />
			</Router>
		</GlobalContextProvider>
	);
};
