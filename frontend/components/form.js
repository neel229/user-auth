const Form = ({
  isMagic,
  errorMessage,
  login,
  setLogin,
  setUsername,
  setEmail,
  setPassword,
  onSubmit,
}) => {
  return (
    <form onSubmit={onSubmit}>
      <div>
        {!login ? (
          <div>
            <div className="selector">
              <button onClick={() => setLogin(!login)}>Login Instead?</button>
            </div>
            <label>
              <span>Username</span>
              <input
                type="text"
                name="username"
                required
                onChange={(e) => setUsername(e.target.value)}
              />
            </label>

            <label>
              <span>Email</span>
              <input
                type="email"
                name="email"
                required
                onChange={(e) => setEmail(e.target.value)}
              />
            </label>

            <label>
              <span>Password</span>
              <input
                type="password"
                name="password"
                required
                onChange={(e) => setPassword(e.target.value)}
              />
            </label>

            <div className="submit">
              <button type="submit">Sign Up</button>
            </div>

            {errorMessage && <p className="error">{errorMessage}</p>}
          </div>
        ) : (
          <div>
            <div className="selector">
              <button onClick={() => setLogin(!login)}>Sign Up</button>
            </div>
            <label>
              <span>Email</span>
              <input
                type="email"
                name="email"
                required
                onChange={(e) => setEmail(e.target.value)}
              />
            </label>

            <label>
              <span>Password</span>
              <input
                type="password"
                name="pasword"
                required
                onChange={(e) => setPassword(e.target.value)}
              />
            </label>

            <div className="submit">
              <button type="submit">Login</button>
            </div>

            {errorMessage && <p className="error">{errorMessage}</p>}
          </div>
        )}
      </div>

      <style jsx>{`
        form,
        label {
          display: flex;
          flex-flow: column;
        }
        label > span {
          font-weight: 600;
        }
        input {
          padding: 8px;
          margin: 0.3rem 0 1rem;
          border: 1px solid #ccc;
          border-radius: 4px;
        }
        .submit {
          display: flex;
          justify-content: flex-end;
          align-items: center;
          justify-content: space-between;
        }
        .submit > a {
          text-decoration: none;
        }
        .submit > button {
          padding: 0.5rem 1rem;
          cursor: pointer;
          background: #fff;
          border: 1px solid #ccc;
          border-radius: 4px;
        }
        .submit > button:hover {
          border-color: #888;
        }
        .error {
          color: brown;
          margin: 1rem 0 0;
        }
        .selector {
          margin-bottom: 5px;
          display: flex;
          flex-direction: row;
        }
        .selector > button {
          padding: 16px;
          cursor: pointer;
        }
      `}</style>
    </form>
  );
};

export default Form;
