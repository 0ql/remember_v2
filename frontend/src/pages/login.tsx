import { createState } from "solid-js";

const [credentials, setCredential] = createState({
  email: "",
  password: "",
});

const login = async () => {
  const res = await fetch("http://localhost:5000/api/v1/login", {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });
  if (res.status === 200) {
    console.log("Successfully Logged In");
  }
};

export default () => {
  return (
    <>
      <input
        placeholder="Email Eingeben"
        onChange={(e: any) => setCredential("email", e.target.value)}
      />
      <input
        placeholder="Passwort Eingeben"
        onChange={(e: any) => setCredential("password", e.target.value)}
      />
      <button class="button" onClick={login}>
        Einloggen
      </button>
    </>
  );
};
