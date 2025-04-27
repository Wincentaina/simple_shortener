import React, {useState} from 'react';
// @ts-ignore
import s from "./ShortLinkPage.module.scss"
import {LinksForm} from "../../components/LinksForm/LinksForm";

export const ShortLinkPage: React.FC = () => {
    return (
        <div className={s.short_page}>
            <LinksForm />
        </div>
    );
};
