import {
  createEffect,
  createResource,
  createState,
  For,
  Suspense,
} from "solid-js";

type UserData = {
  username: string;
  email: string;
  rank: string;
  elo: number;
};

const fetchMyInformation = async () => {
  const res = await fetch("http://localhost:5000/api/v1/getMyInformation", {
    credentials: "include",
  });
  return res.json();
};

type User = {
  email: string;
  username: string;
  rank: string;
  elo: number;
};

const getAllUsers = async () => {
  const res = await fetch("http://localhost:5000/api/v1/getAllUsers", {
    credentials: "include",
  });
  if (res.status === 200) {
    return res.json();
  }
  return new Error(res.status + "|" + res.statusText);
};

const ListUsers = () => {
  const [allUsers] = createResource(getAllUsers);

  return (
    <>
      <table class="striped">
        <thead>
          <tr>
            <th>Email</th>
            <th>Benutzername</th>
            <th>Rang</th>
            <th>Elo</th>
          </tr>
        </thead>
        <For each={allUsers()}>
          {(user: User) => (
            <tr>
              <td>{user.email}</td>
              <td>{user.username}</td>
              <td>{user.rank}</td>
              <td>{user.elo}</td>
            </tr>
          )}
        </For>
      </table>
    </>
  );
};

const CreateUser = () => {
  const [newUser, setNewUser] = createState({
    username: "",
    email: "",
    password: "",
  });

  const createUser = async () => {
    const res = await fetch("http://localhost:5000/api/v1/createUser", {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(newUser),
    });
    if (res.status === 200) {
      console.log("User Created Successfully");
    }
  };

  return (
    <>
      <h4>Benutzer Erstellen</h4>
      <hr />
      <fieldset>
        <legend>BENUTZER ERSTELLEN</legend>
        <p>
          <label>Benutzername</label>
          <input
            onChange={(e: any) => setNewUser("username", e.target.value)}
            type="text"
            placeholder="Benutzername"
          />
        </p>

        <p>
          <label>Email Addresse</label>
          <input
            onChange={(e: any) => setNewUser("email", e.target.value)}
            placeholder="Email"
            type="email"
          />
        </p>

        <label>Elo Zahl</label>
        <p class="grouped">
          <input type="number" placeholder="Elo Zahl" value="100" />
          <input
            onChange={(e: any) => setNewUser("password", e.target.value)}
            type="password"
            placeholder="Passwort"
          ></input>
          <button onClick={createUser}>Erstellen</button>
        </p>
      </fieldset>
    </>
  );
};

const AdminComponent = (props: any) => {
  return (
    <>
      <h2>Admin</h2>
      <hr></hr>
      <ListUsers></ListUsers>
      <CreateUser></CreateUser>
    </>
  );
};

export default () => {
  // TODO: set to global context
  const [getInfo] = createResource(fetchMyInformation, {
    initialValue: {
      username: "",
      elo: 100,
      rank: "",
      email: "",
    },
  });

  const checkAdmin = () => {
    if (getInfo().rank === "admin") {
      return <AdminComponent></AdminComponent>;
    }
    return;
  };

  return (
    <>
      <h2>Mein Account</h2>
      <hr></hr>
      <h4>Benutzername</h4>
      <div>{getInfo().username}</div>
      <h4>Email Addresse</h4>
      <div>{getInfo().email}</div>
      <h4>Rang</h4>
      <div>{getInfo().rank}</div>
      <h4>Elo Rating</h4>
      <div>{getInfo().elo}</div>
      {checkAdmin()}
    </>
  );
};
