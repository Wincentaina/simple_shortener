import React, {useState} from 'react';
// @ts-ignore
import s from "./LinksForm.module.scss"
import axios from "axios";

export const LinksForm: React.FC = () => {
    const [linkToSave, setLinkToSave] = useState("")
    const [alias, setAlias] = useState("")

    const Login = () => {
        axios.post("http://localhost:8082/url/save", {
            "url": linkToSave,
            "alias": alias
        }, {
            headers: {
                Authorization: `Bearer ${window.localStorage.getItem("token")}`
            }
        })
            .then((data) => {
                console.log(data.data)
                setAlias("")
                setLinkToSave("")
                alert(`Создан псевдоним ${data.data.alias}`)
            })
    }

    return (
            <div className={s.form_container}>
                <h3 className={s.title}>Сократить ссылку</h3>
                <div className={s.short_form}>
                    <input className={s.inp_field} value={linkToSave} type="text" placeholder="ссылка формата: http://link.com" onChange={e => setLinkToSave(e.target.value)}/>
                    <input className={s.inp_field} value={alias} type="text" placeholder="псевдоним" onChange={e => setAlias(e.target.value)}/>
                    <button onClick={Login} className={s.short_btn}>Сократить!</button>
                </div>
            </div>
    );
};