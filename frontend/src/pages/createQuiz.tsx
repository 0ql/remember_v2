import { createState } from "solid-js";

const [state, setState] = createState({
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
    body: JSON.stringify(state),
  });
};

export default () => {
  return (
    <>
      <h2>Frage Erstellen</h2>
      <hr></hr>
      <fieldset>
        <legend>FRAGE ERSTELLEN</legend>
        <p>
          <label>Frage</label>
          <input
            type="text"
            placeholder="Frage Eingeben"
            onChange={(e: any) => setState("Question", e.target.value)}
          ></input>
        </p>
        <p>
          <label>Antwort</label>
          <input
            type="text"
            placeholder="Antwort Eingeben"
            onChange={(e: any) => setState("Answer", e.target.value)}
          ></input>
        </p>
        <p>
          <label>1. Falsche Antwort</label>
          <input
            type="text"
            placeholder="Falsche Antwort"
            onChange={(e: any) => setState("Wrong0", e.target.value)}
          ></input>
        </p>
        <p>
          <label>2. Falsche Antwort</label>
          <input
            type="text"
            placeholder="Falsche Antwort"
            onChange={(e: any) => setState("Wrong1", e.target.value)}
          ></input>
        </p>
        <p>
          <label>3. Falsche Antwort</label>
          <input
            type="text"
            placeholder="Falsche Antwort"
            onChange={(e: any) => setState("Wrong2", e.target.value)}
          ></input>
        </p>
        <button onClick={sendQuestion} class="button primary">
          Frage Erstellen
        </button>
      </fieldset>
    </>
  );
};
