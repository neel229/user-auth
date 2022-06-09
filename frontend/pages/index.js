import { useUser } from "../lib/hooks";
import Layout from "../components/layout";
import Login from "./login";
import { useEffect, useState } from "react";
import Router from "next/router";

const Home = () => {
  const user = useUser();
  const [authorized, setAuthorized] = useState();

  useEffect(() => {
    let exists = localStorage.getItem("jwt");
    if (exists) {
      setAuthorized(true);
    }
  }, [authorized]);

  if (authorized && !user) {
    Router.push("/me");
  }

  return (
    <Layout>
      {!user ? (
        <Login setAuthorized={setAuthorized}></Login>
      ) : (
        <>
          <p>Currently logged in as:</p>
          <pre>{JSON.stringify(user, null, 2)}</pre>
          <div>
            <button style={{ padding: "16px", fontWeight: "bold" }}>
              Logout
            </button>
          </div>
        </>
      )}

      <style jsx>{`
        li {
          margin-bottom: 0.5rem;
        }
        pre {
          white-space: pre-wrap;
          word-wrap: break-word;
        }
      `}</style>
    </Layout>
  );
};

export default Home;
