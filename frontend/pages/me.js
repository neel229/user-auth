import Router from "next/router";
import { useEffect, useState } from "react";
import Layout from "../components/layout";
import { me } from "./api/backend";

const Me = () => {
  const [authorized, setAuthorized] = useState(true);
  const [name, setName] = useState("");
  let token;
  useEffect(async () => {
    token = localStorage.getItem("jwt");
    if (token === null) {
      setAuthorized(false);
      return;
    }
    name = await me(token);
    if (name === "Unauthorized") {
      setAuthorized(false);
    }
    setName(name);
    setAuthorized(true);
  });

  if (!authorized) {
    localStorage.removeItem("jwt");
    Router.push("/");
  }

  return (
    <Layout>
      {!authorized ? (
        <div className="content">
          <div>Unauthorized</div>
        </div>
      ) : (
        <div className="content">
          <div>Hello {name}</div>
        </div>
      )}
      <style jsx>{`
        .content {
          position: absolute;
          left: 50%;
          top: 50%;
          transform: translate(-50%, -50%);
          font-size: 72px;
        }
      `}</style>
    </Layout>
  );
};

export default Me;
