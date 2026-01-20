// frontend/components/project-preview.tsx
'use client'
import styles from './project-preview.module.css'
import { Project } from "@/api/project";
import image_bg from "@/public/image_bg.png";
import Image from "next/image";
import { useImage } from "@/hooks/use-image";
import { BUCKETS } from "@/lib/config";

export default function ProjectPreview({project}: {project: Project}) {
    const imageSrc = useImage(
        project.quickPeekPictureFile,
        project.quick_peek_picture_path,
        BUCKETS.PROJECTS,
        image_bg
    );

    const progress = project.wanted_money > 0
        ? Math.min(100, (project.current_money / project.wanted_money) * 100)
        : 0;

    return (
        <div className={styles.preview}>
            <div className={styles.image_container}>
                <Image alt="image" src={imageSrc} fill={true} style={{ borderRadius: '8px' }} />
            </div>

            <div className={styles.project_info}>
                <p className={styles.project_name}>{project.name || 'Безымянный'}</p>
                {/* В будущем здесь можно вывести название организации по creator_id */}
                <a href={"/organisation/" +project.creator_id}>ID Организации: {project.creator_id}</a>
                <br/>
                <p>{project.quick_peek}</p>
            </div>

            <footer className={styles.footer}>
                <div className={styles.progress_bar}>
                    <div className={styles.progress_fill} style={{width: `${progress}%`}}/>
                </div>

                <div className={styles.stats_container}>
                    <div className={styles.stat_block}>
                        прогресс
                        <div className={styles.stat_value}>
                            {progress | 0}%
                        </div>
                    </div>

                    <div className={styles.stat_block}>
                        собрано
                        <div className={styles.stat_value}>
                            {project.current_money}<span>₽</span>
                        </div>
                    </div>

                    <div className={styles.stat_block}>
                        осталось
                        <div className={styles.stat_value}>{project.duration_days} дней</div>
                    </div>
                </div>
            </footer>
        </div>
    )
}