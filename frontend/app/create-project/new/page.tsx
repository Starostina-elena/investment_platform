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
        content: 'Полное описание проекта...',
        category: '',
        location: '',
        duration_days: 30,
        wanted_money: 100000,
        current_money: 0,
        is_public: true,
        is_completed: false,
        monetization_type: "donation",
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

    const CurrentRoute = route < STEPS.length ? STEPS[route].component : ProjectView;

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
                {/* Передаем userOrgs в компонент шага (понадобится только для OrgSelector) */}
                <CurrentRoute
                    project={project}
                    setProject={setProject}
                    userOrgs={userOrgs}
                >
                    <div className={styles.controls}>
                        {route > 0 &&
                            <button onClick={() => setRoute(route - 1)}
                                    className={styles.controls_btn + ' ' + styles.controls_btn_prev}>
                                Назад
                            </button>}

                        {route < STEPS.length - 1 &&
                            <button onClick={e => {
                                // Валидация формы перед переходом
                                // Для OrgSelector проверка creator_id
                                if (route === 0 && project.creator_id <= 0) {
                                    alert("Пожалуйста, выберите организацию");
                                    return;
                                }

                                // Для остальных форм (GenericInfo)
                                const form = document.querySelector('form');
                                if (form && !form.checkValidity()) {
                                    form.reportValidity();
                                    return;
                                }

                                setRoute(route + 1)
                                e.preventDefault()
                            }}
                                    className={styles.controls_btn + ' ' + styles.controls_btn_next}>
                                Вперёд
                            </button>}

                        <MessageComponent message={response}/>

                        {/* Кнопка публикации (на шаге превью или последнем шаге) */}
                        {route == STEPS.length &&
                            <button disabled={requestSent} onClick={e => {
                                if (project.creator_id <= 0) {
                                    setResponse({isError: true, message: "Не выбрана организация-автор"});
                                    return;
                                }
                                setRequestSent(true);
                                PublishProject(project, (msg) => {
                                    setResponse(msg);
                                    setRequestSent(false);
                                    if (!msg.isError) {
                                        // Опционально: редирект на созданный проект через пару секунд
                                        // setTimeout(() => router.push(`/project?id=${newId}`), 2000);
                                    }
                                })
                                e.preventDefault()
                            }}
                                    className={styles.controls_btn + ' ' + styles.controls_btn_next}>
                                {requestSent && <Spinner size={30} style={{margin: "-11px 0 -11px -32px", paddingRight: "32px"}}/>}
                                Опубликовать
                            </button>}
                    </div>
                </CurrentRoute>

                {/* Превью справа (показываем всегда, кроме мобилок) */}
                {/* Скрываем превью на первом шаге выбора орги, чтобы не мешало, или оставляем */}
                {route > 0 && route < STEPS.length && <ProjectPreview project={project}/>}
            </div>
        </>
    );
};

export default CreateProjectPage;