'use client'
import styles from './project-view.module.css';
import Image from "next/image";
import {Project, SetProjectCompleted, SetProjectPublic, UpdateProject} from "@/api/project";
import avatar from "@/public/temp/avatar.png";
import rocket from "@/public/rocket.svg";
import image_bg from "@/public/image_bg.png";
import {useEffect, useRef, useState} from "react";
import {useImage} from "@/hooks/use-image";
import {BUCKETS} from "@/lib/config";
import InvestModal from "@/app/components/invest-modal";
import BannedBanner from "@/app/components/banned-banner";
import AdminBanControl from "@/app/components/admin-ban-control";
import {GetUserById, User} from "@/api/user";
import {toast} from "sonner";
import { Button } from "@/app/components/ui/button";
import { Edit3, CheckCircle, Eye, EyeOff } from "lucide-react";
import MessageComponent from "@/app/components/message";
import { Message } from "@/api/api";

// Функция для маппинга типа монетизации
const getMonetizationLabel = (type: string, percent?: number) => {
    switch (type) {
        case 'charity': return "Безвозмездная поддержка (Благотворительность)";
        case 'fixed_percent': return `Инвестиции: Фиксированный доход ${percent}%`;
        case 'time_percent': return `Инвестиции: ${percent}% годовых`;
        default: return "Смешанный тип финансирования";
    }
}

interface ProjectViewProps {
    project: Project;
    setProject?: (p: Project) => void;
    userOrgs?: any[];
    children?: React.ReactNode;
}

export default function ProjectView({project: initialProject, setProject: externalSetProject, userOrgs}: ProjectViewProps) {
    const [project, setProject] = useState(initialProject);
    const [author, setAuthor] = useState<User | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [editData, setEditData] = useState(initialProject);
    const [message, setMessage] = useState<Message | null>(null);

    // Синхронизация при изменении initialProject
    useEffect(() => {
        setProject(initialProject);
        setEditData(initialProject);
    }, [initialProject]);

    // Загружаем данные автора (организации/юзера) для кнопки "Написать"
    // В текущей модели creator_id - это Organisation ID.
    // Если нужно получить email владельца, надо дернуть /api/org/{id}
    // Для упрощения пока оставим mailto пустым или сделаем заглушку.

    const imageSrc = useImage(
        project.quickPeekPictureFile,
        project.quick_peek_picture_path,
        BUCKETS.PROJECTS,
        image_bg
    );

    const handleShare = (platform: string) => {
        const url = window.location.href;
        const text = `Поддержите проект "${project.name}" на Sipis!`;
        let shareUrl = "";

        switch(platform) {
            case 'telegram': shareUrl = `https://t.me/share/url?url=${url}&text=${text}`; break;
            case 'vk': shareUrl = `https://vk.com/share.php?url=${url}`; break;
            case 'whatsapp': shareUrl = `https://api.whatsapp.com/send?text=${text} ${url}`; break;
            case 'ok': shareUrl = `https://connect.ok.ru/offer?url=${url}`; break;
        }

        if (shareUrl) window.open(shareUrl, '_blank');
    };

    return (
        <div className={styles.project_card}>
            {project.is_banned && <div className="mb-6"><BannedBanner type="Проект" /></div>}

            <div className={styles.content_wrapper}>
                <div className={styles.project_header}>
                    <div className={styles.media_container}>
                        <div className={styles.project_image_wrapper}>
                            <Image fill={true} alt="Фото проекта" className={styles.project_image} src={imageSrc}/>
                        </div>
                        <div className={styles.project_meta}>
                            <div className={styles.category}>{getMonetizationLabel(project.monetization_type, project.percent)}</div>
                        </div>
                    </div>
                    <div className={styles.project_details}>
                        <h1 className={styles.project_title}>{project.name}</h1>
                        <div className={styles.project_description}>{project.quick_peek}</div>

                        <div className={styles.author_info}>
                            <div className={styles.author}>
                                <div className={styles.author_avatar}>
                                    <Image src={avatar} alt="Аватарка" className={styles.avatar_image} />
                                </div>
                                <div className={styles.author_details}>
                                    {/* Ссылка на страницу организации */}
                                    <div className={styles.publisher}>ID Организации: {project.creator_id}</div>
                                    <a href={`mailto:support@Sipis.ru?subject=Вопрос по проекту ${project.id}`}
                                       className={styles.contact_link}>
                                        Написать автору
                                    </a>
                                </div>
                            </div>
                        </div>

                        {/* Исправленный тип инвестиций */}
                        <div className={styles.investment_type} style={{
                            background: 'rgba(0, 208, 255, 0.1)',
                            padding: '10px',
                            borderRadius: '8px',
                            color: '#00D0FF',
                            fontWeight: 'bold',
                            marginTop: '1rem'
                        }}>
                            {getMonetizationLabel(project.monetization_type, project.percent)}
                        </div>

                        <div className={styles.funding_stats}>
                            <div className={styles.funding_amount}>
                                <div className={styles.amount}>{project.current_money} ₽</div>
                                <div className={styles.goal}>из {project.wanted_money} собрано</div>
                            </div>
                            <div className={styles.backers_info}>
                                <div className={styles.time_left}>
                                    {/* Простая логика дней, в реальности нужно считать разницу дат */}
                                    Осталось {project.duration_days} дн.
                                </div>
                            </div>
                        </div>

                        <div className={styles.progress_container}>
                            <div className={styles.progress_bar}>
                                <div className={styles.progress_fill} style={{width: `${Math.min(100, project.current_money / project.wanted_money * 100)}%`}}/>
                            </div>
                        </div>

                        <div className={styles.actions}>
                            {/* Рабочая кнопка инвестирования */}
                            {!project.is_completed && !project.is_banned && (
                                <InvestModal
                                    projectId={project.id}
                                    projectName={project.name}
                                    onSuccess={() => window.location.reload()}
                                />
                            )}

                            {/* Рабочие кнопки шаринга */}
                            <div className={styles.social_sharing}>
                                <div className={styles.share_buttons}>
                                    <button className={styles.share_button} onClick={() => handleShare('telegram')}>
                                        <span className={styles.share_icon + ' ' + styles.share_telegram}/>
                                    </button>
                                    <button className={styles.share_button} onClick={() => handleShare('vk')}>
                                        <span className={styles.share_icon + ' ' + styles.share_vk}/>
                                    </button>
                                    <button className={styles.share_button} onClick={() => handleShare('ok')}>
                                        <span className={styles.share_icon + ' ' + styles.share_ok}/>
                                    </button>
                                    <button className={styles.share_button} onClick={() => handleShare('whatsapp')}>
                                        <span className={styles.share_icon + ' ' + styles.share_whatsapp}/>
                                    </button>
                                </div>
                            </div>
                        </div>

                        {/* Кнопки управления проектом (для владельца) */}
                        {project.creator_id && (
                            <div style={{ display: 'flex', gap: '1rem', marginTop: '1.5rem', flexWrap: 'wrap' }}>
                                <Button
                                    onClick={() => setIsEditing(!isEditing)}
                                    style={{
                                        backgroundColor: isEditing ? '#ff6666' : '#825e9c',
                                        color: 'white',
                                        display: 'flex',
                                        alignItems: 'center',
                                        gap: '0.5rem'
                                    }}
                                >
                                    <Edit3 size={16} />
                                    {isEditing ? 'Отмена' : 'Редактировать'}
                                </Button>

                                {!project.is_completed && (
                                    <Button
                                        onClick={() => {
                                            SetProjectCompleted(project.id, true, (msg) => {
                                                setMessage(msg);
                                                if (!msg.isError) {
                                                    setProject({...project, is_completed: true});
                                                }
                                            });
                                        }}
                                        style={{
                                            backgroundColor: '#2ecc71',
                                            color: 'white',
                                            display: 'flex',
                                            alignItems: 'center',
                                            gap: '0.5rem'
                                        }}
                                    >
                                        <CheckCircle size={16} />
                                        Завершить
                                    </Button>
                                )}

                                <Button
                                    onClick={() => {
                                        SetProjectPublic(project.id, !project.is_public, (msg) => {
                                            setMessage(msg);
                                            if (!msg.isError) {
                                                setProject({...project, is_public: !project.is_public});
                                            }
                                        });
                                    }}
                                    style={{
                                        backgroundColor: project.is_public ? '#3498db' : '#95a5a6',
                                        color: 'white',
                                        display: 'flex',
                                        alignItems: 'center',
                                        gap: '0.5rem'
                                    }}
                                >
                                    {project.is_public ? <Eye size={16} /> : <EyeOff size={16} />}
                                    {project.is_public ? 'Публичный' : 'Приватный'}
                                </Button>
                            </div>
                        )}

                        {message && <div style={{ marginTop: '1rem' }}><MessageComponent message={message} /></div>}

                        {/* Форма редактирования */}
                        {isEditing && (
                            <div style={{ marginTop: '1.5rem', padding: '1.5rem', backgroundColor: 'rgba(0,0,0,0.2)', borderRadius: '8px' }}>
                                <h3 style={{ color: 'white', marginBottom: '1rem', fontWeight: 'bold' }}>Редактирование проекта</h3>
                                <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                                    <div>
                                        <label style={{ color: '#ccc', display: 'block', marginBottom: '0.5rem' }}>Название</label>
                                        <input
                                            type="text"
                                            value={editData.name}
                                            onChange={(e) => setEditData({...editData, name: e.target.value})}
                                            style={{
                                                width: '100%',
                                                padding: '0.75rem',
                                                backgroundColor: '#333',
                                                color: 'white',
                                                border: '1px solid #555',
                                                borderRadius: '4px'
                                            }}
                                        />
                                    </div>
                                    <div>
                                        <label style={{ color: '#ccc', display: 'block', marginBottom: '0.5rem' }}>Короткое описание</label>
                                        <textarea
                                            value={editData.quick_peek}
                                            onChange={(e) => setEditData({...editData, quick_peek: e.target.value})}
                                            style={{
                                                width: '100%',
                                                padding: '0.75rem',
                                                backgroundColor: '#333',
                                                color: 'white',
                                                border: '1px solid #555',
                                                borderRadius: '4px',
                                                minHeight: '80px',
                                                fontFamily: 'inherit'
                                            }}
                                        />
                                    </div>
                                    <div>
                                        <label style={{ color: '#ccc', display: 'block', marginBottom: '0.5rem' }}>Полное описание</label>
                                        <textarea
                                            value={editData.content}
                                            onChange={(e) => setEditData({...editData, content: e.target.value})}
                                            style={{
                                                width: '100%',
                                                padding: '0.75rem',
                                                backgroundColor: '#333',
                                                color: 'white',
                                                border: '1px solid #555',
                                                borderRadius: '4px',
                                                minHeight: '120px',
                                                fontFamily: 'inherit'
                                            }}
                                        />
                                    </div>
                                    <div style={{ display: 'flex', gap: '1rem' }}>
                                        <Button
                                            onClick={() => {
                                                UpdateProject(project.id, editData, (msg) => {
                                                    setMessage(msg);
                                                    if (!msg.isError) {
                                                        setProject(editData);
                                                        setIsEditing(false);
                                                    }
                                                });
                                            }}
                                            style={{ backgroundColor: '#2ecc71', color: 'white' }}
                                        >
                                            Сохранить
                                        </Button>
                                        <Button
                                            onClick={() => setIsEditing(false)}
                                            variant="outline"
                                            style={{ color: '#ff6666', borderColor: '#ff6666' }}
                                        >
                                            Отмена
                                        </Button>
                                    </div>
                                </div>
                            </div>
                        )}

                        <AdminBanControl
                            entityType="project"
                            entityId={project.id}
                            isBanned={project.is_banned}
                            onUpdate={(banned) => setProject({...project, is_banned: banned})}
                        />
                    </div>
                </div>
            </div>
        </div>
    )
}