'use client'

import { useState, useEffect } from 'react';
import { GetAllOrganisationProjects, Project } from '@/api/project';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import { BUCKETS } from '@/lib/config';
import Spinner from './spinner';
import placeholder from '@/public/image_bg.png';

export default function OrganisationProjects({ orgId }: { orgId: number }) {
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        GetAllOrganisationProjects(orgId)
            .then(setProjects)
            .finally(() => setLoading(false));
    }, [orgId]);

    if (loading) return <Spinner />;

    return (
        <div style={{ display: 'grid', gap: '1.5rem' }}>
            {projects.length === 0 ? (
                <p style={{ color: '#ccc', fontSize: '0.95rem' }}>Нет проектов</p>
            ) : (
                projects.map(project => {
                    const imageSrc = project.quick_peek_picture_path 
                        ? `${BUCKETS.PROJECTS}/${project.quick_peek_picture_path}`
                        : placeholder;
                    
                    return (
                        <div
                            key={project.id}
                            onClick={() => router.push(`/project/${project.id}`)}
                            style={{
                                display: 'grid',
                                gridTemplateColumns: '120px 1fr',
                                gap: '1.5rem',
                                padding: '1.5rem',
                                backgroundColor: '#555652',
                                borderRadius: '8px',
                                cursor: 'pointer',
                                transition: 'transform 0.2s, box-shadow 0.2s',
                            }}
                            onMouseEnter={(e) => {
                                (e.currentTarget as HTMLElement).style.transform = 'translateY(-4px)';
                                (e.currentTarget as HTMLElement).style.boxShadow = '0 8px 16px rgba(0,0,0,0.3)';
                            }}
                            onMouseLeave={(e) => {
                                (e.currentTarget as HTMLElement).style.transform = 'translateY(0)';
                                (e.currentTarget as HTMLElement).style.boxShadow = 'none';
                            }}
                        >
                            {/* Изображение */}
                            <div style={{ position: 'relative', borderRadius: '8px', overflow: 'hidden', height: '100px' }}>
                                <Image
                                    src={imageSrc}
                                    alt={project.name}
                                    fill
                                    style={{ objectFit: 'cover' }}
                                />
                            </div>

                            {/* Информация */}
                            <div style={{ display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
                                <div>
                                    <h3 style={{ fontSize: '1.1rem', fontWeight: 'bold', color: 'white', marginBottom: '0.5rem' }}>
                                        {project.name}
                                    </h3>
                                    <p style={{ color: '#ccc', fontSize: '0.9rem', marginBottom: '0.75rem', lineHeight: '1.4' }}>
                                        {project.quick_peek}
                                    </p>
                                </div>

                                {/* Статус и метаинформация */}
                                <div style={{ display: 'flex', gap: '1rem', fontSize: '0.85rem', color: '#aaa' }}>
                                    <span>
                                        {project.is_completed ? '✓ Завершён' : `${Math.round((project.current_money / project.wanted_money) * 100)}% собрано`}
                                    </span>
                                    <span>{project.current_money.toLocaleString()} / {project.wanted_money.toLocaleString()} ₽</span>
                                    {project.is_public && <span style={{ color: '#00aef4' }}>Публичный</span>}
                                    {project.is_banned && <span style={{ color: '#ff6666' }}>Заблокирован</span>}
                                </div>
                            </div>
                        </div>
                    );
                })
            )}
        </div>
    );
}
