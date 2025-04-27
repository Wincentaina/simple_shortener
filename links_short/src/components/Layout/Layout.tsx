import React from 'react';
import {Header} from "../Header/Header";
import {Footer} from "../Footer/Footer";

// @ts-ignore
import s from "./Layout.module.scss"

// TODO: replace any to react type
export const Layout: React.FC = ({children}: any) => {
    return (
        <div className={s.container}>
            <Header />
            <main className={s.main_container}>
                {children}
            </main>
            <div className={s.footer_container}>
                <Footer />
            </div>
        </div>
    );
};
