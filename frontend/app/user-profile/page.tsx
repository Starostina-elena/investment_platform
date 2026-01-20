'use client';

import React, {Suspense, useEffect, useState} from "react";
import {GetUserById, User} from "@/api/user";
import {useSearchParams} from "next/navigation";
import styles from "@/app/user-profile/page.module.css"; // Исправлен путь импорта
import Spinner from "@/app/components/spinner";
import Link from "next/link";
import UserView from "@/app/components/user-view";
import {useUserStore} from "@/context/user-store";

export default function Page(){
    return (
        <Suspense>
            <PageUnwrapped />
        </Suspense>
    )
}

function PageUnwrapped(){
    const [viewedUser, setViewedUser] = useState<User | undefined | null>(undefined);
    const params = useSearchParams();

    // Получаем текущего залогиненного юзера из стора
    const { user: currentUser } = useUserStore();

    useEffect(() => {
        const idParam = params.get('id');

        // Если ID в URL нет, показываем профиль текущего пользователя
        if (!idParam) {
            if (currentUser) {
                setViewedUser(currentUser);
            } else {
                // Если нет ID и не залогинен - редирект на логин (обрабатывается в layout/starter, но можно и тут)
                setViewedUser(null);
            }
        } else {
            // Если ID есть, грузим данные
            GetUserById(+idParam).then(setViewedUser);
        }
    }, [params, currentUser]);

    // Проверка прав на редактирование
    // Можно редактировать, если это мой профиль (нет ID в url ИЛИ ID совпадает)
    const isOwner = currentUser && viewedUser
        ? currentUser.id === viewedUser.id
        : false;

    if (viewedUser === undefined)
        return (
            <div className={styles.container} style={{display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
                <Spinner/>
            </div>
        )

    if (viewedUser === null)
        return (
            <div className={styles.container} style={{padding: '5rem', textAlign: 'center', color: 'white'}}>
                <h1>Пользователь не найден или вы не авторизованы</h1>
                <br/>
                <Link href="/" style={{textDecoration: 'underline'}}>На главную</Link>
            </div>
        )

    return (
        <div className={styles.main}>
            {/* Передаем isOwner, чтобы компонент знал, показывать кнопку "Редактировать" или нет */}
            <UserView user={viewedUser} isOwner={isOwner} />
        </div>
    )
}