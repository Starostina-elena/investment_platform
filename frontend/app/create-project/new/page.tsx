'use client'
import styles from './page.module.css';
import {useEffect, useState} from 'react';
import {Project, PublishProject} from "@/api/project";
import ProjectPreview from "@/app/components/project-preview";
import GenericInfo from "@/app/create-project/new/generic-info";
import NavBar from "@/app/create-project/new/nav-bar";
import ProjectView from "@/app/components/project-view";
import {Message} from "@/api/api";
import MessageComponent from "@/app/components/message";
import Spinner from "@/app/components/spinner";
import {GetMyOrganisations, Organisation} from "@/api/organisation";
import OrgSelector from "@/app/create-project/new/org-selector";


const STEPS = [
    { title: 'Автор', component: OrgSelector },
    { title: 'Основное', component: GenericInfo },
]

const CreateProjectPage = () => {
    const [route, setRoute] = useState<number>(0);
    const [isLoadingOrgs, setIsLoadingOrgs] = useState(true);
    const [userOrgs, setUserOrgs] = useState<Organisation[]>([]);

    const [project, setProject] = useState<Project>({
        created_at: "",
        id: -1,
        creator_id: -1, // Инициализируем как -1, чтобы заставить пользователя выбрать
        name: '',
        quick_peek: '',
        quick_peek_picture_path: "",
        content: '',
        duration_days: 30,
        wanted_money: 100000,
        current_money: 0,
        is_public: true,
        is_completed: false,
        monetization_type: "charity",
        quickPeekPictureFile: null,
        is_banned: false,
        percent: 0,
    });

    const [response, setResponse] = useState<Message | null>(null)
    const [requestSent, setRequestSent] = useState<boolean>(false);

    // Загружаем организации пользователя при маунте
    useEffect(() => {
        GetMyOrganisations()
            .then(data => {
                setUserOrgs(data);
                // Если у пользователя всего одна организация, выбираем её автоматически
                if (data.length === 1) {
                    setProject(prev => ({...prev, creator_id: data[0].id}));
                }
            })
            .finally(() => setIsLoadingOrgs(false));
    }, []);

    if (isLoadingOrgs) {
        return <div className={styles.page_container}><Spinner/></div>;
    }

    return (
        <>
            <NavBar routes={STEPS.map(e => e.title)}
                    currentRoute={route < STEPS.length ? STEPS[route].title : 'preview'}
                    onPreviewClick={() => setRoute(STEPS.length)}
                    onRouteClick={rt => {setRoute(STEPS.findIndex(e => e.title === rt))}}/>

            <div className={styles.page_container}>
                {/* Левая колонка: форма и кнопки */}
                <div className={styles.form_container}>
                    {/* Всегда рендерим все компоненты, скрываем ненужные */}
                    <div style={{ display: route === 0 ? 'block' : 'none' }}>
                        <OrgSelector
                            project={project}
                            setProject={setProject}
                            userOrgs={userOrgs}
                        />
                    </div>
                    <div style={{ display: route === 1 ? 'block' : 'none' }}>
                        <GenericInfo
                            project={project}
                            setProject={setProject}
                        />
                    </div>
                    <div style={{ display: route >= STEPS.length ? 'block' : 'none' }}>
                        <ProjectView
                            project={project}
                            setProject={setProject}
                            userOrgs={userOrgs}
                        />
                    </div>

                    {/* Кнопки навигации */}
                    <div className={styles.controls}>
                        {route > 0 &&
                            <button onClick={() => setRoute(route - 1)}
                                    className={styles.controls_btn + ' ' + styles.controls_btn_prev}>
                                Назад
                            </button>}

                        {route < STEPS.length && (
                            <button onClick={e => {
                                if (e.currentTarget.form?.reportValidity())
                                    setRoute(route + 1)
                                e.preventDefault()
                            }}
                                    className={styles.controls_btn + ' ' + styles.controls_btn_next}>
                                Вперёд
                            </button>
                        )}

                        {/* Кнопка "Опубликовать" (на шаге Превью) */}
                        {route == STEPS.length && (
                            <button disabled={requestSent} onClick={e => {
                                // Валидация перед отправкой
                                if (project.creator_id <= 0) {
                                    setResponse({isError: true, message: "Не выбрана организация"});
                                    return;
                                }
                                setRequestSent(true);
                                PublishProject(project, (msg) => {
                                    setResponse(msg);
                                    setRequestSent(false);
                                })
                                e.preventDefault()
                            }}
                                    className={styles.controls_btn + ' ' + styles.controls_btn_next}>
                                {requestSent && <Spinner size={30} style={{margin: "-11px 0 -11px -32px", paddingRight: "32px"}}/>}
                                Опубликовать
                            </button>
                        )}

                        <MessageComponent message={response}/>
                    </div>
                </div>

                {/* Правая колонка: превью */}
                {route > 0 && route < STEPS.length && <ProjectPreview project={project}/>}
            </div>
        </>
    );
};

export default CreateProjectPage;