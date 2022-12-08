import { useEffect, useState } from "react";
import { useLocation, useParams } from "react-router";
import CryptoJS from "crypto-js";
import { marked } from "marked";
import { useSelector, useDispatch } from "react-redux";
import { setPassword } from "../store";

export default function Post() {
    const url = useLocation()
    const dispatch = useDispatch()
    const locale = useSelector(state => state.locale.value)

    const { folder, id } = useParams()

    useEffect(() => {
        fetch(import.meta.env.VITE_BLOG_PREFIX + "posts/" + folder + "/" + id + ".json").then((response) => {
            response.json().then((json) => {
                var html = marked.parse(json.Content)
                document.getElementById("content").innerHTML = html
            })
        })
    }, [url, locale])

    return  (
        <div>
            <h1>Restricted</h1>
            <input type="password" id="password" />
            <input type="submit" onClick={(e) => {
                // get value from password
                const passwd = document.getElementById("password").value

                dispatch(setPassword(passwd))

                // store in local storage
                localStorage.setItem("password", passwd)
            }} />
            <div id="content" class="content">

            </div>
        </div>
    );
}
