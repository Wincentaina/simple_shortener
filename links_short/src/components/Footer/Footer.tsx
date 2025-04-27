import React from 'react';
// @ts-ignore
import s from "./Footer.module.scss"

export const Footer = () => {
    return (
        <footer>
            <div className={s.footer_container}>
                <div className={s.link_on_networks}>
                    <a className={s.link} href="github.com">GitHub</a>
                </div>
                <div className={s.author}>
                    <p>Created by ...</p>
                </div>
            </div>
        </footer>
    );
};