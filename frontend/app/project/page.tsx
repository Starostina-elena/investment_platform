'use client';
import ProjectView from "@/app/components/project-view";
import React, {Suspense, useEffect, useState} from "react";
import {useSearchParams} from "next/navigation";
import {GetProjectById, Project} from "@/api/project";
import Spinner from "@/app/components/spinner";
import Link from "next/link";
import styles from "./page.module.css"
import CommentsSection from "@/app/components/comments-section"; // Импорт

const TABS = ["Проект", "Новости", "Комментарии", "Спонсоры"]

export default function Page(){
    return (
        <Suspense>
            <Page_unwrapped />
        </Suspense>
    )
}

function Page_unwrapped() {
    const [project, setProject] = useState<Project | undefined | null>(undefined)
    const params = useSearchParams()

    const [activeTab, setActiveTab] = useState<number>(0)

    useEffect(() => {
        const id = params.get('id')
        if (id == null) {
            setProject(null)
            console.log("Null Id")
        }
        else {
            GetProjectById(+id).then(setProject)
        }
    }, []);

    if (project === undefined)
        return (
            <div className={styles.main}>
                <Spinner/>
            </div>
        )

    if (project === null)
        return (
            <div className={styles.main}>
                <h1>К сожалению данный проект не найден</h1>
                <Link href="/" className={styles.main_link}>На главную</Link>
            </div>
        )

    return (
        <div className={styles.main}>
            <ProjectView project={project}/>
            <div className={styles.tabs}>
                {TABS.map((e, i) => (
                    <p key={i} className={i == activeTab ? styles.active : undefined}
                       onClick={() => setActiveTab(i)}>{e}</p>
                ))}
                <Link href='/create-project#faq'>F.A.Q.</Link>
            </div>
            <div className={styles.content}>
                {/* Вкладка Описание */}
                {activeTab == 0 &&
                    <div style={{color: 'white', whiteSpace: 'pre-wrap'}}>
                        {project.content}
                    </div>
                }

                {/* Вкладка Новости (Заглушка) */}
                {activeTab == 1 &&
                    <p style={{color: '#999', textAlign: 'center'}}>Новостей пока нет.</p>
                }

                {/* Вкладка Комментарии */}
                {activeTab == 2 &&
                    <CommentsSection projectId={project.id} />
                }

                {/* Вкладка Спонсоры (Заглушка) */}
                {activeTab == 3 &&
                    <p style={{color: '#999', textAlign: 'center'}}>Список спонсоров скрыт.</p>
                }
            </div>
        </div>
    )
}