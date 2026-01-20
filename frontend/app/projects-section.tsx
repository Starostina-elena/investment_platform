'use client';
import styles from "@/app/page.module.css";
import {useEffect, useState} from "react";
import {GetProjects, Project} from "@/api/project";
import ProjectPreviewNew from "@/app/components/project-preview-new";


export default function ProjectsSection(){
    const [projects, setProjects] = useState<Project[]>([])

    useEffect(() => {
        GetProjects().then(setProjects);
    }, []);

    return (
        <section className={styles.section} id="projects">
            <h2 className={styles.section_title + ' ' + styles.projects_title}>лучшие сегодня!</h2>
            <div className={styles.projects_container}>
                {projects.map((project,i ) => (
                    <ProjectPreviewNew project={project} key={i}/>
                ))}
            </div>
        </section>
    )
}