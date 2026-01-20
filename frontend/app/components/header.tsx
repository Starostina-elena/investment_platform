'use client'
import Link from 'next/link';
import styles from './header.module.css';
import search from "@/public/search.svg"
import Image from "next/image";
import logo from "@/public/logo.svg"
import {useEffect, useRef, useState} from "react";
import cross from "@/public/cross.svg";
import {addBasePath} from "next/dist/client/add-base-path";
import {useUserStore} from "@/context/user-store";
import avatar from "@/public/avatar.svg";
import {useRouter} from "next/navigation";
import { useImage } from "@/hooks/use-image"; // Импортируем хук
import { BUCKETS } from "@/lib/config";

const LINKS = [
    { name: 'Главная', className: styles.main_link, link: '/' },
    { name: 'О нас', className: styles.about_link, link: '/#about' },
    { name: 'Как это работает', className: styles.how_it_works_link, link: '/#features' },
    { name: 'Каталог проектов', link: '/projects' },
    { name: 'Контакты', link: '/#contacts' }
]

const Header = () => {
    const [burger, setBurger] = useState(false);
    const [showComboBox, setShowComboBox] = useState(false);
    const user = useUserStore((state) => state.user);
    const router = useRouter();

    // Загрузка аватарки через хук (чтобы работало с MinIO)
    const avatarSrc = useImage(null, user?.avatar_path, BUCKETS.AVATARS, avatar);

    // Блокировка скролла при открытом меню
    useEffect(() => {
        if (burger) document.body.style.overflow = 'hidden';
        else document.body.style.overflow = 'unset';
    }, [burger]);

    return (
        <>
            <header className={styles.header_container}>
                <Link href="/" className={styles.logo_link}>
                    <Image src={logo} alt="Логотип" fill className="object-contain"/>
                </Link>

                {/* Навигация (Десктоп) */}
                <div className={styles.nav_group}>
                    {LINKS.map(link => (
                        <Link key={link.name} href={link.link} className={styles.nav_link}>
                            {link.name}
                        </Link>
                    ))}
                </div>

                {/* Поиск (Десктоп) */}
                <div className={styles.search_container}>
                    <form
                        onSubmit={(e) => {
                            e.preventDefault();
                            const formData = new FormData(e.currentTarget);
                            router.push(`/projects?q=${formData.get('q')}`);
                        }}
                        className={styles.search_form}
                    >
                        <input className={styles.search_input} type="text" name="q" placeholder="Поиск..." />
                        <Image className={styles.search_icon} src={search} alt="поиск" width={18} height={18}/>
                    </form>
                </div>

                {/* Кнопки справа (Десктоп) */}
                <div className={styles.actions_group}>
                    <Link href="/create-project" className={styles.create_btn}>Создать проект</Link>
                    <Link href="/projects" className={styles.invest_btn}>Инвестировать</Link>

                    {!user ? (
                        <Link href="/login" className={styles.login_link}>Войти</Link>
                    ) : (
                        <div className="relative">
                            <button
                                onClick={() => router.push('/user-profile')}
                                className={styles.profile_link}
                            >
                                <Image src={avatarSrc} alt='Профиль' fill className="object-cover"/>
                            </button>
                        </div>
                    )}
                </div>

                {/* Бургер кнопка (Мобильные) */}
                <div className={styles.burger_button} onClick={() => setBurger(!burger)}>
                    <span></span>
                    <span></span>
                    <span></span>
                </div>
            </header>

            {/* Мобильное меню */}
            <div className={`${styles.burger_menu} ${burger ? styles.burger_open : ''}`}>
                <div className={styles.burger_links}>
                    {/* Кнопка закрытия */}
                    <div className="flex justify-end w-full mb-4">
                        <Image src={cross} alt='Закрыть' width={30} height={30} onClick={() => setBurger(false)} className="cursor-pointer invert"/>
                    </div>

                    {LINKS.map(link => (
                        <Link key={link.name} href={link.link} onClick={() => setBurger(false)}>
                            {link.name}
                        </Link>
                    ))}

                    <hr className="border-gray-700 my-2"/>

                    {!user ? (
                        <Link href="/login" onClick={() => setBurger(false)}>Войти</Link>
                    ) : (
                        <>
                            <Link href="/user-profile" onClick={() => setBurger(false)}>Мой профиль</Link>
                            <Link href="/organisation/my" onClick={() => setBurger(false)}>Мои организации</Link>
                        </>
                    )}

                    <Link href="/create-project" className="text-[#825e9c]" onClick={() => setBurger(false)}>Создать проект</Link>
                </div>
            </div>
        </>
    );
};

export default Header;