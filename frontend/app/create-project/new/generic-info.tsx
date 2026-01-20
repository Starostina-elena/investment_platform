// frontend/app/create-project/new/generic-info.tsx
import { Project } from "@/api/project";
import styles from "./generic-info.module.css";
import Image from "next/image";
import image_bg from "@/public/image_bg.png";
import { Fragment } from "react";
import { useImage } from "@/hooks/use-image";
import { BUCKETS } from "@/lib/config";
import {CATEGORIES} from "@/globals";

export default function GenericInfo({project, setProject, children}: {
    project: Project,
    setProject: (project: Project) => void,
    children: React.ReactNode
}) {
    const coverImageSrc = useImage(
        project.quickPeekPictureFile,
        project.quick_peek_picture_path,
        BUCKETS.PROJECTS,
        image_bg
    );

    return (
        <form>
            <h1 className={styles.page_title}>Основное</h1>
            <p className={styles.welcome_text}>
                Поздравляем! Вы начали новый проект!
            </p>
            <div className={styles.main_content}>
                <div className={styles.form_group}>
                    <label htmlFor="project-name" className={styles.label}>
                        Название проекта
                        <span className={styles.hint}>
                            Должно быть простым, запоминающимся и отражать суть вашего проекта.
                        </span>
                    </label>
                    <input
                        id="project-name"
                        className={styles.input_field}
                        value={project.name} required
                        onChange={(e) => setProject({...project, name: e.target.value})}
                        maxLength={50} minLength={5}
                    />
                    <span className={styles.char_count}>{`${project.name.length}/50`}</span>
                </div>

                <div className={styles.form_group}>
                    <div className={styles.label}>Обложка</div>
                    <label className={styles.cover_wrapper}>
                        <div className={styles.upload_img_placeholder}>
                            <Image src={coverImageSrc} alt="placeholder" fill={true} />
                        </div>
                        <div className={styles.upload_btn}>
                            Загрузить изображение
                            <input
                                type="file"
                                accept=".jpg,.jpeg,.png,.webp,.bmp"
                                onChange={(e) => {
                                    const file = e.currentTarget.files?.[0]
                                    if (!file) return;

                                    if (file.size > 5 * 1024 * 1024){
                                        e.currentTarget.setCustomValidity('Максимальный размер 5МБ!');
                                        return;
                                    }
                                    e.currentTarget.setCustomValidity('');

                                    setProject({
                                        ...project,
                                        quickPeekPictureFile: file
                                    })
                                }}
                            />
                        </div>
                        <p className={styles.file_info}>
                            JPEG, PNG или BMP. Макс. 5 МБ.
                        </p>
                    </label>
                </div>

                <div className={styles.form_group}>
                    <label htmlFor="shortDesc" className={styles.label}>
                        Короткое описание
                    </label>
                    <textarea
                        id="shortDesc"
                        className={styles.textarea_field}
                        value={project.quick_peek}
                        onChange={(e) => setProject({...project, quick_peek: e.target.value})}
                        maxLength={100}
                    />
                    <span className={styles.char_count}>{`${project.quick_peek.length}/100`}</span>
                </div>

                {/* Бэкенд пока не хранит category и location, но мы их держим в стейте для UI */}
                <div className={styles.form_group}>
                    <label htmlFor="category" className={styles.label}>Тип финансирования</label>
                    <select
                        id="category" className={styles.select_field}
                        value={project.monetization_type || ""}
                        required
                        onChange={(e) => {
                            const val = e.target.value;
                            setProject({
                                ...project,
                                monetization_type: val,
                                category: CATEGORIES[val] || val
                            })
                        }}
                    >
                        <option value="" disabled>Выберите тип...</option>
                        {Object.entries(CATEGORIES).map(([key, label]) => (
                            <option key={key} value={key}>{label}</option>
                        ))}
                    </select>
                </div>

                {/* Поле Процент (появляется только для инвест проектов) */}
                {(project.monetization_type === 'fixed_percent' || project.monetization_type === 'time_percent') && (
                    <div className={styles.form_group}>
                        <label className={styles.label}>Процентная ставка (%)</label>
                        <input
                            type="number"
                            className={styles.input_field}
                            value={project.percent || ''}
                            onChange={(e) => setProject({...project, percent: parseFloat(e.target.value)})}
                            placeholder="Например: 10"
                            min="0.1" step="0.1" required
                        />
                    </div>
                )}

                <div className={styles.form_group}>
                    <label htmlFor="location" className={styles.label}>Место реализации</label>
                    <select
                        id="location" className={styles.select_field}
                        value={project.location || ""}
                        required
                        onChange={(e) => setProject({...project, location: e.target.value})}
                    >
                    </select>
                </div>

                <div className={styles.form_group}>
                    <label className={styles.label}>Срок (дней)</label>
                    <div className={styles.days_input}>
                        <input
                            type='number' className={styles.number_field}
                            value={project.duration_days}
                            max={60} min={1} required
                            onChange={(e) => setProject({...project, duration_days: +(e.target.value)})}
                        />
                        <span className={styles.days_label}>дн.</span>
                    </div>
                </div>

                <div className={styles.form_group}>
                    <label htmlFor="amount" className={styles.label}>Необходимая сумма</label>
                    <div className={styles.amount_container}>
                        <input
                            id="amount" min={10_000}
                            className={styles.amount_field}
                            type="number"
                            value={project.wanted_money}
                            onChange={(e) => setProject({...project, wanted_money: parseInt(e.target.value)})}
                        />
                        <span className={styles.currency_symbol}>₽</span>
                    </div>
                </div>

                <div className={styles.form_group}>
                    <div className={styles.checkbox_container}>
                        <label className={styles.checkbox_label}>
                            <input
                                type="checkbox"
                                checked={!project.is_public}
                                onChange={() => setProject({...project, is_public: !project.is_public})}
                            /> Приватный проект
                        </label>
                    </div>
                </div>

                {children}
            </div>
        </form>
    )
}