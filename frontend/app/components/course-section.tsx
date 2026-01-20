import styles from './course-section.module.css';
import {APP_NAME} from "@/app/globals";
import course_bg from "@/public/course_bg.png"

const CourseSection = () => {
    return (
        <>
            <div className={styles.main_container}>
                <div className={styles.left_section}>
                    <div className={styles.heading_primary}>
                        Sipis — место, где вы можете
                        <br/>
                        объединить людей вокруг вашей идеи.
                    </div>

                    <div className={styles.subheading}>
                        Пройдите онлайн-курс из 35 простых видеоуроков длительностью по 3-7 минут и успешно привлекайте
                        деньги через Sipis.
                    </div>

                    <div
                        className={styles.promo_block}
                        style={{backgroundImage: `url(${course_bg.src})`}}
                    >
                        <div className={styles.promo_text}>
                            Бесплатный обучающий курс
                            <br/>
                            «Мастер краудфандинга»
                        </div>
                    </div>
                </div>

                <div className={styles.stats}>
                    <div className={styles.stat_block}>
                        <div className={styles.stat_number}>N</div>
                        <div className={styles.stat_description}>
                            <span className={styles.highlighted_text}>успешных проектов</span>
                            на Sipis
                        </div>
                    </div>

                    <div className={styles.stat_block}>
                        <div className={styles.stat_number}>N</div>
                        <div className={styles.stat_description}>
                            <span className={styles.highlighted_text}>Миллион рублей</span>
                            привлечено проектами
                        </div>
                    </div>
                </div>
            </div>
            <div className={styles.main_container_categories}>
                <div className={styles.content_wrapper}>
                    <div className={styles.left_section}>
                        <div className={styles.section_heading}>Категории проектов на {APP_NAME}</div>

                        <div className={styles.description}>
                            Вы всегда сможете найти подходящую вашему проекту категорию:
                            <br/>
                            <strong className={styles.highlight}>музыка, фильмы, искусство, технологии, дизайн, еда,
                                издательское дело, мода, настольные игры и многие другие.</strong>
                        </div>
                    </div>

                    <div className={styles.info_card}>
                        <div className={styles.card_title}>Есть идея?</div>

                        <div className={styles.card_description}>Создайте черновик проекта прямо сейчас!</div>

                        <button className={styles.primary_button}>Создать проект</button>
                    </div>
                </div>
            </div>
            <div className={styles.hero_section}>
                <div className={styles.content_wrapper_video}>
                    <div className={styles.section_heading_video}>
                        Начало — это эффективный инструмент привлечения финансирования
                    </div>

                    <div className={styles.section_description}>
                        Предложите людям интересную идею или что-то полезное.
                        <br/>
                        Расскажите об этом всем, протестируйте спрос на свой продукт через ранние предзаказы.
                        <br/>
                        Получите деньги и измените мир прямо сейчас!
                    </div>

                    <button className={styles.video_button}>Посмотреть видео</button>
                    <button className={styles.scroll_button} onClick={() => {
                        window.scrollTo(0, (document.querySelector('#type_chooser') as HTMLDivElement).offsetHeight || 0)
                    }}/>
                </div>
            </div>
        </>
    );
};

export default CourseSection;