import React, {useState} from 'react';
// @ts-ignore
import s from "./AuthForm.module.scss"
import axios from "axios";

export const AuthForm: React.FC = () => {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")

    const Login = () => {
        if (username.length != 0 && password.length != 0) {
            axios({
                method: 'post',
                url: 'http://localhost:8082/auth/login',
                data: {
                    "username": username,
                    "password": password
                }
            })
                .then((data) => {
                    window.localStorage.setItem('token', data.data.token)
                    setPassword("")
                    setUsername("")
                    alert("Success!")
                })
        }
    }

    return (
        <div className={s.form_container}>
            <h3 className={s.title}>Войти</h3>
            <div className={s.auth_form}>
                <input className={s.inp_field} value={username} type="text" onChange={e => setUsername(e.target.value)}/>
                <input className={s.inp_field} value={password} type="text" onChange={e => setPassword(e.target.value)}/>
                <button onClick={Login} className={s.login_btn}>Войти</button>
            </div>
        </div>
    );
};