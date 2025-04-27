import React, {useState} from 'react';
// @ts-ignore
import s from "./AuthPage.module.scss"
import {AuthForm} from "../../components/AuthForm/AuthForm";

export const AuthPage: React.FC = () => {


    return (
        <div className={s.auth_page_container}>
            <AuthForm />
        </div>
    );
};
