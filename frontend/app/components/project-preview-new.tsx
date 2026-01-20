import styles from './project-preview-new.module.css'
import {Project} from "@/api/project";
import Image from "next/image";
import image_bg from "@/public/image_bg.png";
import {addBasePath} from "next/dist/client/add-base-path";
import {useEffect, useRef} from "react";
import {useImage} from "@/hooks/use-image";
import {BUCKETS} from "@/lib/config";

export default function ProjectPreviewNew({project}: { project: Project }) {
    const imagePreview = useRef<HTMLImageElement>(null);

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
            <div className={styles.preview_image}>
                <Image src={image_bg} alt={"Фото проекта " + project.name} fill={true} ref={imagePreview}/>
            </div>
            <h3 className={styles.preview_name}>{project.name}</h3>
            <Image alt="image" src={imageSrc} fill={true} style={{ borderRadius: '8px' }} />
        </div>
    )
}