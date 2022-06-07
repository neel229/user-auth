import axios from "axios";

const authURL = process.env.NEXT_PUBLIC_AUTH_URL;
const meURL = process.env.NEXT_PUBLIC_ME_URL;

export const register = async (body) => {
  try {
    const res = await axios.post(
      authURL,
      {
        username: body.username,
        email: body.email,
        password: body.password,
      },
      {
        params: {
          register: "true",
        },
      }
    );
    localStorage.setItem("jwt", res.data);
  } catch (err) {
    console.error(err);
    return err.message;
  }
};

export const login = async (body) => {
  try {
    const res = await axios.post(
      authURL,
      {
        email: body.email,
        password: body.password,
      },
      {
        params: {
          login: "true",
        },
      }
    );
    localStorage.setItem("jwt", res.data);
  } catch (err) {
    console.error(err);
    return err.message;
  }
};

export const me = async (token) => {
  try {
    const res = await axios.get(meURL, {
      headers: {
        authorization: token,
      },
    });
    return res.data;
  } catch (err) {
    return "Unauthorized";
  }
};
