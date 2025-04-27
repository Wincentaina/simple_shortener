import React, {useEffect, useState} from 'react';
// @ts-ignore
import s from "./ManageLinksPage.module.scss"
import axios from "axios";

type ResponseUrl = {
    "Alias": string,
    "Url": string
}

type Url = {
    "alias": string,
    "url": string
}

export const ManagePage = () => {
    const [urls, setUrls] = useState([])
    const [isLoaded, setIsLoaded] = useState(false)

    useEffect(() => {
            if (!isLoaded) {
                axios.get("http://localhost:8082/url/users_all", {
                    headers: {
                        Authorization: `Bearer ${window.localStorage.getItem("token")}`
                    }
                })
                    // .then((res) => {
                    //     let respUrls = res.data.urls
                    //
                    //     // @ts-ignore
                    //     setUrls(() => {
                    //         respUrls.map((item: ResponseUrl): Url => {
                    //             return ({
                    //                 "alias": item.Alias,
                    //                 "url": item.Url
                    //             })
                    //         })
                    //     })
                    // })
                    .then(res => {
                        let results = res.data.urls
                        // @ts-ignore
                        setUrls(results)
                    })
                setIsLoaded(true)
            }
        }
        , [isLoaded])

    return (
        <div className={s.mange_container}>
            <h3>Ваши ссылки:</h3>
            <div className={s.table}>

            </div>
        </div>
    );
};