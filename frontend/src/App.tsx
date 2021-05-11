import { Component, createContext, createState } from "solid-js";
import "./App.css";
import NavConfig from "./NavConfig";
import Login from "./pages/login";

const [state, setState] = createState({
  activeTab: 2,
});

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

const Tabs: Component = () => {
  let tabs = [];
  for (let i = 0; i < NavConfig.length; i++) {
    tabs.push(
      <a
        onClick={() => setState("activeTab", (val) => (val = i))}
        class={state.activeTab === i ? "active" : ""}
      >
        {NavConfig[i].title}
      </a>
    );
  }
  return tabs;
};

export default () => {
  if (getCookie("session_id") === null) {
    return <Login></Login>;
  }

  return (
    <>
      <nav
        class="tabs is-full"
        style="background-color: var(--bg-color);position: sticky; top: 0;"
      >
        <Tabs />
      </nav>
      <main class="container">{NavConfig[state.activeTab].pageComponent}</main>
    </>
  );
};
