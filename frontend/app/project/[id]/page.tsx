'use client'

import { useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { GetProjectById, Project } from '@/api/project';
import ProjectView from '@/app/components/project-view';
import Spinner from '@/app/components/spinner';

export default function ProjectPage() {
    const params = useParams();
    const [project, setProject] = useState<Project | null | undefined>(undefined);

    useEffect(() => {
        if (params.id) {
            GetProjectById(+params.id).then(setProject);
        }
    }, [params.id]);

    if (project === undefined) {
        return (
            <div style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', backgroundColor: '#989694' }}>
                <Spinner />
            </div>
        );
    }

    if (project === null) {
        return (
            <div style={{ minHeight: '100vh', display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', backgroundColor: '#989694', color: 'white' }}>
                <h1 style={{ fontSize: '2rem', fontWeight: 'bold', marginBottom: '1rem' }}>Проект не найден</h1>
                <p style={{ color: '#ccc' }}>Проект с таким ID не существует или был удалён.</p>
            </div>
        );
    }

    return (
        <div style={{ minHeight: '100vh', backgroundColor: '#989694', padding: '2rem' }}>
            <ProjectView project={project} setProject={setProject} />
        </div>
    );
}
