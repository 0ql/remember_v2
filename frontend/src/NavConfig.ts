import Dashboard from "./pages/dashboard"
import Competitive from "./pages/competitive"
import CreateQuiz from "./pages/createQuiz"
import questionList from "./pages/questionList"
import user from "./pages/user"

export default [
	{
		title: "Dashboard",
		pageComponent: Dashboard
	},
	{
		title: "Competitive",
		pageComponent: Competitive
	},
	{
		title: "Create",
		pageComponent: CreateQuiz
	},
	{
		title: "Fragenliste",
		pageComponent: questionList
	},
	{
		title: "User",
		pageComponent: user
	}
]