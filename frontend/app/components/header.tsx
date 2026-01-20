'use client'
import Link from 'next/link';
import styles from './header.module.css';
import search from "@/public/search.svg"
import Image from "next/image";
import logo from "@/public/logo.png"
import {useRef, useState} from "react";
import burger_icon from "@/public/burger.svg";
import cross from "@/public/cross.svg";
import {addBasePath} from "next/dist/client/add-base-path";
import {useUserStore} from "@/context/user-store";
import avatar from "@/public/avatar.svg";
import {useRouter} from "next/navigation";

const LINKS = [
    {
        name: 'Главная',
        className: styles.main_link,
        link: '/'
    },
    {
        name: 'О нас',
        className: styles.about_link,
        link: '/#about'
    },
    {
        name: 'Как это работает',
        className: styles.how_it_works_link,
        link: '/#features'
    },
    {
        name: 'Каталог проектов',
        link: '/#projects'
    },
    {
        name: 'Рекордсмены',
        className: styles.records_link,
        link: '/#records'
    },
    {
        name: 'Контакты',
        link: '/#contacts'
    }
]
const Header = () => {
    const [filtersOpened, setFiltersOpened] = useState(false);
    const [burger, setBurger] = useState(false);
    const [showComboBox, setShowComboBox] = useState(false);
    const header = useRef<HTMLHeadingElement>(null);
    const user = useUserStore((state) => state.user);
    const router = useRouter();

    return (
        <>
            <header className={styles.header_container} ref={header} style={filtersOpened ? {position: 'fixed'} : {}}>
                <Link href="/" className={styles.logo_link}>
                    <Image src={logo} alt="Логотип" fill={true}/>
                </Link>
                <div className={styles.nav_group}>
                    {LINKS.map(link => (
                        <Link key={link.name} href={link.link} className={styles.nav_link + ' ' + link.className}>
                            {link.name}
                        </Link>
                    ))}
                    <div className={styles.search_container}>
                        <form
                            onSubmit={(e) => {
                                e.preventDefault();
                                const formData = new FormData(e.currentTarget);
                                const q = formData.get('q');
                                router.push(`/projects?q=${q}`);
                            }}
                            className={styles.search_form}
                        >
                            <input className={styles.search_input} type="text" name="q" placeholder="Поиск..." />
                            <button type="submit">
                                <Image className={styles.search_icon} src={search} alt="поиск"/>
                            </button>
                        </form>
                    </div>
                </div>

                <div className={styles.actions_group}>
                    <Link href="/create-project" className={styles.create_btn}>
                        Создать проект
                    </Link>
                    <Link href="/projects" className={styles.invest_btn}>
                        Инвестировать
                    </Link>
                    {!user && <Link href="/login" className={styles.login_link}>
                        Войти
                    </Link>}
                    {user && <Link href="/organisation/my" className="text-sm text-gray-300 hover:text-white mr-4">
                        Мои организации
                    </Link>}
                    {user && <Link href="/user-profile" className={styles.profile_link}>
                        Профиль{/*onClick={e => {
                                       e.preventDefault();
                                       setShowComboBox(!showComboBox);
                                   }}*/}
                        <Image src={user.avatar_path || avatar} alt='Профиль' fill={true}/>
                    </Link>}
                    {showComboBox && <div className={styles.combo_box}>
                        <button onClick={() => {
                            useUserStore.getState().Logout();
                            router.push(addBasePath('/login'));
                        }}>Выйти из аккаунта
                        </button>
                    </div>}
                </div>

                <div className={styles.burger_button} onClick={() => setBurger(!burger)}>
                    <Image src={burger_icon} alt='Бургер' fill={true}/>
                </div>
            </header>
            {filtersOpened && header.current &&
                <div style={{height: header.current.getBoundingClientRect().height}}></div>}
            <div className={styles.burger_menu + ' ' + (burger ? styles.burger_open : "")}
                 onClick={e => {
                     if (e.target === e.currentTarget) setBurger(false)
                 }}>
                <div className={styles.burger_links}>
                    <Image className={styles.burger_close} src={cross}
                           alt='Крестик' width={30} height={30} onClick={() => setBurger(false)}/>
                    <h3 className={styles.burger_title}>Навигация</h3>
                    <div className={styles.search_container}>
                        <form action={addBasePath("/search")} className={styles.search_form}>
                            <input
                                className={styles.search_input}
                                type="text"
                                name="q"
                            />
                            <button><Image className={styles.search_icon} src={search} alt="поиск"/></button>
                        </form>
                    </div>
                    {LINKS.map(link => (
                        <Link key={link.name} href={link.link} className={styles.burger_link}>
                            {link.name}
                        </Link>
                    ))}
                    {!user && <Link href="/login" className={styles.login_link}>
                        Войти
                    </Link>}
                    {user && <Link href="/user-profile" className={styles.profile_link}>
                        Профиль{/*onClick={e => {
                                       e.preventDefault();
                                       setShowComboBox(!showComboBox);
                                   }}*/}
                    </Link>}
                    <div className={styles.actions_group}>
                        <Link href="/create-project" className={styles.create_btn}>
                            Создать проект
                        </Link>
                        <Link href="/#" className={styles.invest_btn}>
                            Инвестировать
                        </Link>
                    </div>
                </div>
            </div>
        </>
    );
};

export default Header;