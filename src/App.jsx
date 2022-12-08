import React, { useEffect, useState } from "react";
import { Outlet, Route, Routes } from "react-router";
import Home from "./components/Home";
import Restricted from "./components/Restricted";
import Post from "./components/Post";
import NoMatch from "./components/NoMatch.jsx";
import List from "./components/List.jsx";
import { Link } from "react-router-dom";
import { useSelector, useDispatch } from "react-redux";
import { setLocale, setPassword } from "./store";

function App() {
	const dispatch = useDispatch()

	useEffect(() => {
		// load password from local storage
		const password = localStorage.getItem("password");
		// set to redux
		// if not empty
		if (password) {
			dispatch(setPassword(password));
		}

		// load locale from local storage
		const locale = localStorage.getItem("locale");
		// set to redux
		if (locale) {
			dispatch(setLocale(locale));
		}
	})
  return (
    <div>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />}></Route>
          <Route path="/restricted/list" element={<List />}></Route>
          <Route path="/restricted/:hash/:id" element={<Restricted />}></Route>
          <Route path="/posts/list" element={<List />}></Route>
          <Route path="/posts/:folder/:id" element={<Post />}></Route>
          <Route path="*" element={<NoMatch />}></Route>
        </Route>
      </Routes>
    </div>
  );
}

function Layout() {
	const [locales, setLocales] = useState([])
  const dispatch = useDispatch();

  useEffect(() => {
    const prefix = import.meta.env.VITE_BLOG_PREFIX;
    // fetch README.md
    fetch(prefix + "locales.json").then((response) => {
      // get the body as text
      response.json().then((json) => {
		// sort, make zh-cn first
		const locales = json.sort((a, b) => {
			if (a === "zh-cn") {
				return -1;
			} else if (b === "zh-cn") {
				return 1;
			} else {
				return 0;
			}
		});

		setLocales(json)
      });
    });

  }, []);

  return (
    <div>
      <nav>
        <ul>
          <li>
            <Link to="/">Home</Link>
          </li>
          <li>
            <Link to="/posts/list">Posts</Link>
          </li>
          <li>
            <Link to="/restricted/list">Restricted</Link>
          </li>
          <li>
            <select 
				defaultValue={locales.length == 0 ? "": locales[0]}
				onChange={(e) => {
					// set locale
					dispatch(setLocale(e.target.value));
					// save locale to storage
					localStorage.setItem("locale", e.target.value);
				}}
				>
				{locales.map((locale) => (
					<option key={locale} value={locale}>
						{locale}
					</option>
				))}
            </select>
          </li>
        </ul>
      </nav>
      <Outlet></Outlet>
    </div>
  );
}

export default App;
