'use client'
import styles from './project-preview-new.module.css'
import {Project} from "@/api/project";
import Image from "next/image";
import image_bg from "@/public/image_bg.png"; // Убедитесь, что эта картинка есть
import {addBasePath} from "next/dist/client/add-base-path";
import {useImage} from "@/hooks/use-image";
import {BUCKETS} from "@/lib/config";

export default function ProjectPreviewNew({project}: {project: Project}) {
    // Используем хук для загрузки картинки (или заглушки, если нет)
    const imageSrc = useImage(
        project.quickPeekPictureFile,
        project.quick_peek_picture_path,
        BUCKETS.PROJECTS,
        image_bg
    );

    return (
        <div className={styles.preview} onClick={() => {
            window.location.href = addBasePath('/project?id=' + project.id)
        }}>
            {/* Обертка для картинки */}
            <div className={styles.preview_image}>
                <Image
                    src={imageSrc}
                    alt={project.name}
                    fill={true}
                    // unoptimized // Можно добавить, если картинки с внешнего MinIO не грузятся
                    className={styles.image_element} // Класс для object-fit
                />
            </div>

            <div className={styles.preview_content}>
                <h3 className={styles.preview_name}>{project.name}</h3>
                <p className={styles.preview_description}>{project.quick_peek}</p>

                {/* Можно добавить прогресс бар для красоты */}
                <div className={styles.progress_bar}>
                    <div
                        className={styles.progress_fill}
                        style={{width: `${Math.min(100, (project.current_money / project.wanted_money) * 100)}%`}}
                    />
                </div>
                <div className={styles.money_info}>
                    <span className="text-[#DB935B] font-bold">{project.current_money} ₽</span> из {project.wanted_money} ₽
                </div>
            </div>
        </div>
    )
}