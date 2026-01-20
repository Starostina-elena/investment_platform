import styles from './project-type-chooser.module.css';
import Link from "next/link";

interface ProjectType {
    name: string
    description: string[]
    path: string
}

const PROJECTS: ProjectType[] = [
    {
        name: 'Благотворительный сбор',
        description: [
            'Сумма сбора не ограничена – можно получить сколько угодно, больше 100%.',
            'Можно забрать деньги, если собрано меньше 100% суммы.',
            'Проект завершается по таймеру: максимальный срок проекта – 60 дней Можно продлевать.'
        ],
        path: 'donate'
    },
    {
        name: 'Инвестиции',
        description: [
            'Сумма сбора не ограничена – можно получить сколько угодно, больше 100%.',
            'Можно забрать деньги, если собрано меньше 100% суммы.',
            'Проект завершается по таймеру: максимальный срок проекта – 60 дней Можно продлевать.'
        ],
        path: 'invest'
    }
]

const ProjectTypeChooser = () => {
    return (
        <div className={styles.main_container}>
            <div className={styles.content_wrapper}>
                <div className={styles.heading} id="type_chooser">Выберите вид проекта</div>
                <div className={styles.project_options_container}>
                    {PROJECTS.map(e => (
                        <div className={styles.project_card_container} key={e.name}>
                            <Link href={"/create-project/new?type=" + e.path}
                               className={styles.project_card}>
                                <div className={styles.card_title}>{e.name}</div>

                                {e.description.map((e,i) => (
                                    <div className={styles.card_description} key={i}>
                                        {e}
                                    </div>
                                ))}

                                <button className={styles.primary_button}>Выбрать</button>
                            </Link>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default ProjectTypeChooser;