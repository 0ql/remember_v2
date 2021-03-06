import "solid-js";
import { render } from "solid-js/web";
import "./index.css";
import App from "./App";
import * as serviceWorker from "./serviceWorker";

render(App, document.getElementById("root") as Node);

serviceWorker.unregister();
