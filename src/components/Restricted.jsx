import { useEffect, useState } from "react";
import { useLocation, useParams } from "react-router";
import CryptoJS from "crypto-js";
import { marked } from "marked";
import { useSelector, useDispatch } from "react-redux";
import { setPassword } from "../store";

const decrypt = function (rawKey, encryptedData) {
    // split the key and iv by :
    const splited = encryptedData.split(':');
    const iv = splited[0]
    const encrypted = splited[1]

    // iv to bytes
    var ivBytes = CryptoJS.enc.Base64.parse(iv);

    var key = CryptoJS.enc.Utf8.parse(rawKey);

    // use cryptojs to decrypt data
    const decryptedData = CryptoJS.AES.decrypt(encrypted, key, {
        iv: ivBytes,
        mode: CryptoJS.mode.CBC,
        padding: CryptoJS.pad.Pkcs7
    });

    // return decrypted data
    return decryptedData.toString(CryptoJS.enc.Utf8);
}

export default function Restricted() {
    const url = useLocation()
    const dispatch = useDispatch()
    const password = useSelector(state => state.password.value)
    const locale = useSelector(state => state.locale.value)

    const { hash, id } = useParams()

    useEffect(() => {
        fetch(import.meta.env.VITE_BLOG_PREFIX + "restricted/" + hash + "/" + id + ".json").then((response) => {
            response.text().then((text) => {
                var decrypted = decrypt(password, text)
                var json = JSON.parse(decrypted)
                var html = marked.parse(json.Content)
                document.getElementById("content").innerHTML = html
            })
        })
    }, [url, password, locale])

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
