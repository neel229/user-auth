import { useState } from "react";
import Router from "next/router";
import { useUser } from "../lib/hooks";
import Layout from "../components/layout";
import Form from "../components/form";

import { Magic } from "magic-sdk";
import { register, login as loginFn } from "./api/backend";

const Login = () => {
  useUser({ redirectTo: "/", redirectIfFound: true });

  const [errorMsg, setErrorMsg] = useState("");
  const [isMagic, setIsMagic] = useState(false);
  const [login, setLogin] = useState(false);
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();

    if (errorMsg) setErrorMsg("");

    if (isMagic) {
      const body = {
        email: e.currentTarget.email.value,
      };

      try {
        const magic = new Magic(process.env.NEXT_PUBLIC_MAGIC_PUBLISHABLE_KEY);
        const didToken = await magic.auth.loginWithMagicLink({
          email: body.email,
        });
        const res = await fetch("/api/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: "Bearer " + didToken,
          },
          body: JSON.stringify(body),
        });
        if (res.status === 200) {
          Router.push("/");
        } else {
          throw new Error(await res.text());
        }
      } catch (error) {
        console.error("An unexpected error happened occurred:", error);
        setErrorMsg(error.message);
      }
    }

    try {
      if (login) {
        const body = {
          email: email,
          password: password,
        };
        await loginFn(body);
        return;
      }
      const body = {
        username: username,
        email: email,
        password: password,
      };
      await register(body);
    } catch (err) {
      console.error(err);
    }
  }

  return (
    <Layout>
      <div className="login">
        <Form
          isMagic={isMagic}
          errorMessage={errorMsg}
          login={login}
          setLogin={setLogin}
          setUsername={setUsername}
          setEmail={setEmail}
          setPassword={setPassword}
          onSubmit={handleSubmit}
        />
        <button onClick={() => setIsMagic(!isMagic)}>
          {isMagic ? <span>Use Email</span> : <span>Use Magic Link</span>}
        </button>
      </div>
      <style jsx>{`
        .login {
          max-width: 21rem;
          margin: 0 auto;
          padding: 1rem;
          border: 1px solid #ccc;
          border-radius: 4px;
        }
        button {
          margin: 8px 0px;
          padding: 0.5rem;
        }
      `}</style>
    </Layout>
  );
};

export default Login;
