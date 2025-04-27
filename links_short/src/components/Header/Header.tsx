import React from 'react';
// @ts-ignore
import s from "./Header.module.scss"

export const Header = () => {
    return (
        <header>
            <div className={s.header_container}>
                <div className={s.menu_buttons}>
                    <a href="/" className={s.link}>Сократить ссылку</a>
                    <a href="/links" className={s.link}>Просмотреть мои ссылки</a>
                </div>
                <div className={s.login_area}>
                    <a href="/login"  className={s.link}>Войти</a>
                </div>
            </div>
        </header>
    );
};