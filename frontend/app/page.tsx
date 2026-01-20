import styles from "./page.module.css";
import Link from "next/link";
import Image from "next/image";
import coins from "@/public/coins.svg";
import lamp from "@/public/idea_lamp.svg";
import docs from "@/public/docs.svg";
import about_us from "@/public/about_us.jpg";
import ProjectsSection from "@/app/projects-section";
import avatar from "@/public/temp/avatar.png";
import React from "react";
import Footer from "@/app/components/footer";

export default function Home() {
    return (
        <>
            <main className={styles.main}>
                <div className={styles.main_dark}>
                    <h1 className={styles.page_title}>
                        Краудфандинговая
                        и инвестиционная платформа
                    </h1>
                    <p className={styles.main_description}>Мы объединяем инвесторов и предпринимателей,
                        чтобы помочь реализовать ваши идеи.
                        Начните свой путь к успеху!
                    </p>
                    <div className={styles.main_links}>
                        <Link href="/create-project">Создать проект</Link>
                        <Link href="#">Инвестировать</Link>
                    </div>
                </div>
            </main>
            <section className={styles.section} id="features">
                <h2 className={styles.section_title}>Инвестируйте, создавайте
                    и просматривайте проекты</h2>
                <p className={styles.section_description}>Наша платформа предлагает уникальные возможности для
                    инвестирования и создания проектов.
                    Присоединяйтесь к нам, чтобы поддержать малый бизнес и стартапы.</p>
                <div className={styles.section_cards}>
                    <div>
                        <Image src={coins} alt="Стопка монет"/>
                        <h3>
                            Инвестирование
                            в перспективные проекты
                        </h3>
                        <p>Получите доступ к разнообразным инвестиционным возможностям.</p>
                    </div>
                    <div>
                        <Image src={lamp} alt="Лампочка"/>
                        <h3>
                            Создание и развитие собственных проектов
                        </h3>
                        <p>Запустите свой проект и привлеките финансирование.</p>
                    </div>
                    <div>
                        <Image src={docs} alt="Документ"/>
                        <h3>Просмотр каталога успешных проектов</h3>
                        <p>Изучите успешные проекты, станьте их частью
                            или вдохновитесь на создание своего.</p>
                    </div>
                </div>
                <div className={styles.section_links}>
                    <Link href="/create-project">Создать проект</Link>
                    <Link href="#">Инвестировать
                        <svg xmlns="http://www.w3.org/2000/svg" width="8" height="12" viewBox="0 0 8 12" fill="none">
                            <path
                                d="M1.70697 11.9496L7.41397 6.24264L1.70697 0.535645L0.292969 1.94964L4.58597 6.24264L0.292969 10.5356L1.70697 11.9496Z"
                                fill="#B7FF00"/>
                        </svg>
                    </Link>
                </div>
            </section>
            <section className={styles.section + ' ' + styles.section_gradient_about} id="about">
                <h2 className={styles.section_title}>О нас</h2>
                <p className={styles.section_description}>
                    Мы платформа, которая специализируется на объединении инвесторов и малых предпринимателей, которым
                    необходимы инвестиции на разработку и развитие проектов.
                </p>
                <div className={styles.about_image}>
                    <Image src={about_us} alt='Офис сотрудников' fill={true}/>
                </div>
            </section>
            <ProjectsSection/>
            <section className={styles.section + ' ' + styles.section_choose_way}>
                <h2 className={styles.section_title}>выберите свой путь</h2>
                <div className={styles.section_cards}>
                    <Link href="/projects" className={styles.small_invest}>
                        <p>инвестировать / пожертвовать</p>
                    </Link>
                    {/*<Link href="#" className={styles.big_invest}>*/}
                    {/*    <p>большое инвестирование</p>*/}
                    {/*</Link>*/}
                    <Link href="/create-project" className={styles.start_project}>
                        <p>создать проекта</p>
                    </Link>
                </div>
            </section>
            <section className={styles.section + ' ' + styles.section_benefits}>
                <div className={styles.benefits_block}>
                    <div className={styles.benefit}>
                        <h3 className={styles.benefit_title}>благотворительность</h3>
                        <p className={styles.benefit_description}>Благотворительность - это ваша возможность поучаствовать в создании новой мечты, реализация которой может дать начало великим делам</p>
                    </div>
                    <Link href='#' className={styles.benefit_link}>Пожертвовать</Link>
                </div>
            </section>
            <section className={styles.section} id="reviews">
                <h2 className={styles.section_title}>что говорят наши клиенты</h2>
                <div className={styles.review}>
                    <p className={styles.review_text}>"Эта платформа помогла мне собрать необходимые средства для моего стартапа. Я благодарен за поддержку и профессионализм команды!"</p>
                    <div className={styles.review_avatar}>
                        <Image src={avatar} alt="Фото автора"/>
                    </div>
                    <p className={styles.review_author}>Иван Петров</p>
                    <p className={styles.review_position}>Основатель, Стартап XYZ</p>
                </div>
            </section>
            <section className={styles.section + ' ' + styles.section_start_today}>
                <div>
                    <h2 className={styles.section_title}>Начните свой проект сегодня</h2>
                    <p className={styles.section_description}>Присоединяйтесь к нам и инвестируйте в будущее</p>
                </div>
                <div className={styles.main_links}>
                    <Link href="/create-project">Создать проект</Link>
                    <Link href="#">Инвестировать</Link>
                </div>
            </section>
            <Footer/>
        </>
    );
}
