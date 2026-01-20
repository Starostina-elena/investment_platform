import styles from "./page.module.css";
import Link from "next/link";
import Footer from "@/app/components/footer";
import ProjectsSection from "@/app/projects-section";
import { Coins, Lightbulb, FileText } from "lucide-react"; // Используем векторные иконки

export default function Home() {
    return (
        <>
            <main className={styles.main}>
                <div className={styles.main_dark}>
                    <h1 className={styles.page_title}>
                        Система инвестиционных<br/>проектов и сборов
                    </h1>
                    <p className={styles.main_description}>
                        Мы — самый крутой сайт для сборов и венчурных инвестиций.
                        Позволим вам собрать средства для ваших идей в два счёта или найти проект,
                        который принесет вам прибыль.
                    </p>
                    <div className={styles.main_links}>
                        <Link href="/create-project">Создать сбор</Link>
                        <Link href="/projects">Я - меценат!</Link>
                        {/* Пример третьей кнопки как на скрине */}
                        <Link href="/projects" style={{background: '#934f4f'}}>Выгодно вложиться</Link>
                    </div>
                </div>
            </main>

            <section className={styles.section} id="features">
                <h2 className={styles.section_title}>Инвестируйте, создавайте, просматривайте</h2>
                <p className={styles.section_description}>
                    Наша платформа предлагает уникальные возможности для всех участников рынка.
                </p>
                <div className={styles.section_cards}>
                    {/* Карточка 1 */}
                    <div>
                        {/* Красим иконку в бледно-коричневый */}
                        <Coins size={64} color="#DB935B" strokeWidth={1.5} />
                        <h3>Инвестирование</h3>
                        <p>Получите доступ к разнообразным инвестиционным возможностям в перспективные проекты.</p>
                    </div>
                    {/* Карточка 2 */}
                    <div>
                        <Lightbulb size={64} color="#825e9c" strokeWidth={1.5} />
                        <h3>Создание</h3>
                        <p>Запустите свой проект, привлеките финансирование и найдите партнеров.</p>
                    </div>
                    {/* Карточка 3 */}
                    <div>
                        <FileText size={64} color="#934f4f" strokeWidth={1.5} />
                        <h3>Каталог</h3>
                        <p>Изучите успешные кейсы, вдохновитесь идеями и отслеживайте тренды рынка.</p>
                    </div>
                </div>
            </section>

            <ProjectsSection/>

            <section className={`${styles.section} ${styles.section_choose_way}`}>
                <h2 className={styles.section_title}>ВЫБЕРИТЕ СВОЙ ПУТЬ</h2>
                <p className={styles.section_description} style={{color: '#ccc'}}>
                    Кем вы хотите быть сегодня?
                </p>
                <div className={styles.section_cards}>
                    <Link href="/projects?type=invest">
                        <p>инвестировать / пожертвовать</p>
                    </Link>
                    <Link href="/create-project">
                        <p>создать проект</p>
                    </Link>
                </div>
            </section>

            <Footer/>
        </>
    );
}